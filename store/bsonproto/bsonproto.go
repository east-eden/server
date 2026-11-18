// Package bsonproto 在 protobuf 消息与 BSON 文档之间做双向转换。
//
// 设计目标：持久化模型统一用 protobuf 定义，一份 schema 同时服务网络协议、
// Redis 缓存和 MongoDB 落库；但 MongoDB 里保留“按字段可读、可建索引”的 BSON
// 文档结构，而不是把整个对象压成一个二进制 blob。
//
// 字段名映射规则：
//   - proto 字段名统一转成 snake_case，例如 TypeId -> type_id；
//   - 顶层主键字段（Id / id）映射为 MongoDB 的 _id；
//   - 嵌套 message 转成子文档，repeated message 转成数组，map 转成文档。
package bsonproto

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// IDField MongoDB 主键字段名
const IDField = "_id"

var (
	ErrNilMessage  = errors.New("bsonproto: nil message")
	ErrTypeMisuse  = errors.New("bsonproto: unexpected bson value type")
	ErrUnsupported = errors.New("bsonproto: unsupported proto field kind")
)

// Marshal 把 proto 消息编码成 BSON 文档
func Marshal(msg proto.Message) (bson.D, error) {
	if msg == nil {
		return nil, ErrNilMessage
	}

	return encodeMessage(msg.ProtoReflect(), true)
}

// Unmarshal 把 BSON 文档解码进 proto 消息
func Unmarshal(doc bson.D, msg proto.Message) error {
	if msg == nil {
		return ErrNilMessage
	}

	return decodeMessage(doc, msg.ProtoReflect(), true)
}

// Key 返回 proto 字段名对应的 BSON 键：顶层主键为 _id，其余为 snake_case
func Key(name protoreflect.Name, topLevel bool) string {
	if topLevel && isIDName(name) {
		return IDField
	}

	return SnakeCase(string(name))
}

// SnakeCase 把 CamelCase 的 proto 字段名转换为 snake_case
func SnakeCase(s string) string {
	var b strings.Builder
	b.Grow(len(s) + 4)

	for i, r := range s {
		if isUpper(r) {
			if i > 0 {
				b.WriteByte('_')
			}

			b.WriteRune(r + ('a' - 'A'))
			continue
		}

		b.WriteRune(r)
	}

	return b.String()
}

func isUpper(r rune) bool { return r >= 'A' && r <= 'Z' }

func isIDName(name protoreflect.Name) bool {
	return name == "Id" || name == "id" || name == "ID"
}

func encodeMessage(m protoreflect.Message, topLevel bool) (bson.D, error) {
	fields := m.Descriptor().Fields()
	doc := make(bson.D, 0, fields.Len())

	var encodeErr error
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		value, err := encodeValue(fd, v)
		if err != nil {
			encodeErr = fmt.Errorf("encode field %q: %w", fd.FullName(), err)
			return false
		}

		doc = append(doc, bson.E{Key: Key(fd.Name(), topLevel), Value: value})
		return true
	})

	if encodeErr != nil {
		return nil, encodeErr
	}

	return doc, nil
}

func encodeValue(fd protoreflect.FieldDescriptor, v protoreflect.Value) (any, error) {
	switch {
	case fd.IsMap():
		sub := make(bson.D, 0, v.Map().Len())

		var encodeErr error
		v.Map().Range(func(k protoreflect.MapKey, mv protoreflect.Value) bool {
			value, err := encodeSingular(fd.MapValue(), mv)
			if err != nil {
				encodeErr = err
				return false
			}

			sub = append(sub, bson.E{Key: k.String(), Value: value})
			return true
		})

		return sub, encodeErr

	case fd.IsList():
		list := v.List()
		arr := make(bson.A, 0, list.Len())

		for i := 0; i < list.Len(); i++ {
			value, err := encodeSingular(fd, list.Get(i))
			if err != nil {
				return nil, err
			}

			arr = append(arr, value)
		}

		return arr, nil

	default:
		return encodeSingular(fd, v)
	}
}

func encodeSingular(fd protoreflect.FieldDescriptor, v protoreflect.Value) (any, error) {
	// map 的 value 描述符不是 field，但 Kind/Message 依然可用
	kind := fd.Kind()

	switch kind {
	case protoreflect.BoolKind:
		return v.Bool(), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return int32(v.Int()), nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return v.Int(), nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return uint32(v.Uint()), nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return v.Uint(), nil
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return v.Float(), nil
	case protoreflect.StringKind:
		return v.String(), nil
	case protoreflect.BytesKind:
		return primitive.Binary{Data: append([]byte(nil), v.Bytes()...)}, nil
	case protoreflect.EnumKind:
		return int32(v.Enum()), nil
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return encodeMessage(v.Message(), false)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupported, kind)
	}
}

func decodeMessage(doc bson.D, m protoreflect.Message, topLevel bool) error {
	fields := m.Descriptor().Fields()

	lookup := make(map[string]protoreflect.FieldDescriptor, fields.Len())
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		lookup[Key(fd.Name(), topLevel)] = fd
	}

	for _, e := range doc {
		fd, ok := lookup[e.Key]
		if !ok {
			// 未知字段（历史遗留或已删除的字段）直接忽略
			continue
		}

		if err := decodeValue(m, fd, e.Value); err != nil {
			return fmt.Errorf("decode field %q: %w", fd.FullName(), err)
		}
	}

	return nil
}

func decodeValue(m protoreflect.Message, fd protoreflect.FieldDescriptor, raw any) error {
	switch {
	case fd.IsList():
		list := m.Mutable(fd).List()

		items, err := asSlice(raw)
		if err != nil {
			return err
		}

		list.Truncate(0)
		for _, item := range items {
			// repeated message 需要先构造元素再递归解码，否则 Append 会拿到零值
			if fd.Kind() == protoreflect.MessageKind || fd.Kind() == protoreflect.GroupKind {
				entries, err := asDocument(item)
				if err != nil {
					return err
				}

				elem := list.NewElement()
				if err := decodeMessage(entries, elem.Message(), false); err != nil {
					return err
				}

				list.Append(elem)
				continue
			}

			value, err := decodeSingular(fd, item)
			if err != nil {
				return err
			}

			list.Append(value)
		}

		return nil

	case fd.IsMap():
		entries, err := asDocument(raw)
		if err != nil {
			return err
		}

		mp := m.Mutable(fd).Map()
		for _, entry := range entries {
			key, err := decodeMapKey(fd.MapKey(), entry.Key)
			if err != nil {
				return fmt.Errorf("decode map key %q: %w", entry.Key, err)
			}

			// map<string, message> 的 value 同样需要构造子消息
			if fd.MapValue().Kind() == protoreflect.MessageKind || fd.MapValue().Kind() == protoreflect.GroupKind {
				entries, err := asDocument(entry.Value)
				if err != nil {
					return err
				}

				elem := mp.NewValue()
				if err := decodeMessage(entries, elem.Message(), false); err != nil {
					return err
				}

				mp.Set(key, elem)
				continue
			}

			value, err := decodeSingular(fd.MapValue(), entry.Value)
			if err != nil {
				return err
			}

			mp.Set(key, value)
		}

		return nil

	default:
		value, err := decodeSingular(fd, raw)
		if err != nil {
			return err
		}

		if fd.Kind() == protoreflect.MessageKind || fd.Kind() == protoreflect.GroupKind {
			sub := m.Mutable(fd).Message()
			entries, err := asDocument(raw)
			if err != nil {
				return err
			}

			return decodeMessage(entries, sub, false)
		}

		m.Set(fd, value)
		return nil
	}
}

func decodeSingular(fd protoreflect.FieldDescriptor, raw any) (protoreflect.Value, error) {
	switch fd.Kind() {
	case protoreflect.BoolKind:
		v, err := asBool(raw)
		return protoreflect.ValueOfBool(v), err

	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		v, err := asInt64(raw)
		return protoreflect.ValueOfInt32(int32(v)), err

	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		v, err := asInt64(raw)
		return protoreflect.ValueOfInt64(v), err

	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		v, err := asUint64(raw)
		return protoreflect.ValueOfUint32(uint32(v)), err

	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		v, err := asUint64(raw)
		return protoreflect.ValueOfUint64(v), err

	case protoreflect.FloatKind, protoreflect.DoubleKind:
		v, err := asFloat64(raw)
		return protoreflect.ValueOfFloat64(v), err

	case protoreflect.StringKind:
		v, err := asString(raw)
		return protoreflect.ValueOfString(v), err

	case protoreflect.BytesKind:
		v, err := asBytes(raw)
		return protoreflect.ValueOfBytes(v), err

	case protoreflect.EnumKind:
		v, err := asInt64(raw)
		return protoreflect.ValueOfEnum(protoreflect.EnumNumber(v)), err

	case protoreflect.MessageKind, protoreflect.GroupKind:
		// 由调用方负责填充子消息，这里只需要一个占位值
		return protoreflect.Value{}, nil

	default:
		return protoreflect.Value{}, fmt.Errorf("%w: %s", ErrUnsupported, fd.Kind())
	}
}

// decodeMapKey 按 proto map key 的类型把 BSON 文档键还原成对应标量。
// proto map 的键只能是整数、bool 或 string，而 BSON 文档的键恒为字符串。
func decodeMapKey(fd protoreflect.FieldDescriptor, key string) (protoreflect.MapKey, error) {
	switch fd.Kind() {
	case protoreflect.StringKind:
		return protoreflect.ValueOfString(key).MapKey(), nil

	case protoreflect.BoolKind:
		v, err := strconv.ParseBool(key)
		if err != nil {
			return protoreflect.MapKey{}, err
		}

		return protoreflect.ValueOfBool(v).MapKey(), nil

	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		v, err := strconv.ParseInt(key, 10, 32)
		if err != nil {
			return protoreflect.MapKey{}, err
		}

		return protoreflect.ValueOfInt32(int32(v)).MapKey(), nil

	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		v, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			return protoreflect.MapKey{}, err
		}

		return protoreflect.ValueOfInt64(v).MapKey(), nil

	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		v, err := strconv.ParseUint(key, 10, 32)
		if err != nil {
			return protoreflect.MapKey{}, err
		}

		return protoreflect.ValueOfUint32(uint32(v)).MapKey(), nil

	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		v, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			return protoreflect.MapKey{}, err
		}

		return protoreflect.ValueOfUint64(v).MapKey(), nil

	default:
		return protoreflect.MapKey{}, fmt.Errorf("%w: map key %s", ErrUnsupported, fd.Kind())
	}
}

func asDocument(raw any) (bson.D, error) {
	switch v := raw.(type) {
	// bson.D / bson.M / bson.A 都是 primitive 包类型的别名，这里只需列出一种
	case bson.D:
		return v, nil
	case bson.M:
		doc := make(bson.D, 0, len(v))
		for k, val := range v {
			doc = append(doc, bson.E{Key: k, Value: val})
		}

		return doc, nil
	default:
		return nil, fmt.Errorf("%w: want document, got %T", ErrTypeMisuse, raw)
	}
}

func asSlice(raw any) ([]any, error) {
	switch v := raw.(type) {
	case bson.A:
		return v, nil
	case []any:
		return v, nil
	default:
		return nil, fmt.Errorf("%w: want array, got %T", ErrTypeMisuse, raw)
	}
}

func asBool(raw any) (bool, error) {
	if v, ok := raw.(bool); ok {
		return v, nil
	}

	return false, fmt.Errorf("%w: want bool, got %T", ErrTypeMisuse, raw)
}

func asString(raw any) (string, error) {
	if v, ok := raw.(string); ok {
		return v, nil
	}

	return "", fmt.Errorf("%w: want string, got %T", ErrTypeMisuse, raw)
}

func asBytes(raw any) ([]byte, error) {
	switch v := raw.(type) {
	case primitive.Binary:
		return v.Data, nil
	case []byte:
		return v, nil
	default:
		return nil, fmt.Errorf("%w: want binary, got %T", ErrTypeMisuse, raw)
	}
}

func asInt64(raw any) (int64, error) {
	switch v := raw.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("%w: want integer, got %T", ErrTypeMisuse, raw)
	}
}

func asUint64(raw any) (uint64, error) {
	switch v := raw.(type) {
	case uint:
		return uint64(v), nil
	case uint8:
		return uint64(v), nil
	case uint16:
		return uint64(v), nil
	case uint32:
		return uint64(v), nil
	case uint64:
		return v, nil
	case int:
		return uint64(v), nil
	case int8:
		return uint64(v), nil
	case int16:
		return uint64(v), nil
	case int32:
		return uint64(v), nil
	case int64:
		return uint64(v), nil
	case float32:
		return uint64(v), nil
	case float64:
		return uint64(v), nil
	default:
		return 0, fmt.Errorf("%w: want unsigned integer, got %T", ErrTypeMisuse, raw)
	}
}

func asFloat64(raw any) (float64, error) {
	switch v := raw.(type) {
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("%w: want float, got %T", ErrTypeMisuse, raw)
	}
}

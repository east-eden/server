package cache

// import (
// 	"errors"
// 	"flag"
// 	"testing"

// 	"github.com/google/go-cmp/cmp"
// 	"github.com/urfave/cli/v2"
// )

// type EmbededObject struct {
// 	ObjId        int64 `json:"obj_id"`
// 	EmbededDepth int32 `json:"embeded_depth"`
// }

// // test object
// type Object struct {
// 	EmbededObject `json:",inline"`
// 	Id            int64 `json:"_id"`
// 	OwnerId       int64 `json:"owner_id"`
// 	TypeId        int32 `json:"type_id"`
// 	Exp           int64 `json:"exp"`
// 	Level         int32 `json:"level"`
// }

// func (o *Object) GetStoreIndex() int64 {
// 	return o.OwnerId
// }

// func TestCache(t *testing.T) {
// 	set := flag.NewFlagSet("cache", flag.ContinueOnError)
// 	set.String("redis_addr", "localhost:6379", "redis address")
// 	ctx := cli.NewContext(nil, set, nil)
// 	cc := NewGoRedis(ctx)

// 	obj := &Object{
// 		EmbededObject: EmbededObject{
// 			ObjId:        1111000001,
// 			EmbededDepth: 1,
// 		},
// 		Id:      1111000001,
// 		OwnerId: 1111000001,
// 		TypeId:  1,
// 		Exp:     2,
// 		Level:   10,
// 	}

// 	err := cc.SaveObject("test_obj", obj.Id, obj)
// 	if err != nil {
// 		t.Fatalf("TestCache SaveObject failed: %s", err.Error())
// 	}

// 	var newObj Object
// 	err = cc.LoadObject("test_obj", 1111000001, &newObj)
// 	if err != nil {
// 		t.Fatalf("TestCache LoadObject hit failed: %s", err.Error())
// 	}

// 	var newObj2 Object
// 	err = cc.LoadObject("test_obj", 1111000001, &newObj2)
// 	if err != nil && !errors.Is(err, ErrNoResult) && !errors.Is(err, ErrObjectNotFound) {
// 		t.Fatalf("TestCache LoadObject not hit failed: %s", err.Error())
// 	}

// 	diff := cmp.Diff(*obj, newObj)
// 	if diff != "" {
// 		t.Fatalf("TestCache Compare failed: %s", diff)
// 	}
// }

// func BenchmarkRejson(b *testing.B) {
// 	set := flag.NewFlagSet("cache", flag.ContinueOnError)
// 	set.String("redis_addr", "localhost:6379", "redis address")
// 	ctx := cli.NewContext(nil, set, nil)
// 	cc := NewGoRedis(ctx)

// 	for n := 0; n < b.N; n++ {
// 		obj := &Object{
// 			EmbededObject: EmbededObject{
// 				ObjId:        int64(n),
// 				EmbededDepth: int32(n),
// 			},
// 			Id:      int64(n),
// 			OwnerId: int64(n),
// 			TypeId:  1,
// 			Exp:     2,
// 			Level:   10,
// 		}

// 		err := cc.SaveObject("obj_rejson", obj.Id, obj)
// 		if err != nil {
// 			b.Fatal(err)
// 		}
// 	}

// 	for n := 0; n < b.N; n++ {
// 		var obj Object
// 		err := cc.LoadObject("obj_rejson", n, &obj)
// 		if err != nil {
// 			b.Fatal(err)
// 		}
// 	}
// }

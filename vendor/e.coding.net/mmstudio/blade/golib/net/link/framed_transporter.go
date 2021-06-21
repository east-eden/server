package link

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	SizeOfLen              = 4
	DefaultRetryInterval   = time.Millisecond
	DefaultRecvBufferWait  = time.Second
	DefaultWriteBufferWait = time.Second
)

// per session per transporter
type FramedTransporter struct {
	conn net.Conn

	bufferReadHead []byte
	bufferSendHead []byte

	// 上层主动设置超时的情况下在内部尝试设定超时时间，否则直接返回底层错误
	readDeadlineSet  bool
	writeDeadlineSet bool
	opts             *ProtocolOptions
}

func NewFramedTransporter(conn net.Conn, opts *ProtocolOptions) *FramedTransporter {
	if opts == nil {
		panic("opts should not ne nil")
	}
	return &FramedTransporter{
		opts:           opts,
		conn:           conn,
		bufferReadHead: make([]byte, SizeOfLen),
		bufferSendHead: make([]byte, SizeOfLen),
	}
}

func (f *FramedTransporter) ReadBuffer(buffer []byte, length int, zeroRetry bool, tags ...string) (readCount int, err error) {
	if length == 0 {
		return 0, errors.New("read buffer length should not be zero")
	}
	retryCount := f.opts.RetryCountWhenTempError
	retry := false
	for {
		if retry {
			log.Info().Int("retry_left", retryCount).Int("read", readCount).Int("total", length).Strs("tags", tags).Bool("deadline_set", f.readDeadlineSet).Msg("ReadBuffer with retry")
			if f.readDeadlineSet {
				// 读取超时，允许进行buffer wait内的超时设置
				if errSetDeadline := f.SetReadDeadline(time.Now().Add(f.opts.RecvBufferWait)); errSetDeadline != nil {
					return readCount, err
				}
			} else {
				// 系统层未设置超时
				return readCount, err
			}
		}
		var nCurRead int
		nCurRead, err = f.conn.Read(buffer[readCount:])
		readCount += nCurRead
		if readCount == length {
			if err != nil {
				log.Error().Int("length", length).Int("read_count", readCount).Err(err).Msg("ReadBuffer length = index, but retrun error")
			}
			break
		}
		if err != nil {
			if !isTemporaryErr(err) {
				// 连接关闭或非临时性错误
				return readCount, err
			}
			if !zeroRetry && readCount == 0 {
				// 防止无意义重试，交给上层处理
				return readCount, err
			}
			if retryCount == 0 {
				return readCount, err
			}
			retryCount--
			time.Sleep(f.opts.ReadRetryIntervalWhenTempError)
			retry = true
		}
	}
	return readCount, err
}

type ErrorCanNotRecover struct {
	err string
}

func (ec *ErrorCanNotRecover) Error() string { return ec.err }

func (f *FramedTransporter) Receive() (interface{}, error) {
	readCount, err := f.ReadBuffer(f.bufferReadHead, SizeOfLen, false, "header")
	if err != nil {
		if readCount == 0 {
			// 未读取任何字节，如果是temp error应该允许上层恢复这个错误
			return nil, err
		}
		return nil, &ErrorCanNotRecover{err: fmt.Sprintf("error got while reading header,error:%s", err.Error())}
	}

	size := int(f.opts.ByteOrder.Uint32(f.bufferReadHead))
	if size > f.opts.MaxRecv {
		return nil, &ErrorCanNotRecover{err: fmt.Sprintf("error got while parse header,packet too large with:%d limit:%d", size, f.opts.MaxRecv)}
	}
	var frame []byte
	if f.opts.BytesPool == nil {
		frame = make([]byte, size)
	} else {
		frame = f.opts.BytesPool.Alloc(size)
	}
	// header读取成功，但读取body出错不可恢复，内部已经做了重试机制
	sStart := time.Now()
	if readCount, err = f.ReadBuffer(frame, size, true, "body"); err != nil {
		err = &ErrorCanNotRecover{err: fmt.Sprintf("error got while reading body,error:%s,readCount:%d size:%d", err.Error(), readCount, size)}
		log.Error().Err(err).Dur("cost", time.Since(sStart)).Msg("########## header ok,body got error")
	}
	return frame, err
}

func (f *FramedTransporter) WriteBuffer(p interface{}, length int, tags ...string) (err error) {
	var index int
	retryCount := f.opts.RetryCountWhenTempError
	for {
		var nCurWrite int
		if index > 0 {
			log.Info().Int("retry_left", retryCount).Int("write", index).Int("total", length).Strs("tags", tags).Bool("deadline_set", f.writeDeadlineSet).Msg("WriteBuffer with retry")
			if f.writeDeadlineSet {
				// 读取超时，允许进行buffer wait内的超时设置
				if errSetDeadline := f.SetWriteDeadline(time.Now().Add(f.opts.WriteBufferWait)); errSetDeadline != nil {
					return err
				}
			} else {
				return err
			}
		}
		switch buf := p.(type) {
		case []byte:
			nCurWrite, err = f.conn.Write(buf[index:])
		case net.Buffers:
			var writeLen int64
			writeLen, err = io.Copy(f.conn, &buf)
			nCurWrite = int(writeLen)
		default:
			return errors.New("type not supported")
		}
		index += nCurWrite
		if index == length {
			if err != nil {
				log.Error().Int("length", length).Int("size", nCurWrite).Err(err).Msg("WriteBuffer length = index,but return error")
			}
			break
		}
		if err != nil {
			if !isTemporaryErr(err) {
				// 连接关闭或非临时性错误
				return err
			}
			if retryCount == 0 {
				return err
			}
			retryCount--
			time.Sleep(f.opts.WriteRetryIntervalWhenTempError)
		}
	}
	return err
}

// Send data, not goroutine safe
// use writev &  avoid bytes copy
// https://golang.org/pkg/net/#Buffers
// https://github.com/golang/go/issues/13451
func (f *FramedTransporter) Send(msg interface{}) (err error) {
	bodyLen := 0
	switch body := msg.(type) {
	case [][]byte:
		for i := 0; i < len(body); i++ {
			bodyLen += len(body[i])
		}
		f.opts.ByteOrder.PutUint32(f.bufferSendHead, uint32(bodyLen))

		err := f.WriteBuffer(f.bufferSendHead, len(f.bufferSendHead), "[][]byte")
		if err != nil {
			return err
		}
		return f.WriteBuffer(net.Buffers(body), bodyLen)
	case []byte:
		bodyLen = len(body)
		buffs := [][]byte{f.bufferSendHead, body}
		f.opts.ByteOrder.PutUint32(buffs[0], uint32(bodyLen))
		return f.WriteBuffer(net.Buffers(buffs), bodyLen+SizeOfLen, "net.Buffer")
	default:
		if s, ok := msg.(interface {
			WriteTo(w io.Writer) (n int64, err error)
		}); ok {
			_, err = s.WriteTo(f.conn)
		} else {
			err = errors.New("must give [][]byte or []byte, the transporter do not have codec")
		}
	}
	return
}

func (f *FramedTransporter) Close() error { return f.conn.Close() }
func (f *FramedTransporter) SetReadDeadline(t time.Time) error {
	err := f.conn.SetReadDeadline(t)
	if err == nil {
		f.readDeadlineSet = true
	}
	return err
}

func (f *FramedTransporter) SetWriteDeadline(t time.Time) error {
	err := f.conn.SetWriteDeadline(t)
	if err == nil {
		f.writeDeadlineSet = true
	}
	return err
}

func (f *FramedTransporter) RemoteAddr() net.Addr { return f.conn.RemoteAddr() }
func (f *FramedTransporter) LocalAddr() net.Addr  { return f.conn.LocalAddr() }

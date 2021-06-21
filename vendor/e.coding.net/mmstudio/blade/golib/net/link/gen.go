package link

import (
	"encoding/binary"
	"time"

	"e.coding.net/mmstudio/blade/golib/time2"
)

//go:generate gotemplate -outfmt gen_%v "e.coding.net/mmstudio/blade/golib/container/templates/cmap" "ConcurrentMapUint64Session(uint64,*_session,func(key uint64)uint64 {return key})"

//go:generate optiongen --option_with_struct_name=true --v=true
func ProtocolOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"BytesPool":                       BytesPool(nil),
		"MaxRecv":                         int(maxRecv),
		"MaxSend":                         int(maxSend),
		"ByteOrder":                       binary.ByteOrder(binary.LittleEndian),
		"RetryCountWhenTempError":         3,
		"RecvBufferWait":                  time.Duration(DefaultRecvBufferWait),
		"WriteBufferWait":                 time.Duration(DefaultWriteBufferWait),
		"ReadRetryIntervalWhenTempError":  time.Duration(DefaultRetryInterval),
		"WriteRetryIntervalWhenTempError": time.Duration(DefaultRetryInterval),
	}
}

//go:generate optiongen --option_with_struct_name=true --v=true
func WsProtocolOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"RetryCountWhenTempError":         3,
		"RecvBufferWait":                  time.Duration(DefaultRecvBufferWait),
		"WriteBufferWait":                 time.Duration(DefaultWriteBufferWait),
		"ReadRetryIntervalWhenTempError":  time.Duration(DefaultRetryInterval),
		"WriteRetryIntervalWhenTempError": time.Duration(DefaultRetryInterval),
	}
}

const (
	tcpReadBuffer  = 128 * 1024 * 1024
	tcpWriteBuffer = 128 * 1024 * 1024
	maxRecv        = 64 * 1024 * 1024
	maxSend        = 64 * 1024 * 1024
)

//go:generate optiongen --option_with_struct_name=true --v=true
func SpecOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Proto":   "tcp",
		"Address": "",
		//TCPReadBuffer tcp read buffer length
		"TCPReadBuffer": int(tcpReadBuffer),
		//TCPWriteBuffer tcp write buffer length
		"TCPWriteBuffer": int(tcpWriteBuffer),
		// true for single user connection(better performance)
		"TCPNoDelay":      true,
		"TCPLingerSecond": 0,
		"AcceptTimeout":   time.Duration(500 * time.Millisecond),
		//ReadTimeout zero for not set read deadline for client (better  performance)
		"ReadTimeout": time.Duration(0 * time.Millisecond),
		//WriteTimeout zero for not set write deadline for client (better performance)
		"WriteTimeout": time.Duration(0 * time.Millisecond),
		//IdleTimeout idle timeout
		"IdleTimeout": time.Duration(600000 * time.Millisecond),
		// https://tools.ietf.org/html/rfc1122#page-101
		// http://codearcana.com/posts/2015/08/28/tcp-keepalive-is-a-lie.html
		"KeepAlivePeriod": time.Duration(0 * time.Millisecond),
		"SendChanSize":    16,
		//<0： 用户连接设定，不独立起goroutine处理，0：每个请求独立goroutine，>0:goroutine pool处理
		"InvokerPerSession": int(-1),
		//QueuePerInvoker queue gap
		"QueuePerInvoker":          5000,
		"FlushTimeout":             time.Duration(time.Second * 5),
		"timingWheelForFlushInner": (*time2.TimingWheel)(nil),
		"WebSocketPattern":         "/ws",
		"TLSCertFile":              "",
		"TLSKeyFile":               "",
	}
}

func init() {
	InstallSpecWatchDog(func(cc *Spec) {
		cc.timingWheelForFlushInner = time2.NewTimewheel(time.Second, int(cc.FlushTimeout.Seconds()))
	})
	InstallProtocolOptionsWatchDog(func(cc *ProtocolOptions) {
		if cc.ReadRetryIntervalWhenTempError == 0 {
			cc.ReadRetryIntervalWhenTempError = DefaultRetryInterval
		}
		if cc.WriteRetryIntervalWhenTempError == 0 {
			cc.WriteRetryIntervalWhenTempError = DefaultRetryInterval
		}
		if cc.RecvBufferWait == 0 {
			cc.RecvBufferWait = DefaultRecvBufferWait
		}
	})
}

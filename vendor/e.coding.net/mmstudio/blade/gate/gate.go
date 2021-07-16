package gate

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"

	"e.coding.net/mmstudio/blade/gate/msg"
	"e.coding.net/mmstudio/blade/gate/selector"
	"github.com/rs/zerolog/log"

	"e.coding.net/mmstudio/blade/golib/net/link"
	"e.coding.net/mmstudio/blade/golib/reuseport"
	"e.coding.net/mmstudio/blade/golib/sync2"
)

type Gate struct {
	spec     *Spec
	stopChan chan struct{}
	listener net.Listener

	historyConnCount int64

	gatePipe sync2.AtomicInt64
	wsServer *http.Server
	wsProto  link.Protocol

	pluginMgr *pluginManager
}

func New(spec *Spec) (gate *Gate, err error) {
	gate = new(Gate)
	gate.spec = spec
	gate.stopChan = make(chan struct{})
	gate.pluginMgr = newPluginManager(gate.spec.PatchPath)

	portStr := ":" + strconv.Itoa(spec.PortForClient)
	if spec.PortForClientReuse {
		gate.listener, err = reuseport.NewReusablePortListener("tcp", portStr)
	} else {
		gate.listener, err = net.Listen("tcp", portStr)
	}

	if spec.EnableXListener {
		gate.listener, err = gate.XListener()
	}
	setMaxOpenFile()
	return
}

func (g *Gate) Run(closeChan chan struct{}) {
	g.stopChan = closeChan
	go g.listenAndServe()
	go g.webServe()
	go g.pluginMgr.watchPlugin(g.stopChan)
	<-g.stopChan
}

func (g *Gate) OnInit() {}
func (g *Gate) OnClose() {
	_ = g.listener.Close()
	if g.wsServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		_ = g.wsServer.Shutdown(ctx)
		cancel()
	}
}
func (g *Gate) Name() string { return "gate" }

func (g *Gate) listenAndServe() {
	defer func() { _ = g.listener.Close() }()
	for {
		conn, err := link.Accept(g.listener)
		if err != nil {
			log.Warn().Err(err).Msg("gate accept error")
			return
		}
		g.historyConnCount++
		go g.handleConn(conn)
	}
}

func (g *Gate) handshake(conn net.Conn, ws bool) (net.Conn, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Warn().Interface("recover", r).Str("stack", string(debug.Stack())).Msg("recover in handshake")
		}
	}()
	// 连接进来第一个协议必须是 握手 否则关闭
	// 新建立链接: 随机获取后端节点路由
	// 重连时连接: 协议中后断节点
	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	defer func() {
		// 这里的超时设置不要影响后序的处理，所以要取消掉
		_ = conn.SetReadDeadline(time.Time{})
	}()
	var transferFront link.Transporter
	var err error
	if ws {
		transferFront, err = g.wsProto.NewTransporter(conn)
	} else {
		transferFront, err = g.spec.TransferProvider(conn)
	}
	if err != nil {
		log.Warn().Err(err).Msg("get transfer for front fail")
		return nil, err
	}
	respFront := func(code msg.ErrorCode, desc string) {
		resp := msg.HandshakeResp{
			Code: code,
			Desc: desc,
		}
		if data, err := g.spec.MessageCodec.Marshal(&resp); err == nil {
			log.Debug().Str("resp", resp.String()).Str("remote", conn.RemoteAddr().String()).Msg("response handshake")

			_ = transferFront.Send(data)
		}
	}

	payload, err := transferFront.Receive()
	if err != nil {
		log.Warn().Err(err).Msg("get data from front fail")
		respFront(msg.ErrorCode_Unauthorized, fmt.Sprintf("read from client fail:%s", err.Error()))
		return nil, err
	}
	hsFront := msg.Handshake{}
	err = g.spec.MessageCodec.Unmarshal(payload.([]byte), &hsFront)
	if err != nil {
		log.Warn().Err(err).Msg("unmarshal front data fail")
		respFront(msg.ErrorCode_Forbidden, fmt.Sprintf("unmarshal client handshake fail:%s", err.Error()))
		return nil, err
	}
	if g.spec.HandshakeSecret {
		if err = g.checkHandshake(&hsFront); err != nil {
			log.Warn().Err(err).Msg("handshake secret check fail")
			respFront(msg.ErrorCode_Unauthorized, fmt.Sprintf("handshake secret check fail:%s", err.Error()))
			return nil, err
		}
	}

	log.Debug().Str("handshake", hsFront.String()).Str("remote", conn.RemoteAddr().String()).Msg("recv handshake from front")

	entry, err := g.spec.Selector.Select(hsFront.ServiceName,
		selector.WithStrategy(g.spec.SelectorStrategy),
		selector.WithFilters(g.MakeFilter(&hsFront)...),
	)
	if err != nil {
		log.Warn().Err(err).Str("servicename", hsFront.ServiceName).Str("serverId", hsFront.ServerId).Msg("get backend entry fail")
		respFront(msg.ErrorCode_ServiceUnavailable, fmt.Sprintf("get entry fail:%s", err.Error()))
		return nil, err
	}
	backend, err := dialTimeout("tcp", fmt.Sprintf("%s:%s", entry.Host, entry.Port), g.spec.DispatcherDialTimeout)
	if err != nil {
		log.Warn().Err(err).Interface("backend", entry).Msg("dial backend fail")
		respFront(msg.ErrorCode_ServiceUnavailable, fmt.Sprintf("dial entry fail:%s", err.Error()))
		return nil, err
	}

	if g.spec.BackEndHandshake {
		transferBackend, err := g.spec.TransferProvider(backend)
		if err != nil {
			log.Warn().Err(err).Interface("backend", entry).Msg("get transfer for backend fail")
			respFront(msg.ErrorCode_ServiceUnavailable, fmt.Sprintf("backend transfer fail:%s", err.Error()))
			_ = backend.Close()
			return nil, err
		}

		hsBackend := hsFront
		hsBackend.Src = msg.SrcType_GATE
		hsBackend.ClientAddr = conn.RemoteAddr().String()
		hsBackend.ServerId = entry.Identifier
		dataBack, err := g.spec.MessageCodec.Marshal(&hsBackend)

		if err != nil {
			log.Warn().Err(err).Interface("backend", entry).Msg("marshal handshake for backend fail")
			respFront(msg.ErrorCode_ServiceUnavailable, fmt.Sprintf("marshal backend handshake fail:%s", err.Error()))
			_ = backend.Close()
			return nil, err
		}
		err = transferBackend.Send(dataBack)
		if err != nil {
			log.Warn().Err(err).Interface("backend", entry).Msg("send handshake to backedn fail")
			respFront(msg.ErrorCode_ServiceUnavailable, fmt.Sprintf("write to backend fail:%s", err.Error()))
			_ = backend.Close()
			return nil, err
		}
	}

	respFront(msg.ErrorCode_Success, entry.Identifier)
	return backend, nil
}

// must have `ts`,`random`,`sign` in handshake meta
// ts: timestamp seconds
// random: random string
const (
	tsKey     = "ts"
	randomKey = "random"
	signKey   = "sign"
)

var (
	SignFailTimeFormatError = errors.New("timestamp format error")
	SignFailMissingField    = errors.New("missing required fields")
	SignFailTimeError       = errors.New("time error")
	SignFailSignError       = errors.New("sign error")
)

func (g *Gate) checkHandshake(h *msg.Handshake) error {
	var ts int64
	var random string
	var sign string
	kv := make([]string, 0)
	for _, m := range h.Meta {
		if m.Key == signKey {
			sign = m.Value
			continue
		} else if m.Key == randomKey {
			random = m.Value
		} else if m.Key == tsKey {
			t, err := strconv.ParseInt(m.Value, 10, 64)
			if err != nil {
				return SignFailTimeFormatError
			}
			ts = t
		}
		kv = append(kv, fmt.Sprintf("%s=%s", m.Key, m.Value))
	}
	if random == "" || ts == 0 || sign == "" {
		return SignFailMissingField
	}
	if g.spec.HandshakeCheckTime {
		now := time.Now().Unix()
		if now-ts < -20 || now-ts > 20 {
			// 20秒有效
			return SignFailTimeError
		}
	}
	// md5(k1=v1&k2=v2+secret)
	sort.Strings(kv)
	bs := fmt.Sprintf("%s%s", strings.Join(kv, "&"), g.spec.SecretKey)
	hasher := md5.New()
	hasher.Write([]byte(bs))
	ms := hex.EncodeToString(hasher.Sum(nil))
	if ms != sign {
		log.Warn().Str("beforeSign", bs).Str("mySign", ms).Str("recvSign", sign).Msg("check sign fail")
		return SignFailSignError
	}
	return nil
}

func (g *Gate) handleConn(conn net.Conn) {

	backend, err := g.handshake(conn, false)
	if err != nil {
		_ = conn.Close()
		return
	}

	// TODO pipe不参与协议解析，后端节点下线更新时，客户端是有感知的，会触发断线
	if g.spec.EnableZeroCopy {
		go pipeZeroCopy(backend, conn)
		pipeZeroCopy(conn, backend)
	} else {
		go g.pipe(backend, conn)
		g.pipe(conn, backend)
	}
}

func (g *Gate) pipe(destConn, srcConn net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			log.Warn().Interface("recover", r).Str("stack", string(debug.Stack())).Msg("recover in pipe")
		}
	}()
	defer func() { _ = destConn.Close() }()

	buf := make([]byte, 2*4096)
	// TODO 有frame解析，封装不太好，暂时注释
	//onBytes := onBytesInOut(g.spec.ByteOrderForLog, tag, g.spec.OnFrameForLog)
	for {
		select {
		case <-g.stopChan:
			return
		default:
		}
		_ = srcConn.SetReadDeadline(time.Now().Add(g.spec.ConnReadTimeout))
		nr, er := srcConn.Read(buf)
		if nr > 0 {
			index := 0
			for {
				_ = destConn.SetWriteDeadline(time.Now().Add(g.spec.ConnWriteTimeout))
				nw, ew := destConn.Write(buf[index:nr])
				if ew != nil && !isTemporaryErr(ew) {
					log.Error().Err(ew).Str("addr", destConn.LocalAddr().String()).Msg("destConn write fail")
					break
				}
				index += nw
				if nr == nw {
					break
				}
			}
			//onBytes(buf[0:nr])
		}
		if isTemporaryErr(er) {
			continue
		}
		if er == io.EOF {
			log.Error().Err(er).Str("addr", srcConn.LocalAddr().String()).Msg("srcConn read fail EOF")
			break
		}
		if er != nil {
			log.Error().Err(er).Str("addr", srcConn.LocalAddr().String()).Msg("srcConn read fail")
			break
		}
	}
}

func isTemporaryErr(err error) bool {
	if err == nil {
		return false
	}
	netErr, ok := err.(net.Error)
	if ok && netErr.Timeout() && netErr.Temporary() {
		return true
	}
	return false
}

func dialTimeout(network, address string, timeout time.Duration) (conn net.Conn, err error) {
	m := int(timeout / time.Second)
	for i := 0; i < m; i++ {
		conn, err = net.DialTimeout(network, address, timeout)
		if err == nil || !strings.Contains(err.Error(), "can't assign requested address") {
			break
		}
		time.Sleep(time.Second)
	}
	return
}

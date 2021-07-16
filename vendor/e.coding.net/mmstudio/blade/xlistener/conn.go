package xlistener

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"time"

	"e.coding.net/mmstudio/blade/xlistener/internal/dh64"
)

var _ net.Conn = &xConn{}

type Dialer func() (net.Conn, error)

type xConn struct {
	net.Conn
	connId    uint64
	secretKey [8]byte
}

func newXConn(conn net.Conn, connId uint64) (*xConn, error) {
	c := &xConn{
		Conn:   conn,
		connId: connId,
	}
	return c, nil
}

func (x *xConn) handshake(serverAddr []byte, handshakeTimeout time.Duration) (err error) {
	if handshakeTimeout > 0 {
		_ = x.Conn.SetDeadline(time.Now().Add(handshakeTimeout))
		defer func() { _ = x.Conn.SetDeadline(time.Time{}) }()
	}
	var (
		addressLen     = 1 + len(serverAddr)
		buf            = make([]byte, addressLen+24)
		addressField   = buf[0:addressLen]
		publicKeyField = buf[addressLen : addressLen+8]
		connIdField    = buf[addressLen+8 : addressLen+16]
		randField      = buf[addressLen+16 : addressLen+24]
	)
	// read client public secret
	if _, err = io.ReadFull(x.Conn, publicKeyField); err != nil {
		return
	}
	clientPubKey := binary.LittleEndian.Uint64(publicKeyField)
	if clientPubKey == 0 {
		err = errors.New("client public key is zero")
		return
	}
	serverPriKey, serverPubKey := dh64.KeyPair()
	secret := dh64.Secret(serverPriKey, clientPubKey)
	binary.LittleEndian.PutUint64(x.secretKey[:], secret)
	addressField[0] = byte(addressLen - 1)
	copy(addressField[1:], serverAddr)
	binary.LittleEndian.PutUint64(publicKeyField, serverPubKey)
	binary.LittleEndian.PutUint64(connIdField, x.connId)
	_, _ = rand.Read(randField)
	if _, err = x.Conn.Write(buf[:]); err != nil {
		err = errors.New("send handshake response failed: " + err.Error())
		return
	}

	var buf2 [16]byte
	if _, err = io.ReadFull(x.Conn, buf2[:]); err != nil {
		err = errors.New("read twice handshake failed: " + err.Error())
		return
	}

	hash := md5.New()
	hash.Write(randField)
	hash.Write(x.secretKey[:])
	md5sum := hash.Sum(nil)
	if !bytes.Equal(buf2[:], md5sum) {
		err = fmt.Errorf("twice handshake not equals: %x, %x", buf2[:], md5sum)
		return
	}
	return
}

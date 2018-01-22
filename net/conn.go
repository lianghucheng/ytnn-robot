package net

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/name5566/leaf/log"
	"net"
	"reflect"
	"sync"
)

type MyConn struct {
	sync.Mutex
	conn      *websocket.Conn
	writeChan chan []byte
	maxMsgLen uint32
	closeFlag bool
}

func newMyConn(conn *websocket.Conn, pendingWriteNum int, maxMsgLen uint32) *MyConn {
	wsConn := new(MyConn)
	wsConn.conn = conn
	wsConn.writeChan = make(chan []byte, pendingWriteNum)
	wsConn.maxMsgLen = maxMsgLen
	go func() {
		for b := range wsConn.writeChan {
			if b == nil {
				break
			}
			err := conn.WriteMessage(websocket.BinaryMessage, b)
			if err != nil {
				break
			}
		}
		conn.Close()
		wsConn.Lock()
		wsConn.closeFlag = true
		wsConn.Unlock()
	}()
	return wsConn
}

func (wsConn *MyConn) doDestroy() {
	wsConn.conn.UnderlyingConn().(*net.TCPConn).SetLinger(0)
	wsConn.conn.Close()

	if !wsConn.closeFlag {
		close(wsConn.writeChan)
		wsConn.closeFlag = true
	}
}

func (wsConn *MyConn) Destroy() {
	wsConn.Lock()
	defer wsConn.Unlock()

	wsConn.doDestroy()
}

func (wsConn *MyConn) Close() {
	wsConn.Lock()
	defer wsConn.Unlock()
	if wsConn.closeFlag {
		return
	}

	wsConn.doWrite(nil)
	wsConn.closeFlag = true
}

func (wsConn *MyConn) doWrite(b []byte) {
	if len(wsConn.writeChan) == cap(wsConn.writeChan) {
		log.Debug("close conn: channel full")
		wsConn.doDestroy()
		return
	}
	wsConn.writeChan <- b
}

func (wsConn *MyConn) LocalAddr() net.Addr {
	return wsConn.conn.LocalAddr()
}

func (wsConn *MyConn) RemoteAddr() net.Addr {
	return wsConn.conn.RemoteAddr()
}

// goroutine not safe
func (wsConn *MyConn) ReadMsg() ([]byte, error) {
	_, b, err := wsConn.conn.ReadMessage()
	return b, err
}

// args must not be modified by the others goroutines
func (wsConn *MyConn) WriteMsg(msg interface{}) error {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		return errors.New("json message pointer required")
	}
	msgID := msgType.Elem().Name()
	m := map[string]interface{}{msgID: msg}
	data, err := json.Marshal(m)
	if err != nil {
		return errors.New("marshal message error")
	}
	wsConn.Lock()
	defer wsConn.Unlock()
	if wsConn.closeFlag {
		return nil
	}

	// get len
	msgLen := uint32(len(data))

	// check len
	if msgLen > wsConn.maxMsgLen {
		return errors.New("message too long")
	} else if msgLen < 1 {
		return errors.New("message too short")
	}

	wsConn.doWrite(data)
	return nil
}

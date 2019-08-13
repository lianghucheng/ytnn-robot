package robot

import (
	"sync"
	"ytnn-robot/net"

	"time"

	"encoding/json"
	"ytnn-robot/msg"

	"github.com/gorilla/websocket"
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/network"
)

var gameMu *sync.Mutex
var gamecount int

func init() {
	gameMu = new(sync.Mutex)
}
func InitGame(addr string) {
	client := new(net.Client)
	client.Addr = "ws://" + addr
	client.ConnNum = 1
	client.ConnectInterval = 3 * time.Second
	client.HandshakeTimeout = 10 * time.Second
	client.PendingWriteNum = 100
	client.MaxMsgLen = 4096
	client.NewAgent = newAgentGame

	client.Start()
	clients = append(clients, client)
}

func DestroyGame() {
	for _, client := range clients {
		client.Close()
	}
}

type AgentGame struct {
	conn       *net.MyConn
	playerData *PlayerData
}

func newAgentGame(conn *net.MyConn) network.Agent {
	a := new(AgentGame)
	a.conn = conn
	a.playerData = newPlayerData()
	return a
}

func (a *AgentGame) writeMsg(msg interface{}) {
	err := a.conn.WriteMsg(msg)
	if err != nil {
		log.Debug("write message: %v", err)
	}
	return
}

func (a *AgentGame) readMsg() {
	for {
		msg, err := a.conn.ReadMsg()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Debug("error: %v", err)
			}
			break
		}

		jsonMap := map[string]interface{}{}
		err = json.Unmarshal(msg, &jsonMap)
		if err == nil {
			a.handleMsg(jsonMap)
		} else {
			log.Error("%v", err)
		}
	}
}

func (a *AgentGame) Run() {
	go func() {
		for {
			(<-dispatcher.ChanTimer).Cb()
		}
	}()

	go a.robotLogin()
	a.readMsg()
}

func (a *AgentGame) OnClose() {

}

func (a *AgentGame) sendHeartbeat() {
	a.writeMsg(&msg.C2S_Heartbeat{})
}

func (a *AgentGame) robotLogin() {
	gameMu.Lock()
	defer gameMu.Unlock()

	a.writeMsg(&msg.C2S_TokenAuthorize{
		Token: tokenMap[gamecount],
	})
	gamecount++
}

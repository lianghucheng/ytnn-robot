package robot

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"sync"
	"time"
	"ytnn-robot/msg"
	"ytnn-robot/net"

	"github.com/gorilla/websocket"
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/network"
	"github.com/name5566/leaf/timer"
)

var (
	//addr = "ws://niuniu.shenzhouxing.com:3661"
	addr = "ws://localhost:3653"
	//addr = "ws://47.93.13.22:2000"
	//addr        = "ws://139.199.180.94:3661"
	clients     []*net.Client
	unionids    []string
	nicknames   []string
	headimgurls []string
	loginIPs    []string
	hallCount   = 0
	loginCount  = 0
	mu          sync.Mutex
	loginMu     sync.Mutex
	robotNumber = 100 // 机器人数量

	dispatcher *timer.Dispatcher
	tokenMap   =make(map[int]string)
)

func init() {
	rand.Seed(time.Now().UnixNano())

	dispatcher = timer.NewDispatcher(0)
}

func InitHall() {
	client := new(net.Client)
	client.Addr = addr
	client.ConnNum = robotNumber
	client.ConnectInterval = 3 * time.Second
	client.HandshakeTimeout = 10 * time.Second
	client.PendingWriteNum = 100
	client.MaxMsgLen = 4096
	client.NewAgent = newAgent

	client.Start()
	clients = append(clients, client)
}

func DestroyHall() {
	for _, client := range clients {
		client.Close()
	}
}

type Agent struct {
	conn       *net.MyConn
	playerData *PlayerData
}

func newAgent(conn *net.MyConn) network.Agent {
	a := new(Agent)
	a.conn = conn
	a.playerData = newPlayerData()
	return a
}

func newPlayerData() *PlayerData {
	playerData := new(PlayerData)
	return playerData
}

func (a *Agent) writeMsg(msg interface{}) {
	err := a.conn.WriteMsg(msg)
	if err != nil {
		log.Debug("write message: %v", err)
	}
	return
}

func (a *Agent) readMsg() {
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
			a.handleMsgHall(jsonMap)
		} else {
			log.Error("%v", err)
		}
	}
}

func (a *Agent) Run() {
	go func() {
		for {
			(<-dispatcher.ChanTimer).Cb()
		}
	}()

	go a.robotLogin()
	a.readMsg()
}

func (a *Agent) OnClose() {

}

/*
func (a *Agent) handleMsgHall(jsonMap map[string]interface{}) {
	for k, v := range jsonMap {
		switch k {
		case "H2C_Heartbeat":
			a.sendHallHeartbeat()

		case "H2C_GameAddr":
			nNAddr := v.(map[string]interface{})["NNAddr"].(string)
			log.Debug("牛牛的地址:%v", nNAddr)
			//InitGame(nNAddr)

		case "H2C_Login":
			token := v.(map[string]interface{})["Token"].(string)
			a.playerData.Token = token
			a.playerData.Accountid = int64(v.(map[string]interface{})["AccountID"].(float64))
			mu.Lock()
			tokenMap[hallCount] = token
			hallCount++
			mu.Unlock()
		case "H2C_UpdateUserChips":
		case "H2C_GameIngAddr":
			{
				gameAddr := v.(map[string]interface{})["Addr"].(string)
				InitGame(gameAddr)
			}
		default:
			log.Release("message: <%v> not deal", k)
		}
	}
}
*/
func (a *Agent) handleMsgHall(jsonMap map[string]interface{}) {
	for k, v := range jsonMap {
		switch k {
		case "H2C_Heartbeat":
			a.sendHallHeartbeat()

		case "H2C_GameAddr":
			nNAddr := v.(map[string]interface{})["NN"].(map[string]interface{})["Addr"].(string)
			log.Debug("牛牛的地址:%v", nNAddr)
			InitGame(nNAddr)

		case "H2C_Login":
			token := v.(map[string]interface{})["Token"].(string)
			a.playerData.Token = token
			a.playerData.Accountid = int64(v.(map[string]interface{})["AccountID"].(float64))
			mu.Lock()
			log.Debug("token:%v   hallCount:%v",token,hallCount)
			tokenMap[hallCount] = token
			hallCount++
			mu.Unlock()
		case "H2C_UpdateUserChips":
		default:
			log.Release("message: <%v> not deal", k)
		}
	}
}
func (a *Agent) robotLogin() {
	loginMu.Lock()
	defer loginMu.Unlock()
	account:=""
	if loginCount > 0 &&loginCount<9{
		account = "root40"
	}else{
		account = "root4"
	}
	account += strconv.Itoa(loginCount + 1)
	if loginCount == 0 {
		account = "root401"
	}
	a.writeMsg(&msg.C2H_AccountLogin{
		Account:  account,
		Password: "123456",
	})
	loginCount++
	log.Debug("account:%v,password:%v", account, "123456")

}

func (a *Agent) sendHallHeartbeat() {
	a.writeMsg(&msg.C2H_Heartbeat{})
}

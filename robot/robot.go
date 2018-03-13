package robot

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/network"
	"github.com/name5566/leaf/timer"
	"math/rand"
	"strconv"
	"sync"
	"time"
	"ytnn-robot/common"
	"ytnn-robot/net"
)

var (
	addr = "ws://niuniu.shenzhouxing.com:3661"
	//addr        = "ws://192.168.1.34:3661"
	//addr        = "ws://139.199.180.94:3661"
	clients     []*net.Client
	unionids    []string
	nicknames   []string
	headimgurls []string
	loginIPs    []string
	count       = 0
	mu          sync.Mutex

	robotNumber = 100 // 机器人数量

	dispatcher *timer.Dispatcher
)

func init() {
	rand.Seed(time.Now().UnixNano())

	names, ips := make([]string, 0), make([]string, 0)
	var err error
	names, err = common.ReadFile("conf/robot_nickname.txt")
	names = common.Shuffle2(names)

	ips, _ = common.ReadFile("conf/robot_ip.txt")
	ips = common.Shuffle2(ips)
	if err == nil {
		nicknames = append(nicknames, names[:robotNumber]...)
		loginIPs = append(loginIPs, ips[:robotNumber]...)
	} else {
		log.Debug("read file error: %v", err)
	}
	temp := rand.Perm(robotNumber)

	for i := 0; i < robotNumber; i++ {
		unionids = append(unionids, strconv.Itoa(i))
		headimgurls = append(headimgurls, "https://www.shenzhouxing.com/robot/"+strconv.Itoa(temp[i])+".jpg")
	}

	dispatcher = timer.NewDispatcher(0)
}

func Init() {
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

func Destroy() {
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
	a.playerData.PlayTimes = rand.Intn(4) + 2
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
			a.handleMsg(jsonMap)
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

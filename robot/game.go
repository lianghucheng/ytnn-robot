package robot

import (
	"math/rand"
	"time"
	"ytnn-robot/msg"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (a *Agent) sendHeartbeat() {
	a.writeMsg(&msg.C2S_Heartbeat{})
}

func (a *Agent) robotLogin() {
	mu.Lock()
	defer mu.Unlock()
	a.playerData.Unionid = unionids[count]
	a.playerData.Nickname = nicknames[count]
	a.writeMsg(&msg.C2S_RobotLogin{
		UnionID:    unionids[count],
		Nickname:   nicknames[count],
		Headimgurl: headimgurls[count],
		LoginIP:    loginIPs[count],
	})
	//log.Debug("UnionID: %v - IP: %v - Nickname: %v", unionids[count], loginIPs[count], nicknames[count])
	count++
}

func (a *Agent) enterRoom() {
	a.writeMsg(&msg.C2S_Matching{
		RoomType:      a.playerData.RoomType,
		BaseScore:     a.playerData.BaseScore,
		RedPacketType: a.playerData.RedPacketType,
	})
}

func (a *Agent) getAllPlayer() {
	a.writeMsg(&msg.C2S_GetAllPlayers{})
}

func (a *Agent) reconnect() {
	a.writeMsg(&msg.C2S_EnterRoom{})
}

func (a *Agent) doBid(bid int) {
	switch bid {
	case 0, 1, 2, 3, 4:
	default:
		bid = 0
	}
	a.writeMsg(&msg.C2S_Bid{
		Bid: bid,
	})
}

func (a *Agent) doDouble(double int) {
	switch double {
	case 5, 10, 15, 20, 25:
	default:
		double = 0
	}
	a.writeMsg(&msg.C2S_Double{
		Double: double,
	})
}

func (a *Agent) show() {
	a.writeMsg(&msg.C2S_Show{})
}

func (a *Agent) exit() {
	a.writeMsg(&msg.C2S_ExitRoom{})
}

func (a *Agent) wxFake(fee int) {
	a.writeMsg(&msg.C2S_FakeWXPay{
		TotalFee: fee, // 100 = 1 元
	})
}

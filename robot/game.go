package robot

import (
	"github.com/name5566/leaf/log"
	"math/rand"
	"ytnn-robot/msg"
)

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
	log.Debug("UnionID: %v - IP: %v - Nickname: %v", unionids[count], loginIPs[count], nicknames[count])
	count++
}

func (a *Agent) bid() {
	a.writeMsg(&msg.C2S_Bid{
		Bid: rand.Intn(5),
	})
}

func (a *Agent) double() {
	double := []int{5, 10, 15, 20, 25}
	a.writeMsg(&msg.C2S_Double{
		Double: double[rand.Intn(len(double))],
	})
}

func (a *Agent) show() {
	a.writeMsg(&msg.C2S_Show{})
}

func (a *Agent) exit() {
	a.writeMsg(&msg.C2S_ExitRoom{})
}

func (a *Agent) Fake() {
	a.writeMsg(&msg.C2S_FakeWXPay{
		TotalFee: 100, // 100 = 1 元
	})
}

package robot

import (
	"math/rand"
	"strconv"
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
	baseScore := 0
	redPacketType := 0
	index, _ := strconv.Atoi(a.playerData.Unionid)
	switch {
	case index > -1 && index < 50:
		a.playerData.RoomType = roomBaseScoreMatching
		switch {
		case a.playerData.Chips >= 50000:
			baseScore = 400
		case a.playerData.Chips >= 200000:
			baseScore = 1000
		default:
			baseScore = 100
		}
	case index > 49 && index < 75:
		a.playerData.RoomType = roomBaseScoreMatching
		switch {
		case a.playerData.Chips >= 200000:
			baseScore = 1000
		default:
			baseScore = 400
		}
	case index > 74 && index < 100:
		a.playerData.RoomType = roomRedPacketMatching
		switch {
		case a.playerData.Chips >= 80000:
			redPacketType = 10
		default:
			redPacketType = 1
		}
		CronFunc("10 0 19 * * *", a.enterRoom)
	default:
		return
	}
	a.writeMsg(&msg.C2S_Matching{
		RoomType:      a.playerData.RoomType,
		BaseScore:     baseScore,
		RedPacketType: redPacketType,
	})
}

func (a *Agent) reconnect() {
	a.writeMsg(&msg.C2S_EnterRoom{})
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

func (a *Agent) wxFake(fee int) {
	a.writeMsg(&msg.C2S_FakeWXPay{
		TotalFee: fee, // 100 = 1 元
	})
}

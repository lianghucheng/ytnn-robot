package robot

import (
	"math/rand"
	"strconv"
	"ytnn-robot/msg"
)

func (a *Agent) isMe(pos int) bool {
	return a.playerData.Position == pos
}

func (a *Agent) sendHeartbeat() {
	a.writeMsg(&msg.C2S_Heartbeat{})
}

func (a *Agent) wechatLogin() {
	mu.Lock()
	defer mu.Unlock()
	a.playerData.Unionid = unionids[count]
	a.playerData.Nickname = nicknames[count]
	a.writeMsg(&msg.C2S_WeChatLogin{
		Unionid:    unionids[count],
		NickName:   nicknames[count],
		Headimgurl: headimgurls[count],
	})
	count++
}

func (a *Agent) setRobotData() {
	index, _ := strconv.Atoi(a.playerData.Unionid)
	a.writeMsg(&msg.C2S_SetRobotData{
		LoginIP: loginIPs[index],
	})
}

func (a *Agent) enterRoom() {
	a.writeMsg(&msg.C2S_EnterRoom{})
}

func (a *Agent) enterRandRoom() {
	a.playerData.getRandRoom()
	a.startMatching(roomBaseScoreMatching, a.playerData.BaseScore, 0)
}

func (a *Agent) startMatching(roomType int, baseScore int, redPacketType int) {
	a.writeMsg(&msg.C2S_Matching{
		RoomType:      roomType,
		BaseScore:     baseScore,
		RedPacketType: redPacketType,
	})
}

func (a *Agent) getAllPlayer() {
	a.writeMsg(&msg.C2S_GetAllPlayers{})
}

func (a *Agent) prepare() {
	a.writeMsg(&msg.C2S_Prepare{})
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

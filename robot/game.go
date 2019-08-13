package robot

import (
	"math/rand"
	"time"
	"ytnn-robot/msg"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (a *AgentGame) enterRoom() {
	a.writeMsg(&msg.C2NN_Matching{
		RoomType:      a.playerData.RoomType,
		MinChips:      a.playerData.BaseScore,
		RedPacketType: a.playerData.RedPacketType,
	})
}
func (a *AgentGame) joinRoom() {
	a.writeMsg(&msg.C2NN_EnterRoom{})
}
func (a *AgentGame) getAllPlayer() {
	a.writeMsg(&msg.C2NN_GetAllPlayers{})
}

func (a *AgentGame) reconnect() {
	a.writeMsg(&msg.C2NN_EnterRoom{})
}

func (a *AgentGame) doBid(bid int) {
	switch bid {
	case 0, 1, 2, 3, 4:
	default:
		bid = 0
	}
	a.writeMsg(&msg.C2NN_Bid{
		Bid: bid,
	})
}

func (a *AgentGame) doDouble(double int) {
	switch double {
	case 5, 10, 15, 20, 25:
	default:
		double = 0
	}
	a.writeMsg(&msg.C2NN_Double{
		Double: double,
	})
}

func (a *AgentGame) show() {
	a.writeMsg(&msg.C2NN_Show{})
}

func (a *AgentGame) exit() {
	a.writeMsg(&msg.C2NN_ExitRoom{})
}

func (a *AgentGame) wxFake(fee int) {
	a.writeMsg(&msg.C2S_FakeWXPay{
		TotalFee: fee, // 100 = 1 元
	})
}

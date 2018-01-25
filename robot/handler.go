package robot

import (
	"github.com/name5566/leaf/log"
	"math/rand"
	"time"
)

const (
	S2C_EnterRoom_LackOfChips = 6
	S2C_LeaveRoom_LackOfChips = 1
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (a *Agent) handleMsg(jsonMap map[string]interface{}) {
	for k, v := range jsonMap {
		switch k {
		case "S2C_Heartbeat":
			a.sendHeartbeat()
		case "S2C_ActionBid":
			Delay(time.Duration(rand.Intn(6)+3)*time.Second, a.bid)
		case "S2C_ActionDouble":
			Delay(time.Duration(rand.Intn(6)+3)*time.Second, a.double)
		case "S2C_ShowFifthCard":
			Delay(time.Duration(rand.Intn(6)+3)*time.Second, a.show)
		case "S2C_ShowWinnersAndLosers":
			Delay(time.Duration(rand.Intn(6)+3)*time.Second, a.exit)
		case "S2C_EnterRoom":
			if int(v.(map[string]interface{})["Error"].(float64)) == S2C_EnterRoom_LackOfChips {
				a.Fake()
			}
		case "S2C_LeaveRoom":
			if int(v.(map[string]interface{})["Error"].(float64)) == S2C_LeaveRoom_LackOfChips {
				a.Fake()
			}
		default:
			if k == "S2C_PayOK" {
				log.Debug("message: <%v> ", k)
			}
		}
	}
}

func Delay(d time.Duration, cb func()) {
	if cb == nil {
		return
	}
	time.AfterFunc(d, func() {
		cb()
	})
}

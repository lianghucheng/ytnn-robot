package robot

import (
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/timer"
	"math/rand"
	"time"
)

const (
	roomBaseScoreMatching = 1 // 底分匹配
	roomRedPacketMatching = 4 // 红包匹配

	S2C_EnterRoom_OK          = 0
	S2C_EnterRoom_NotCreated  = 1 // "房间: " + S2C_EnterRoom.RoomNumber + " 未创建"
	S2C_EnterRoom_Full        = 2 // "房间: " + S2C_EnterRoom.RoomNumber + " 玩家人数已满"
	S2C_EnterRoom_Unknown     = 4 // 进入房间出错，请稍后重试
	S2C_EnterRoom_LackOfChips = 6 // 需要 + S2C_EnterRoom.MinChips + 筹码才能进入
	S2C_EnterRoom_NotRightNow = 7 // 比赛暂未开始，请到时再来

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
		case "S2C_UpdateUserChips":
			// 机器人筹码会先于登陆信息传，如果还没收到就先为0，进入最低等级房
			a.playerData.Chips = int64(v.(map[string]interface{})["Chips"].(float64))
		case "S2C_Login":
			// 触发进入房间
			if v.(map[string]interface{})["AnotherRoom"].(bool) {
				a.reconnect()
				return
			}
			DelayDo(time.Duration(10)*time.Second, a.enterRoom)
		case "S2C_EnterRoom":
			switch int(v.(map[string]interface{})["Error"].(float64)) {
			case S2C_EnterRoom_OK:
			case S2C_EnterRoom_Unknown:
				// 机器人进入房间不会创建，如果没有一人房或者两人房就返回这条错误
				DelayDo(time.Duration(10)*time.Second, a.enterRoom)
			case S2C_EnterRoom_LackOfChips:
				a.wxFake(100)
			case S2C_EnterRoom_NotRightNow:
				// 红包比赛场未开始
			}
			a.playerData.PlayTimes = rand.Intn(9) + 2
		case "S2C_PayOK":
		case "S2C_SitDown":
		case "S2C_GameStart":
			a.playerData.PlayTimes--
		case "S2C_ActionBid":
			DelayDo(time.Duration(rand.Intn(2)+3)*time.Second, a.bid)
		case "S2C_ActionDouble":
			DelayDo(time.Duration(rand.Intn(2)+3)*time.Second, a.double)
		case "S2C_ShowFifthCard":
			DelayDo(time.Duration(rand.Intn(2)+3)*time.Second, a.show)
		case "S2C_ShowWinnersAndLosers":
			if a.playerData.PlayTimes <= 0 {
				DelayDo(time.Duration(rand.Intn(5)+10)*time.Second, a.exit)
			}
		case "S2C_ExitRoom", "S2C_LeaveRoom":
			// 退出房间
		default:
			log.Debug("message: <%v> ", k)
		}
	}
}

func DelayDo(d time.Duration, cb func()) {
	if cb == nil {
		return
	}
	time.AfterFunc(d, cb)
}

func CronFunc(expr string, cb func()) {
	cronExpr, _ := timer.NewCronExpr(expr)
	dispatcher.CronFunc(cronExpr, func() {
		if cb != nil {
			cb()
		}
	})
}

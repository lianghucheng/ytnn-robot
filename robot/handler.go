package robot

import (
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/timer"
	"math/rand"
	"strconv"
	"time"
)

const (
	roomBaseScoreMatching = 1 // 底分匹配
	roomRedPacketMatching = 4 // 红包匹配

	S2C_EnterRoom_OK            = 0
	S2C_EnterRoom_Full          = 2 // "房间: " + S2C_EnterRoom.RoomNumber + " 玩家人数已满"
	S2C_EnterRoom_Unknown       = 4 // 进入房间出错，请稍后重试
	S2C_EnterRoom_LackOfChips   = 6 // 需要 + S2C_EnterRoom.MinChips + 筹码才能进入
	S2C_EnterRoom_NotRightNow   = 7 // 比赛暂未开始，请到时再来
	S2C_EnterRoom_MaxChipsLimit = 8 // 进入房间最大携带金币限制

	S2C_CreateRoom_InOtherRoom = 3 // 正在其他房间对局，是否回去？
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (a *Agent) handleMsg(jsonMap map[string]interface{}) {
	for k, v := range jsonMap {
		switch k {
		case "S2C_Heartbeat":
			a.sendHeartbeat()
		case "S2C_UpdateRedPacketTaskList":
		case "S2C_UpdateChipTaskList":
		case "S2C_UpdateTaskProgress":
		case "S2C_UpdateUserChips":
			a.playerData.Chips = int64(v.(map[string]interface{})["Chips"].(float64))
		case "S2C_UpdatePlayerChips":
		case "S2C_Login":
			index, _ := strconv.Atoi(a.playerData.Unionid)
			switch {
			case index > -1 && index < 25:
				a.playerData.RoomType = roomBaseScoreMatching
				a.playerData.BaseScore = 100
			case index > 24 && index < 50:
				a.playerData.RoomType = roomBaseScoreMatching
				a.playerData.BaseScore = 400
			case index > 49:
				a.playerData.RoomType = roomRedPacketMatching
				a.playerData.RedPacketType = 1
			}
			// 断线重连
			if v.(map[string]interface{})["AnotherRoom"].(bool) {
				a.reconnect()
				log.Debug("unionid: %v 断线重连", a.playerData.Unionid)
				return
			}
			// 触发进入房间
			if index > 49 {
				a.enterRoom()
				CronFunc("10 0 19 * * *", a.enterRoom)
			} else {
				a.enterRoom()
			}
		case "S2C_CreateRoom":
			switch int(v.(map[string]interface{})["Error"].(float64)) {
			case S2C_CreateRoom_InOtherRoom:
			default:
				log.Debug("unionid: %v 创建房间 error: %v", a.playerData.Unionid, int(v.(map[string]interface{})["Error"].(float64)))
			}
		case "S2C_EnterRoom":
			switch int(v.(map[string]interface{})["Error"].(float64)) {
			case S2C_EnterRoom_OK:
				log.Debug("unionid: %v 进入房间", a.playerData.Unionid)
				//a.playerData.PlayTimes = rand.Intn(9) + 2
				a.playerData.PlayTimes = 1
			case S2C_EnterRoom_Full:
				DelayDo(time.Duration(10)*time.Second, a.enterRoom)
			case S2C_EnterRoom_Unknown:
				// 机器人进入房间不会创建，如果没有一人房或者两人房就返回这条错误
				DelayDo(time.Duration(10)*time.Second, a.enterRoom)
			case S2C_EnterRoom_LackOfChips:
				log.Debug("unionid: %v 请求充钱", a.playerData.Unionid)
				a.wxFake(100)
			case S2C_EnterRoom_NotRightNow:
				// 不处理等待定时器自动触发
			case S2C_EnterRoom_MaxChipsLimit:
				log.Debug("unionid: %v 携带金币超过上限", a.playerData.Unionid)
				a.wxFake(-100)
			default:
				log.Debug("unionid: %v 进入房间 error: %v", a.playerData.Unionid, int(v.(map[string]interface{})["Error"].(float64)))
			}
		case "S2C_GameStop":
		case "S2C_StandUp":
		case "S2C_PayOK":
			log.Debug("unionid: %v 充值成功", a.playerData.Unionid)
			DelayDo(time.Duration(10)*time.Second, a.enterRoom)
		case "S2C_SitDown":
		case "S2C_GameStart":
			a.playerData.PlayTimes--
		case "S2C_ActionBid":
			DelayDo(time.Duration(rand.Intn(2)+3)*time.Second, a.bid)
		case "S2C_Bid":
		case "S2C_Grab":
		case "S2C_DecideDealer":
		case "S2C_ActionDouble":
			DelayDo(time.Duration(rand.Intn(2)+3)*time.Second, a.double)
		case "S2C_Double":
		case "S2C_UpdatePokerHands":
		case "S2C_ShowFifthCard":
			DelayDo(time.Duration(rand.Intn(2)+3)*time.Second, a.show)
		case "S2C_OxResult":
		case "S2C_ShowAllResults":
		case "S2C_ShowWinnersAndLosers":
			if a.playerData.PlayTimes < 1 {
				DelayDo(time.Duration(rand.Intn(4)+11)*time.Second, a.exit)
			}
		case "S2C_ClearAction":
		case "S2C_AddPlayerChips":
		case "S2C_AddPlayerRedPacket":
		case "S2C_ActionStart":
		case "S2C_LeaveRoom":
			DelayDo(time.Duration(10)*time.Second, a.enterRoom)
		case "S2C_ExitRoom":
			DelayDo(time.Duration(10)*time.Second, a.enterRoom)
		default:
			log.Debug("unionid: %v message: <%v> not deal", a.playerData.Unionid, k)
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

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

	S2C_ExitRoom_OK          = 0
	S2C_ExitRoom_GamePlaying = 1 // 游戏进行中，不能退出房间
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
			log.Debug("金币数: %v", a.playerData.Chips)
		case "S2C_UpdatePlayerChips":
		case "S2C_Login":
			index, _ := strconv.Atoi(a.playerData.Unionid)
			switch {
			case index >= 0 && index < 25:
				a.playerData.RoomType = roomBaseScoreMatching
				a.playerData.BaseScore = 100
			case index >= 25 && index < 50:
				a.playerData.RoomType = roomBaseScoreMatching
				a.playerData.BaseScore = 400
			case index >= 50:
				a.playerData.RoomType = roomRedPacketMatching
				a.playerData.RedPacketType = 1
			}
			log.Debug("登陆成功")
			// 断线重连
			if v.(map[string]interface{})["AnotherRoom"].(bool) {
				a.reconnect()
				log.Debug("断线重连")
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
				log.Debug("创建房间 error: %v", int(v.(map[string]interface{})["Error"].(float64)))
			}
		case "S2C_EnterRoom":
			switch int(v.(map[string]interface{})["Error"].(float64)) {
			case S2C_EnterRoom_OK:
				log.Debug("进入房间")
				a.playerData.PositionHands = make(map[int][]int)
				a.playerData.PlayTimes = 1
				a.playerData.Position = int(v.(map[string]interface{})["Position"].(float64))
				log.Debug("座位号: %v", a.playerData.Position)
				a.getAllPlayer()
			case S2C_EnterRoom_Full:
				log.Debug("房间已满")
				DelayDo(time.Duration(10)*time.Second, a.enterRoom)
			case S2C_EnterRoom_Unknown:
				log.Debug("无单人房")
				DelayDo(time.Duration(10)*time.Second, a.enterRoom)
			case S2C_EnterRoom_LackOfChips:
				log.Debug("请求充钱")
				a.wxFake(100)
			case S2C_EnterRoom_NotRightNow:
				// 不处理等待定时器自动触发
			case S2C_EnterRoom_MaxChipsLimit:
				log.Debug("金币过多")
				a.wxFake(-100)
			default:
				log.Debug("进入房间 error: %v", int(v.(map[string]interface{})["Error"].(float64)))
			}
		case "S2C_ActionStart":
			log.Debug("开始倒计时")
		case "S2C_GameStart":
			log.Debug("游戏开始")
			a.playerData.PlayTimes--
		case "S2C_GameStop":
			log.Debug("游戏终止")
		case "S2C_PayOK":
			log.Debug("充值成功")
			DelayDo(time.Duration(10)*time.Second, a.enterRoom)
		case "S2C_SitDown":
			pos := int(v.(map[string]interface{})["Position"].(float64))
			a.playerData.PositionHands[pos] = []int{}
			log.Debug("座位号: %v 玩家就位", pos)
		case "S2C_StandUp":
			pos := int(v.(map[string]interface{})["Position"].(float64))
			if a.playerData.Position == pos {
				log.Debug("自己起立")
			} else {
				log.Debug("座位号: %v 玩家起立", pos)
			}
			delete(a.playerData.PositionHands, pos)
		case "S2C_ActionBid":
			DelayDo(time.Duration(rand.Intn(2)+3)*time.Second, func() {
				log.Debug("叫庄: 0")
				a.doBid(0)
			})
		case "S2C_Bid":
		case "S2C_Grab":
		case "S2C_DecideDealer":
			a.playerData.DealerPos = int(v.(map[string]interface{})["Position"].(float64))
			log.Debug("座位号: %v 玩家做庄", a.playerData.DealerPos)
		case "S2C_ActionDouble":
			DelayDo(time.Duration(rand.Intn(2)+3)*time.Second, func() {
				log.Debug("叫倍: 5")
				a.doDouble(5)
			})
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
		case "S2C_LeaveRoom":
			DelayDo(time.Duration(10)*time.Second, a.enterRoom)
		case "S2C_ExitRoom":
			switch int(v.(map[string]interface{})["Error"].(float64)) {
			case S2C_ExitRoom_OK:
				pos := int(v.(map[string]interface{})["Position"].(float64))
				if a.playerData.Position == pos {
					log.Debug("自己退出房间")
					DelayDo(time.Duration(10)*time.Second, a.enterRoom)
				} else {
					log.Debug("%v 号退出房间", pos)
					log.Debug("房间剩余 %v 人", len(a.playerData.PositionHands))
				}
			}
		default:
			log.Debug("message: <%v> not deal", k)
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
	if cb == nil {
		return
	}
	cronExpr, _ := timer.NewCronExpr(expr)
	dispatcher.CronFunc(cronExpr, cb)
}

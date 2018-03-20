package robot

import (
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/timer"
	"math/rand"
	"sort"
	"strconv"
	"time"
	"ytnn-robot/poker"
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
			//log.Debug("登陆成功")
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
				log.Release("创建房间 error: %v", int(v.(map[string]interface{})["Error"].(float64)))
			}
		case "S2C_EnterRoom":
			switch int(v.(map[string]interface{})["Error"].(float64)) {
			case S2C_EnterRoom_OK:
				log.Debug(" %v 进入房间", a.playerData.Unionid)
				a.playerData.PositionHands = make(map[int]*CardsDetail)
				//a.playerData.PlayTimes = 1
				a.playerData.Position = int(v.(map[string]interface{})["Position"].(float64))
				//log.Debug("座位号: %v", a.playerData.Position)
				a.getAllPlayer()
			case S2C_EnterRoom_Full:
				//log.Debug("房间已满")
				DelayDo(10*time.Second, a.enterRoom)
			case S2C_EnterRoom_Unknown:
				//log.Debug("无单人房")
				DelayDo(10*time.Second, a.enterRoom)
			case S2C_EnterRoom_LackOfChips:
				log.Release("%v 请求充钱", a.playerData.Unionid)
				a.wxFake(100)
			case S2C_EnterRoom_NotRightNow:
				// 不处理等待定时器自动触发
			case S2C_EnterRoom_MaxChipsLimit:
				log.Release("%v 金币过多", a.playerData.Unionid)
				a.wxFake(-100)
			default:
				log.Release("进入房间 error: %v", int(v.(map[string]interface{})["Error"].(float64)))
			}
		case "S2C_ActionStart":
			//log.Debug("开始倒计时")
		case "S2C_GameStart":
			//log.Debug("游戏开始")
			//a.playerData.PlayTimes--
		case "S2C_UpdatePokerHands":
			pos := int(v.(map[string]interface{})["Position"].(float64))
			hands := ArrayInterfaceToInt(v.(map[string]interface{})["Hands"].([]interface{}))
			sort.Sort(sort.Reverse(sort.IntSlice(hands)))
			if s, ok := a.playerData.PositionHands[pos]; ok {
				s.Cards = hands
				s.Type = poker.GetWinType(s.Cards)
			}
		case "S2C_GameStop":
			//log.Debug("游戏终止")
		case "S2C_PayOK":
			//log.Debug("充值成功")
			DelayDo(10*time.Second, a.enterRoom)
		case "S2C_SitDown":
			pos := int(v.(map[string]interface{})["Position"].(float64))
			a.playerData.PositionHands[pos] = &CardsDetail{}
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
			bid := rand.Intn(5) // 先随机出倍数，然后比不过就置最小值
			mine, ok := a.playerData.PositionHands[a.playerData.Position]
			if !ok {
				return
			}
			for pos, s := range a.playerData.PositionHands {
				if a.playerData.Position == pos {
					continue
				}
				if len(a.playerData.PositionHands[pos].Cards) < 1 {
					continue
				}
				if s.Type > mine.Type ||
					(s.Type == mine.Type && s.Cards[0] > mine.Cards[0]) {
					bid = 0
					break
				}
			}
			DelayDo(time.Duration(rand.Intn(3)+2)*time.Second, func() {
				log.Debug("叫庄: %v", bid)
				a.doBid(bid)
			})
		case "S2C_Bid":
		case "S2C_Grab":
		case "S2C_DecideDealer":
			a.playerData.DealerPos = int(v.(map[string]interface{})["Position"].(float64))
			//log.Debug("座位号: %v 玩家做庄", a.playerData.DealerPos)
		case "S2C_ActionDouble":
			if a.playerData.Position == a.playerData.DealerPos {
				log.Debug("等待其他玩家叫倍")
				return
			}
			canDouble := []int{5, 10, 15, 20, 25}
			double := canDouble[rand.Intn(5)] // 先随机出倍数，然后比不过就置最小值
			mine, ok := a.playerData.PositionHands[a.playerData.Position]
			if !ok || len(mine.Cards) < 1 {
				return
			}
			if s, ok := a.playerData.PositionHands[a.playerData.DealerPos]; ok {
				if s.Type > mine.Type ||
					(s.Type == mine.Type && s.Cards[0] > mine.Cards[0]) {
					double = canDouble[0]
				}
			}
			DelayDo(time.Duration(rand.Intn(3)+2)*time.Second, func() {
				log.Debug("叫倍: %v", double)
				a.doDouble(double)
			})
		case "S2C_Double":
		case "S2C_ShowFifthCard":
			DelayDo(time.Duration(rand.Intn(2)+3)*time.Second, a.show)
		case "S2C_OxResult":
		case "S2C_ShowAllResults":
		case "S2C_ShowWinnersAndLosers":
			DelayDo(time.Duration(rand.Intn(4)+11)*time.Second, a.exit)
			/*switch a.playerData.RoomType {
			case roomBaseScoreMatching:
				DelayDo(time.Duration(rand.Intn(4)+11)*time.Second, func() {
					if len(a.playerData.PositionHands) > 2 || len(a.playerData.PositionHands) == 1 {
						log.Debug("人数: %v, 退出房间", len(a.playerData.PositionHands))
						a.exit()
					}
				})
			case roomRedPacketMatching:
				DelayDo(time.Duration(rand.Intn(4)+11)*time.Second, a.exit)
			default:
				log.Release("a.playerData.RoomType - default not deal")
			}*/
		case "S2C_ClearAction":
		case "S2C_AddPlayerChips":
		case "S2C_AddPlayerRedPacket":
		case "S2C_LeaveRoom":
			DelayDo(10*time.Second, a.enterRoom)
		case "S2C_ExitRoom":
			switch int(v.(map[string]interface{})["Error"].(float64)) {
			case S2C_ExitRoom_OK:
				pos := int(v.(map[string]interface{})["Position"].(float64))
				if a.playerData.Position == pos {
					log.Debug("自己退出房间")
					DelayDo(10*time.Second, a.enterRoom)
				} else {
					log.Debug("%v 号退出房间", pos)
				}
			}
		default:
			log.Release("message: <%v> not deal", k)
		}
	}
}

func DelayDo(d time.Duration, cb func()) {
	if cb == nil {
		return
	}
	time.AfterFunc(d, cb)
}

func ArrayInterfaceToInt(array []interface{}) []int {
	var temp []int
	for _, v := range array {
		temp = append(temp, int(v.(float64)))
	}
	return temp
}

func CronFunc(expr string, cb func()) {
	if cb == nil {
		return
	}
	cronExpr, _ := timer.NewCronExpr(expr)
	dispatcher.CronFunc(cronExpr, cb)
}

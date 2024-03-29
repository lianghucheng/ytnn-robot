package robot

import (
	"math/rand"
	"sort"
	"strconv"
	"time"
	"ytnn-robot/poker"

	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/timer"
)

const (
	roomBaseScoreMatching = 1 // 底分匹配
	roomRedPacketMatching = 4 // 红包匹配

	NN2C_EnterRoom_OK            = 0
	NN2C_EnterRoom_Full          = 2 // "房间: " + S2C_EnterRoom.RoomNumber + " 玩家人数已满"
	NN2C_EnterRoom_Unknown       = 4 // 进入房间出错，请稍后重试
	NN2C_EnterRoom_LackOfChips   = 6 // 需要 + S2C_EnterRoom.MinChips + 筹码才能进入
	NN2C_EnterRoom_NotRightNow   = 7 // 比赛暂未开始，请到时再来
	NN2C_EnterRoom_MaxChipsLimit = 8 // 进入房间最大携带金币限制

	NN2C_CreateRoom_InOtherRoom = 3 // 正在其他房间对局，是否回去？

	NN2C_ExitRoom_OK          = 0
	NN2C_ExitRoom_GamePlaying = 1 // 游戏进行中，不能退出房间
)

func RandInt() int {
	return rand.Intn(3)
}
func init() {
	rand.Seed(time.Now().UnixNano())
}
func (a *AgentGame) handleMsg(jsonMap map[string]interface{}) {
	for k, v := range jsonMap {
		switch k {

		case "S2C_Heartbeat":
			a.sendHeartbeat()
		case "NN2C_UpdateRedPacketTaskList":
		case "NN2C_UpdateChipTaskList":
		case "NN2C_UpdateTaskProgress":
		case "NN2C_UpdateUserChips":
			a.playerData.Chips = int64(v.(map[string]interface{})["Chips"].(float64))
			log.Debug("金币数: %v", a.playerData.Chips)
		case "NN2C_UpdatePlayerChips":
		case "S2C_Authorize":
			minChips := 0
			index := RandInt()
			if index == 0 {
				minChips = 2000
			}
			if index == 1 {
				minChips = 50000
			}
			if index == 2 {
				minChips = 500000
			}
			a.playerData.RoomType = roomBaseScoreMatching
			a.playerData.BaseScore = minChips
			a.enterRoom()
		case "S2C_Login":
			//log.Debug("登陆成功")
			// 断线重连
			if v.(map[string]interface{})["AnotherRoom"].(bool) {
				a.reconnect()
				log.Debug("断线重连")
				return
			}
			index, _ := strconv.Atoi(a.playerData.Unionid)
			switch {
			case index < 1000:
				a.playerData.RoomType = roomBaseScoreMatching
				a.playerData.BaseScore = 100
				a.enterRoom()
			case index < 50:
				a.playerData.RoomType = roomBaseScoreMatching
				a.playerData.BaseScore = 400
				a.enterRoom()
			case index < 100:
				a.playerData.RoomType = roomRedPacketMatching
				a.playerData.RedPacketType = 1
				CronFunc("10 0 19 * * *", a.enterRoom)
			default:
				log.Release("有剩余机器人未处理")
			}
		case "NN2C_CreateRoom":
			switch int(v.(map[string]interface{})["Error"].(float64)) {
			case NN2C_CreateRoom_InOtherRoom:
				a.joinRoom()
			default:
				log.Release("创建房间 error: %v", int(v.(map[string]interface{})["Error"].(float64)))
			}
		case "NN2C_EnterRoom":
			switch int(v.(map[string]interface{})["Error"].(float64)) {
			case NN2C_EnterRoom_OK:
				log.Debug(" %v 进入房间", a.playerData.Unionid)
				a.playerData.PositionHands = make(map[int]*CardsDetail)
				a.playerData.PlayTimes = rand.Intn(5) + 5
				a.playerData.Position = int(v.(map[string]interface{})["Position"].(float64))
				//log.Debug("座位号: %v", a.playerData.Position)
				a.getAllPlayer()
			case NN2C_EnterRoom_Full:
				//log.Debug("房间已满")
				DelayDo(10*time.Second, a.enterRoom)
			case NN2C_EnterRoom_Unknown:
				log.Debug("无单人房")
				DelayDo(10*time.Second, a.enterRoom)
			case NN2C_EnterRoom_LackOfChips:
				log.Release("%v 请求充钱", a.playerData.Unionid)
				a.wxFake(100)
			case NN2C_EnterRoom_NotRightNow:
				// 不处理等待定时器自动触发
			case NN2C_EnterRoom_MaxChipsLimit:
				log.Release("%v 金币过多", a.playerData.Unionid)
				a.wxFake(-100)
			default:
				log.Release("进入房间 error: %v", int(v.(map[string]interface{})["Error"].(float64)))
			}
		case "NN2C_ActionStart":
			//log.Debug("开始倒计时")
		case "NN2C_GameStart":
			//log.Debug("游戏开始")
			a.playerData.PlayTimes--
		case "NN2C_UpdatePokerHands":
			pos := int(v.(map[string]interface{})["Position"].(float64))
			hands := ArrayInterfaceToInt(v.(map[string]interface{})["Hands"].([]interface{}))
			sort.Sort(sort.Reverse(sort.IntSlice(hands)))
			if s, ok := a.playerData.PositionHands[pos]; ok {
				s.Cards = hands
				s.Type = poker.GetWinType(s.Cards)
			}
		case "NN2C_GameStop":
			//log.Debug("游戏终止")
		case "NN2C_PayOK":
			//log.Debug("充值成功")
			DelayDo(10*time.Second, a.enterRoom)
		case "NN2C_SitDown":
			pos := int(v.(map[string]interface{})["Position"].(float64))
			a.playerData.PositionHands[pos] = &CardsDetail{}
			log.Debug("座位号: %v 玩家就位", pos)
		case "NN2C_StandUp":
			pos := int(v.(map[string]interface{})["Position"].(float64))
			if a.playerData.Position == pos {
				log.Debug("自己起立")
			} else {
				log.Debug("座位号: %v 玩家起立", pos)
			}
			delete(a.playerData.PositionHands, pos)
		case "NN2C_ActionBid":
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
		case "NN2C_Bid":
		case "NN2C_Grab":
		case "NN2C_DecideDealer":
			a.playerData.DealerPos = int(v.(map[string]interface{})["Position"].(float64))
			//log.Debug("座位号: %v 玩家做庄", a.playerData.DealerPos)
		case "NN2C_ActionDouble":
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
		case "NN2C_Double":
		case "NN2C_ShowFifthCard":
			DelayDo(time.Duration(rand.Intn(2)+3)*time.Second, a.show)
		case "NN2C_OxResult":
		case "NN2C_ShowAllResults":

		case "NN2C_ShowWinnersAndLosers":
			//DelayDo(time.Duration(rand.Intn(4)+11)*time.Second, a.exit)
			switch a.playerData.RoomType {
			case roomBaseScoreMatching:
				DelayDo(time.Duration(rand.Intn(4)+11)*time.Second, func() {
					if a.playerData.PlayTimes < 1 {
						//log.Debug("人数: %v, 退出房间", len(a.playerData.PositionHands))
						a.exit()
					}
					log.Debug("****")
				})
			case roomRedPacketMatching:
				DelayDo(time.Duration(rand.Intn(4)+11)*time.Second, a.exit)
			default:
				log.Release("a.playerData.RoomType - default not deal")
			}
		case "NN2C_ClearAction":
		case "NN2C_AddPlayerChips":
			pos := int(v.(map[string]interface{})["Position"].(float64))
			chips := int(v.(map[string]interface{})["Chips"].(float64))
			log.Debug("position:%v chips:%v", pos, chips)
		case "NN2C_AddPlayerRedPacket":

		case "NN2C_LeaveRoom":
			minChips := 0
			index := RandInt()
			if index == 0 {
				minChips = 2000
			}
			if index == 1 {
				minChips = 50000
			}
			if index == 2 {
				minChips = 500000
			}
			a.playerData.RoomType = roomBaseScoreMatching
			a.playerData.BaseScore = minChips
			log.Debug("自己退出房间")
			DelayDo(10*time.Second, a.enterRoom)
		case "NN2C_ExitRoom":
			pos := int(v.(map[string]interface{})["Position"].(float64))
			switch int(v.(map[string]interface{})["Error"].(float64)) {
			case NN2C_ExitRoom_OK:
				if a.playerData.Position == pos {
					minChips := 0
					index := RandInt()
					if index == 0 {
						minChips = 2000
					}
					if index == 1 {
						minChips = 50000
					}
					if index == 2 {
						minChips = 500000
					}
					a.playerData.RoomType = roomBaseScoreMatching
					a.playerData.BaseScore = minChips
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
	cb()
	cronExpr, _ := timer.NewCronExpr(expr)
	dispatcher.CronFunc(cronExpr, cb)
}

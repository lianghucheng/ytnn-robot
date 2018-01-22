package robot

import (
	"github.com/name5566/leaf/log"
	"math/rand"
	"time"
	"ytnn-robot/poker"
)

func (a *Agent) handleMsg(jsonMap map[string]interface{}) {
	if _, ok := jsonMap["S2C_Heartbeat"]; ok {
		a.sendHeartbeat()
	} else if res, ok := jsonMap["S2C_Login"].(map[string]interface{}); ok {
		a.playerData.AccountID = int(res["AccountID"].(float64))
		log.Debug("accID: %v 登录", a.playerData.AccountID)
		a.playerData.Role = int(res["Role"].(float64))
		if a.playerData.Role == 1 {
			a.setRobotData()
			return
		}
		if res["AnotherRoom"].(bool) {
			a.enterRoom()
		} else {
			//a.enterRandRoom()
		}
	} else if res, ok := jsonMap["S2C_CreateRoom"].(map[string]interface{}); ok {
		err := res["Error"].(float64)
		switch err {
		case 6:
			log.Debug("accID: %v 需要%v筹码才能游戏", a.playerData.AccountID, res["MinChips"].(float64))
		}
	} else if res, ok := jsonMap["S2C_EnterRoom"].(map[string]interface{}); ok {
		err := res["Error"].(float64)
		switch err {
		case 0:
			a.playerData.Position = int(res["Position"].(float64))
			a.playerData.RoomType = int(res["RoomType"].(float64))
			roomNumber := res["RoomNumber"].(string)
			switch a.playerData.RoomType {
			case roomBaseScoreMatching:
				a.playerData.BaseScore = int(res["BaseScore"].(float64))
				log.Debug("accID: %v 进入房间:%v 底分: %v", a.playerData.AccountID, roomNumber, a.playerData.BaseScore)
			}
			a.getAllPlayer()
		case 6:
			log.Debug("accID: %v 需要%v筹码才能进入", a.playerData.AccountID, res["MinChips"].(float64))
			Delay(func() {
				a.enterRandRoom()
			})
		}
	} else if res, ok := jsonMap["S2C_SitDown"].(map[string]interface{}); ok {
		if a.isMe(int(res["Position"].(float64))) {
			a.prepare()
		}
	} else if res, ok := jsonMap["S2C_UpdatePokerHands"].(map[string]interface{}); ok {
		a.playerData.hands = To1DimensionalArray(res["Hands"].([]interface{}))
		log.Debug("hands: %v", poker.ToCardsString(a.playerData.hands))
	} else if _, ok := jsonMap["S2C_ActionBid"].(map[string]interface{}); ok {
		Delay(func() {
			a.bid()
		})
	} else if _, ok := jsonMap["S2C_ActionDouble"].(map[string]interface{}); ok {
		Delay(func() {
			a.double()
		})
	} else if _, ok := jsonMap["S2C_ShowWinnersAndLosers"]; ok {
		Delay(func() {
			a.enterRandRoom()
		})
	}
}

func To1DimensionalArray(array []interface{}) []int {
	var newArray []int
	for _, v := range array {
		newArray = append(newArray, int(v.(float64)))
	}
	return newArray
}

func Delay(cb func()) {
	time.AfterFunc(time.Duration((rand.Intn(2))+3)*time.Second, func() {
		if cb != nil {
			cb()
		}
	})
}

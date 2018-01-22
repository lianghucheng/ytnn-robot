package robot

import "math/rand"

// 房间类型
const (
	roomBaseScoreMatching = 1 // 1 底分匹配
)

var (
	roomType  = roomBaseScoreMatching
	baseScore = []int{100, 400, 1000}
)

type PlayerData struct {
	Unionid       string
	Nickname      string
	AccountID     int
	RoomType      int
	BaseScore     int
	RedPacketType int
	Position      int
	Role          int

	hands []int
}

func (playerData *PlayerData) getRandRoom() {
	playerData.RoomType = roomType
	playerData.BaseScore = baseScore[rand.Intn(len(baseScore))]
}

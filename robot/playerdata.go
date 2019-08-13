package robot

type PlayerData struct {
	Token    string
	Unionid  string
	Nickname string

	// 在登陆的时候处理第一次要打的随机局数
	// 之后每次接受开始游戏的时候把生成好的随机局数进行减一操作
	// 在退出的函数处理下一次要打的随机局数
	PlayTimes     int   // 进入房间机器人玩的局数
	Chips         int64 // 金币数
	RoomType      int   // 进入房间类型
	BaseScore     int
	RedPacketType int
	Position      int                  // 自己的位置
	DealerPos     int                  // 庄的位置
	PositionHands map[int]*CardsDetail // 玩家的手牌
	Bid           int                  // 叫庄
	Double        int                  // 叫倍
	Accountid     int64
}

type CardsDetail struct {
	Cards []int
	Type  int
}

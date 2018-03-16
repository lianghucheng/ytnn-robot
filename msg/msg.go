package msg

type C2S_Heartbeat struct {
}

type C2S_RobotLogin struct {
	UnionID    string
	Nickname   string
	Headimgurl string
	LoginIP    string
}

type C2S_Matching struct {
	RoomType      int // 房间类型: 1 底分匹配、4 红包匹配
	BaseScore     int // 底分: 100、400、1000
	RedPacketType int // 红包种类(元): 1、5、10、50
}

type C2S_GetAllPlayers struct {
}

type C2S_EnterRoom struct {
}

type C2S_Bid struct {
	Bid int
}

type C2S_Double struct {
	Double int // 5、10、 15、 20、 25
}

type C2S_Show struct {
}

type C2S_ExitRoom struct {
}

type C2S_FakeWXPay struct {
	TotalFee int
}

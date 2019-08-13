package msg

type C2S_Heartbeat struct {
}

type C2H_Heartbeat struct {
}
type C2H_AccountLogin struct {
	Account  string
	Password string
}
type C2S_TokenAuthorize struct {
	Token string
}
type C2NN_Matching struct {
	RoomType      int // 房间类型: 1 底分匹配、4 红包匹配
	MinChips      int // 底分: 100、400、1000
	RedPacketType int // 红包种类(元): 1、5、10、50
}

type C2NN_GetAllPlayers struct {
}

type C2NN_EnterRoom struct {
}

type C2NN_Bid struct {
	Bid int
}

type C2NN_Double struct {
	Double int // 5、10、 15、 20、 25
}

type C2NN_Show struct {
}

type C2NN_ExitRoom struct {
}

type C2S_FakeWXPay struct {
	TotalFee int
}

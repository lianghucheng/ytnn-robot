package msg

type C2S_Heartbeat struct {
}

type C2S_WeChatLogin struct {
	NickName   string
	Headimgurl string
	Unionid    string
}

type C2S_SetRobotData struct {
	LoginIP string
}

type C2S_EnterRoom struct {
}

type C2S_Matching struct {
	RoomType      int // 房间类型: 0 练习、1 底分匹配、4 红包匹配
	BaseScore     int // 底分: 100、400、1000
	RedPacketType int // 红包种类(元): 1、5、10、50
}

type C2S_GetAllPlayers struct {
}

type C2S_Prepare struct {
}

type C2S_Bid struct {
	Bid int
}

type C2S_Double struct {
	Double int // 5、10、 15、 20、 25
}

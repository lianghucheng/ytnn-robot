package msg

type C2S_Heartbeat struct {
}

type C2S_RobotLogin struct {
	UnionID    string
	Nickname   string
	Headimgurl string
	LoginIP    string
}

type C2S_Prepare struct {
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

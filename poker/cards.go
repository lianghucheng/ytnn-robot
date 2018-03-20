package poker

const (
	_           = iota
	DiamondCard // 方块
	ClubCard    // 梅花
	HeartCard   // 红桃
	SpadeCard   // 黑桃
	JokerCard   // 王
)

// 游戏结果
const (
	ResultLose = iota // 0 失败
	ResultWin         // 1 胜利
	ResultDraw        // 2 流局
)

// 牛牛牌型
const (
	Untreated    = -1 // -1 未处理
	NoOx         = 0  //  0 没牛
	OneOx        = 1  //  1 牛一
	TwoOx        = 2  //  2 牛二
	ThreeOx      = 3  //  3 牛三
	FourOx       = 4  //  4 牛四
	FiveOx       = 5  //  5 牛五
	SixOx        = 6  //  6 牛六
	SevenOx      = 7  //  7 牛七
	EightOx      = 8  //  8 牛八
	NineOx       = 9  //  9 牛九
	OxOx         = 10 // 10 牛牛
	FourBomb     = 11 // 11 四炸
	FiveFlowerOx = 12 // 12 五花牛
	FiveLittleOx = 13 // 13 五小牛
)

var (
	Diamonds = []int{0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48}  // 方块A到方块K
	Clubs    = []int{1, 5, 9, 13, 17, 21, 25, 29, 33, 37, 41, 45, 49}  // 梅花A到梅花K
	Hearts   = []int{2, 6, 10, 14, 18, 22, 26, 30, 34, 38, 42, 46, 50} // 红桃A到红桃K
	Spades   = []int{3, 7, 11, 15, 19, 23, 27, 31, 35, 39, 43, 47, 51} // 黑桃A到黑桃K

	Jokers = []int{52, 53} // 小王、大王

	CardType = []int{
		DiamondCard, ClubCard, HeartCard, SpadeCard, DiamondCard, ClubCard, HeartCard, SpadeCard,
		DiamondCard, ClubCard, HeartCard, SpadeCard, DiamondCard, ClubCard, HeartCard, SpadeCard,
		DiamondCard, ClubCard, HeartCard, SpadeCard, DiamondCard, ClubCard, HeartCard, SpadeCard,
		DiamondCard, ClubCard, HeartCard, SpadeCard, DiamondCard, ClubCard, HeartCard, SpadeCard,
		DiamondCard, ClubCard, HeartCard, SpadeCard, DiamondCard, ClubCard, HeartCard, SpadeCard,
		DiamondCard, ClubCard, HeartCard, SpadeCard, DiamondCard, ClubCard, HeartCard, SpadeCard,
		DiamondCard, ClubCard, HeartCard, SpadeCard, JokerCard, JokerCard,
	}

	CardString = []string{
		"方块A", "梅花A", "红桃A", "黑桃A", "方块2", "梅花2", "红桃2", "黑桃2",
		"方块3", "梅花3", "红桃3", "黑桃3", "方块4", "梅花4", "红桃4", "黑桃4",
		"方块5", "梅花5", "红桃5", "黑桃5", "方块6", "梅花6", "红桃6", "黑桃6",
		"方块7", "梅花7", "红桃7", "黑桃7", "方块8", "梅花8", "红桃8", "黑桃8",
		"方块9", "梅花9", "红桃9", "黑桃9", "方块10", "梅花10", "红桃10", "黑桃10",
		"方块J", "梅花J", "红桃J", "黑桃J", "方块Q", "梅花Q", "红桃Q", "黑桃Q",
		"方块K", "梅花K", "红桃K", "黑桃K", "小王", "大王",
	}
	CardValue = []int{
		1, 1, 1, 1, 2, 2, 2, 2,
		3, 3, 3, 3, 4, 4, 4, 4,
		5, 5, 5, 5, 6, 6, 6, 6,
		7, 7, 7, 7, 8, 8, 8, 8,
		9, 9, 9, 9, 10, 10, 10, 10,
		10, 10, 10, 10, 10, 10, 10, 10,
		10, 10, 10, 10,
	}
)

func ToCardsString(cards []int) []string {
	s := []string{}
	for _, v := range cards {
		if v > -1 {
			s = append(s, CardString[v])
		}
	}
	return s
}

func ToCardsTypeString(cardsType int) string {
	switch cardsType {
	case Untreated: // -1 未处理
		return "未处理"
	case NoOx: // 0 没牛
		return "没牛"
	case OneOx: // 1 牛一
		return "牛一"
	case TwoOx: // 2 牛二
		return "牛二"
	case ThreeOx: // 3 牛三
		return "牛三"
	case FourOx: // 4 牛四
		return "牛四"
	case FiveOx: // 5 牛五
		return "牛五"
	case SixOx: // 6 牛六
		return "牛六"
	case SevenOx: // 7 牛七
		return "牛七"
	case EightOx: // 8 牛八
		return "牛八"
	case NineOx: // 9 牛九
		return "牛九"
	case OxOx: // 10 牛牛
		return "牛牛"
	case FourBomb: // 10 四炸
		return "四炸"
	case FiveFlowerOx: // 11 五花牛
		return "五花牛"
	case FiveLittleOx: // 12 五小牛
		return "五小牛"
	default:
		return "牌型错误"
	}
}

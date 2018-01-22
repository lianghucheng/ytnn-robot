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

var (
	Diamonds = []int{0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48}  // 方块A到方块K
	Clubs    = []int{1, 5, 9, 13, 17, 21, 25, 29, 33, 37, 41, 45, 49}  // 梅花A到梅花K
	Hearts   = []int{2, 6, 10, 14, 18, 22, 26, 30, 34, 38, 42, 46, 50} // 红桃A到红桃K
	Spades   = []int{3, 7, 11, 15, 19, 23, 27, 31, 35, 39, 43, 47, 51} // 黑桃A到黑桃K
	Jokers   = []int{52, 53}                                           // 小王、大王

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

func CountCardValue(cards []int, card int) int {
	count, value := 0, CardValue[card]
	for _, v := range cards {
		if CardValue[v] == value {
			count++
		}
	}
	return count
}

func ToCardsString(cards []int) []string {
	s := []string{}
	for _, v := range cards {
		if v > -1 {
			s = append(s, CardString[v])
		}
	}
	return s
}

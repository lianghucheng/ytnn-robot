package poker

import (
	"ytnn-robot/common"
)

// cards 为数组，从大到小排序
func GetWinType(cards []int) int {
	if len(cards) == 0 || len(cards) > 5 {
		return Untreated
	}
	// 总和
	sum := 0
	for i := 0; i < len(cards); i++ {
		sum += CardValue[cards[i]]
	}
	// 五小牛
	if cards[0] < 16 {
		if sum <= 10 {
			return FiveLittleOx
		}
	}
	// 五花牛
	if cards[4] > 39 {
		return FiveFlowerOx

	}
	// 四炸
	if HaveBomb(cards) {
		return FourBomb
	}
	// 牛牛、牛一 ~ 牛九
	// 从小到大循环找出最大的拼牛组合
	for i := 4; i > 0; i-- {
		for j := i - 1; j > -1; j-- {
			if (sum-CardValue[cards[i]]-CardValue[cards[j]])%10 == 0 {
				if (CardValue[cards[i]]+CardValue[cards[j]])%10 == 0 {
					return OxOx
				} else {
					return (CardValue[cards[i]] + CardValue[cards[j]]) % 10
				}
			}
		}
	}
	// 没牛
	return NoOx
}

func HaveBomb(cards []int) bool {
	// 五张检测两张就可判定是否有炸
	indexGroup := []int{} // 同值牌型索引除4应该在同一组里
	for i := 0; i < 5; i++ {
		indexGroup = append(indexGroup, cards[i]/4)
	}
	if common.Count(indexGroup, cards[0]/4) == 4 {
		return true
	}
	if common.Count(indexGroup, cards[1]/4) == 4 {
		return true
	}
	return false
}

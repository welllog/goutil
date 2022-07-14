package base

import "testing"

func TestRandIncrGenerate_MaxCount(t *testing.T) {
	incr := &RandIncrGenerate{}
	// 前30分钟，每3秒增加1~2次，半小时后增加1~100次
	incr.AddRule(1800, 3, 2, 100)
	// 第一天内，每15秒增加1~3次，一天后增加1~300次
	incr.AddRule(86400, 15, 3, 300)
	// 前5天内，每180秒增加1~4次
	incr.AddRule(5*86400, 180, 4, 0)
	// 前60天内，每600秒增加1~5次
	incr.AddRule(60*86400, 600, 5, 0)

	id := "test"
	t.Log(incr.MaxCount(50 * 86400))
	t.Log(incr.Calculate(id, 50*86400))
	t.Log(incr.MinCount(50 * 86400))
}

package base

import (
	"sort"
)

type RandIncrGenerate struct {
	rules []rule
}

type rule struct {
	gradient     int // 时间梯度，例：3600
	interval     int // 间隔多少秒，例：10
	randMaxBase  int // 达成该梯度最大增加值 例：100
	randMaxMulti int // 每间隔interval增加的最大值 例:2
	quickNum     int
	// 如：前一小时内，每10秒最大增加(1~2),达到一小时增加(1~100)
}

func (r *RandIncrGenerate) AddRule(gradient, interval, maxMulti, maxBase int) {
	r.rules = append(r.rules, rule{
		gradient:     gradient,
		interval:     interval,
		randMaxBase:  maxBase,
		randMaxMulti: maxMulti,
	})
	sort.Slice(r.rules, func(i, j int) bool {
		return r.rules[i].gradient < r.rules[j].gradient
	})
}

func (r *RandIncrGenerate) Calculate(id string, diff int) int {
	if diff <= 0 {
		return 0
	}

	randN := r.fnv32(id)
	var count, lastGradient int
	for _, v := range r.rules {
		multi := r.getRand(randN, v.randMaxMulti)
		if diff < v.gradient {
			return (diff-lastGradient)/v.interval*multi + count
		}
		count += (v.gradient-lastGradient)/v.interval*multi + r.getRand(randN, v.randMaxBase)
		lastGradient = v.gradient
	}
	return count
}

func (r *RandIncrGenerate) MaxCount(diff int) int {
	if diff <= 0 {
		return 0
	}

	var count, lastGradient int
	for _, v := range r.rules {
		if diff < v.gradient {
			return (diff-lastGradient)/v.interval*v.randMaxMulti + count
		}
		count += (v.gradient-lastGradient)/v.interval*v.randMaxMulti + v.randMaxBase
		lastGradient = v.gradient
	}
	return count
}

func (r *RandIncrGenerate) MinCount(diff int) int {
	if diff <= 0 {
		return 0
	}

	var count, lastGradient int
	for _, v := range r.rules {
		if diff < v.gradient {
			return (diff-lastGradient)/v.interval + count
		}
		count += (v.gradient-lastGradient)/v.interval + 1
		lastGradient = v.gradient
	}
	return count
}

func (r *RandIncrGenerate) fnv32(str string) uint32 {
	const preme32 = uint32(16777619)
	hash := uint32(2166136261)
	for i := 0; i < len(str); i++ {
		hash *= preme32
		hash ^= uint32(str[i])
	}
	return hash
}

func (r *RandIncrGenerate) getRand(randn uint32, max int) int {
	if max == 0 {
		return 0
	}
	return int(randn%uint32(max) + 1)
}

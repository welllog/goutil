package base

func ArrayDiff[T comparable](a, b []T) []T {
	bl := len(b)
	if bl == 0 {
		return a
	}

	bm := make(map[T]struct{}, bl)
	for _, bv := range b {
		bm[bv] = struct{}{}
	}

	na := make([]T, 0, len(a))
	for _, av := range a {
		if _, ok := bm[av]; !ok {
			na = append(na, av)
		}
	}
	return na
}

func ArrayDiffReuse[T comparable](a, b []T) []T {
	bl := len(b)
	if bl == 0 {
		return a
	}

	bm := make(map[T]struct{}, bl)
	for _, bv := range b {
		bm[bv] = struct{}{}
	}

	var remain int
	for ai, av := range a {
		if _, ok := bm[av]; !ok {
			a[remain], a[ai] = a[ai], a[remain]
			remain++
		}
	}
	return a[:remain]
}

func ArrayUnique[T comparable](s []T) []T {
	sl := len(s)
	if sl == 0 {
		return s
	}

	keys := make(map[T]struct{}, sl)
	newS := make([]T, 0, sl)
	var keysLength int
	for _, v := range s {
		keys[v] = struct{}{}
		kl := len(keys)
		if keysLength < kl {
			newS = append(newS, v)
			keysLength = kl
		}
	}
	return newS
}

func ArrayUniqueReuse[T comparable](s []T) []T {
	sl := len(s)
	if sl == 0 {
		return s
	}

	keys := make(map[T]struct{}, sl)
	var remain, keysLength int
	for i, v := range s {
		keys[v] = struct{}{}
		kl := len(keys)
		if keysLength < kl {
			s[remain], s[i] = s[i], s[remain]
			keysLength = kl
			remain++
		}
	}
	return s[:remain]
}

// 结果集会包含重复元素
func ArrayIntersect[T comparable](a, b []T) []T {
	al, bl := len(a), len(b)
	if bl == 0 {
		return b
	}

	if al == 0 {
		return a
	}

	min := al
	if al > bl {
		min = bl
		a, b = b, a
	}

	m := make(map[T]struct{}, min)
	for _, v := range a {
		m[v] = struct{}{}
	}

	ns := make([]T, 0, min)
	for _, v := range b {
		if _, ok := m[v]; ok {
			ns = append(ns, v)
		}
	}

	return ns
}

func ArrayIntersectReuse[T comparable](a, b []T) []T {
	al, bl := len(a), len(b)
	if bl == 0 {
		return b
	}

	if al == 0 {
		return a
	}

	min := al
	if al > bl {
		min = bl
		a, b = b, a
	}

	m := make(map[T]struct{}, min)
	for _, v := range a {
		m[v] = struct{}{}
	}

	var remain int
	for i, v := range b {
		if _, ok := m[v]; ok {
			b[remain], b[i] = b[i], b[remain]
			remain++
		}
	}
	return b[:remain]
}

func ArrayEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	b = b[:len(a)]
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

func InArray[T comparable](search T, arr []T) bool {
	for _, v := range arr {
		if search == v {
			return true
		}
	}
	return false
}

func ArraySum[T Number](arr []T) T {
	var sum T
	for _, v := range arr {
		sum += v
	}
	return sum
}

func ArrayChunk[T any](args []T, chunk int) [][]T {
	if chunk < 1 {
		return [][]T{args}
	}

	l := len(args)
	if l <= chunk {
		return [][]T{args}
	}

	n := l / chunk
	arr := make([][]T, 0, n+1)
	var start, end int
	for i := 0; i < n; i++ {
		end = start + chunk
		arr = append(arr, args[start:end])
		start = end
	}
	if l > start {
		arr = append(arr, args[start:])
	}
	return arr
}

func ArrayChunkFunc[T any](args []T, chunk int, fn func([]T) error) error {
	l := len(args)
	if chunk < 1 || l <= chunk {
		return fn(args)
	}

	n := l / chunk
	var start, end int
	for i := 0; i < n; i++ {
		end = start + chunk
		if err := fn(args[start:end]); err != nil {
			return err
		}
		start = end
	}
	if l > start {
		return fn(args[start:])
	}
	return nil
}

func ArrayFilter[T any](args []T, fn func(T) bool) []T {
	var remain int
	for i, v := range args {
		if fn(v) {
			args[remain], args[i] = args[i], args[remain]
			remain++
		}
	}
	return args[:remain]
}

func ArrayCopy[T any](arr []T, start, length int) []T {
	var result []T
	if arr == nil {
		return result
	}

	result = []T{}
	l := len(arr)
	if l == 0 {
		return result
	}
	if start >= l {
		return result
	}
	if length == 0 {
		return result
	}

	if start < 0 {
		start = 0
	}

	max := l - start
	if length < 0 || length >= max {
		result = make([]T, max)
		copy(result, arr[start:])
		return result
	}

	result = make([]T, length)
	copy(result, arr[start:start+length])
	return result
}

package slices

type Number interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8 | ~uint | ~uint64 | ~uint32 | ~uint16 | ~uint8 | float32 | float64
}

// Diff returns the elements in s1 that aren't in s2.
func Diff[T comparable](s1, s2 []T) []T {
	if len(s1) == 0 {
		return nil
	}

	m := make(map[T]struct{}, len(s2))
	for _, v := range s2 {
		m[v] = struct{}{}
	}

	var diff []T
	for _, v := range s1 {
		if _, ok := m[v]; !ok {
			diff = append(diff, v)
		}
	}
	return diff
}

// DiffReuse returns the elements in s1 that aren't in s2. And it will reuse s1's memory.
func DiffReuse[T comparable](s1, s2 []T) []T {
	if len(s2) == 0 {
		return s1
	}

	m := make(map[T]struct{}, len(s2))
	for _, v := range s2 {
		m[v] = struct{}{}
	}

	var remain int
	for i, v := range s1 {
		if _, ok := m[v]; !ok {
			s1[remain], s1[i] = s1[i], s1[remain]
			remain++
		}
	}
	return s1[:remain]
}

// Unique returns the unique elements in s.
func Unique[T comparable](s []T) []T {
	if len(s) == 0 {
		return nil
	}

	seen := make(map[T]struct{}, len(s))
	var unique []T
	var uniqueCount int
	for _, v := range s {
		seen[v] = struct{}{}
		if uniqueCount < len(seen) {
			unique = append(unique, v)
			uniqueCount = len(seen)
		}
	}
	return unique
}

// UniqueInPlace returns the unique elements in s. And it will reuse s's memory.
func UniqueInPlace[T comparable](s []T) []T {
	if len(s) == 0 {
		return s
	}

	seen := make(map[T]struct{}, len(s))
	var remain, uniqueCount int
	for i, v := range s {
		seen[v] = struct{}{}
		if uniqueCount < len(seen) {
			s[remain], s[i] = s[i], s[remain]
			uniqueCount = len(seen)
			remain++
		}
	}
	return s[:remain]
}

// Intersect returns the elements in s1 elements that are also in s2.
func Intersect[T comparable](s1, s2 []T) []T {
	if len(s1) == 0 || len(s2) == 0 {
		return nil
	}

	m := make(map[T]struct{}, len(s2))
	for _, v := range s2 {
		m[v] = struct{}{}
	}

	var intersect []T
	for _, v := range s1 {
		if _, ok := m[v]; ok {
			intersect = append(intersect, v)
		}
	}

	return intersect
}

// IntersectReuse returns the elements in s1 elements that are also in s2. And it will reuse s1's memory.
func IntersectReuse[T comparable](s1, s2 []T) []T {
	if len(s1) == 0 {
		return s1
	}

	m := make(map[T]struct{}, len(s2))
	for _, v := range s2 {
		m[v] = struct{}{}
	}

	var remain int
	for i, v := range s1 {
		if _, ok := m[v]; ok {
			s1[remain], s1[i] = s1[i], s1[remain]
			remain++
		}
	}
	return s1[:remain]
}

// Filter returns a slice holding only the elements of s that satisfy predicate.
func Filter[T any](s []T, predicate func(T) bool) []T {
	var filtered []T
	for _, v := range s {
		if predicate(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

// FilterInPlace returns a slice holding only the elements of s that satisfy predicate, And it will reuse s's memory.
func FilterInPlace[T any](s []T, predicate func(T) bool) []T {
	var remain int
	for i, v := range s {
		if predicate(v) {
			s[remain], s[i] = s[i], s[remain]
			remain++
		}
	}
	return s[:remain]
}

// Equal reports whether s1 and s2 are equal.
func Equal[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}

	if len(s1) == 0 {
		return true
	}

	s2 = s2[:len(s1)]
	for i, v := range s1 {
		if v != s2[i] {
			return false
		}
	}

	return true
}

// Index returns the index of the first instance of v in s, or -1 if v is not present in s.
func Index[T comparable](s []T, v T) int {
	for i, e := range s {
		if e == v {
			return i
		}
	}
	return -1
}

// Contains reports whether v is present in s.
func Contains[T comparable](s []T, v T) bool {
	return Index(s, v) >= 0
}

// Sum returns the sum of all elements in s.
func Sum[T Number](s []T) T {
	var sum T
	for _, v := range s {
		sum += v
	}
	return sum
}

// Max returns the maximum value in s.
func Max[T Number](s ...T) T {
	if len(s) == 0 {
		return 0
	}

	max := s[0]
	for _, v := range s[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// Min returns the minimum value in s.
func Min[T Number](s ...T) T {
	if len(s) == 0 {
		return 0
	}

	min := s[0]
	for _, v := range s[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// Chunk returns a slice holding the chunked elements of s.
func Chunk[T any](s []T, chunkSize int) [][]T {
	if chunkSize < 1 {
		return [][]T{s}
	}

	if len(s) <= chunkSize {
		return [][]T{s}
	}

	n := len(s) / chunkSize
	chunks := make([][]T, 0, n+1)
	var start, end int
	for i := 0; i < n; i++ {
		end = start + chunkSize
		chunks = append(chunks, s[start:end])
		start = end
	}
	if len(s) > start {
		chunks = append(chunks, s[start:])
	}
	return chunks
}

// ChunkFunc calls fn for each chunk of s.
func ChunkFunc[T any](s []T, chunkSize int, fn func([]T) error) error {
	if chunkSize < 1 || len(s) <= chunkSize {
		return fn(s)
	}

	n := len(s) / chunkSize
	var start, end int
	for i := 0; i < n; i++ {
		end = start + chunkSize
		if err := fn(s[start:end]); err != nil {
			return err
		}
		start = end
	}
	if len(s) > start {
		return fn(s[start:])
	}
	return nil
}

// Clone returns a new slice holding the elements of s.
func Clone[T any](s []T, start, length int) []T {
	l := len(s)
	if l == 0 || start >= l || length == 0 {
		return []T{}
	}

	if start < 0 {
		start = 0
	}

	max := l - start
	if length < 0 || length > max {
		length = max
	}

	return append([]T(nil), s[start:start+length]...)
}

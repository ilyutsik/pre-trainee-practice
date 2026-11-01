package pipeline

func Double(in []int) []int {
	if in == nil {
		return nil
	}
	resSlice := make([]int, len(in))
	for i, v := range in {
		resSlice[i] = v * 2
	}
	return resSlice
}

func FilterMin(in []int, min int) []int {
	if in == nil {
		return nil
	}
	resSlice := make([]int, 0, len(in))
	for _, v := range in {
		if v >= min {
			resSlice = append(resSlice, v)
		}
	}
	return resSlice
}

func UniqueStable(in []int) []int {
	if in == nil {
		return nil
	}
	uniqueConstrains := make(map[int]struct{}, len(in))
	resSlice := make([]int, 0, len(in))
	for _, v := range in {
		if _, exist := uniqueConstrains[v]; exist {
			continue
		}
		uniqueConstrains[v] = struct{}{}
		resSlice = append(resSlice, v)
	}
	return resSlice
}

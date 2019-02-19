package via

// Size provides sorter interface for Plans slice type
type Size PlanSlice

func (s Size) Len() int {
	return len(s)
}
func (s Size) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s Size) Less(i, j int) bool {
	return s[i].Size > s[j].Size
}

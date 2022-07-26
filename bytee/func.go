package bytee

func copyss() {
	var s1 = [2][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9}}
	var s = make([]int, 9)
	var l int
	for _, v := range s1 {
		copy(s[l:], v)
		l = len(v)
	}
}

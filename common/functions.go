package common

//传入顺序切片值和查找数，返回index，不存在返回0
func BinarySearch(i []int, j int) int {
	res := i[len(i)/2]
	index := len(i) / 2
	if res == j {
		return index
	} else if res > j {
		return BinarySearch(i[:index-1], j) + index
	} else if res < j {
		return BinarySearch(i[index:], j) + index
	}
	return 0
}
func test() {

}

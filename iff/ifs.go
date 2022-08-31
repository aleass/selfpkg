package iff

func ff() {
	//&&优先
	//if a() && b() || c() && d() {} //(a() && b()) || (c() && d())

	//if a() || (b() && c()) {} //a() || (b() && c())
	if a() && b() || c() { //(a() && b())|| c()
	}

}
func a() bool {
	print(1)
	return true
}
func b() bool {
	print(2)
	return false
}
func c() bool {
	print(3)
	return true
}
func d() bool {
	print(4)
	return false
}

package arrayy

var ms = [3]string{
	1: "gzl",
	2: "zmj",
}
var str = [12]uint8{0, 3, 6}
var se = "gzlzmj"

//数组字符串组合时，减少内存分配
//以及 c8791538315e633a41461173964a99bd90e103e3 cmd/compile/internal/syntax/scanner.go:397
/*
cmd/compile/internal/syntax: use stringer for operators and tokens

With its new -linecomment flag, it is now possible to use stringer on
values whose strings aren't valid identifiers. This is the case with
tokens and operators in Go.

Operator alredy had inline comments with each operator's string
representation; only minor modifications were needed. The inline
comments were added to each of the token names, using the same strategy.

Comments that were previously inline or part of the string arrays were
moved to the line immediately before the name they correspond to.

Finally, declare tokStrFast as a function that uses the generated arrays
directly. Avoiding the branch and strconv call means that we avoid a
performance regression in the scanner, perhaps due to the lack of
mid-stack inlining.

Performance is not affected. Measured with 'go test -run StdLib -fast'
on an X1 Carbon Gen2 (i5-4300U @ 1.90GHz, 8GB RAM, SSD), the best of 5
runs before and after the changes are:
*/
func array() {
	var tok = uint8(1)
	println(se[str[tok-1]:str[tok]])
	println(ms[tok])
	tok = 2
	println(se[str[tok-1]:str[tok]])
	println(ms[tok])
}

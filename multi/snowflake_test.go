package multi

import (
	"testing"
)

func BenchmarkSnowFlake(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SnowFlake()
	}
}

/*
goos: darwin
goarch: arm64
pkg: selfpkg/multi
BenchmarkSnowFlake
BenchmarkSnowFlake-8   	19841487	        59.29 ns/op
PASS
*/

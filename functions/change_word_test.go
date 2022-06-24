package functions

import (
	"fmt"
	"testing"
)

func BenchmarkChangeWord(b *testing.B) {
	b.ReportAllocs()
	b.Run("TestNameObjectPost", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ChangeWord("TestNameObjectPost") //BenchmarkChangeWord-4            9419270               119.9 ns/op            72 B/op          2 allocs/op
		}
	})
	b.Run("test_name_object_post", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ChangeWord("test_name_object_post") //BenchmarkChangeWord-4            5911789               196.2 ns/op           136 B/op          3 allocs/op
		}
	})
	b.Run("testnameobjectpost", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ChangeWord("testnameobjectpost") //BenchmarkChangeWord-4           10577607               105.5 ns/op            88 B/op          3 allocs/op
		}
	})
}

func TestChangeWord(t *testing.T) {
	fmt.Println(ChangeWord("test_name_object_post"))
}

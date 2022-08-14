package testt

import (
	"fmt"
	"testing"
)

func TestT1(t *testing.T) {
	fmt.Println("t11")
}

func TestT22(t *testing.T) {
	fmt.Println("t22")
}

func BenchmarkName(b *testing.B) {
	b.ReportAllocs()
	b.Run("test1", func(b *testing.B) {

	})
	b.Run("test1", func(b *testing.B) {

	})
}
func TestDo(t *testing.T) {

}

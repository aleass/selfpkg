package main

import (
	"testing"
)

func BenchmarkLinkSlice(b *testing.B) {
	b.Run("slice", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var g2 = guard2{}
			for j := 0; j < 10000; j++ {
				g2.Add(&sse{
					data: j,
				})
			}
			for id, v := range g2.data {
				if v.data%2 == 0 {
					g2.Del(id)
				}
			}

			for id, v := range g2.data {
				if v == nil {
					continue
				}
				g2.Del(id)
			}
		}

	})

	b.Run("link", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var g = guard{}
			for j := 0; j < 10000; j++ {
				g.Add(&sse{
					data: i,
				})
			}
			res := g.data
			for ; res != nil; res = res.next {
				if res.data%2 == 0 {
					g.Del(res)
				}
			}
			res = g.data
			for ; res != nil; res = res.next {
				g.Del(res)
			}
		}

	})

}

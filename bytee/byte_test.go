package bytee

import (
	"bytes"
	"testing"
	"unsafe"
)

// go test -fuzz FuzzGuids
func FuzzGuids(f *testing.F) {
	f.Fuzz(func(t *testing.T, ids int64) {
		guid := EncodeGuid2(ids)
		_ids := DecodeGuid(guid)
		if ids != _ids {
			t.Errorf("guid:%d != ids:%d,res:%d", guid, _ids, ids-_ids)
		}
	})
}

func BenchmarkAddUrl(b *testing.B) {
	data := []byte(`["/aaaaa/bbbbb/20220718/cccc/.JPG", "/dddd/eee/20220718/eeeee.jpeg"]`)
	b.Run("a", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			AddOssUrlFast(data, i)
		}
	})
	b.Run("a3", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			AddOssUrlSlow(data)
		}
	})
}
func BenchmarkLinkSlice(b *testing.B) {
	var bb = []byte{1, 2, 3, 4, 5, 6, 2, 7, 8, 9}
	b.Run("bb", func(b *testing.B) {
		b.ReportAllocs()
		sb := bytes.Buffer{}
		for i := 0; i < b.N; i++ {
			sb.Write(bb)
		}
	})

	b.Run("aa", func(b *testing.B) {
		b.ReportAllocs()
		sb := bytes.Buffer{}
		for i := 0; i < b.N; i++ {
			var p = (*Buffer)(unsafe.Pointer(&sb))
			(*p).buf = bb
		}
	})

}

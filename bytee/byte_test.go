package bytee

import "testing"

// go test -fuzz FuzzGuids
func FuzzGuids(f *testing.F) {
	f.Fuzz(func(t *testing.T, ids int64) {
		guid := EncodeGuid(ids)
		_ids := DecodeGuid(guid)
		if ids != _ids {
			t.Errorf("guid:%d != ids:%d,res:%d", guid, _ids, ids-_ids)
		}
	})
}


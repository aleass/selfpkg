package functions

import "strings"

//a - z 97 122
//A - Z 65 90
func ChangeWord(word string) string {
	if res := strings.Split(word, "_"); len(res) > 0 { // _
		words := make([]byte, 0, len(word)+20)
		for _, v := range res {
			if v == "" { //_xxx_xxx
				continue
			}
			if v[0] < 65 || v[0] > 122 || (v[0] > 90 && v[0] < 97) {
				word += v
				continue
			}
			words = append(words, v[0]-32)
			words = append(words, v[1:]...)
		}
		return string(words)
	} else if word[0] >= 97 && word[0] <= 122 { //a
		word = string(word[0]-32) + word[1:]
	} else { //A
		bytes := []byte(word)
		_vb := make([]byte, 0, len(bytes)+20)
		for i, v := range bytes {
			if v < 65 || v > 122 || (v > 90 && v < 97) {
				_vb = append(_vb, v)
				continue
			}
			if i == 0 {
				v = v + 32
			}
			if v < 97 {
				v = v + 32
				_vb = append(_vb, '_')
			}
			_vb = append(_vb, v)
		}
		word = string(_vb)
	}
	return word
}

package hashh

// hash is a perfect hash function for keywords.
// It assumes that s has at least length 2.
var keywordMap [1 << 6]int

func hash(s []byte) uint {
	return (uint(s[0])<<4 ^ uint(s[1]) + uint(len(s))) & uint(len(keywordMap)-1)
}

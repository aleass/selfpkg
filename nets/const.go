package nets

const (
	//pack size
	PackSize      = 4 //自身header和包长信息
	VerSize       = 2 //version
	RawHeaderSize = VerSize + PackSize

	HeaderOffset = 0
	VerOffset    = 4
)

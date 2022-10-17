package hashh

/*
// 数据库exchange相关转换
var exchangeNameCode = map[string]int{
	//股票类型，1上证，2深证，3北证
	"sh":     1,
	"sz":     2,
	"bj":     3,
	"hu":     4,
	"comex":  5,
	"global": 6,
}

var exchangeNameCode2 = [...]string{
	//股票类型，1上证，2深证，3北证
	1: "sh",
	2: "sz",
	3: "bj",
	4: "hu",
	5: "comex",
	6: "global",
}


func GetMysqlCodeByExcName(exchangeName string) int {
	return exchangeNameCode[exchangeName]
}

func GetMysqlExcNameByCode(code int) string {
	return exchangeNameCode2[code]
}
*/

type stock int

var (
	_stock_name  = "shszbjhucomexglobal"
	_stock_index = [...]int{0, 2, 4, 6, 8, 13, 19}
	stockMap     [1 << 3]stock
)

const (
	_ stock = iota
	sh
	sz
	bj
	hu
	comex
	global
	_max
)

// bkdr hash
func hash(str string) (hash int) {
	for i := 0; i < len(str); i++ {
		hash = (hash * 131) + int(str[i])
	}
	return hash & (len(stockMap) - 1)
}

//减少内存
func (i stock) String() string {
	if i < 1 || i > stock(len(_stock_index)-1) {
		return ""
	}
	return _stock_name[_stock_index[i-1]:_stock_index[i]]
}

func init() {
	for i := sh; i < _max; i++ {
		stockMap[hash(i.String())] = i
	}
}

// 获取股票id
func GetMysqlCodeByExcName(exchangeName string) int {
	return int(stockMap[hash(exchangeName)])
}

// 获取股票
func GetMysqlExcNameByCode(code int) string {
	return stock(code).String()
}

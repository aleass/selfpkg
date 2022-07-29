package pprof

import (
	"encoding/json"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/go-redis/redis"
	"go.uber.org/zap/buffer"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const (
	frontSize         = 13.0 // 顶部和底部字体大小
	numHigSize        = 10   //默认数字的高度
	logoSize          = 30   //长方体内部logo大小
	baseStockDrawSize = 240  //长方体内部股票数据的基础宽度像素值，默认：240(4小时的分钟)
	spaceSize         = 7    //空格size
	baseTimeHig       = 5    //底部x轴的时间的轴长度
	frontFile         = "1.ttc"
	dataSize          = 241 //min = 4 * 60 +1 (1为930)
	frontFile2        = "11.ttf"

	outFile       = "stock.png"
	textInterSize = 3.5 //文字间隔size
	numSize       = 12  //左右侧数字文字大小
	RootPath      = ""
)

//颜色
const (
	white     = "#FFFFFF"
	black     = "#000000"
	lightGray = "#C1CDCD"
	blue      = "#1084CB"
	lightBlue = "#E1FFFF"
	yellow    = "#EBCE85"
	red       = "#fd211c"
	green     = "#008000"
)

//左侧k线颜色
var color = map[int]string{
	0: red,
	1: red,
	2: black,
	3: green,
	4: green,
}

//服务器json文件
type sortData struct {
	data string
	date int
}

type StockSimple struct {
	Id           string `json:"id" gorm:"-;primary_key;AUTO_INCREMENT"`
	Name         string `json:"name" bson:"name"`
	Code         string `json:"code" bson:"code"`
	SimpleName   string `json:"simple_name" bson:"simple_name"`
	ExchangeCode int64  `json:"exchange_code"`
	CreateTime   int64  `json:"create_time" bson:"create_time"`
	UpdateTime   int64  `json:"update_time" bson:"update_time"`
	Kind         int64  `json:"kind" bson:"kind"`
}

var Redis *redis.Client

func init() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("192.168.11.60:6379"),
		Password: "AaB123456./", // no password set
		DB:       0,             // use default DB
		Network:  "tcp",
		PoolSize: 50,
	})
}

type historyStock struct {
	Code int      `json:"code"`
	Data []string `json:"data"`
}

//max, min, dayStore, avgStore
func dayStock(code, date string, model *StockSimple) (
	float64, float64, []float64, []float64, float64, float64, float64, error) {
	var sortSlice []sortData
	if date == "" {
		key := fmt.Sprintf("trend:%s", strings.Replace(code, "-", ":", 1))
		data := Redis.HGetAll(key)
		sortSlice = make([]sortData, 0, len(data.Val()))
		for k, v := range data.Val() {
			k = k[8:]
			_date, _ := strconv.Atoi(k)
			sortSlice = append(sortSlice, sortData{
				data: v,
				date: _date,
			})
		}

		sort.Slice(sortSlice, func(i, j int) bool {
			if sortSlice[i].date > sortSlice[j].date {
				return false
			}
			return true
		})
	} else {
		hostname := "https://api.dev.jz3377.com/"
		fullPath := hostname + path.Join("trend", code, date+".json")
		res, err := http.Get(fullPath)
		if err != nil {
			return 0, 0, nil, nil, 0, 0, 0, err
		}
		if res.StatusCode != 200 {
			return 0, 0, nil, nil, 0, 0, 0, err
		}
		bytedata, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return 0, 0, nil, nil, 0, 0, 0, err
		}
		var data historyStock
		json.Unmarshal(bytedata, &data)

		if len(data.Data) == 0 {
			return 0, 0, nil, nil, 0, 0, 0, err
		}
		sortSlice = make([]sortData, len(data.Data))
		for i := 0; i < len(data.Data); i++ {
			if index := strings.Index(data.Data[i], ","); index != -1 {
				sortSlice[i] = sortData{data: data.Data[i][index+1:]}
			}
		}
	}

	var (
		minTop, minLow, turnover, volume, max, min, maxVolume, maxTurnover float64
		minVolume, minTurnover, nowVolume, nowTurnover                     float64
		dayStore, avgStore                                                 = make([]float64, dataSize), make([]float64, dataSize)
		i                                                                  int
	)

	index := 0 //同样的数据json字符串比redis字符串需要少偏移1
	if date == "" {
		index = 1
	}

	for ; i < len(sortSlice) && i < dataSize; i++ {
		v := sortSlice[i]
		res := strings.Split(v.data, ",")
		stocks, _ := strconv.ParseFloat(res[0], 64)           //价格
		minTop, _ = strconv.ParseFloat(res[1+index], 64)      //分钟内最高
		minLow, _ = strconv.ParseFloat(res[2+index], 64)      //分钟内最低
		nowVolume, _ = strconv.ParseFloat(res[3+index], 64)   //成交量
		nowTurnover, _ = strconv.ParseFloat(res[4+index], 64) //成交额

		var _avg float64
		if model.Kind == 1 {
			if nowTurnover != 0 && nowVolume != 0 {
				turnover += nowTurnover
				volume += nowVolume
				_avg = turnover / volume / 100
			}
		} else {
			turnover += stocks
			_avg = stocks
			if i != 0 {
				_avg = turnover / float64(i+1)
			}
		}
		avgStore[i] = _avg

		//获取最值
		if max < minTop || max == 0 {
			max = minTop
		}
		if min > minLow || min == 0 {
			min = minLow
		}
		if nowVolume > maxVolume {
			maxVolume = nowVolume
		}
		if nowTurnover > maxTurnover {
			maxTurnover = nowTurnover
		}

		if nowVolume < minVolume || minVolume == 0 {
			minVolume = nowVolume
		}
		if nowTurnover < minTurnover || minTurnover == 0 {
			minTurnover = nowTurnover
		}

		if max < _avg && _avg != 0 {
			max = _avg
		}
		if min > _avg && _avg != 0 {
			min = _avg
		}

		dayStore[i] = stocks
	}
	dayStore = dayStore[:i]
	avgStore = avgStore[:i]
	return max, min, dayStore, avgStore, turnover, volume, nowVolume, nil
}

func Draw(code string, open float64, stock *StockSimple) ([]byte, error) {
	size := float64(480)
	var imageWidth, imageHeight = size, size / 1.6 //比例 16:10
	var dayStore, avgStore []float64

	max, min, dayStore, avgStore, turnover, volume, nowVolume, _ := dayStock(code, "", stock)

	if len(dayStore) == 0 {
		return nil, nil
	}
	return run(imageWidth, imageHeight, dayStore, avgStore, max, min, open, turnover, volume, nowVolume, stock, "")
	//	return  outFile
}

/*
* @DrawImage 绘制图片
* @param imageWidth 长方形宽度
* @param imageHeight 长方形高度
* @param dayStore 日指数数组
* @param avgStore 平均指数数组
* @param max 最大指数
* @param min 最小指数
* @param Opening 开盘值
 */
func run(imageWidth, imageHeight float64, dayStore, avgStore []float64, max, min, opening, turnover, volume,
	nowVolume float64, model *StockSimple, date string) ([]byte, error) {
	c := opening / 100
	_max, _min := max, min
	if opening > 10 {
		c /= 10
	}
	//指数值较大需要循环直到<10
	for c > 10 {
		c /= 10
	}
	max += c
	min -= c

	/*--------------------------------------------------创建画布
	 */
	dc := gg.NewContext(int(imageWidth), int(imageHeight))
	dc.SetHexColor(white) // 设置颜色 白
	dc.Clear()            // 使用当前颜色设置背景色
	//设置长方体的起点
	size, _ := dc.MeasureString(fmt.Sprintf("%0.f", opening))
	//字体的长高度
	var originX, originY = size + textInterSize*5, frontSize * 2 //这里的高度 = 信息的条数 * spaceSize

	//设置长方形的长宽高
	size, _ = dc.MeasureString("-00.00%") //右侧涨幅长度
	informWidth := imageWidth - size - originX
	informHeight := imageHeight - frontSize*2 - spaceSize - originY - numHigSize

	/*--------------------------------------------------设置画布的底部时间
	 */
	spacing := informWidth / 2
	high := originY + informHeight + numHigSize + baseTimeHig*2 //numHigSize + spaceSize
	dc.SetHexColor(black)
	var tickArr = [5]string{"09:30", "10:30", "11:30/13:00", "14:00", "15:00"}
	spacing = informWidth / 4
	var strX float64
	var x = originX
	for i, v := range tickArr {
		strX, _ = dc.MeasureString(v)
		switch i {
		case 0:
			x = originX
		case 4:
			x = originX + float64(i)*spacing - strX
		default:
			x = originX + float64(i)*spacing - strX/2
		}
		dc.DrawString(v, x, high)
	}

	/*--------------------------------------------------底部x轴时间的轴线和长方体内时间的y轴
	 */
	spacing = informWidth / 4 //设置间距，此处4等分长方体的长
	for i := 0.0; i < 5; i++ {
		dc.SetHexColor(lightGray)
		dc.SetLineWidth(1) // 设置画笔宽度为5像素
		//定义一个像素表示一分钟，一个时间轴间隔是一个钟，因此此处*60
		y := informHeight + originY
		dc.DrawLine(originX+i*spacing+textInterSize, y, originX+i*spacing+textInterSize, y+baseTimeHig)
		dc.Stroke()

		dc.SetLineWidth(0.5)      // 设置画笔宽度像素
		dc.SetHexColor(lightGray) //
		y = originY
		dc.DrawLine(originX+i*spacing+textInterSize, y, originX+i*spacing+textInterSize, originY+informHeight)
		dc.Stroke()
	}
	/*--------------------------------------------------设置左侧的y轴指数
	 */
	var ranges float64
	var index = make([]float64, 5) //[2]int{指数值,对应的颜色}
	if opening >= max {            //open >= max >= min
		ranges = (opening - min) / 2
	} else if opening <= min { //max >= min >= open
		ranges = float64(max-opening) / 2
	} else { //情况1
		res := max - opening
		res2 := opening - min
		if res > res2 {
			ranges = res / 2
		} else {
			ranges = res2 / 2
		}
	}
	index[0] = opening + ranges*2
	index[1] = opening + ranges
	index[2] = opening
	index[3] = opening - ranges
	index[4] = opening - ranges*2

	//开始画指数
	var brokenLine float64
	spacing = informHeight / 4 //设置间距，此处4等分长方体的长
	dc.SetLineWidth(1)         // 设置画笔宽度像素
	var beforeFunc [5]func()
	var val = "+" //给正数加 + 符号
	for i, v := range index {
		dc.SetHexColor(color[i]) //
		brokenLine = originY + spacing*float64(i)
		if i == 2 {
			_brokenLine := brokenLine
			beforeFunc[i] = func() {
				var x = originX
				dc.DrawLine(x+textInterSize, _brokenLine, informWidth+originX+textInterSize, _brokenLine)
				dc.Stroke()
			}
			val = ""
		} else {
			_brokenLine := brokenLine
			beforeFunc[i] = func() {
				var x = originX
				for ; x < informWidth+originX; x += 6 { //
					dc.SetHexColor(lightGray) //
					dc.DrawLine(x+textInterSize, _brokenLine, x+textInterSize+2, _brokenLine)
				}
				dc.Stroke()
			}
		}
		if i == 2 {
			dc.DrawString(" 0.00%", informWidth+originX+textInterSize*2, originY+spacing*float64(i)+5) //右侧信息
		} else {
			dc.DrawString(val+strconv.FormatFloat((v-opening)/opening*100, 'f', 2, 64)+"%",
				informWidth+originX+textInterSize*2, originY+spacing*float64(i)+5) //右侧信息
		}
		dc.DrawString(fmt.Sprintf("%0.1f", v), textInterSize, originY+spacing*float64(i)+5) //左侧信息
	}

	min, max = index[4], index[0]
	single := informHeight / (max - min)

	//设置左侧的y轴指数对应的y轴线
	for _, f := range beforeFunc {
		f()
	}

	/*--------------------------------------------------指数折线信息区域
	 */
	dc.SetLineWidth(1) // 设置画笔宽度像素

	//计算填充线高度
	x = originX
	spacing = (informWidth - 2) / baseStockDrawSize //减去两边宽度
	pre := dayStore[0]
	bottomY := make([]float64, 0, len(dayStore))
	bottomY = append(bottomY, informHeight+originY+(min-pre)*single)
	for i := 0; i < len(dayStore); i++ {
		trend := dayStore[i]
		y := informHeight + originY + (min-pre)*single
		bottomY = append(bottomY, y)
		pre = trend
	}

	//填充行情线折线图
	x = originX
	dc.SetLineWidth(1)   // 设置画笔宽度像素
	dc.SetHexColor(blue) //
	pre = dayStore[0]
	var YList = make([]float64, len(dayStore))
	for i := 0; i < len(dayStore); i++ {
		trend := dayStore[i]
		y := informHeight + originY + (min-pre)*single
		dc.DrawLine(x+textInterSize+1, y, x+spacing+textInterSize, informHeight+originY+(min-trend)*single) //+1 减去宽度
		YList[i] = y
		x += spacing
		pre = trend
	}
	dc.Stroke()

	//行情线下的蓝色区域填色
	yLine := informHeight + originY //长方形的y轴值
	dc.SetHexColor(lightBlue)
	x = originX
	for _, y := range YList {
		dc.DrawLine(x+textInterSize+1, y+1, x+textInterSize, yLine)
		x += spacing
	}
	dc.Stroke()

	//填充均价折线图
	dc.SetLineWidth(1)     // 设置画笔宽度像素
	dc.SetHexColor(yellow) //
	pre = avgStore[0]
	x = originX
	for i := 0; i < len(avgStore); i++ {
		trend := avgStore[i]
		if trend == 0 {
			continue
		}
		y := informHeight + originY + (min-pre)*single
		dc.DrawLine(x+textInterSize+1, y, x+spacing+textInterSize, informHeight+originY+(min-trend)*single)
		x += spacing
		pre = trend
	}
	dc.Stroke()

	/*--------------------------------------------------外面的长方形
	 */
	dc.SetHexColor(lightGray)                                                   //
	dc.DrawRectangle(originX+textInterSize, originY, informWidth, informHeight) //
	dc.SetLineWidth(1)                                                          // 设置画笔宽度为5像素
	dc.StrokePreserve()                                                         // 使用当前颜色（红）描出当前路径（矩形），但不删除当前路径
	dc.Stroke()

	/*--------------------------------------------------设置字体
	 */
	//设置公司名称
	path := frontFile
	if err := dc.LoadFontFace(path, logoSize); err != nil {
	}
	strX, strY := dc.MeasureString("盛世创富")
	var ttfSize = informWidth/2 + originX - strX/2
	dc.DrawString("盛世创富", ttfSize, strY/2+originY+informHeight/2) // 直接将文字贴入画布中

	//设置股票对应的产品
	if err := dc.LoadFontFace(frontFile, frontSize); err != nil {
	}

	/*--------------------------------------------------顶部信息
	 */
	dc.SetHexColor(color[2]) //改为黑色
	var latest = dayStore[len(dayStore)-1]
	info := fmt.Sprintf("%s[%s] %0.2f", model.Name, model.Code, latest)
	strX, strY = dc.MeasureString(info)
	dc.DrawString(info, textInterSize, strY/3*2+spaceSize)

	rate := (latest - opening) / opening * 100
	runes := ""
	if rate > 0 {
		dc.SetHexColor(color[1])
		runes = "+"
	} else if rate == 0 {
		dc.SetHexColor(color[2])
	} else {
		dc.SetHexColor(color[3])
	}
	info = runes + fmt.Sprintf("%0.2f%%", rate)

	x = textInterSize + strX + spaceSize
	strX, strY = dc.MeasureString(info)
	dc.DrawString(info, x, strY/3*2+spaceSize)

	dc.SetHexColor(color[2]) //改为黑色
	info = fmt.Sprintf("成交量:%s", ChangeUtil(nowVolume))
	x += strX + spaceSize
	_, strY = dc.MeasureString(info)

	dc.DrawString(info, x, strY/3*2+spaceSize)

	/*--------------------------------------------------设置顶部右边的时间
	 */
	dc.SetHexColor(black)
	strTime := date
	if strTime == "" {
		now := time.Now()
		if now.Hour() >= 15 { //15:00以后
			strTime = now.Format("2006-01-02 ") + "15:00" // 直接将文字贴入画布中
		} else if now.Hour() > 9 || (now.Hour() == 9 && now.Minute() >= 25) { //股市时间内
			//在11:30-13:00范围内
			if now.Hour() < 13 && ((now.Hour() == 11 && now.Minute() > 30) || now.Hour() > 11) {
				strTime = now.Format("2006-01-02 ") + "11:30" // 直接将文字贴入画布中
			} else {
				strTime = now.Format("2006-01-02 15:04") // 直接将文字贴入画布中
			}
		} else { //9:30以前
			strTime = now.AddDate(0, 0, -1).Format("01-02 ") + "15:00"
		}
	}

	strX, strY = dc.MeasureString(strTime)
	dc.DrawString(strTime, imageWidth-strX-textInterSize, strY/3*2+spaceSize) // 直接将文字贴入画布中

	/*--------------------------------------------------底部信息
	 */
	//if err := dc.LoadFontFace(frontFile, frontSize); err != nil {
	//}
	dc.SetHexColor(black)

	topRate := (_max - opening) / opening * 100
	lowRate := (_min - opening) / opening * 100

	info = fmt.Sprintf("最高：%.2f", _max)
	residual := high + numHigSize + spaceSize
	x = textInterSize
	dc.DrawString(info, x, residual)

	runes = ""
	if topRate > 0 {
		dc.SetHexColor(color[1])
		runes = "+"
	} else if rate == 0 {
		dc.SetHexColor(color[2])
	} else {
		dc.SetHexColor(color[3])
	}

	strX, _ = dc.MeasureString(info)
	x += strX + spaceSize
	info = runes + fmt.Sprintf("%.2f", topRate)
	dc.DrawString(info, x, residual)

	dc.SetHexColor(black)
	strX, _ = dc.MeasureString(info)
	x += strX + spaceSize
	info = fmt.Sprintf("最低：%.2f", _min)
	dc.DrawString(info, x, residual)

	runes = ""
	if lowRate > 0 {
		dc.SetHexColor(color[1])
		runes = "+"
	} else if rate == 0 {
		dc.SetHexColor(color[2])
	} else {
		dc.SetHexColor(color[3])
	}
	strX, _ = dc.MeasureString(info)
	x += strX + spaceSize
	info = runes + fmt.Sprintf("%.2f", lowRate)
	dc.DrawString(info, x, residual)

	dc.SetHexColor(black)
	strX, _ = dc.MeasureString(info)
	x += strX + spaceSize
	info = fmt.Sprintf("成交量：%s手 成交额：%s", ChangeUtil(volume), ChangeUtil(turnover))
	dc.DrawString(info, x, residual)

	/*--------------------------------------------------左右侧数字
	 */

	//避免频繁生成删除io文件操作，此处直接调用内部函数生成byte流返回到前端
	//dc.SavePNG(RootPath + outFile)
	var f = buffer.Buffer{}
	p := unsafe.Pointer(dc)
	im := *(**image.RGBA)(unsafe.Add(p, uintptr(24)))
	png.Encode(&f, im)
	return f.Bytes(), nil
}

func ChangeUtil(v float64) string {
	if v < 1e4 {
		return fmt.Sprintf("%0.2f", v)
	}
	if v < 1e8 {
		v = v / 1e4
		return fmt.Sprintf("%0.2f%s", v, "万")
	}
	v = v / 1e8
	return fmt.Sprintf("%0.2f%s", v, "亿")
}

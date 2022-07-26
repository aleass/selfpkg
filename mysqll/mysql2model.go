package mysqll

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"os"
	"strings"
)

var db *gorm.DB

const (
	dir = "model/"
)

func dbfile() {
	var (
		host     = "192.168.11.60:3306"
		user     = "stock_user"
		pass     = "sfsctlrw92ijywi"
		database = "stock_spider"
		pack     = "model"
		table    = ""
		err      error
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=True&loc=Local", user, pass, host, database)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	if table == "" {
		type Field struct {
			Field string `gorm:"column:Tables_in_stock_spider"`
		}
		var fieldObj []Field
		sql := fmt.Sprintf("show tables")
		db.Raw(sql).Scan(&fieldObj)
		for _, v := range fieldObj {
			ToAFile(pack, v.Field)
		}
		return
	}
	ToAFile(pack, table)
}
func ToAFile(pack, table string) {
	type Field struct {
		Field   string `gorm:"column:Field"`
		Type    string `gorm:"column:Type"`
		Comment string `gorm:"column:Comment"`
	}
	var fieldObj []Field
	file := ""
	sql := fmt.Sprintf("SELECT TABLE_NAME, COLUMN_COMMENT as Comment, COLUMN_NAME as Field, DATA_TYPE as Type from INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='%s';", table)
	db.Raw(sql).Scan(&fieldObj)
	for _, value := range fieldObj {
		file += fmt.Sprintf("\n\t %s %s `json:\"%s\"` // %s", toUp(value.Field), getType(value.Type), value.Field, value.Comment)
	}
	infos := fmt.Sprintf(`package %s
type %s struct {%s
}
`, pack, toUp(table), file)

	infos += fmt.Sprintf(`
func (%s) TableName() string {
	return "%s"
}`, toUp(table), table)

	fileName := fmt.Sprintf(dir+"%s.go", table)
	f, e := os.Create(fileName)
	if e != nil {
		fmt.Println("打开文件错误", e)
		return
	}
	_, we := io.WriteString(f, infos)
	if we != nil {
		fmt.Println("写入文件错误", we)
		return
	}

}

func toUp(field string) string {
	var nextUp bool
	str := ""
	for key, value := range field {
		if key == 0 {
			str = str + strings.ToUpper(string(value))
			continue
		}
		if string(value) == "_" {
			nextUp = true
			continue
		}
		if nextUp {
			str = str + strings.ToUpper(string(value))
			nextUp = false
		} else {
			str = str + string(value)
		}
	}

	return str

}

var m = map[string]string{
	"tinyint":    "int64",
	"smallint":   "int64",
	"mediumint":  "int64",
	"int":        "int64",
	"bigint":     "int64",
	"float":      "float64",
	"decimal":    "string",
	"bit":        "string",
	"year":       "string",
	"time":       "string",
	"date":       "string",
	"datetime":   "string",
	"timestamp":  "string",
	"char":       "string",
	"varchar":    "string",
	"tinytext":   "string",
	"text":       "string",
	"mediumtext": "string",
	"longtext":   "string",
	"enum":       "string",
	// 其他类型默认转字符
}

func getType(typeString string) string {

	if val, ex := m[typeString]; ex {
		return val
	} else {
		return "string"
	}
}

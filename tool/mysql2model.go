package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"os"
	"os/exec"
	"strings"
)

var db *gorm.DB
var Configer *Config

type App struct {
	Host     string `toml:"host"`
	User     string `toml:"user"`
	Pass     string `toml:"pass"`
	Database string `toml:"database"`
	Pack     string `toml:"pack"`
	SaveDir  string `toml:"saveDir"`
}
type Config struct {
	App App `toml:"app"`
}

func init() {
	Configer = &Config{}
	_, err := toml.DecodeFile("config.toml", Configer)
	if err != nil {
		panic("读取配置文件失败!,原因:" + err.Error())
	}
}

func dbfile() {
	var (
		host     = Configer.App.Host
		user     = Configer.App.User
		pass     = Configer.App.Pass
		database = Configer.App.Database
		pack     = Configer.App.Pack
		path     = Configer.App.SaveDir
		err      error
		table    string
	)

	flag.StringVar(&table, "t", "", "表名字")
	flag.Parse()
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=True&loc=Local", user, pass, host, database)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	type tabField struct {
		Name    string `gorm:"column:Name"`
		Comment string `gorm:"column:Comment"`
	}
	sql := "show table status;"
	var tableObj []tabField
	db.Raw(sql).Scan(&tableObj)
	var tableName = make(map[string]string, len(tableObj))
	for _, v := range tableObj {
		if v.Comment == "" {
			continue
		}
		tableName[v.Name] = "\n" + "//" + v.Comment
	}
	if table == "" {
		type Field struct {
			Field string `gorm:"column:table_name"`
		}
		var fieldObj []Field
		sql = fmt.Sprintf("select table_name  from INFORMATION_SCHEMA.TABLES where table_schema = '%s'", database)
		db.Raw(sql).Scan(&fieldObj)
		for _, v := range fieldObj {
			ToAFile(pack, v.Field, database, path, tableName)
		}
		return
	}
	ToAFile(pack, table, database, path, tableName)
	c := exec.Command("go", "fmt", "-n", path)
	c.String()
	err = c.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
	c.Wait()
}
func ToAFile(pack, table, dbs, path string, tabMap map[string]string) {
	if table == "comment" {
		println()
	}
	type Field struct {
		Field   string `gorm:"column:Field"`
		Type    string `gorm:"column:Type"`
		Comment string `gorm:"column:Comment"`
	}
	var fieldObj []Field
	file := ""
	sql := fmt.Sprintf("SELECT TABLE_NAME, COLUMN_COMMENT as Comment, COLUMN_NAME as Field, DATA_TYPE as Type "+
		"from INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='%s' and TABLE_SCHEMA = '%s';", table, dbs)
	db.Raw(sql).Scan(&fieldObj)
	var length int
	for _, value := range fieldObj {
		if len(value.Field) > length {
			length = len(value.Field)
		}
	}
	var space = "                                                                                                      "
	for _, value := range fieldObj {
		spaces := space[:length-len(value.Field)]
		file += fmt.Sprintf("\n\t %s %s `gorm:\"column:%s\" %sdesc:\"%s\"`", toUp(value.Field), getType(value.Type), value.Field, spaces, value.Comment)
	}

	infos := fmt.Sprintf(`package %s
%s
type %s struct {%s
}
`, pack, tabMap[table], toUp(table), file)

	infos += fmt.Sprintf(`
func (%s) TableName() string {
	return "%s"
}`, toUp(table), table)

	fileName := fmt.Sprintf(path+"%s.go", table)
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
	f.Close()
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

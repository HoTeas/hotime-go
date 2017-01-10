package hotime

import (

	//"time"
	//	"fmt"
	"reflect"
	"strconv"
	"strings"
	//"errors"
	//"math/rand"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"net/http"
	"math"
)

//字符串截取
func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

//获取最后出现字符串的下标
//return  找不到返回 -1
func IndexLastStr(str, sep string) int {
	sepSlice := []rune(sep)
	strSlice := []rune(str)
	if len(sepSlice) > len(strSlice) {
		return -1
	}

	v := sepSlice[len(sepSlice)-1]

	for i := len(strSlice) - 1; i >= 0; i-- {
		vs := strSlice[i]
		if v == vs {
			j := len(sepSlice) - 2
			for ; j >= 0; j-- {
				vj := sepSlice[j]
				vsj := strSlice[i-(len(sepSlice)-j-1)]
				if vj != vsj {
					break
				}
			}
			if j < 0 {
				return i - len(sepSlice) + 1
			}
		}

	}
	return -1
}

func ObjToStr(obj interface{}) string {
	//	fmt.Println(reflect.ValueOf(obj).Type().String() )
	str := ""
	switch reflect.ValueOf(obj).Type().String() {
	case "int":
		str = strconv.Itoa(obj.(int))
	case "uint8":
		str = obj.(string)
	case "int64":
		str = strconv.FormatInt(obj.(int64), 10)
	case "[]byte":
		str = string(obj.([]byte))
	case "string":
		str = obj.(string)
	}

	return str
}

//md5
func Md5(req string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(req))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

//随机数
func Rand(count int) int {

	res := 0
	for i := 0; i < count; i++ {
		res = res * 10
		res = res + R.Intn(10)
		if i == 0 && res == 0 {
			for {
				res = res + R.Intn(10)
				if res != 0 {
					break
				}
			}
		}
	}
	return res
}

//随机数范围
func RandX(small int, max int) int {
	res := 0
	if small == max {
		return small
	}

	for {
		res = R.Intn(max)
		if res >= small {
			break
		}
	}
	return res
}

//字符串转int
func StrToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	return i, err
}

//路由
func Router(ctr CtrInterface) {

	str := reflect.ValueOf(ctr).Type().String()
	a := strings.IndexByte(str, '.')

	c := len(str) - len("Ctr") - 1

	app := Substr(str, 0, a)    //属于哪个app
	ct := Substr(str, a+1, c-a) //属于哪个控制器
	var x = map[string]CtrInterface{}
	if _, ok := Proj[app]; ok {
		//存在APP
		x = Proj[app]
	} else {
		x = map[string]CtrInterface{}
	}
	x[ct] = ctr //将控制器存入APP
	//	fmt.Println(c)
	Proj[app] = x //将APP存入测试
}

func RunMethodListener(test func(app []string)){
	RunMethodListenerFunc=test
}

func SetDb(db *sql.DB) {
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.Ping()
	SqlDB = &*db
	GetDb()
}

//获取数据库
func GetDb() (HoTimeDB, error) {
	Db.DB = &*SqlDB
	return Db, nil
}

//初始化方法
func Init() {


	http.HandleFunc("/", PublicCore.myHandler)
	http.ListenAndServe(":"+Config["port"].(string), nil)
}

//设置Config
func SetCfg(tData Map) {
	for k, v := range tData {
		Config[k] = v
	}
}

//浮点数四舍五入保留小数
func Round(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}

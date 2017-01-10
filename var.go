package hotime

import (
	"time"
	"math/rand"
	"database/sql"
)


//数据库实例
var SqlDB *sql.DB
var Db = HoTimeDB{}

//实例化核心库
var PublicCore Core = Core{}

//随机对象
var R = rand.New(rand.NewSource(time.Now().UnixNano()))
var RunMethodListenerFunc func(app []string)
//控制器接口
type CtrInterface interface {
}
//hotime的常用map
type Map map[string]interface{}
type Slice []interface{}

//程序的项目
var Proj = map[string]map[string]CtrInterface{}

type ModelInterface interface {

}

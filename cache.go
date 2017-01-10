package hotime

import (
	"time"
	"sync"
	"reflect"
)

type cache struct {
	time int64
	data interface{}
}

var cacheMap = map[interface{}]cache{}
//读写锁
var M sync.RWMutex;
//获取Cache键只能为string类型
func getCache(key string) interface{} {
	M.RLock();
	defer M.RUnlock();
	data, ok := cacheMap[key];

	if !ok {
		return nil
	}
	//data:=cacheMap[key];
	if data.time <= time.Now().Unix() {
		delete(cacheMap, key)
		return nil;
	}

	return data.data
}
//key value ,时间为时间戳
func setCache(key string, value interface{}, time int64) {
	M.Lock();
	defer M.Unlock();
	data := cache{time:time, data:value}
	cacheMap[key] = data
}

func Cache(key interface{}, data ...interface{}) interface{} {
	if (len(data) == 0) {
		return getCache(ObjToStr(key))
	}
	tim := time.Now().Unix()
	if (len(data) == 1) {
		tim += int64(Config["cacheTime"].(int))
	}
	if (len(data) == 2) {
		tim += int64(data[1].(int))
	}
	setCache(ObjToStr(key), data[0], tim)
	return nil
}

//剔除等译uid的缓存
func DeleteByUid(uid int64) {
	M.Lock();
	defer M.Unlock();
	for k,v := range cacheMap{
		if reflect.TypeOf(v.data).String() != "hotime.Map" {
			continue;
		}

		if v.data.(Map)["uid"] == nil {
			continue;
		}
		if v.data.(Map)["uid"].(int64) == uid {
			delete(cacheMap, k);
			Db.Delete("cached",Map{"ckey":k});
			break;
		}
	}
}


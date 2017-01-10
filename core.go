package hotime

import (
	//"fmt"
	"os"
	//"runtime"
	"reflect"
	"strings"
	"strconv"
	"io/ioutil"
	"net/http"

	"encoding/json"
	"time"
)

type Core struct{
	
}


var sessionTemp string
//网络错误
func (this *Core) webError(code int, w http.ResponseWriter) {
	w.WriteHeader(code)
}

//路由处理
func (this *Core) Router(w *http.ResponseWriter,req *http.Request)bool{

	q := strings.Index(req.RequestURI, "?")
	if q == -1 {
		q = len(req.RequestURI)
	}
	o := Substr(req.RequestURI, 0, q)

	r := strings.SplitN(o, "/", -1)

	var s =make([]string,0)

	for i := 0; i < len(r); i++ {
		if !strings.EqualFold("", r[i]) {
			s=append(s,r[i])
		}
	}

	if len(s) != 3 {
		return false
	}

	if _, ok := Proj[s[0]]; !ok {
		return false
	}

	//存在APP
	x := Proj[s[0]]
	if _, ok := x[strings.Title(s[1])]; !ok {

		return false
	}
	m := x[strings.Title(s[1])]

	n := reflect.ValueOf(m)

	n=reflect.New(n.Type())
	//fmt.Println(n.MethodByName("Test").Type().String())
	//for i:=0;i<n.NumMethod();i++{
	//
	//	fmt.Println(runtime.FuncForPC(n.Pointer()).Name())
	//}

	//n=reflect.MakeChan(n.Type,1)
	//fmt.Println(n)
	y := n.MethodByName(strings.Title(s[2]))

	if !y.IsValid() {
		return false
	}
	if RunMethodListenerFunc!=nil{
		RunMethodListenerFunc(s)
	}

	//获取cookie
	// 如果cookie存在直接将sessionId赋值为cookie.Value
	// 如果cookie不存在就查找传入的参数中是否有token
	// 如果token不存在就生成随机的sessionId
	// 如果token存在就判断token是否在Session中有保存
	// 如果有取出token并复制给cookie
	// 没有保存就生成随机的session
	cookie,err:=req.Cookie((Config["sessionName"]).(string))
	sessionId:=Md5(strconv.Itoa(Rand(10)))
	token:= req.FormValue("token")

	if err!=nil ||(len(token)==32&&cookie.Value!=token){
		if len(token)==32{
			sessionId=token
		}
		http.SetCookie(*w,&http.Cookie{Name:Config["sessionName"].(string),Value:sessionId,Path:"/"})
	}else{
		sessionId=cookie.Value
	}


	sessionObj:=Cache(Config["sessionName"].(string)+":"+sessionId)
	if sessionObj==nil{
		//获取最新的json数据
		sessionJson:=Db.Get("cached","cvalue",Map{"ckey":Config["sessionName"].(string)+":"+sessionId})
		sessionObj=map[string]interface{}{}

		if sessionJson!=nil{
			json.Unmarshal([]byte(sessionJson["cvalue"].(string)),&sessionObj)
			tt:=sessionObj.(map[string]interface{})
			for k,v:=range tt{
				switch reflect.TypeOf(v).String() {
				case "int":tt[k]=int64(v.(int))
				case "float64":dt,_:=StrToInt(strconv.FormatFloat(v.(float64),'f',0,32));tt[k]=int64(dt)
				}
			}
			sessionObj=tt
		}



	}

	data,_:=json.Marshal(sessionObj)
	sessionTemp=string(data)
	n.Elem().FieldByName("SessionId").Set(reflect.ValueOf(sessionId))
	n.Elem().FieldByName("Session").Set(reflect.ValueOf(sessionObj))
	n.Elem().FieldByName("Request").Set(reflect.ValueOf(&*req))
	n.Elem().FieldByName("ResponseWriter").Set(reflect.ValueOf(*w))
	n.MethodByName("Init").Call(nil)

	y.Call(nil)
	//dd := n.Elem().FieldByName("Data").Bytes()

	Db.Delete("cached",Map{"endtime[<]":time.Now().Unix()})
	sessionObj=n.Elem().FieldByName("Session").Interface().(Map)
	data,_=json.Marshal(sessionObj)
	tp:=string(data)
	if tp!=sessionTemp{
		res:=Db.Update("cached",Map{"cvalue":tp,"endtime":time.Now().Unix()+int64(Config["sessionTime"].(int))},Map{"ckey":Config["sessionName"].(string)+":"+sessionId})
		if res==int64(0){
			Db.Insert("cached",Map{"ckey":Config["sessionName"].(string)+":"+sessionId,"cvalue":tp,"time":time.Now().Unix(),"endtime":time.Now().Unix()+int64(Config["sessionTime"].(int))})
		}
	}
	Cache(Config["sessionName"].(string)+":"+sessionId,n.Elem().FieldByName("Session").Interface(),Config["sessionTime"])

	return true
}

//网络请求的Handler
func (this *Core) myHandler(w http.ResponseWriter, req *http.Request) {

	result := this.Router(&w,req)
	if result != false {
		//mtd := result.([]byte)
		//w.Write(mtd)
		return
	}

	q := strings.Index(req.RequestURI, "?")
	if q == -1 {
		q = len(req.RequestURI)
	}
	o := Substr(req.RequestURI, 0, q)
	//url赋值
	path := Config["tpt"].(string) + o



	//判断是否为默认
	if path[len(path)-1] == '/' {
		defFile := Config["defFile"].([]string)
		for i := 0; i < len(defFile); i++ {
			temp := path + defFile[i]
			_, err := os.Stat(temp)

			if err == nil {
				path = temp
				break
			}

		}
		if path[len(path)-1] == '/' {
			this.webError(403, w)
			return ;
		}
	}
	if strings.Contains(path, "/.") {
		this.webError(403, w)
		return
	}



	var data []byte
	var err error
	var ok bool
	if data, ok = Config["tpt:"+path].([]byte); !ok {
		//不存在缓存
		data, err = ioutil.ReadFile(path)
		if err != nil {
			this.webError(404, w)
			return
		}
		if Config["cached"].(bool) {
			Config["tpt:"+path] = data
		}




	}
	header:=w.Header()

	//类型判断并设置Content-Type
	if strings.Contains(path,".css"){
		header.Set("Content-Type","text/css")
		header.Get(Config["sessionName"].(string))

	}

	w.Write(data)
}




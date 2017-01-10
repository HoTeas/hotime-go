package hotime

var Config = Map{
	"dbHost":"localhost",
	"dbName":"test",
	"dbUser":"root",
	"dbPwd":"root",
	"dbPort":    "3306",
	"port":"8080",
	"tpt":     "tpt",
	"cacheTime":2*60*60,//缓存时间
	"sessionName":"HOTIME",
	"sessionTime":14*24*60*60,//session保存时间
	"defFile": []string{"index.html", "index.htm"},
}
//session
//var Session=map[string]map[string]interface{}{}
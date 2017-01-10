package hotime

import (
	"net/http"
	"encoding/json"
)

type Ctr struct {
	Mdl
	//Data    []byte
	Request *http.Request
	ResponseWriter http.ResponseWriter
	SessionId string
	Session Map
}

func (this *Ctr)Display(statu int, data interface{}) {

	resp:=Map{"statu": statu}
	if statu!=0{
			temp:=map[string]interface{}{}
			temp["type"]=Config["error"].(map[int]string)[statu]
			temp["msg"]=data
			resp["result"]= temp
	}else{
		resp["result"]=data
	}


	d, err := json.Marshal(resp)
	if err != nil {
		return
	}

	this.ResponseWriter.Write(d)
	//this.Data=d;
}





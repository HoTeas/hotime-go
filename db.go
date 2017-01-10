package hotime

import (
	"reflect"
	"database/sql"
	//"fmt"
	//"strconv"
	"strings"
	//"sort"

)






type HoTimeDB struct {
	*sql.DB
	LastQuery string
	LastData []interface{}
	LastErr error
	limit Slice
}

func (this *HoTimeDB)Page(page,pageRow int)  *HoTimeDB{
	page=(page-1)*pageRow
	if page<0{
		page=1
	}

	this.limit=Slice{page,pageRow}
	return this
}

func (this *HoTimeDB)PageSelect(table string,qu ...interface{})[]Map{

	if len(qu)==1{
		qu=append(qu,Map{"LIMIT":this.limit})
	}
	if len(qu)==2{
		temp:=qu[1].(Map)
		temp["LIMIT"]=this.limit
		qu[1]=temp
	}
	if len(qu)==3{
		temp:=qu[2].(Map)
		temp["LIMIT"]=this.limit
		qu[2]=temp
	}
	//fmt.Println(qu)
	data:=this.Select(table,qu...)

	return data
}
//数据库数据解析
func (this *HoTimeDB)Row(resl *sql.Rows) []Map{
	dest:=make([]Map,0)
	strs, _ := resl.Columns()

	for i := 0; resl.Next(); i++ {
		lis:=make(Map,0)
		a:=make([]interface{},len(strs))

		b:=make([]interface{},len(a))
		for j:=0;j<len(a);j++{
			b[j]=&a[j]
		}
		resl.Scan(b...)
		for j:=0;j<len(a);j++{
			if a[j]!=nil&&reflect.ValueOf(a[j]).Type().String()=="[]uint8"{
				lis[strs[j]]=string(a[j].([]byte))
			}else{
				lis[strs[j]]=a[j]
			}

		}


		dest=append(dest,lis)

	}

	return dest
}

func (this *HoTimeDB)Exec(query string,args... interface{})(sql.Result,error){
	resl,err:=this.DB.Exec(query,args...)
	this.LastErr=err
	return resl,err
}
func (this *HoTimeDB)Query(query string,args... interface{})[]Map{
	var err error
	var resl *sql.Rows


	resl,err=this.DB.Query(query,args...)
	this.LastErr=err

	if err!=nil {
		return nil
	}


	return this.Row(resl)
}

func (this *HoTimeDB)Select(table string,qu ...interface{})[]Map{

	query:="SELECT"
	where :=Map{}
	qs:=make([]interface{},0)
	intQs,intWhere:=0,1
	join:=false
	if(len(qu)==3){
		intQs=1
		intWhere=2
		join=true

	}


	if len(qu)>0{
		if reflect.ValueOf(qu[intQs]).Type().String()=="string"{
			query+=" "+qu[intQs].(string)
		}else{
			for i:=0;i<len(qu[intQs].(Slice));i++{
				if i+1!=len(qu[intQs].(Slice)){
					query+=" `"+qu[intQs].(Slice)[i].(string)+"`,"
				}else{
					query+=" `"+qu[intQs].(Slice)[i].(string)+"`"
				}

			}
		}

		query+=" FROM "+table
	}

	if join{
		for k,v:=range  qu[0].(Map){
			switch Substr(k,0,3) {
			case "[>]":query+=" LEFT JOIN "+Substr(k,3,len(k)-3)+" ON "+v.(string)
			case "[<]":query+=" RIGHT JOIN "+Substr(k,3,len(k)-3)+" ON "+v.(string)
			}
			switch Substr(k,0,4) {
			case "[<>]":query+=" FULL JOIN "+Substr(k,4,len(k)-4)+" ON "+v.(string)
			case "[><]":query+=" INNER JOIN "+Substr(k,4,len(k)-4)+" ON "+v.(string)
			}
		}
	}


	if len(qu)>1{
		where=qu[intWhere].(Map)
	}




	temp,resWhere:=this.where(where)


	query+=temp
	qs=append(qs,resWhere...)

	//fmt.Println(query)
	//fmt.Println(qs)
	this.LastQuery=query
	this.LastData=qs

	res:=this.Query(query,qs...)

	//return 0
	//rows,_:=res.RowsAffected()
	return res


}

func (this *HoTimeDB)Get(table string,qu ...interface{})Map{
	//fmt.Println(qu)
	if len(qu)==1{
		qu=append(qu,Map{"LIMIT":1})
	}
	if len(qu)==2{
		temp:=qu[1].(Map)
		temp["LIMIT"]=1
		qu[1]=temp
	}
	if len(qu)==3{
		temp:=qu[2].(Map)
		temp["LIMIT"]=1
		qu[2]=temp
	}
	//fmt.Println(qu)
	data:=this.Select(table,qu...)
	if len(data)==0{
		return nil
	}
	return data[0]
}

/**
** 计数
 */
func (this *HoTimeDB)Count(table string,qu ...interface{})int{
	req:=[]interface{}{}

	if len(qu)==2{
		req=append(req,qu[0])
		req=append(req,"COUNT(*)")
		req=append(req,qu[1])
	}else{
		req=append(req,"COUNT(*)")
		req=append(req,qu...)
	}

	//req=append(req,qu...)
	data:=this.Select(table,req...)
	//fmt.Println(data)
	if len(data)==0{
		return 0
	}
	//res,_:=StrToInt(data[0]["COUNT(*)"].(string))
	res:=ObjToStr(data[0]["COUNT(*)"])
	count,_ := StrToInt(res)
	return count


}
//更新数据
func  (this *HoTimeDB)Update(table string,data Map,where Map)int64{
	query:="UPDATE "+table+" SET "
	//UPDATE Person SET Address = 'Zhongshan 23', City = 'Nanjing' WHERE LastName = 'Wilson'
	qs:=make([]interface{},0)
	tp:=len(data)

	for k,v:=range data{
		vstr:="?"
		if Substr(k,len(k)-3,3)=="[#]"{
			k=strings.Replace(k,"[#]","",-1)

			vstr=ObjToStr(v)
		}else{
			qs = append(qs, v)
		}
		query+="`"+k+"`="+vstr+""
		if tp--;tp!=0{
			query+=", "
		}
	}
	temp,resWhere:=this.where(where)
	//fmt.Println(resWhere)

	query+=temp
	qs=append(qs,resWhere...)

	//fmt.Println(query)
	//fmt.Println(qs)
	this.LastQuery=query
	this.LastData=qs

	res,err:=this.Exec(query,qs...)

	if err!=nil{
		return 0
	}

	//return 0
	rows,_:=res.RowsAffected()
	return rows
}

var condition=[]string{"AND","OR"}
var vcond=[]string{"GROUP","ORDER","LIMIT"}


//where语句解析
func (this *HoTimeDB)where(data Map) (string,[]interface{}){



	where:=""

	res:=make([]interface{},0)
	//AND OR判断
	for k, v := range data {
		x:=0
		for i:=0;i<len(condition);i++{

			if condition[i]==k{

				tw,ts:=this.cond(k,v.(Map))

				where+=tw
				res=append(res,ts...)
				break
			}
			x++
		}

		y:=0
		for j:=0;j<len(vcond);j++{

			if vcond[j]==k{
				break
			}
			y++
		}
		if x==len(condition)&&y==len(vcond){
			tv,vv:=this.varCond(k,v)

			where+=tv;
			res=append(res,vv...)
		}



	}


	if(len(where)!=0){
		where=" WHERE "+where
	}

	//特殊字符
	for j:=0;j<len(vcond);j++{
		for k,v:=range data{
			if vcond[j]==k{
				if k=="ORDER"{
					where+=" "+k+" BY "
					//fmt.Println(reflect.ValueOf(v).Type())

					//break
				}else if k=="GROUP"{

					where+=" "+k+" BY "
				}else{

					where+=" "+k
				}

				if(reflect.ValueOf(v).Type().String()=="hotime.Slice"){
					for i:=0;i<len(v.(Slice));i++{
						where+=" "+ObjToStr(v.(Slice)[i])

						if len(v.(Slice))!=i+1{
							where+=","
						}

					}
				}else{
					//fmt.Println(v)
					where+=" "+ObjToStr(v)
				}



				break
			}
		}



	}


	return where,res


}

func (this *HoTimeDB)varCond(k string,v interface{}) (string,[]interface{}){
	where:=""
	res:=make([]interface{},0)
	length:=len(k)
	if length>4{
		def:= false
		switch Substr(k,length-3,3) {
		case "[>]":k=strings.Replace(k,"[>]","",-1);where+="`"+k+"`>? ";res=append(res,v)
		case "[<]":k=strings.Replace(k,"[<]","",-1);where+="`"+k+"`<? ";res=append(res,v)
		case "[!]":k=strings.Replace(k,"[!]","",-1);where,res=this.notIn(k,v,where,res)
		case "[#]":k=strings.Replace(k,"[#]","",-1);where+="`"+k+"`="+ObjToStr(v)
		case "[~]":k=strings.Replace(k,"[~]","",-1);where+="`"+k+"` LIKE ? ";v="%"+v.(string)+"%";res=append(res,v)
		default:def=true

		}
		if def{
			switch Substr(k,length-4,4) {
			case "[>=]":k=strings.Replace(k,"[>=]","",-1);where+="`"+k+"`>=? ";res=append(res,v)
			case "[<=]":k=strings.Replace(k,"[<]","",-1);where+="`"+k+"`<=? ";res=append(res,v)
			case "[><]":k=strings.Replace(k,"[><]","",-1);where+="`"+k+"` NOT BETWEEN ? AND ? ";res=append(res,v.(Slice)[0]);res=append(res,v.(Slice)[1]);
			case "[<>]":k=strings.Replace(k,"[<>]","",-1);where+="`"+k+"` BETWEEN ? AND ? ";res=append(res,v.(Slice)[0]);res=append(res,v.(Slice)[1]);
			default:
				if reflect.ValueOf(v).Type().String()=="hotime.Slice"{
					where+="`"+k+"` IN ("
					res=append(res,(v.(Slice))...)
					for i:=0;i<len(v.(Slice));i++{
						if i+1!=len(v.(Slice)){
							where+="?,"
						}else{
							where+="?) "
						}
						//res=append(res,(v.(Slice))[i])
					}
				}else{
					where+="`"+k+"`=? "
					res=append(res,v)
				}

			}
		}



	}else {
		//fmt.Println(reflect.ValueOf(v).Type().String())
		if reflect.ValueOf(v).Type().String()=="hotime.Slice"{
			where+="`"+k+"` IN ("
			res=append(res,(v.(Slice))...)
			for i:=0;i<len(v.(Slice));i++{
				if i+1!=len(v.(Slice)){
					where+="?,"
				}else{
					where+="?) "
				}
				//res=append(res,(v.(Slice))[i])
			}
		}else if reflect.ValueOf(v).Type().String()=="[]interface {}"{
			where+="`"+k+"` IN ("
			res=append(res,(v.([]interface{}))...)
			for i:=0;i<len(v.([]interface{}));i++{
				if i+1!=len(v.([]interface{})){
					where+="?,"
				}else{
					where+="?) "
				}
				//res=append(res,(v.(Slice))[i])
			}
		}else{
			where+="`"+k+"`=? "
			res=append(res,v)
		}
	}

	return where,res
}
//	this.Db.Update("user",hotime.Map{"ustate":"1"},hotime.Map{"AND":hotime.Map{"OR":hotime.Map{"uid":4,"uname":"dasda"}},"ustate":1})
func (this *HoTimeDB)notIn(k string,v interface{},where string,res []interface{})(string,[]interface{}){
	//where:=""
	//fmt.Println(reflect.ValueOf(v).Type().String())
	if reflect.ValueOf(v).Type().String()=="hotime.Slice"{
		where+="`"+k+"` NOT IN ("
		res=append(res,(v.(Slice))...)
		for i:=0;i<len(v.(Slice));i++{
			if i+1!=len(v.(Slice)){
				where+="?,"
			}else{
				where+="?) "
			}
			//res=append(res,(v.(Slice))[i])
		}
	}else if reflect.ValueOf(v).Type().String()=="[]interface {}"{
		where+="`"+k+"` NOT IN ("
		res=append(res,(v.([]interface{}))...)
		for i:=0;i<len(v.([]interface{}));i++{
			if i+1!=len(v.([]interface{})){
				where+="?,"
			}else{
				where+="?) "
			}
			//res=append(res,(v.(Slice))[i])
		}
	}else{
		where+="`"+k+"` !=? "
		res=append(res,v)
	}

	return  where,res
}

func (this *HoTimeDB)cond(tag string,data Map)(string,[]interface{}){
	where:=" "
	res:=make([]interface{},0)
	lens:=len(data)
	//fmt.Println(lens)
	for k,v :=range data  {
		x:=0
		for i:=0;i<len(condition);i++{

			if condition[i]==k{

				tw,ts:=this.cond(k,v.(Map))
				if lens--;lens<=0{
					//fmt.Println(lens)
					where+="("+tw+") "
				}else{
					where+="("+tw+") "+tag+" "
				}

				res=append(res,ts...)
				break
			}
			x++
		}

		if x==len(condition){

			tv,vv:=this.varCond(k,v)

			res=append(res,vv...)
			if lens--;lens<=0{
				where+=tv+""
			}else{
				where+=tv+" "+tag+" "
			}
		}

	}

	return where,res
}

func (this *HoTimeDB)Delete(table string,data map[string]interface{})  int64{
	query:="DELETE FROM "+table+" "

	temp,resWhere:=this.where(data)
	query+=temp

	this.LastQuery=query
	this.LastData=resWhere

	res,err:=this.Exec(query,resWhere...)

	if err!=nil{
		return 0
	}

	//return 0
	rows,_:=res.RowsAffected()
	return rows
}

//插入新数据
func (this *HoTimeDB)Insert(table string,data map[string]interface{})int64{

	if this.DB.Ping()!=nil{
		return 0
	}

	values:=make([]interface{},0)
	queryString:=" ("
	valueString:=" ("

	lens:=len(data)
	tempLen:=0
	for k, v := range data {
		tempLen++
		values=append(values,v)
		if tempLen<lens{
			queryString+="`"+k+"`,"
			valueString+="?,"
		}else{
			queryString+="`"+k+"`) "
			valueString+="?);"
		}

	}
	query:="INSERT INTO "+table+queryString+"VALUES"+valueString
	res,err:=this.DB.Exec(query,values...)

	this.LastQuery=query
	this.LastData=values

	this.LastErr=err
	if(err!=nil){
		return 0
	}
	id,err:=res.LastInsertId()
	if err!=nil {
		return 0
	}
	//fmt.Println(id)
	return id
}
package hotime

type  Mdl struct {
	Db HoTimeDB
}
//初始化mdl
func (this *Mdl)Init(){

	this.Db=Db

}


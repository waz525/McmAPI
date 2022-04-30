/*
数据库操作类
*/
package Utils

import (
	"fmt"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql" //mysql驱动
	_ "github.com/bmizerany/pq"       //postgresql驱动 
	_ "github.com/mattn/go-sqlite3"   //sqlite3驱动
	"github.com/jmoiron/sqlx"
//	"database/sql"
	seelog "github.com/cihub/seelog"
)

//数据库配置参数
type DBMethed struct {
		DBType		string
		Host		string
		Port		string
		Username	string
		Password	string
		Database	string
		DBFile		string
		DBHandle	*sqlx.DB
}



//连接数据库
func (this *DBMethed) ConnDataBase() {
	//默认是mysql数据库
	dbInfo := ""+this.Username+":"+this.Password+"@tcp("+this.Host+":"+this.Port+")/"+this.Database+"" ;
	//postgres 数据库
	if this.DBType == "postgres" {
		dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",this.Host, this.Port, this.Username, this.Password, this.Database )
	}
	//sqlite3数据库
	if this.DBType == "sqlite3" {
		dbInfo = this.DBFile
	}

	DBHandle, err := sqlx.Open(this.DBType, dbInfo)
	if err != nil {
		seelog.Critical("Open "+this.DBType+" failed,", err) ; seelog.Flush()
		return
	}

	DBHandle.SetMaxOpenConns(200)
	DBHandle.SetMaxIdleConns(20)
	this.DBHandle = DBHandle
}

//关闭数据库连接
func (this *DBMethed) Close() {
	if this.DBHandle != nil	{
		this.DBHandle.Close()
	}
}


//查询并转换成Map
func (this *DBMethed) Query(sqlString string)  []map[string]interface{} {
	tableData := make([]map[string]interface{}, 0)
	if this.DBHandle == nil {
		this.ConnDataBase()
	}
	stmt, err := this.DBHandle.Prepare(sqlString)
	if err != nil {
		errinfo := make(map[string]interface{})
		errinfo["ERROR"] = err.Error()
		tableData = append(tableData, errinfo)
		return tableData

	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		errinfo := make(map[string]interface{})
		errinfo["ERROR"] = err.Error()
		tableData = append(tableData, errinfo)
		return tableData

	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		errinfo := make(map[string]interface{})
		errinfo["ERROR"] = err.Error()
		tableData = append(tableData, errinfo)
		return tableData
	}
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	return tableData ;
}

//查询结果转成Json字符串
func (this *DBMethed) Map2Json(tableData []map[string]interface{}, oneFlag bool) string {
	rst := ""
	if oneFlag && len(tableData) > 0 {
		jsonData, err := json.Marshal(tableData[0])
		if err != nil {
			return "{\"ERROR\":\""+err.Error()+"\"}"
		}
		rst = string(jsonData)
	} else {
		jsonData, err := json.Marshal(tableData)
		if err != nil {
			return "{\"ERROR\":\""+err.Error()+"\"}"
		}
		rst = string(jsonData)
	}
	return rst
}

//普通查询数据
func (this *DBMethed) Query2Json(sql string) string {
	return this.Map2Json(this.Query(sql), false)

}

func (this *DBMethed) Query2JsonOne(sql string) string {
	return this.Map2Json(this.Query(sql), true)
}

//根据id查询单条数据
func (this *DBMethed) QueryById(TableName, ObjectID string) string {
	sql := "select * from "+TableName+" where id = '"+ObjectID+"' "
	seelog.Debug("Query Sql:", sql) ; seelog.Flush()
	return this.Map2Json(this.Query(sql), true)
}

//查询数据量
func (this *DBMethed) QueryCount(TableName, WhereStr string) string {
	sql := "select count(1) as COUNT from "+TableName+ " where "+WhereStr ;
	seelog.Debug("QueryCount Sql:", sql); seelog.Flush()
	return this.Map2Json(this.Query(sql), true)

}



//执行sql，包括insert,update,delete,create,alter,drop等
func (this *DBMethed) ExecSql(sql string) string {
	if this.DBHandle == nil {
		this.ConnDataBase()
	}
	seelog.Debug("Exec Sql:", sql); seelog.Flush()
	stmt,err := this.DBHandle.Prepare(sql)
	if err != nil {
		return "{\"ERROR\":\""+err.Error()+"\"}"
	}
	rst, err := stmt.Exec()
	if err != nil {
		return "{\"ERROR\":\""+err.Error()+"\"}"
	}
	num, err := rst.RowsAffected()
	if err != nil {
		return "{\"ERROR\":\""+err.Error()+"\"}"
	}
	return "{\"ChageLine\":"+Itoa(num)+"}"
}

//查询表字段，返回字段数组
func (this *DBMethed) QueryFields(TableName string) string {

	if this.DBHandle == nil {
		this.ConnDataBase()
	}

	sqlString := "select * from "+TableName+" where 1 = 2 "

	seelog.Debug(sqlString) ; seelog.Flush();
	stmt, err := this.DBHandle.Prepare(sqlString)
	if err != nil {
		return "{\"ERROR\":\""+err.Error()+"\"}"
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return "{\"ERROR\":\""+err.Error()+"\"}"
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return "{\"ERROR\":\""+err.Error()+"\"}"
	}

	var rst map[string][]string
	rst = make(map[string][]string)
	rst["Fields"] = columns

	fmt.Println(rst)
	jsonData, err := json.Marshal(rst)
	if err != nil {
		return "{\"ERROR\":\""+err.Error()+"\"}"
	}
	return string(jsonData)

}

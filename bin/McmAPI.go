/*
Name:		McmAPI.go
Create:		20210512
Modify:		20210512
Version:	1.0.0
Auth:		Worden
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
//	"regexp"
	filepath "path"
	"io/ioutil"
	seelog "github.com/cihub/seelog"
	"net/http"
	"net/http/cgi"
//	"encoding/json"
	"BinInfo"
	"./Utils" //通用类

)

var nServerConf Utils.ServerConf
var nDBMethed Utils.DBMethed
var nDBProcess DBProcess

//默认页面返回，404
func HttpDefaultInterface(rw http.ResponseWriter, req *http.Request) {
	url_path := req.URL.Path
	seelog.Info("Receive "+req.Method+" from "+req.RemoteAddr+" ,URL.Path: "+url_path) ; seelog.Flush()
	filePath := Utils.GetRealPath(nServerConf.FileStaticPath+"/"+url_path[len(nServerConf.FileStaticPrefix)-1:])

	//默认路径加上index.html
	if filePath[len(filePath)-1] == '/' {
		filePath = filePath + "index.html"
	}

	if Utils.IsFileExist(filePath) {
		rw.Header().Set("Content-Type", Utils.GetContentTypeByFilePath(filePath))
		fmt.Fprintln(rw, Utils.GetFileContent(filePath) )
	} else {
		seelog.Error("File Path ( "+filePath+" ) not Exist !") ; seelog.Flush()
		http.NotFound(rw, req)
	}
}

//处理cgi-bin
func CgibinInterface(rw http.ResponseWriter, req *http.Request) {
	url_path := req.URL.Path
	seelog.Info("Receive "+req.Method+" from "+req.RemoteAddr+" ,URL.Path: "+url_path) ; seelog.Flush()
	handler := new(cgi.Handler)
	handler.Path = Utils.GetRealPath(nServerConf.CigBinPath)+"/"+url_path[len(nServerConf.CigBinPrefix)-1:]
	//fmt.Println("handler.Path: "+handler.Path)
	handler.Dir = Utils.GetRealPath(nServerConf.CigBinPath)
	handler.ServeHTTP(rw, req)
}

//处理登录验证
func McmUserLoginInterface(rw http.ResponseWriter, req *http.Request) {
	url_path := req.URL.Path
	seelog.Info("Receive "+req.Method+" from "+req.RemoteAddr+" ,URL.Path: "+url_path) ; seelog.Flush()
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Content-Type", "application/json;charset=utf-8")
	if strings.EqualFold(req.Method,"post") {
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(rw, "{\"ERROR\":\"ParseForm() err: %v\"}", err)
			return
		}
		postForm := req.PostForm ;
		seelog.Debug( "req.PostForm: ", postForm) ; seelog.Flush();
		rst := nDBProcess.DoUserLogin(postForm["username"][0], postForm["password"][0])
		seelog.Debug("UserLogin rst: "+rst) ; seelog.Flush()
		fmt.Fprintln(rw, rst)
	}  else {
		http.NotFound(rw, req)
	}

}

//处理McmApi
func McmApiInterface(rw http.ResponseWriter, req *http.Request) {
	url_path := req.URL.Path
	seelog.Info("Receive "+req.Method+" from "+req.RemoteAddr+" ,URL.Path: "+url_path) ; seelog.Flush()
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Content-Type", "application/json;charset=utf-8")
	url_split := strings.Split(url_path, "/")
	var TableName, ObjectID string
	rst := "[]" 
	//获取TableName
	if len(url_split) > 3  {
		TableName = url_split[3]
		seelog.Debug("TableName: "+TableName) ; seelog.Flush();
	}
	//获取ObjectID
	if len(url_split) > 4  {
		ObjectID = url_split[4]
		seelog.Debug("ObjectID: "+ObjectID) ; seelog.Flush();
	}
	if strings.EqualFold(req.Method,"get") {
		//获取filter
		filter := req.URL.Query().Get("filter")
		seelog.Debug("filter: "+filter) ; seelog.Flush();

		if strings.EqualFold(ObjectID, "count") {					//查询数据量
			rst = nDBProcess.DoGetTableCount(TableName, filter)
		} else if strings.EqualFold(ObjectID, "TableFields") {		//查询表结构
			rst = nDBProcess.DoTableFields(TableName)
		} else {													//查询数据
			rst = nDBProcess.DoGetQuery(TableName, ObjectID, filter)
		}
	} else if strings.EqualFold(req.Method,"post") {
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(rw, "{\"ERROR\":\"ParseForm() err: %v\"}", err)
			return
		}
		postForm := req.PostForm ;
		seelog.Debug( "req.PostForm: ", postForm) ; seelog.Flush();
		if strings.EqualFold(ObjectID, "CreateTable") {			//查询表
			rst = nDBProcess.DoCreateTable(TableName, postForm)
		} else if strings.EqualFold(ObjectID, "DropTable") {	//删除表
			rst = nDBProcess.DoDropTable(TableName, postForm)
		} else if strings.EqualFold(ObjectID, "AddField") {		//增加字段
			rst = nDBProcess.DoAddField(TableName, postForm)
		} else if strings.EqualFold(ObjectID, "DelField") {		//删除字段
			rst = nDBProcess.DoDelField(TableName, postForm)
		} else {													//对表数据的增删改
			rst = nDBProcess.DoTablePost(TableName, ObjectID, postForm)
		}

	}
	seelog.Debug("API Result: "+rst) ; seelog.Flush()
	fmt.Fprintln(rw, rst)
}

//文件删除服务接口
func McmFileDeleteInterface(rw http.ResponseWriter, req *http.Request) {
	url_path := req.URL.Path
	seelog.Info("Receive "+req.Method+" from "+req.RemoteAddr+" ,URL.Path: "+url_path) ; seelog.Flush()
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Content-Type", "application/json;charset=utf-8")
	url_split := strings.Split(url_path, "/")
	rst := ""

	//获取ObjectID
	if len(url_split) > 4  {
		ObjectID := url_split[4]
		seelog.Debug("ObjectID: "+ObjectID) ; seelog.Flush();
		rst = nDBProcess.DoDeleteFile("file", ObjectID)
	} else {
		rst = `{"ERROR":"INVALID FILE ID"}`
	}

	fmt.Fprintln(rw, rst)
}

//文件上传服务接口
func McmFileUploadInterface(rw http.ResponseWriter, req *http.Request) {

	url_path := req.URL.Path
	seelog.Info("Receive "+req.Method+" from "+req.RemoteAddr+" ,URL.Path: "+url_path) ; seelog.Flush()

	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Content-Type", "application/json;charset=utf-8")


	var maxUploadSize int64
	maxUploadSize = 200 * 1024 * 1024 // 200 MB 
	uploadPath := nServerConf.FileStaticPath+"UploadFile/"
	req.Body = http.MaxBytesReader(rw, req.Body, maxUploadSize)
	if err := req.ParseMultipartForm(maxUploadSize); err != nil {
		fmt.Fprintln(rw, `{"ERROR":"FILE TOO BIG, UP TO 200M"}` )
		return
	}
	file, handler, err := req.FormFile("file_upload")
	if err != nil {
		fmt.Fprintln(rw, `{"ERROR":"INVALID FILE"}` )
		return
	}
	defer file.Close()
	//获取文件名
	fileName := handler.Filename
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintln(rw, `{"ERROR":"INVALID FILE"}` )
		return
	}
	/*
	filetype := http.DetectContentType(fileBytes)

	//根据文件类型获取文件后缀
	fileEndings, err := mime.ExtensionsByType(filetype)
	if err != nil {
		renderError(rw, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
		return
	}
	seelog.Debug("Upload FileType:",filetype,", fileEndings:",fileEndings) ; seelog.Flush() ;
	*/
	newPath := Utils.GetRealPath( filepath.Join(uploadPath, fileName) )
	newFile, err := os.Create(newPath)
	if err != nil {
		fmt.Fprintln(rw, `{"ERROR":"CANT WRITE FILE"}` )
		return
	}
	defer newFile.Close()
	if _, err := newFile.Write(fileBytes); err != nil {
		fmt.Fprintln(rw, `{"ERROR":"CANT WRITE FILE"}` )
		return
	}
	seelog.Info("Upload File: ",newPath," success !") ; seelog.Flush() ;
	url := "http://"+req.Host+nServerConf.FileStaticPrefix+"UploadFile/"+fileName
	rst := nDBProcess.DoUploadFile("file", fileName, newPath, url)
	fmt.Fprintln(rw, rst )

}

//开启Http服务
func StartHttpService() {
	http.HandleFunc(nServerConf.CigBinPrefix, CgibinInterface)								//配置cgi处理接口
	http.HandleFunc(nServerConf.FileStaticPrefix+"api/user/login", McmUserLoginInterface)	//配置文件上传处理接口
	http.HandleFunc(nServerConf.FileStaticPrefix+"api/fileUpload", McmFileUploadInterface)	//配置文件上传处理接口
	http.HandleFunc(nServerConf.FileStaticPrefix+"api/fileDelete/", McmFileDeleteInterface)	//配置文件删除处理接口
	http.HandleFunc(nServerConf.FileStaticPrefix+"api/", McmApiInterface)					//配置api处理接口
	http.HandleFunc(nServerConf.FileStaticPrefix, HttpDefaultInterface)						//配置一般http文件处理接口

	//设置静态文件路径，但没有日志，不再使用
	//fs := http.FileServer(http.Dir(nServerConf.FileStaticPath))
	//http.Handle(nServerConf.FileStaticPrefix, http.StripPrefix(nServerConf.FileStaticPrefix, fs))

	seelog.Info("Server Listen on "+nServerConf.ServerPort+" ...")
	seelog.Flush()
	//开启http侦听
	err := http.ListenAndServe(":"+nServerConf.ServerPort, nil)
	if err != nil {
		seelog.Error("ListenAndServe:", err)
		seelog.Flush()
	}
}


func main() {
	var conf_file_name	string
	//读取命令行,-c 配置文件路径
	flag.StringVar(&conf_file_name, "c", "../conf/McmAPI.conf", "Config File Path")
	v := flag.Bool("v", false, "Show Version Info")
	flag.Parse()

	//显示版本信息
	if *v {
		_, _ = fmt.Fprint(os.Stderr, BinInfo.StringifyMultiLine())
		os.Exit(1)
	}


	//读取主配置文件
	ini_parser := Utils.IniParser{}
	if err := ini_parser.Load(conf_file_name); err != nil {
		fmt.Printf("try load config file[%s] error[%s]\n", conf_file_name, err.Error())
		return
	}
	//读取主配置，包括日志配置/侦听端口
	nServerConf.LogConfigPath = ini_parser.GetString("system", "LogConfigPath")
	nServerConf.ServerPort = ini_parser.GetString("system", "ServerPort")
	nServerConf.DBType =ini_parser.GetString("system", "DBType")

	//读取数据库配置
	nDBMethed.DBType = nServerConf.DBType 
	nDBMethed.Host = ini_parser.GetString(nServerConf.DBType, "Host")
	nDBMethed.Port = ini_parser.GetString(nServerConf.DBType, "Port")
	nDBMethed.Username = ini_parser.GetString(nServerConf.DBType, "Username")
	nDBMethed.Password = ini_parser.GetString(nServerConf.DBType, "Password")
	nDBMethed.Database = ini_parser.GetString(nServerConf.DBType, "Database")
	nDBMethed.DBFile = ini_parser.GetString(nServerConf.DBType, "DBFile")
	nDBMethed.ConnDataBase()
	//将nDBMethed赋给nDBProcess作数据处理
	nDBProcess.SetDBMethed(&nDBMethed)

	//读取http配置
	nServerConf.FileStaticPrefix =ini_parser.GetString("http", "FileStaticPrefix")
	nServerConf.FileStaticPath =ini_parser.GetString("http", "FileStaticPath")
	nServerConf.CigBinPrefix = ini_parser.GetString("http", "CigBinPrefix")
	nServerConf.CigBinPath =ini_parser.GetString("http", "CigBinPath")



	//加载日志配置
	logger, err := seelog.LoggerFromConfigAsFile(nServerConf.LogConfigPath)
	if err != nil {
		seelog.Critical("err parsing config log file", err)
		return
	}
	seelog.ReplaceLogger(logger)

	//启动服务侦听
	StartHttpService()
}

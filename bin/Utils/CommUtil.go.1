/*
通用工具
1 常用类型
2 常用函数
*/
package Utils

import (
		"fmt"
		seelog "github.com/cihub/seelog"
		"os"
		"io"
		"time"
		"bytes"
		"strings"
		"strconv"
		"reflect"
		"net/http"
		"os/exec"
		"io/ioutil"
		"encoding/json"
)

/* 获取绝对路径 */
func GetRealPath(path string) string {
		if strings.Index(path,"/") == 0  {
				return path
		} else {
				pwd,_ := os.Getwd()
				return pwd+"/"+path
		}
}

/* 字符串转整型 */
func Atoi(s string) int {
		v, err := strconv.Atoi(s)
		if err != nil {
				seelog.Error("Atoi Error: ",err) ; seelog.Flush()
				return 0
		}
		return v
}

/* 整型转字符型 */
func Itoa( v int ) string {
		return strconv.Itoa(v)
}


/* 执行shell命令 */
func RunShellCmd(cmd string) string{
		//fmt.Println("Running Shell cmd:" , cmd)
		result, err := exec.Command("/bin/sh", "-c", cmd).Output()
		if err != nil {
				seelog.Error("Exec Command Error: ",err) ; seelog.Flush()
		}
		return strings.TrimSpace(string(result))
}

/* 判断文件是否存在  存在返回 true 不存在返回false */
func IsFileExist(filename string) bool {
		var exist = true
		if _, err := os.Stat(filename); os.IsNotExist(err) {
				exist = false
		}
		return exist
}

/* 读取文本文件内容 */
func GetFileContent( filepath string) string {
		bytes, err := ioutil.ReadFile(filepath)

		if err != nil {
				seelog.Error("read file: ",filepath,", error:", err) ; seelog.Flush()
				return ""
		}
		return string(bytes)
}

/* 写入文件 */
func WriteFileContent(filepath, content string) {
		var d1 = []byte(content)
		err := ioutil.WriteFile(filepath, d1, 0666) //写入文件(字节数组)
		if err != nil {
				seelog.Error("Write file: ",filepath,", error:", err) ; seelog.Flush()
		}
}


/* 将字符串转换成二维字符串数组 */
func Str2List(rows string, line_fd string, str_fd string) [][]string {
		rst := [][]string{}
		lines := strings.Split(rows, line_fd)
		for ind := 0 ; ind<len(lines) ;ind++ {
				rst = append(rst, strings.Split(strings.TrimSpace(lines[ind]), str_fd ))
		}
		return rst
}

/* 打印二维数组 */
func PrintList(rows [][]string ) {
		for ind := 0 ; ind<len(rows) ;ind++ {
				line :=  rows[ind]
				for j := 0 ; j <len(line) ; j++ {
						fmt.Printf(""+strconv.Itoa(j)+":"+line[j]+"\t")
				}
				fmt.Println()
		}
}

/* 从二维数组中查找对应数据 */
func FindListCell(rows [][]string, key string , f_index int, r_index int) string {
		for ind := 0 ; ind<len(rows) ;ind++ {
				line :=  rows[ind]
				if len(line) > f_index && len(line) > r_index  {
						if line[f_index] == key {
								return line[r_index]
						}
				}
		}
		return ""
}


/* 查找字符是否在数组中 */
func InArray(obj interface{}, target interface{}) (bool) {
		targetValue := reflect.ValueOf(target)
		switch reflect.TypeOf(target).Kind() {
				case reflect.Slice, reflect.Array:
					for i := 0; i < targetValue.Len(); i++ {
						if targetValue.Index(i).Interface() == obj {
								return true
						}
					}
				case reflect.Map:
						if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
								return true
				}
		}

		return false
}

/* 打印类的所有属性 */
func GetClassAttribute(body interface{}) {
		var prop []string
		refType := reflect.TypeOf(body)
		if refType.Kind() != reflect.Struct {
				fmt.Println("Not a structure type.")
		}
		for i:=0;i<refType.NumField();i++ {
				field := refType.Field(i)
				if field.Anonymous {
						prop = append(prop, field.Name)
						for j := 0;j<field.Type.NumField();j++ {
								prop = append(prop, field.Type.Field(j).Name)
						}
						continue
				}
				prop = append(prop, field.Name)
		}
		fmt.Printf("%v\n", prop)
}

/* 将Struct转成字符串 */
func Struct2String(v interface{}) string {
		res, _ := json.Marshal(v)
		return string(res)
}


/* 发送GET请求 */
func HttpGet(url string) string {
	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
			seelog.Error("HttpGet Error:",err) ; seelog.Flush()
			return ""
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			seelog.Error("HttpGet Error:",err) ; seelog.Flush()
			return ""
		}
	}

	return result.String()
}

/* 发送POST请求 */
func HttpPost(url string, data string, contentType string) string {
	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	//jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(url, contentType, strings.NewReader(data))
	if err != nil {
		seelog.Error("HttpPost Error:",err) ; seelog.Flush()
		return ""
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return string(result)
}



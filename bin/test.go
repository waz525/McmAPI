package main

import (
		"fmt"
		"./Utils"
)


func main() {
	filePath := "/mcm/UploadFile.css"
	rst := Utils.GetContentTypeByFilePath(filePath)
	fmt.Println(rst)
}



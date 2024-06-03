/*
Package gtbox_files 文件处理工具
*/
package gtbox_files

import (
	"github.com/george012/gtbox/gtbox_string"
	"os"
	"path"
)

func GTToolsFileRead(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	return string(content), err
}
func GTToolsFileWrite(filePath, content string) bool {
	err := os.WriteFile(filePath, gtbox_string.GTString2Bytes(content), 0755)

	if err != nil {
		return false
	}
	return true
}

func GTToolsFileRemoveAllInDir(dirPath string) {
	dir, err := os.ReadDir(dirPath)

	if err != nil {
		GTCheckDirisNoneToCreate(dirPath)
		return
	}

	for _, d := range dir {
		filePath := path.Join([]string{dirPath, d.Name()}...)
		os.Remove(filePath)
	}
}

func GTToolsFileRemoveAllKeepRunLog(dirPath string) {
	//dir, err := ioutil.ReadDir(dirPath)
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		GTCheckDirisNoneToCreate(dirPath)
		return
	}

	for _, d := range dir {
		if d.Name() != "run.log" {
			filePath := path.Join([]string{dirPath, d.Name()}...)
			os.Remove(filePath)
		}
	}
}

func GTCheckDirisNoneToCreate(dirName string) {
	_, exist := os.Stat(dirName)
	if os.IsNotExist(exist) {
		os.Mkdir(dirName, os.ModePerm)
	}
}

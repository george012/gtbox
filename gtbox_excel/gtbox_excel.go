/*
Package gtbox_excel Excel处理工具
*/
package gtbox_excel

import (
	"fmt"
	"github.com/george012/gtbox/gtbox_files"
	"github.com/xuri/excelize/v2"
)

// GTToolDataExcelCreateExcelFile
// tableTitle := map[string]string{
func GTToolDataExcelCreateExcelFile(dir_path string, file_name string, tableTitle map[string]string) {
	gtbox_files.GTCheckDirisNoneToCreate(dir_path)
	f := excelize.NewFile()
	// 创建一个工作表
	index, _ := f.NewSheet("Sheet1")
	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	for k, v := range tableTitle {
		f.SetCellValue("Sheet1", k, v)
	}
	// 根据指定路径保存文件
	if err := f.SaveAs(dir_path + "/" + file_name + ".xlsx"); err != nil {
		fmt.Println(err)
	}
}
func GTToolDataExcelAppenData(dir_path string, file_name string, valuesMap map[string]string) {
	excleFile, erra := excelize.OpenFile(dir_path + "/" + file_name + ".xlsx")
	if erra != nil {
		return
	}
	for k, v := range valuesMap {
		excleFile.SetCellValue("Sheet1", k, v)
	}

	// Save spreadsheet by the given path.
	if errb := excleFile.SaveAs(dir_path + "/" + file_name + ".xlsx"); errb != nil {
	}
}

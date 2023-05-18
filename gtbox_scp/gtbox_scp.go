package gtbox_scp

import (
	"context"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/george012/gtbox/gtbox_encryption"
	"github.com/george012/gtbox/gtbox_files"
	"github.com/george012/gtbox/gtbox_log"
	"golang.org/x/crypto/ssh"
	"os"
)

// GTToolsScpSetup 初始化SCP连接
func GTToolsScpSetup(addressIPAndPort, userName, password string) scp.Client {
	clientConfig, _ := auth.PasswordKey(userName, password, ssh.InsecureIgnoreHostKey())
	// Create a new SCP client.
	client := scp.NewClient(addressIPAndPort, &clientConfig)

	err := client.Connect()
	if err != nil {
		gtbox_log.LogErrorf("无法连接服务器: %s", err)
	}
	return client
}

// GTToolsScpDownloadFileAndResult 使用SCP 下载文件
// Return [bool] 是否成功
func GTToolsScpDownloadFileAndResult(addressIPAndPort, userName, password, srcFilePath, saveFilePath string) bool {
	client := GTToolsScpSetup(addressIPAndPort, userName, password)

	// Create a local file to write to.
	f, err := os.OpenFile(saveFilePath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		gtbox_log.LogErrorf("保存文件失败")
		return false
	}
	defer f.Close()

	// Use a file name with exotic characters and spaces in them.
	// If this test works for this, simpler files should not be a problem.
	err = client.CopyFromRemote(context.Background(), f, srcFilePath)
	if err != nil {
		gtbox_log.LogErrorf("下载文件失败：%s", err.Error())
		return false
	}

	noenStr, _ := gtbox_files.GTToolsFileRead(saveFilePath)
	enStr := gtbox_encryption.GTEncryptionGo(noenStr, "iPollo")

	if len(enStr) > 0 {
		isOK := gtbox_files.GTToolsFileWrite(saveFilePath, enStr)
		if isOK == false {
			return false
		}
	}

	defer client.Close()
	return true
}

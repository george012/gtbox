/*
Package gtbox_ssh SSH工具库
*/
package gtbox_ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SSHLoginType int

const (
	SSHLoginTypePWD SSHLoginType = iota
	SSHLoginTypeCert
)

func (logT SSHLoginType) String() string {
	switch logT {
	case SSHLoginTypePWD:
		return "password"
	case SSHLoginTypeCert:
		return "key-cert"
	default:
		return "密码登录"
	}
}

// SSHController 连接Detail信息 struct
type SSHController struct {
	Address     string
	Port        int64
	User        string
	Pwd         string       //如果login_type为证书登录pwd为私钥字符串
	LoginType   SSHLoginType `json:"login_type"`
	Client      *ssh.Client
	Session     *ssh.Session
	LastResult  string
	ConnectLock *sync.RWMutex
	RunCmdLock  *sync.RWMutex
}

// CLIConnectInfo 地址信息
type CLIConnectInfo struct {
	Address   string       `json:"address"`
	Port      int64        `json:"port"`
	User      string       `json:"user"`
	Pwd       string       `json:"pwd"`
	LoginType SSHLoginType `json:"login_type"`
}

// CmdResult 命令行返回值
type CmdResult struct {
	Result string    `json:"Cmd_result"`
	Time   time.Time `json:"time"`
}

// SSHResultData 执行信息
type SSHResultData struct {
	CmdResult  []*CmdResult `json:"cmd_result"`
	CmdStr     string       `json:"cmd_str"`
	ResultTime time.Time    `json:"result_time"`
}

// SSHResultInfo 返回信息Struct
type SSHResultInfo struct {
	Address    string           `json:"address"`
	Port       int64            `json:"port"`
	ResultData []*SSHResultData `json:"result_data"`
}

// Connect 连接远端
func (c *SSHController) Connect() (*SSHController, error) {
	c.ConnectLock.Lock()
	defer c.ConnectLock.Unlock()

	config := &ssh.ClientConfig{}
	config.SetDefaults()
	config.User = c.User

	if c.LoginType == SSHLoginTypeCert {

		signer, s_err := ssh.ParsePrivateKey([]byte(c.Pwd))
		if s_err != nil {
			return c, s_err
		}

		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else if c.LoginType == SSHLoginTypePWD {
		config.Auth = []ssh.AuthMethod{ssh.Password(c.Pwd)}
	}

	config.HostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", c.Address, strconv.FormatInt(c.Port, 10)), config)
	if nil != err {
		return c, err
	}
	c.Client = client
	return c, nil
}

// Run 执行shell
func (c SSHController) Run(shell string) (string, error) {
	c.RunCmdLock.Lock()
	defer c.RunCmdLock.Unlock()
	if c.Client == nil {
		if _, err := c.Connect(); err != nil {
			return "", err
		}
	}
	session, err := c.Client.NewSession()
	if err != nil {
		return "", err
	}
	// 关闭会话
	defer session.Close()
	buf, err := session.CombinedOutput(shell)

	c.LastResult = string(buf)
	return c.LastResult, err

}

// GTSSHClientRun 对一个地址多个命令并发执行并返回
func GTSSHClientRun(cliInfo *CLIConnectInfo, cmds []string, endFunc func(result *SSHResultInfo, err error)) {
	cli := &SSHController{
		Address:     cliInfo.Address,
		Port:        cliInfo.Port,
		User:        cliInfo.User,
		Pwd:         cliInfo.Pwd,
		LoginType:   cliInfo.LoginType,
		RunCmdLock:  &sync.RWMutex{},
		ConnectLock: &sync.RWMutex{},
	}
	// 建立连接对象
	c, aerr := cli.Connect()
	if aerr != nil {
		endFunc(nil, aerr)
		return
	}
	// 退出时关闭连接
	defer c.Client.Close()

	wg := &sync.WaitGroup{}

	aReData := &SSHResultInfo{
		Address: cli.Address,
		Port:    cli.Port,
	}

	for _, acmd := range cmds {
		wg.Add(1)
		aResult := &SSHResultData{
			CmdStr: acmd,
		}
		go func() {
			resSrc, _ := c.Run(acmd)
			for _, va := range strings.Split(resSrc, "\n") {
				if va != "" {
					aResult.CmdResult = append(aResult.CmdResult, &CmdResult{
						Result: va,
						Time:   time.Now(),
					})
				}
			}
			wg.Done()
			aResult.ResultTime = time.Now()
		}()
		aReData.ResultData = append(aReData.ResultData, aResult)
	}
	wg.Wait()

	endFunc(aReData, nil)

}

// GTSSHClientRunDualAddress 对多个地址多个命令并发执行并返回 使用
func GTSSHClientRunDualAddress(cliInfos []*CLIConnectInfo, cmds []string, endFunc func(results []*SSHResultInfo, err error)) {

	var aResultInfo []*SSHResultInfo

	for _, aCliInfo := range cliInfos {
		GTSSHClientRun(aCliInfo, cmds, func(result *SSHResultInfo, err error) {
			aResultInfo = append(aResultInfo, result)
		})

	}

	endFunc(aResultInfo, nil)
}

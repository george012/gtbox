package gtbox_sms

import (
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type GTBoxSMS struct {
	SMSToolsAliYunAccessKeyId     string
	SMSToolsAliYunAccessKeySecret string
	SMSToolsAliYunTemplateCode    string
	SMSToolsAliYunSignName        string
}

var (
	SMSTool *GTBoxSMS
)

// SetupSMSConfigWithAli 初始化阿里云
// accessKeyId  accessKeySecret  templateCode:短信模版 signName 签名
func (aSMS *GTBoxSMS) SetupSMSConfigWithAli(accessKeyId string, accessKeySecret string, templateCode string, signName string) {
	SMSTool = &GTBoxSMS{
		SMSToolsAliYunAccessKeyId:     accessKeyId,
		SMSToolsAliYunAccessKeySecret: accessKeySecret,
		SMSToolsAliYunTemplateCode:    templateCode,
		SMSToolsAliYunSignName:        signName,
	}
}

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func GTCreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

func GTSendSMS(phone string, msgCode string) error {
	// 初始化 Client，采用 AK&SK 鉴权访问的方式，此方式可能会存在泄漏风险，建议使用 STS 方式。鉴权访问方式请参考：https://help.aliyun.com/document_detail/378661.html
	// 获取 AK 链接：https://usercenter.console.aliyun.com
	client, _err := GTCreateClient(tea.String(SMSTool.SMSToolsAliYunAccessKeyId), tea.String(SMSTool.SMSToolsAliYunAccessKeySecret))
	if _err != nil {
		return _err
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(phone),
		SignName:      tea.String(SMSTool.SMSToolsAliYunSignName),
		TemplateCode:  tea.String(SMSTool.SMSToolsAliYunTemplateCode),
		TemplateParam: tea.String(fmt.Sprintf("{\"code\":\"%s\"}", msgCode)),
	}

	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		_, _err = client.SendSmsWithOptions(sendSmsRequest, &util.RuntimeOptions{})
		if _err != nil {
			return _err
		}

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 如有需要，请打印 error
		_, _err = util.AssertAsString(error.Message)
		if _err != nil {
			return _err
		}
	}
	return _err
}

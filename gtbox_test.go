package gtbox

import (
	"github.com/george012/gtbox/gtbox_encryption"
	"github.com/george012/gtbox/gtbox_log"
	"github.com/george012/gtbox/gtbox_time"
	"testing"
)

func TimeToolsTest(t *testing.T) {
	aNow := gtbox_time.NowUTC()
	gtbox_log.LogDebugf("%v", aNow)
	gtbox_log.LogDebugf("11 digits: %d", aNow.Unix())
	gtbox_log.LogDebugf("13 digits: %d", aNow.UnixMilli())
	gtbox_log.LogDebugf("16 digits: %d", aNow.UnixMicro())
	gtbox_log.LogDebugf("19 digits: %d", aNow.UnixNano())
}

func EncryptionToolsTest(t *testing.T) {

	pre_en_str := "test wait encryption string"
	en_key := "test"
	alen := len(pre_en_str)

	gtbox_log.LogDebugf("wait enc str : [%s]%d", pre_en_str, alen)

	en_Str := gtbox_encryption.GTEnc(pre_en_str, en_key)

	gtbox_log.LogDebugf("en_str[%s]%d", en_Str, alen)

	de_str := gtbox_encryption.GTDec(en_Str, en_key)
	b_len := len(de_str)

	gtbox_log.LogDebugf("de_str[%s]%d", de_str, b_len)

	cd_str := gtbox_encryption.GTDecryptionGo(en_Str, en_key)
	cd_len := len(cd_str)
	gtbox_log.LogDebugf("use other func src_str[%s]%d de_str:%s%d", de_str, b_len, cd_str, cd_len)

}

func TestGTBoxFunctions(t *testing.T) {
	//	TODO 初始化gtbox及log分片
	//SetupGTBox("test_gtbox", RunModeRelease, "", 3, gtbox_log.GTLogSaveHours, 60)
	SetupGTBox("test_gtbox", RunModeDebug, "", 3, gtbox_log.GTLogSaveHours, 60)

	gtbox_log.LogDebugf("level  %s", "aleve")
	gtbox_log.LogInfof("aaaaa%s", "qqq")
	gtbox_log.LogErrorf("bbb%s", "ccc")

	TimeToolsTest(t)

	EncryptionToolsTest(t)
}

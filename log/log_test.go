package log

import "testing"

var (
	log = GetLog("test")
)

func Test_Error(t *testing.T) {
	log.Error("ooo")
}

func Test_Info(t *testing.T) {
	log.Info("ooo")
}

func Test_Notice(t *testing.T) {
	log.Notice("ooo")
}

func Test_Warn(t *testing.T) {
	log.Warn("ooo")
}

func Test_Debug(t *testing.T) {
	log.Debug("ooo")
}

package cache

import "fmt"

const (
	SIGNAL_BASE_PREFIX          = "signaling"
	SIGNAL_IDENTITY_CODE_PREFIX = "idcode"
)

func GetSignalingIdentifyCode(deviceId string) string {
	return fmt.Sprintf("%s:%s:%s", SIGNAL_BASE_PREFIX, SIGNAL_IDENTITY_CODE_PREFIX, deviceId)
}

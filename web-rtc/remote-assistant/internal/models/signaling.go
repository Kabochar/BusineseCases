package models

import "github.com/gorilla/websocket"

const (
	SIGNAL_FLAG_RA_CONNECT     = "ra-connected"
	SIGNAL_FLAG_START_CONTROL  = "start-control"
	SIGNAL_FLAG_APPLY_CONTROL  = "apply-control"
	SIGNAL_FLAG_FAILED_CONTROL = "control-failed"
	SIGNAL_FLAG_AGREE_CONTROL  = "agree-control"
	SIGNAL_FLAG_DENY_CONTROL   = "deny-control"
	SIGNAL_FLAG_CANCEL_CONTROL = "cancel-control"
	SIGNAL_FLAG_FORWARD_MSG    = "forward-msg"
	SIGNAL_FLAG_HEART_BEAT     = "heart-beat"
)

type AssistantWsConn struct {
	WConn        *websocket.Conn // ws连接
	IdentityCode string          // 具体协助码
	State        string          // 当前状态
}

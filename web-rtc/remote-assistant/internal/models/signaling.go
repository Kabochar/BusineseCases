package models

import "github.com/gorilla/websocket"

const (
	SIGNAL_FLAG_WS_CONNECTED   = "ws-connected"  //
	SIGNAL_FLAG_START_CONTROL  = "start-control" // 发起协助
	SIGNAL_FLAG_APPLY_CONTROL  = "apply-control" // 申请协助
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
	ConnectTime  int64           // 连接时间
}

type SignalingMsgInfo struct {
	MsgEvent string `json:"msgEvent"`
	Body     string `json:"body"`
	FromUser string `json:"fromUser"`
	ToUser   string `json:"toUser"`

	// 错误信息
	ErrorCode int    `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
}

// 获取 Turn/Stun 服务地址信息
type GetServerInfoReply struct {
	TurnAddr    string `json:"turnAddr"`
	StunAddr    string `json:"stunAddr"`
	TurnAccount string `json:"turnAccount"`
	TurnAccPwd  string `json:"turnAccPwd"`
}

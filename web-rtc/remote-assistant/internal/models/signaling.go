package models

import "github.com/gorilla/websocket"

const (
	SIGNAL_FLAG_START_CONTROL  = "start-control"  // 发起协助
	SIGNAL_FLAG_APPLY_CONTROL  = "apply-control"  // 申请协助
	SIGNAL_FLAG_CONTROLLING    = "controlling"    // 协助中
	SIGNAL_FLAG_FAILED_CONTROL = "control-failed" // 协助失败，eg 对方不在线
	SIGNAL_FLAG_AGREE_CONTROL  = "agree-control"  // 同意协助
	SIGNAL_FLAG_DENY_CONTROL   = "deny-control"   // 拒绝协助
	SIGNAL_FLAG_CANCEL_CONTROL = "cancel-control" // 取消协助
	SIGNAL_FLAG_FORWARD_MSG    = "forward-msg"    // 透传消息
	SIGNAL_FLAG_HEART_BEAT     = "heart-beat"     // 心跳信息

	// 状态可选
	ASSISTANT_STATE_ONLINE    = "online"    // 在线/空闲
	ASSISTANT_STATE_ASSISTING = "assisting" // 协助中
	ASSISTANT_STATE_OFFLINE   = "offline"   // 离线-占位
)

type AssistantWsConn struct {
	WConn             *websocket.Conn // ws连接
	IdentityCode      string          // 具体连接码
	State             string          // 当前状态
	ConnectTime       int64           // 连接时间
	AssistOtherIdCode string          // 协助时对方的连接码
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

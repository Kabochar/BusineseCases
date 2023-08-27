package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"remote-assistant/internal/cache"
	"remote-assistant/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	concurrentHashMap "github.com/hfdpx/concurrent-hash-map"
	"github.com/spf13/cast"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	} // use default options
	signalHandleFunc   = map[string]func(msgInfo *models.SignalingMsgInfo, wsConn *models.AssistantWsConn){} // 事件处理函数mapping
	assistantWsConnHub = concurrentHashMap.New()                                                             // 协助ws连接池
)

func ServerInfo(c *gin.Context) {
	reply := models.GetServerInfoReply{
		TurnAddr:    fmt.Sprintf("%s:%s", os.Getenv("TURN_ADDR"), os.Getenv("TURN_PORT")),
		StunAddr:    fmt.Sprintf("%s:%s", os.Getenv("TURN_ADDR"), os.Getenv("TURN_PORT")),
		TurnAccount: os.Getenv("TURN_ACCOUNT"),
		TurnAccPwd:  os.Getenv("TURN_ACC_PWD"),
	}
	c.JSON(http.StatusOK, reply)
}

// SignalingServer 信令
func SignalingServer(c *gin.Context) {
	// websocket 协议升级
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	assistantWsConn := &models.AssistantWsConn{
		WConn: conn,
	}
	// 加入连接池
	idCode, _, err := handleConnect(c)
	if err != nil {
		log.Println("handel connect event failed", err)
		return
	}
	assistantWsConn.IdentityCode = idCode
	assistantWsConn.ConnectTime = time.Now().Unix()
	assistantWsConnHub.Set(idCode, assistantWsConn)

	// 监听和处理msg
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read msg from ws conn error", err)
			break
		}
		log.Printf("ws conn received: %s", msg)

		go func(msg []byte, assistantWsConn *models.AssistantWsConn) {
			handleMsgEvent(msg, assistantWsConn)
		}(msg, assistantWsConn)
	}

	// 释放连接
	handleReleaseWsConn(assistantWsConn)
}

func init() {
	signalHandleFunc[models.SIGNAL_FLAG_START_CONTROL] = handleStartControl
	signalHandleFunc[models.SIGNAL_FLAG_AGREE_CONTROL] = HandleActionControl
	signalHandleFunc[models.SIGNAL_FLAG_DENY_CONTROL] = HandleActionControl
	signalHandleFunc[models.SIGNAL_FLAG_CANCEL_CONTROL] = HandleActionControl
	signalHandleFunc[models.SIGNAL_FLAG_FORWARD_MSG] = handleForwardMsg
	signalHandleFunc[models.SIGNAL_FLAG_HEART_BEAT] = handleHeartbeat
}

func handleMsgEvent(msg []byte, assistantWsConn *models.AssistantWsConn) {
	msgInfo := &models.SignalingMsgInfo{}
	err := json.Unmarshal(msg, msgInfo)
	if err != nil {
		log.Println("get origin msg body failed", err)
		return
	}

	// 丢弃到对应的
	fn, isExit := signalHandleFunc[msgInfo.MsgEvent]
	if !isExit {
		log.Println("cannot handle this msg event", assistantWsConn.IdentityCode, msgInfo.MsgEvent)
		return
	}
	fn(msgInfo, assistantWsConn)
}

func handleConnect(c *gin.Context) (idCode, reloginToken string, err error) {
	log.Println("start connect")

	if c.Query("deviceId") == "" {
		return "", "", errors.New("cannot get deviceId")
	}

	cacheKey := cache.GetSignalingIdentifyCode(c.Query("deviceId"))
	val, err := cache.GetValueFromRedis(c, cacheKey, "")
	if err != nil {
		return "", "", err
	}
	if val != "" {
		return val, "", nil
	}

	// 随机生成code
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return "", "", err
	}
	idCode = cast.ToString(newUUID.ID() / 10)
	log.Println("generate new idcode", newUUID)
	cache.GetRdbClient().Set(c, cacheKey, idCode, 0)

	return idCode, "", nil
}

// 发起协助
func handleStartControl(msgInfo *models.SignalingMsgInfo, fromUserConn *models.AssistantWsConn) {
	log.Println("start control")
	// 判断是否在线
	rawToUserConn, isExit := assistantWsConnHub.Get(msgInfo.ToUser)
	if !isExit {
		log.Println("other side not online", msgInfo.ToUser)
		handleReplyMsg(fromUserConn, &models.SignalingMsgInfo{
			MsgEvent:  models.SIGNAL_FLAG_FAILED_CONTROL,
			ErrorCode: -1,
			ErrorMsg:  "other side not online",
		})
		return
	}
	toUserConn := rawToUserConn.(*models.AssistantWsConn)

	// 发送具体的协助请求
	handleReplyMsg(toUserConn, &models.SignalingMsgInfo{
		MsgEvent: models.SIGNAL_FLAG_APPLY_CONTROL,
		FromUser: msgInfo.FromUser,
		ToUser:   msgInfo.ToUser,
	})
	handleReplyMsg(fromUserConn, &models.SignalingMsgInfo{
		MsgEvent: models.SIGNAL_FLAG_START_CONTROL,
		FromUser: msgInfo.FromUser,
		ToUser:   msgInfo.ToUser,
	})

	// 锁定双方的协助状态
}

// HandleActionControl 协助中各种命令转发
func HandleActionControl(msgInfo *models.SignalingMsgInfo, fromUserConn *models.AssistantWsConn) {
	log.Println("start action control")
	// 判断是否在线
	rawToUserConn, isExit := assistantWsConnHub.Get(msgInfo.ToUser)
	if !isExit {
		log.Println("other side not online", msgInfo.ToUser)
		handleReplyMsg(fromUserConn, &models.SignalingMsgInfo{
			MsgEvent:  models.SIGNAL_FLAG_FAILED_CONTROL,
			ErrorCode: -1,
			ErrorMsg:  "other side not online",
		})
		return
	}
	toUserConn := rawToUserConn.(*models.AssistantWsConn)

	// 转发具体消息
	handleReplyMsg(toUserConn, &models.SignalingMsgInfo{
		MsgEvent: msgInfo.MsgEvent,
		FromUser: msgInfo.FromUser,
		ToUser:   msgInfo.ToUser,
	})
	handleReplyMsg(fromUserConn, &models.SignalingMsgInfo{
		MsgEvent: msgInfo.MsgEvent,
		FromUser: msgInfo.FromUser,
		ToUser:   msgInfo.ToUser,
	})

	// 判断是否解除双方的协助状态
}

func handleForwardMsg(msgInfo *models.SignalingMsgInfo, fromUserConn *models.AssistantWsConn) {
	log.Println("start forward msg")
	// 判断是否在线
	rawToUserConn, isExit := assistantWsConnHub.Get(msgInfo.ToUser)
	if !isExit {
		log.Println("other side not online", msgInfo.ToUser)
		handleReplyMsg(fromUserConn, &models.SignalingMsgInfo{
			MsgEvent:  models.SIGNAL_FLAG_FAILED_CONTROL,
			ErrorCode: -1,
			ErrorMsg:  "other side not online",
		})
		return
	}
	toUserConn := rawToUserConn.(*models.AssistantWsConn)

	// 转发具体消息
	handleReplyMsg(toUserConn, &models.SignalingMsgInfo{
		MsgEvent: msgInfo.MsgEvent,
		FromUser: msgInfo.FromUser,
		ToUser:   msgInfo.ToUser,
	})
	handleReplyMsg(fromUserConn, &models.SignalingMsgInfo{
		MsgEvent: msgInfo.MsgEvent,
		FromUser: msgInfo.FromUser,
		ToUser:   msgInfo.ToUser,
	})
}

// 处理心跳相关
func handleHeartbeat(msgInfo *models.SignalingMsgInfo, fromUserConn *models.AssistantWsConn) {
	log.Println("start heart-beat")
	// 心跳包的处理和响应
	handleReplyMsg(fromUserConn, &models.SignalingMsgInfo{
		MsgEvent: models.SIGNAL_FLAG_HEART_BEAT,
		Body:     "pong",
	})
}

// 释放连接
func handleReleaseWsConn(assistantConn *models.AssistantWsConn) {
	// 移除连接
	assistantWsConnHub.Remove(assistantConn.IdentityCode)

	// 解除协助状态
}

// handleReplyMsg 响应具体消息
func handleReplyMsg(wsConn *models.AssistantWsConn, msgBody *models.SignalingMsgInfo) {
	wsConn.WConn.WriteJSON(msgBody)
}

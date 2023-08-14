package handler

import (
	"errors"
	"log"
	"net/http"

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
	signalHandleFunc   = map[string]func(ctx *gin.Context, wsConn *models.AssistantWsConn){} // 事件处理函数mapping
	assistantWsConnHub = concurrentHashMap.New()                                             // 协助ws连接池
)

func ServerInfo(c *gin.Context) {
	c.JSON(http.StatusOK, "server.info")
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
	idCode, _, err := Connect(c, assistantWsConn)
	if err != nil {
		log.Println("handel connect event failed", err)
		return
	}
	assistantWsConn.IdentityCode = idCode
	assistantWsConnHub.Set(idCode, assistantWsConn)

	// 监听和处理msg
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}
		log.Printf("Received: %s", msg)

		err = conn.WriteMessage(msgType, msg)
		if err != nil {
			log.Println("Error during message writing:", err)
			break
		}
	}

	HandleReleaseWsConn(assistantWsConn)
}

func init() {
	signalHandleFunc[models.SIGNAL_FLAG_START_CONTROL] = StartControl
	signalHandleFunc[models.SIGNAL_FLAG_AGREE_CONTROL] = ActionControl
	signalHandleFunc[models.SIGNAL_FLAG_DENY_CONTROL] = ActionControl
	signalHandleFunc[models.SIGNAL_FLAG_CANCEL_CONTROL] = ActionControl
	signalHandleFunc[models.SIGNAL_FLAG_FORWARD_MSG] = ForwardMsg
}

func Connect(c *gin.Context, wsConn *models.AssistantWsConn) (idCode, reloginToken string, err error) {
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

func StartControl(c *gin.Context, wsConn *models.AssistantWsConn) {
	log.Println("start control")
}

func ActionControl(c *gin.Context, wsConn *models.AssistantWsConn) {
	log.Println("start action control")
}

func ForwardMsg(c *gin.Context, wsConn *models.AssistantWsConn) {
	log.Println("start forward msg")
}

// 处理心跳相关
func HandleHeartbeat() {
	log.Println("start heart-beat")
}

// 释放连接
func HandleReleaseWsConn(assistantConn *models.AssistantWsConn) {
	assistantWsConnHub.Remove(assistantConn.IdentityCode)
}

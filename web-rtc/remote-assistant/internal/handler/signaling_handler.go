package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"remote-assistant/internal/model"
)

var (
	upgrader         = websocket.Upgrader{} // use default options
	signalHandleFunc = map[string]func(){}
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

	// 加入连接池

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

	// 释放连接
}

func init() {
	signalHandleFunc[model.SIGNAL_FLAG_RA_CONNECT] = Connect
	signalHandleFunc[model.SIGNAL_FLAG_START_CONTROL] = StartControl
	signalHandleFunc[model.SIGNAL_FLAG_AGREE_CONTROL] = ActionControl
	signalHandleFunc[model.SIGNAL_FLAG_DENY_CONTROL] = ActionControl
	signalHandleFunc[model.SIGNAL_FLAG_CANCEL_CONTROL] = ActionControl
	signalHandleFunc[model.SIGNAL_FLAG_FORWARD_MSG] = ForwardMsg
}

func Connect() {
	log.Println("start handle")
}

func StartControl() {
	log.Println("start handle")
}

func ActionControl() {
	log.Println("start handle")
}

func ForwardMsg() {
	log.Println("start handle")
}

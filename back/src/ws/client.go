package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Conn *websocket.Conn
}

func (cl *Client) HandlerWebsocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("%s, error while Upgrading websocket connection\n", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	cl.Conn = conn

	// for {
	// 	messageType, p, err := conn.ReadMessage()
	// 	if err != nil {
	// 		log.Printf("%s, error while reading message\n", err.Error())
	// 		c.AbortWithError(http.StatusInternalServerError, err)
	// 		break
	// 	}

	// 	err = conn.WriteMessage(messageType, append([]byte("server echo: "), p...))
	// 	if err != nil {
	// 		log.Printf("%s, error while writing message\n", err.Error())
	// 		c.AbortWithError(http.StatusInternalServerError, err)
	// 		break
	// 	}
	// }
}

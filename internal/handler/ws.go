package handler

import (
	"net/http"
	"quiz-svc/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type WSHandler struct {
	netService *service.NetService
	logger     *zap.Logger
}

func NewWSHandler(netService *service.NetService, logger *zap.Logger) *WSHandler {
	return &WSHandler{netService: netService, logger: logger}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *WSHandler) Ws(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("Failed to upgrade connection", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	h.logger.Info("New WebSocket connection established")

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			h.logger.Info("WebSocket connection closed", zap.Error(err))
			h.netService.OnDisconnect(conn)
			break
		}

		h.netService.OnIncomingMessage(conn, mt, msg)
	}
}

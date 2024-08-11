package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	service_shared "github.com/sepcon/quizprob/cmd/shared"
	"github.com/sepcon/quizprob/internal/events"
	"io"
	"net/http"
	"net/url"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	bs := events.NewChannelManager()
	r := gin.Default()

	r.GET("/connect", func(c *gin.Context) {
		subscriberID := c.Query("subscriber_id")
		if subscriberID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing subscriber_id"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		}
		bs.Connect(conn, subscriberID)
		go func() {
			closeHandler := func(code int, text string) error {
				fmt.Println("WebSocket connection closed:", code, text)
				// Unsubscribe the user from all channels
				bs.CleanupObsoleteConnection(conn, subscriberID)
				return nil
			}

			// Handle disconnection
			conn.SetCloseHandler(closeHandler)

			for {
				_, bytes, err := conn.ReadMessage()
				if err != nil {
					bs.CleanupObsoleteConnection(conn, subscriberID)
					break
				}
				var subscribeMessage map[string]string
				err = json.Unmarshal(bytes, &subscribeMessage)
				if err != nil {
					conn.WriteJSON(gin.H{"error": "invalid json format"})
					continue
				}
				operationType := subscribeMessage["type"]
				if operationType == "" {
					conn.WriteJSON(gin.H{"error": "missing operation type"})
					continue
				}
				channel := subscribeMessage["channel"]

				switch operationType {
				case "subscribe":
					if channel != "" {
						bs.Subscribe(subscriberID, channel)
					}
				case "unsubscribe":
					if channel != "" {
						bs.Unsubscribe(subscriberID, channel)
					}
				default:
					conn.WriteJSON(gin.H{"error": "invalid operation type"})
				}
			}
		}()
	})

	r.POST("/publish/:channel", func(c *gin.Context) {
		channel, _ := url.QueryUnescape(c.Param("channel"))
		channel, _ = url.PathUnescape(channel)
		event, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = bs.Publish(channel, event)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "published"})
	})

	r.Run(":" + service_shared.CONST_EVENT_SERVICE_PORT)
}

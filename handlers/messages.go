package handlers

import (
	"net/http"

	"hakathon/models"
	"hakathon/repository"

	"github.com/gin-gonic/gin"
)

func GetMe(c *gin.Context) {
	profile, _ := c.Get("profile")
	role, _ := c.Get("Role")
	c.JSON(http.StatusOK, gin.H{
		"Role":    role,
		"Profile": profile,
	})
}

func SendMessage(c *gin.Context) {
	var req models.RequestSendMessage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	authID := getAuthID(c)
	if err := repository.CreateMessage(req, authID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Message sent"})
}

func GetMessages(c *gin.Context) {
	recipientID, err := parseParamID(c, "recipient_id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid recipient ID"})
		return
	}
	senderID := getAuthID(c)
	messages, err := repository.GetMessages(senderID, recipientID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}
	c.JSON(http.StatusOK, messages)
}

func EditMessage(c *gin.Context) {
	msgID, err := parseParamID(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}
	var req models.RequestUpdateMessage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	req.MessageID = msgID
	senderID := getAuthID(c)
	if err := repository.UpdateMessage(&req, senderID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update message"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Message updated"})
}

func DeleteMessageMsg(c *gin.Context) {
	id, err := parseParamID(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := repository.DeleteMessage(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Message deleted"})
}

package handlers

import (
	"hakathon/models"
	"hakathon/redis"
	"hakathon/repository"
	"hakathon/security"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RefreshTokens(c *gin.Context) {
	var Token models.RefreshToken
	c.ShouldBindJSON(&Token)
	RfToken, err := security.ValidateRefreshToken(Token.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is not valid!"})
		return
	}
	err = redis.CheckSession(RfToken.RtID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token does not exist!"})
		return
	}
	auth, err := repository.ReadAuth(RfToken.AuthID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not Found!"})
		return
	}
	DataTransfer, err := security.GenerateToken(RfToken.AuthID, auth.Role)
	redis.DeleteSession(RfToken.RtID)
	redis.SaveSession(
		DataTransfer.RefreshTokenID, DataTransfer.RtExpires)
	c.JSON(http.StatusOK, gin.H{
		"refresh": DataTransfer.RefreshToken,
		"access":  DataTransfer.AccessToken,
	})
}

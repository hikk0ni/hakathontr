package handlers

import (
	"net/http"

	"hakathon/models"
	"hakathon/redis"
	"hakathon/repository"
	"hakathon/security"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var req models.RequestAuth
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "The request does not meet the standart model!"})
		return
	}
	authid, err := repository.FindAuthByEmail(req.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "AuthID not find!"})
		return
	}
	Auth, err := repository.ReadAuth(authid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Auth not find!"})
		return
	}
	err = security.CheckHashPassword(req.Password, Auth.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error password Check!"})
		return
	}

	td, err := security.GenerateToken(Auth.ID, Auth.Role)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error Generate Tokens!"})
		return
	}
	redis.SaveSession(td.RefreshTokenID, td.RtExpires)
	c.JSON(http.StatusOK, gin.H{
		"refresh": td.RefreshToken,
		"access":  td.AccessToken,
	})
}

func RegisterPatient(c *gin.Context) {
	var req models.RequestCreatePatient
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "The request does not meet the standart model!"})
		return
	}
	err = repository.CreatePatient(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error to CreatePatient"})
		return
	}
	authid, err := repository.FindAuthByEmail(req.Email)
	td, err := security.GenerateToken(authid, "patient")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error Generate Tokens!"})
		return
	}
	redis.SaveSession(td.RefreshTokenID, td.RtExpires)
	c.JSON(http.StatusOK, gin.H{
		"refresh": td.RefreshToken,
		"access":  td.AccessToken,
	})
}

func RegisterClinic(c *gin.Context) {
	var req models.RequestCreateClinic
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "The request does not meet the standart model!"})
		return
	}
	err = repository.CreateClinic(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error to Create Clinic"})
		return
	}
	authid, err := repository.FindAuthByEmail(req.Admin.Email)
	td, err := security.GenerateToken(authid, "Doctor")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error Generate Tokens!"})
		return
	}
	redis.SaveSession(td.RefreshTokenID, td.RtExpires)
	c.JSON(http.StatusOK, gin.H{
		"refresh": td.RefreshToken,
		"access":  td.AccessToken,
	})
}

func RegisterDoctor(c *gin.Context) {
	var req models.RequestCreateDoctor
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "The request does not meet the standart model!"})
		return
	}
	err = repository.CreateDoctor(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error to CreatePatient"})
		return
	}
	authid, err := repository.FindAuthByEmail(req.Email)
	td, err := security.GenerateToken(authid, "Doctor")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error Generate Tokens!"})
		return

	}
	redis.SaveSession(td.RefreshTokenID, td.RtExpires)
	c.JSON(http.StatusOK, gin.H{
		"refresh": td.RefreshToken,
		"access":  td.AccessToken,
	})
}

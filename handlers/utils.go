package handlers

import (
	"strconv"

	"hakathon/models"

	"github.com/gin-gonic/gin"
)

// Получает AuthID из контекста, который положил AuthRequired
func getAuthID(c *gin.Context) uint {
	role, _ := c.Get("Role")
	if role == "Doctor" {
		if doc, ok := c.MustGet("profile").(models.Doctor); ok {
			return doc.AuthID
		}
	} else {
		if pat, ok := c.MustGet("profile").(models.Patient); ok {
			return pat.AuthID
		}
	}
	return 0
}

// Получает внутренний ID профиля (Врача или Пациента)
func getProfileID(c *gin.Context) uint {
	role, _ := c.Get("Role")
	if role == "Doctor" {
		if doc, ok := c.MustGet("profile").(models.Doctor); ok {
			return doc.ID
		}
	} else {
		if pat, ok := c.MustGet("profile").(models.Patient); ok {
			return pat.ID
		}
	}
	return 0
}

func parseParamID(c *gin.Context, param string) (uint, error) {
	id, err := strconv.Atoi(c.Param(param))
	return uint(id), err
}

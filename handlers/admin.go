package handlers

import (
	"net/http"

	"hakathon/models"
	"hakathon/repository"

	"github.com/gin-gonic/gin"
)

func UpdateClinicInfo(c *gin.Context) {
	var req models.RequestUpdateClinic
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	docProfile := c.MustGet("profile").(models.Doctor)
	if err := repository.UpdateClinic(&req, docProfile.ClinicID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update clinic"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Clinic updated"})
}

func ApproveDoctor(c *gin.Context) {
	docID, err := parseParamID(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	authIDJWT := getAuthID(c)
	if err := repository.MakeDoctorIsAvailable(docID, authIDJWT); err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Doctor approved"})
}

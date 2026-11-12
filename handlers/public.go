package handlers

import (
	"hakathon/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllClinicsList(c *gin.Context) {
	ApproveFilter := c.Query("ApproveFilter")
	clinics, err := repository.GetAllClinics(ApproveFilter)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get clinics"})
		return
	}
	c.JSON(http.StatusOK, clinics)
}

func GetClinicInfo(c *gin.Context) {
	id, err := parseParamID(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	clinic, err := repository.ReadClinic(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Clinic not found"})
		return
	}
	c.JSON(http.StatusOK, clinic)
}

func GetClinicDocs(c *gin.Context) {
	id, err := parseParamID(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ApproveFilter := c.Query("ApproveFilter")

	doctors, err := repository.GetDoctorsByClinic(id, ApproveFilter)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get doctors"})
		return
	}
	c.JSON(http.StatusOK, doctors)
}

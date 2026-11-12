package handlers

import (
	"net/http"

	"hakathon/models"
	"hakathon/repository"

	"github.com/gin-gonic/gin"
)

func GetPatientProfile(c *gin.Context) {
	id, err := parseParamID(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	patient, err := repository.ReadPatient(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}
	c.JSON(http.StatusOK, patient)
}

func UpdatePatientProfile(c *gin.Context) {
	var req models.RequestUpdatePatient
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	patientID := getProfileID(c)
	if err := repository.UpdatePatient(&req, patientID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
}

func DeleteMyPatientProfile(c *gin.Context) {
	patientID := getProfileID(c)
	if err := repository.DeletePatient(patientID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profile deleted"})
}

func GetMyMedicalRecords(c *gin.Context) {
	patientID := getProfileID(c)
	records, err := repository.GetAllRecordsByPatient(patientID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get records"})
		return
	}
	c.JSON(http.StatusOK, records)
}

func CreateAppointmentReq(c *gin.Context) {
	var req models.RequestAppointment
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	req.PatientID = getProfileID(c)
	if err := repository.CreateAppointment(req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create appointment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Appointment created"})
}

func GetPatientQueries(c *gin.Context) {
	patientID := getProfileID(c)
	queries, err := repository.GetQueryForPatient(patientID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get appointments"})
		return
	}
	c.JSON(http.StatusOK, queries)
}

func CancelAppoint(c *gin.Context) {
	id, err := parseParamID(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := repository.DeleteAppointment(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel appointment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Appointment canceled"})
}

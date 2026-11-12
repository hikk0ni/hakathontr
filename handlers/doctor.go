package handlers

import (
	"net/http"
	"strconv"

	"hakathon/models"
	"hakathon/repository"

	"github.com/gin-gonic/gin"
)

func GetDoctorProfile(c *gin.Context) {
	id, err := parseParamID(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	doctor, err := repository.ReadDoctor(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}
	c.JSON(http.StatusOK, doctor)
}

func UpdateDoctorProfile(c *gin.Context) {
	var req models.RequestUpdateDoctor
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	doctorID := getProfileID(c)
	if err := repository.UpdateDoctor(&req, doctorID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
}

func GetClinicPatientsList(c *gin.Context) {
	docProfile := c.MustGet("profile").(models.Doctor)
	patients, err := repository.GetPatientsByClinic(docProfile.ClinicID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get patients"})
		return
	}
	c.JSON(http.StatusOK, patients)
}

func GetDoctorQueries(c *gin.Context) {
	doctorID := getProfileID(c)
	queries, err := repository.GetQueryForDoctor(doctorID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get appointments"})
		return
	}
	c.JSON(http.StatusOK, queries)
}

func CreateRecord(c *gin.Context) {
	var req models.RequestCreateRecord
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	doctorID := getProfileID(c)
	if err := repository.CreateMedicalRecord(&req, doctorID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create record"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Medical record created"})
}

func UpdateRecord(c *gin.Context) {
	var req models.RequestUpdateRecord
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	req.DoctorID = getProfileID(c)
	if err := repository.UpdateMedicalRecord(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Record updated"})
}

func DeleteRecord(c *gin.Context) {
	recordID, err := parseParamID(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Record ID"})
		return
	}

	patientIDStr := c.Query("patient_id")
	patientID, _ := strconv.Atoi(patientIDStr)

	if err := repository.DeleteMedicalRecord(recordID, uint(patientID)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete record"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Record deleted"})
}

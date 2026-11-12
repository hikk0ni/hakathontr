package handlers

import (
	"net/http"

	"hakathon/repository"

	"github.com/gin-gonic/gin"
)

func ApproveClinic(c *gin.Context) {
	clinicID, err := parseParamID(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	authIDJWT := getAuthID(c)
	if err := repository.MakeClinicIsAvailable(clinicID, authIDJWT); err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Clinic approved"})
}

func PromoteStaff(c *gin.Context) {
	docID, err := parseParamID(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	authIDJWT := getAuthID(c)
	if err := repository.PromoteDoctor(docID, authIDJWT); err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Staff promoted"})
}

func DemoteStaff(c *gin.Context) {
	docID, err := parseParamID(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	authIDJWT := getAuthID(c)
	if err := repository.DemoteDoctor(docID, authIDJWT); err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Staff demoted"})
}

func DeleteClinicGlobal(c *gin.Context) {
	id, err := parseParamID(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := repository.DeleteClinic(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete clinic"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Clinic deleted"})
}

func DeleteDoctorGlobal(c *gin.Context) {
	id, err := parseParamID(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := repository.DeleteDoctor(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete doctor"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Doctor deleted"})
}

func DeletePatientGlobal(c *gin.Context) {
	id, err := parseParamID(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := repository.DeletePatient(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete patient"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Patient deleted"})
}

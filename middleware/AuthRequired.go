package middleware

import (
	"hakathon/repository"
	"hakathon/security"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token does not exist!"})
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		AcToken, err := security.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}
		c.Set("Role", AcToken.Role)
		if AcToken.Role == "Doctor" {
			DoctorID, err := repository.FindDoctorByAuthID(AcToken.AuthID)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User does not exist!"})
				return
			}
			Doctor, err := repository.ReadDoctor(DoctorID)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User does not exist!"})
				return
			}
			if !Doctor.AccountAvailable {
				log.Println("Doctor is not Availible!")
				return
			}
			c.Set("profile", Doctor)
			return
		} else {
			PatientID, err := repository.FindPatientByAuthID(AcToken.AuthID)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User does not exist!"})
				return
			}
			Patient, err := repository.ReadPatient(PatientID)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User does not exist!"})
				return
			}
			c.Set("profile", Patient)
			return
		}
	}
}

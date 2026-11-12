package main

import (
	"log"

	"hakathon/handlers"
	"hakathon/middleware"
	"hakathon/redis"
	"hakathon/repository"

	"github.com/gin-gonic/gin"
)

func main() {
	redis.InitRedis()
	repository.Initdb()

	r := gin.Default()
	r.Use(middleware.CORS())

	public := r.Group("/api")
	{
		public.POST("/login", handlers.Login)
		public.POST("/register/patient", handlers.RegisterPatient) // CreatePatient
		public.POST("/register/doctor", handlers.RegisterDoctor)   // CreateDoctor
		public.POST("/register/clinic", handlers.RegisterClinic)   // CreateClinic

		public.GET("/clinics", handlers.GetAllClinicsList)         // GetAllClinics
		public.GET("/clinics/:id", handlers.GetClinicInfo)         // ReadClinic
		public.GET("/clinics/:id/doctors", handlers.GetClinicDocs) // GetDoctorsByClinic
	}

	private := r.Group("/api")
	private.Use(middleware.AuthRequired())
	{
		private.GET("/me", handlers.GetMe) // checkroot: FindDoctorByAuthID / FindPatientByAuthID

		private.POST("/messages", handlers.SendMessage)              // CreateMessage
		private.GET("/messages/:recipient_id", handlers.GetMessages) // GetMessages
		private.PATCH("/messages/:id", handlers.EditMessage)         // UpdateMessage
		private.DELETE("/messages/:id", handlers.DeleteMessageMsg)   // DeleteMessage

		patient := private.Group("/patient")
		{
			patient.GET("/profile/:id", handlers.GetPatientProfile)     // ReadPatient
			patient.PATCH("/profile", handlers.UpdatePatientProfile)    // UpdatePatient
			patient.DELETE("/profile", handlers.DeleteMyPatientProfile) // DeletePatient (самоудаление)

			patient.GET("/records", handlers.GetMyMedicalRecords) // GetAllRecordsByPatient

			patient.POST("/appointments", handlers.CreateAppointmentReq) // CreateAppointment
			patient.GET("/appointments", handlers.GetPatientQueries)     // GetQueryForPatient
			patient.DELETE("/appointments/:id", handlers.CancelAppoint)  // DeleteAppointment
		}

		doctor := private.Group("/doctor")
		{
			doctor.GET("/profile/:id", handlers.GetDoctorProfile)  // ReadDoctor
			doctor.PATCH("/profile", handlers.UpdateDoctorProfile) // UpdateDoctor

			// Работа с пациентами
			doctor.GET("/patients", handlers.GetClinicPatientsList) // GetPatientsByClinic

			// Очередь записей к этому врачу (Quere)
			doctor.GET("/appointments", handlers.GetDoctorQueries) // GetQueryForDoctor

			// Управление медицинскими картами (MedicalRecord)
			doctor.POST("/records", handlers.CreateRecord)       // CreateMedicalRecord
			doctor.PATCH("/records/:id", handlers.UpdateRecord)  // UpdateMedicalRecord
			doctor.DELETE("/records/:id", handlers.DeleteRecord) // DeleteMedicalRecord
		}

		admin := private.Group("/admin")
		{
			admin.PATCH("/clinic", handlers.UpdateClinicInfo)          // UpdateClinic
			admin.POST("/doctors/:id/approve", handlers.ApproveDoctor) // MakeDoctorIsAvailable
		}

		superadmin := private.Group("/superadmin")
		{
			superadmin.POST("/clinics/:id/approve", handlers.ApproveClinic) // MakeClinicIsAvailable

			superadmin.POST("/doctors/:id/promote", handlers.PromoteStaff) // PromoteDoctor
			superadmin.POST("/doctors/:id/demote", handlers.DemoteStaff)   // DemoteDoctor

			superadmin.DELETE("/clinics/:id", handlers.DeleteClinicGlobal)   // DeleteClinic
			superadmin.DELETE("/doctors/:id", handlers.DeleteDoctorGlobal)   // DeleteDoctor
			superadmin.DELETE("/patients/:id", handlers.DeletePatientGlobal) // DeletePatient
		}
	}

	log.Println("Сервер API запущен на порту 8080...")
	r.Run(":8080")
}

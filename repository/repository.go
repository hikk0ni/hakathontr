package repository

import (
	"errors"
	"log"
	"time"

	"hakathon/models"
	"hakathon/security"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// database
var db *gorm.DB

func Initdb() {
	var err error

	db, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		log.Println("Ошибка подключения к базе данных!: ", err)
	}
	err = db.AutoMigrate(
		&models.Auth{},
		&models.Clinic{},
		&models.Doctor{},
		&models.Patient{},
		&models.MedicalRecord{},
	)
	if err != nil {
		log.Println("Ошибка при миграции таблиц!", err)
	}
	log.Println("База успешно инициализирована!")
}

// CURD!

// Create

func CreateDoctor(rcd *models.RequestCreateDoctor) error {
	return db.Transaction(func(tx *gorm.DB) error {
		PasswordHash, err := security.HashPassword(rcd.Password)
		if err != nil {
			log.Println("Panic!")
			return err
		}
		auth := models.Auth{
			Email:    rcd.Email,
			Password: PasswordHash,
			Role:     "Doctor",
		}
		if err := tx.Create(&auth).Error; err != nil {
			return err
		}
		doctor := models.Doctor{
			ClinicID:   rcd.ClinicID,
			Name:       rcd.Name,
			Surname:    rcd.Surname,
			Speciality: rcd.Speciality,
			AuthID:     auth.ID,
		}
		if err := tx.Create(&doctor).Error; err != nil {
			return err
		}
		return nil
	})
}

func CreatePatient(req *models.RequestCreatePatient) error {
	return db.Transaction(func(tx *gorm.DB) error {
		PasswordHash, err := security.HashPassword(req.Password)
		if err != nil {
			log.Println("Panic!")
			return err
		}
		auth := models.Auth{
			Email:    req.Email,
			Password: PasswordHash,
			Role:     "Patient",
		}
		if err := tx.Create(&auth).Error; err != nil {
			log.Println("Panic! Could not create auth!")
			return err
		}
		patient := models.Patient{
			ClinicID: req.ClinicID,
			Name:     req.Name,
			Surname:  req.Surname,
			Disease:  req.Disease,
			AuthID:   auth.ID,
		}
		if err := tx.Create(&patient).Error; err != nil {
			log.Println("Panic! Could not create patient!")
			return err
		}
		return nil
	})
}

func CreateClinic(req *models.RequestCreateClinic) error {
	return db.Transaction(func(tx *gorm.DB) error {
		clinic := models.Clinic{
			Name:    req.Name,
			Address: req.Address,
		}
		if err := tx.Create(&clinic).Error; err != nil {
			log.Println("Panic! Could not create clinic!")
			return err
		}
		PasswordHash, err := security.HashPassword(req.Admin.Password)
		if err != nil {
			log.Println("Panic!")
			return err
		}
		auth := models.Auth{
			Email:    req.Admin.Email,
			Password: PasswordHash,
			Role:     "Doctor",
		}
		if err := tx.Create(&auth).Error; err != nil {
			log.Println("Panic! Could not create auth for admin!")
			return err
		}
		admin := models.Doctor{
			ClinicID:   clinic.ID,
			Name:       req.Admin.Name,
			Surname:    req.Admin.Surname,
			Speciality: req.Admin.Speciality,
			Root:       1,
			AuthID:     auth.ID,
		}
		if err := tx.Create(&admin).Error; err != nil {
			log.Println("Panic! Could not create admin!")
			return err
		}
		return nil
	})
}

func CreateMedicalRecord(req *models.RequestCreateRecord, doctorID uint) error {
	rec := models.MedicalRecord{
		PatientID:    req.PatientID,
		DoctorID:     doctorID,
		Diagnosis:    req.Diagnosis,
		Prescription: req.Prescription,
		Date:         time.Now().Format("02.01.2006"),
	}
	return db.Create(&rec).Error
}

func CreateMessage(req models.RequestSendMessage, IDJWT uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		AuthID, err := FindAuthByEmail(req.RecipientEmail)
		if err != nil {
			log.Println("Error to get Account ID from Recipient!")
			return err
		}
		m := models.Message{
			SenderID:         IDJWT,
			RecipientID:      AuthID,
			MessageContainer: req.MessageContainer,
		}
		return db.Create(&m).Error
	})
}

func CreateAppointment(req models.RequestAppointment) error {
	return db.Transaction(func(tx *gorm.DB) error {
		pq := models.Queue{
			PatientID: req.PatientID,
			DoctorID:  req.DoctorID,
		}
		err := tx.Create(&pq).Error
		if err != nil {
			log.Println("Error to create appointment!")
			return err
		}
		return nil
	})
}

// Update

func UpdateDoctor(req *models.RequestUpdateDoctor, ID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		d := models.Doctor{}
		err := tx.Model(models.Doctor{}).Select("AuthID").First(&d, ID).Error
		if err != nil {
			log.Println("Panic! UpdateDoctor Failed!")
			return err
		}
		PasswordHash, err := security.HashPassword(req.Auth.Password)
		if err != nil {
			log.Println("Panic!")
			return err
		}
		err = tx.Model(models.Auth{}).Where("ID = ?", d.AuthID).Updates(models.Auth{
			Email:    req.Auth.Email,
			Password: PasswordHash,
		}).Error
		if err != nil {
			log.Println("panic! updatedoctor failed")
			return err
		}
		err = tx.Model(models.Doctor{}).Where("ID = ?", ID).Updates(models.Doctor{
			ClinicID:   req.ClinicID,
			Name:       req.Name,
			Surname:    req.Surname,
			Speciality: req.Speciality,
		}).Error
		if err != nil {
			log.Println("panic! updatedoctor failed")
			return err
		}
		return nil
	})
}

func UpdatePatient(req *models.RequestUpdatePatient, ID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		p := models.Patient{}
		err := tx.Model(models.Patient{}).Select("AuthID").First(&p, ID).Error
		if err != nil {
			log.Println("Panic! UpdatePatient Failed!")
			return err
		}
		PasswordHash, err := security.HashPassword(req.Auth.Password)
		if err != nil {
			log.Println("Panic!")
			return err
		}
		err = tx.Model(models.Auth{}).Where("ID = ?", p.AuthID).Updates(models.Auth{
			Email:    req.Auth.Email,
			Password: PasswordHash,
		}).Error
		if err != nil {
			log.Println("Panic! UpdatePatient Failed!")
			return err
		}
		err = tx.Model(models.Patient{}).Where("ID = ?", ID).Updates(models.Patient{
			ClinicID: req.ClinicID,
			Name:     req.Name,
			Surname:  req.Surname,
			Disease:  req.Disease,
		}).Error
		if err != nil {
			log.Println("Panic! UpdatePatient Failed!")
			return err
		}
		return nil
	})
}

func UpdateClinic(req *models.RequestUpdateClinic, ID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(models.Clinic{}).Where("ID = ?", ID).Updates(models.Clinic{
			Name:    req.Name,
			Address: req.Address,
		}).Error
		if err != nil {
			log.Println("Panic! UpdateClinic Failed!")
			return err
		}
		return nil
	})
}

func UpdateMedicalRecord(req *models.RequestUpdateRecord) error {
	err := db.Model(&models.MedicalRecord{}).Where("id = ?", req.RecordID).Updates(models.MedicalRecord{
		DoctorID:     req.DoctorID,
		Prescription: req.Prescription,
		Date:         req.Date,
	}).Error
	if err != nil {
		log.Println("Ошибка при обновлении записи:", err)
		return err
	}
	return nil
}

func UpdateMessage(req *models.RequestUpdateMessage, SenderID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		err := db.Model(models.Message{}).Where("id = ? AND sender_id = ?", req.MessageID, SenderID).Updates(models.Message{
			MessageContainer: req.MessageContainer,
		}).Error
		if err != nil {
			log.Println("Error to update message!")
			return err
		}
		return nil
	})
}

// Read

func ReadDoctor(ID uint) (models.Doctor, error) {
	d := models.Doctor{}
	err := db.Model(models.Doctor{}).Where("ID = ?", ID).First(&d).Error
	if err != nil {
		log.Println("Panic! Impossible Operation!")
		return models.Doctor{}, err
	}
	return d, nil
}

func GetDoctorsByClinic(clinicID uint) ([]models.Doctor, error) {
	var doctors []models.Doctor
	err := db.Where("clinic_id = ?", clinicID).Find(&doctors).Error
	if err != nil {
		log.Println("Error Get Doctors!")
	}
	return doctors, err
}

func ReadPatient(ID uint) (models.Patient, error) {
	p := models.Patient{}
	err := db.Model(models.Patient{}).Where("ID = ?", ID).First(&p).Error
	if err != nil {
		log.Println("Panic! Impossible Operation!")
		return models.Patient{}, err
	}
	return p, nil
}

func GetPatientsByClinic(clinicID uint) ([]models.Patient, error) {
	var patients []models.Patient
	err := db.Where("clinic_id = ?", clinicID).Find(&patients).Error
	if err != nil {
		log.Println("Error Get Patients!")
	}
	return patients, err
}

func ReadClinic(ID uint) (models.Clinic, error) {
	c := models.Clinic{}
	err := db.Model(models.Clinic{}).Where("ID = ?", ID).First(&c).Error
	if err != nil {
		log.Println("Panic! Impossible Operation!")
		return models.Clinic{}, err
	}
	return c, nil
}

func GetAllClinics() ([]models.Clinic, error) {
	var clinics []models.Clinic
	err := db.Find(&clinics).Error
	return clinics, err
}

func GetAllRecordsByPatient(patientID uint) ([]models.MedicalRecord, error) {
	var records []models.MedicalRecord
	err := db.Where("patient_id = ?", patientID).Find(&records).Error
	if err != nil {
		log.Println("Error Get Records!")
	}
	return records, err
}

func GetMessages(SenderID uint, RecipientID uint) ([]models.Message, error) {
	var Messages []models.Message
	err := db.Where("(sender_id = ? AND recipient_id = ?) OR (recipient_id = ? AND sender_id = ?)", SenderID, RecipientID, SenderID, RecipientID).Find(&Messages).Error
	if err != nil {
		log.Println("Error find messages!")
		return Messages, err
	}
	return Messages, err
}

func GetQueryForPatient(ID uint) ([]models.Queue, error) {
	quere := []models.Queue{}
	err := db.Model(models.Queue{}).Where("patient_id", ID).Find(&quere).Error
	if err != nil {
		log.Println("Error to get Query")
		return quere, err
	}
	return quere, err
}

func GetQueryForDoctor(ID uint) ([]models.Queue, error) {
	quere := []models.Queue{}
	err := db.Model(models.Queue{}).Where("doctor_id", ID).Find(&quere).Error
	if err != nil {
		log.Println("Error to get Query")
		return quere, err
	}
	return quere, err
}

// Delete

func DeleteDoctor(ID uint) error {
	err := db.Delete(&models.Doctor{}, ID).Error
	if err != nil {
		log.Println("Panic! Impossible Operation!")
		return err
	}
	return nil
}

func DeletePatient(ID uint) error {
	err := db.Delete(&models.Patient{}, ID).Error
	if err != nil {
		log.Println("Panic! Impossible Operation!")
		return err
	}
	return nil
}

func DeleteClinic(ID uint) error {
	err := db.Delete(&models.Clinic{}, ID).Error
	if err != nil {
		log.Println("Panic! Impossible Operation!")
		return err
	}
	return nil
}

func DeleteMedicalRecord(recordID uint, patientID uint) error {
	err := db.Where("id = ? AND patient_id = ?", recordID, patientID).Delete(&models.MedicalRecord{}).Error
	if err != nil {
		log.Println("Error to delete record!:", err)
		return err
	}
	return nil
}

func DeleteMessage(ID uint) error {
	err := db.Delete(models.Message{}, ID).Error
	if err != nil {
		log.Println("Error to delete Message", err)
		return err
	}
	return nil
}

func DeleteAppointment(ID uint) error {
	err := db.Delete(models.Queue{}, ID).Error
	if err != nil {
		log.Println("Panic! Impossible Operation!")
		return err
	}
	return nil
}

// Special func

func FindAuthByEmail(Email string) (uint, error) {
	var AuthID uint
	err := db.Model(models.Auth{}).Select("ID").Where("Email = ?", Email).First(&AuthID).Error
	if err != nil {
		log.Println("Error! Not Find Account!")
		return AuthID, err
	}
	return AuthID, nil
}

func FindDoctorByAuthID(AuthID uint) (uint, error) {
	var DoctorID uint
	err := db.Model(models.Doctor{}).Select("ID").Where("auth_id = ?", AuthID).First(&DoctorID).Error
	if err != nil {
		log.Println("Error! Not Find Doctor")
		return DoctorID, err
	}
	return DoctorID, err
}

func FindPatientByAuthID(AuthID uint) (uint, error) {
	var PatientID uint
	err := db.Model(models.Patient{}).Select("ID").Where("auth_id = ?", AuthID).First(&PatientID).Error
	if err != nil {
		log.Println("Error! Not Find Patient")
		return PatientID, err
	}
	return PatientID, err
}

func MakeClinicIsAvailable(CLinicID uint, AuthIDJWT uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var root uint

		err := tx.Model(models.Doctor{}).Where("auth_id = ?", AuthIDJWT).Select("Root").First(&root).Error
		if err != nil {
			log.Println("Error get root-status!")
			return err
		}
		if root == 2 {
			err := tx.Model(models.Clinic{}).Where("id = ?", CLinicID).Update("ClinicIsAvailable", true).Error
			if err != nil {
				log.Println("Error to make clinic is Available!")
				return err
			}
			return nil
		}
		log.Println("User is not a SuperAdmin!")
		return errors.New("недостаточно прав: требуется уровень суперадмина")
	})
}

func MakeDoctorIsAvailable(DoctorID uint, AuthIDJWT uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var root uint

		err := tx.Model(models.Doctor{}).Where("auth_id = ?", AuthIDJWT).Select("Root").First(&root).Error
		if err != nil {
			log.Println("Error get root-status!")
			return err
		}
		if root >= 1 {
			err := tx.Model(models.Doctor{}).Where("id = ?", DoctorID).Update("AccountAvailable", true).Error
			if err != nil {
				log.Println("Error to make doctor is Available!")
				return err
			}
			return nil
		}
		log.Println("User is not a SuperAdmin!")
		return errors.New("недостаточно прав: требуется уровень суперадмина")
	})
}

func PromoteDoctor(DoctorID uint, AuthIDJWT uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var requester models.Doctor

		if err := tx.Where("auth_id = ?", AuthIDJWT).First(&requester).Error; err != nil {
			return errors.New("ошибка: инициатор не найден")
		}

		if requester.Root < 1 {
			return errors.New("недостаточно прав")
		}

		var target models.Doctor
		if err := tx.Where("id = ?", DoctorID).First(&target).Error; err != nil {
			return errors.New("ошибка: целевой врач не найден")
		}
		if requester.Root == 1 {
			if requester.ClinicID != target.ClinicID {
				return errors.New("отказано: вы можете управлять только врачами своей клиники")
			}
			if target.Root >= 1 {
				return errors.New("отказано: вы не можете выдать права выше локального администратора")
			}
		}

		if requester.Root == 2 {
			if target.Root >= 2 {
				return errors.New("у врача уже максимальный уровень прав (2)")
			}
		}
		return tx.Model(&models.Doctor{}).Where("id = ?", DoctorID).Update("Root", target.Root+1).Error
	})
}

func DemoteDoctor(DoctorID uint, AuthIDJWT uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var requester models.Doctor
		if err := tx.Where("auth_id = ?", AuthIDJWT).First(&requester).Error; err != nil {
			return errors.New("ошибка: инициатор не найден")
		}
		if requester.Root < 1 {
			return errors.New("недостаточно прав")
		}
		var target models.Doctor
		if err := tx.Where("id = ?", DoctorID).First(&target).Error; err != nil {
			return errors.New("ошибка: целевой врач не найден")
		}

		if target.Root <= 0 {
			return errors.New("у врача уже минимальный уровень прав (0)")
		}
		if requester.Root == 1 {
			if requester.ClinicID != target.ClinicID {
				return errors.New("отказано: вы можете управлять только врачами своей клиники")
			}
			if target.Root > 1 {
				return errors.New("отказано: вы не можете понизить суперадмина")
			}
		}
		return tx.Model(&models.Doctor{}).Where("id = ?", DoctorID).Update("Root", target.Root-1).Error
	})
}

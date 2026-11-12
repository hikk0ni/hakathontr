package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Таблица аккаунтов

type Auth struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"not null"`
}

// Таблица врачей

type Doctor struct {
	ID               uint `gorm:"primaryKey"`
	AuthID           uint // Ключ аккаунта
	Auth             Auth `gorm:"constraint:OnDelete:CASCADE;"`
	ClinicID         uint `gorm:"not null"` // Ключ клиники
	Name             string
	Surname          string
	Speciality       string
	Root             int  `gorm:"default:0"` // 0 - врач, 1 - админ, 2 - суперпользователь
	AccountAvailable bool `gorm:"default:false"`
}

// Таблица пациентов

type Patient struct {
	ID             uint `gorm:"primaryKey"`
	AuthID         uint
	Auth           Auth `gorm:"constraint:OnDelete:CASCADE;"`
	ClinicID       uint `gorm:"not null"`
	Name           string
	Surname        string
	Disease        string
	MedicalRecords []MedicalRecord `gorm:"constraint:OnDelete:CASCADE;"`
}

// Таблица клиник

type Clinic struct {
	ID                uint `gorm:"primaryKey"`
	Name              string
	Address           string
	Doctors           []Doctor  `gorm:"constraint:OnDelete:CASCADE;"`
	Patients          []Patient `gorm:"constraint:OnDelete:CASCADE;"`
	ClinicIsAvailable bool      `gorm:"default:false"`
}

// Таблицуа Медицинских Карт

type MedicalRecord struct {
	ID           uint   `gorm:"primaryKey"`
	PatientID    uint   `gorm:"not null"`
	DoctorID     uint   `gorm:"not null"`
	Diagnosis    string `gorm:"not null"`
	Prescription string
	Date         string
}

type Message struct {
	ID               uint   `gorm:"primaryKey"`
	SenderID         uint   `gorm:"not null"`
	RecipientID      uint   `gorm:"not null"`
	MessageContainer string `gorm:"not null"`
	CreatedAt        time.Time
}

type Queue struct {
	ID        uint `gorm:"primaryKey"`
	PatientID uint `gorm:"not null"`
	DoctorID  uint `gorm:"not null"`
}

// gin обработчики

type RequestCreateDoctor struct {
	Email      string `json:"Email" binding:"required,email"`
	Password   string `json:"Password" binding:"required,min=6"`
	ClinicID   uint   `json:"ClinicID" binding:"required"`
	Name       string `json:"Name" binding:"required"`
	Surname    string `json:"Surname" binding:"required"`
	Speciality string `json:"Speciality" binding:"required"`
}

type RequestAuth struct {
	Email    string `json:"Email" binding:"required,email"`
	Password string `json:"Password" binding:"required,min=6"`
}

type RequestCreateAdmin struct {
	Email      string `json:"Email" binding:"required,email"`
	Password   string `json:"Password" binding:"required,min=6"`
	Name       string `json:"Name" binding:"required"`
	Surname    string `json:"Surname" binding:"required"`
	Speciality string `json:"Speciality" binding:"required"`
}

type RequestCreateClinic struct {
	Name    string `json:"Name" binding:"required"`
	Address string `json:"Address" binding:"required"`
	Admin   RequestCreateAdmin
}

type RequestCreatePatient struct {
	Email    string `json:"Email" binding:"required,email"`
	Password string `json:"Password" binding:"required,min=6"`
	ClinicID uint   `json:"ClinicID" binding:"required"`
	Name     string `json:"Name" binding:"required"`
	Surname  string `json:"Surname" binding:"required"`
	Disease  string `json:"Disease" binding:"required"`
}

type RequestCreateRecord struct {
	PatientID    uint   `json:"PatientID" binding:"required"`
	Diagnosis    string `json:"Diagnosis" binding:"required"`
	Prescription string `json:"Prescription"`
}

type RequestSendMessage struct {
	RecipientEmail   string `json:"RecipientEmail" binding:"required,email"`
	MessageContainer string `json:"MessageContainer" binding:"required"`
}

type RequestUpdateAuth struct {
	Email    string `json:"Email" binding:"required,email"`
	Password string `json:"Password" binding:"required,min=6"`
}

type RequestUpdateDoctor struct {
	ClinicID   uint   `json:"ClinicID"`
	Name       string `json:"Name"`
	Surname    string `json:"Surname"`
	Speciality string `json:"Speciality"`
	Auth       RequestUpdateAuth
}

type RequestUpdatePatient struct {
	ClinicID uint   `json:"ClinicID"`
	Name     string `json:"Name"`
	Surname  string `json:"Surname"`
	Disease  string `json:"Disease"`
	Auth     RequestUpdateAuth
}

type RequestUpdateClinic struct {
	Name    string `json:"Name" binding:"required"`
	Address string `json:"Address" binding:"required"`
}

type RequestUpdateRecord struct {
	RecordID     uint   `json:"RecordID" binding:"required"`
	DoctorID     uint   `json:"DoctorID" binding:"required"`
	Prescription string `json:"Prescription" binding:"required"`
	Date         string `json:"Date" binding:"required"`
}

type RequestUpdateMessage struct {
	MessageID        uint   `json:"Message_ID" binding:"required"`
	MessageContainer string `json:"MessageContainer" binding:"required"`
}

type RequestAppointment struct {
	PatientID uint `gorm:"not null"`
	DoctorID  uint `gorm:"not null"`
}

// DTO-SENDER

type PublicDoctor struct {
	DoctorID   uint   `json:"DoctorID"`
	ClinicID   uint   `json:"ClinicID"`
	Name       string `json:"Name"`
	Surname    string `json:"Surname"`
	Speciality string `json:"Speciality"`
}

type PublicPatient struct {
	PatientID uint   `json:"PatientID"`
	ClinicID  uint   `json:"ClinicID"`
	Name      string `json:"Name"`
	Surname   string `json:"Surname"`
	Disease   string `json:"Disease"`
}

type PublicClinic struct {
	ID      uint   `json:"ID"`
	Name    string `json:"Name"`
	Address string `json:"Address"`
}

// JWT-ТОКЕН and DTO JWT

type TokenDataTransfer struct {
	AccessToken    string
	RefreshToken   string
	RefreshTokenID string
	AtExpires      int64
	RtExpires      int64
}

type AccessTokenData struct {
	AuthID uint   `json:"AuthID"`
	Role   string `json:"Role"`
	jwt.RegisteredClaims
}

type RefreshTokenData struct {
	RtID   string `json:"RtID"`
	AuthID uint   `json:"AuthID"`
	jwt.RegisteredClaims
}

type RefreshToken struct {
	RefreshToken string `json:"RefreshToken"`
}

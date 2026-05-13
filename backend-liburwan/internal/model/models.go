package model

import (
	"time"

	"github.com/google/uuid"
)

type Toko struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Nama      string    `gorm:"type:varchar(255);not null" json:"nama"`
	IsPusat   bool      `gorm:"default:false" json:"is_pusat"`
	CreatedAt time.Time `json:"created_at"`
}

func (Toko) TableName() string {
	return "toko"
}

type Karyawan struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Nama      string    `gorm:"type:varchar(255);not null" json:"nama"`
	Email     string    `gorm:"type:varchar(255);unique" json:"email"`
	GoogleID  string    `gorm:"type:varchar(255);unique" json:"google_id"`
	Role      string    `gorm:"type:role_type;not null" json:"role"`
	TokoID    uuid.UUID `gorm:"type:uuid" json:"toko_id"`
	Toko      Toko      `gorm:"foreignKey:TokoID" json:"toko,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (Karyawan) TableName() string {
	return "karyawan"
}

type JadwalLibur struct {
	ID               uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	KaryawanID       uuid.UUID         `gorm:"type:uuid" json:"karyawan_id"`
	Karyawan         Karyawan          `gorm:"foreignKey:KaryawanID" json:"karyawan,omitempty"`
	Tanggal          time.Time         `gorm:"type:date;not null" json:"tanggal"`
	Tipe             string            `gorm:"type:leave_type;not null" json:"tipe"`
	BackupAssignment *BackupAssignment `gorm:"foreignKey:JadwalLiburID" json:"backup_assignment,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

func (JadwalLibur) TableName() string {
	return "jadwal_libur"
}

type BackupAssignment struct {
	ID               uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	JadwalLiburID    uuid.UUID    `gorm:"type:uuid" json:"jadwal_libur_id"`
	JadwalLibur      *JadwalLibur `gorm:"foreignKey:JadwalLiburID;constraint:OnDelete:CASCADE" json:"jadwal_libur,omitempty"`
	BackupKaryawanID uuid.UUID    `gorm:"type:uuid" json:"backup_karyawan_id"`
	BackupKaryawan   Karyawan     `gorm:"foreignKey:BackupKaryawanID" json:"backup_karyawan,omitempty"`
	AssignedBy       uuid.UUID    `gorm:"type:uuid" json:"assigned_by"`
	Assigner         Karyawan     `gorm:"foreignKey:AssignedBy" json:"assigner,omitempty"`
	CreatedAt        time.Time    `json:"created_at"`
}

func (BackupAssignment) TableName() string {
	return "backup_assignment"
}

type Konfigurasi struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Key       string    `gorm:"type:varchar(255);unique;not null" json:"key"`
	Value     string    `gorm:"type:varchar(255);not null" json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuditLog struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	KaryawanID *uuid.UUID `gorm:"type:uuid" json:"karyawan_id"`
	Action     string    `gorm:"type:varchar(255);not null" json:"action"`
	Entity     string    `gorm:"type:varchar(255);not null" json:"entity"`
	EntityID   uuid.UUID `gorm:"type:uuid;not null" json:"entity_id"`
	Payload    string    `gorm:"type:jsonb" json:"payload"`
	CreatedAt  time.Time `json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_log"
}

func (Konfigurasi) TableName() string {
	return "konfigurasi"
}

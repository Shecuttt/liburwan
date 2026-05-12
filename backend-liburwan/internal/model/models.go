package model

import (
	"time"

	"github.com/google/uuid"
)

type Toko struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Nama      string    `gorm:"type:varchar(255);not null"`
	IsPusat   bool      `gorm:"default:false"`
	CreatedAt time.Time
}

func (Toko) TableName() string {
	return "toko"
}

type Karyawan struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Nama      string    `gorm:"type:varchar(255);not null"`
	Role      string    `gorm:"type:role_type;not null"` // Using the custom enum type
	TokoID    uuid.UUID `gorm:"type:uuid"`
	Toko      Toko      `gorm:"foreignKey:TokoID"`
	CreatedAt time.Time
}

func (Karyawan) TableName() string {
	return "karyawan"
}

type JadwalLibur struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	KaryawanID uuid.UUID `gorm:"type:uuid"`
	Karyawan   Karyawan  `gorm:"foreignKey:KaryawanID"`
	Tanggal    time.Time `gorm:"type:date;not null"`
	Tipe       string    `gorm:"type:leave_type;not null"` // Using the custom enum type
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (JadwalLibur) TableName() string {
	return "jadwal_libur"
}

type BackupAssignment struct {
	ID               uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	JadwalLiburID    uuid.UUID   `gorm:"type:uuid"`
	JadwalLibur      JadwalLibur `gorm:"foreignKey:JadwalLiburID;constraint:OnDelete:CASCADE"`
	BackupKaryawanID uuid.UUID   `gorm:"type:uuid"`
	BackupKaryawan   Karyawan    `gorm:"foreignKey:BackupKaryawanID"`
	AssignedBy       uuid.UUID   `gorm:"type:uuid"`
	Assigner         Karyawan    `gorm:"foreignKey:AssignedBy"`
	CreatedAt        time.Time
}

func (BackupAssignment) TableName() string {
	return "backup_assignment"
}

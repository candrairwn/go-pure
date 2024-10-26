package entities

import (
	"github.com/google/uuid"
)

type User struct {
	Id         uuid.UUID `gorm:"column:id;type:uuid;primaryKey;not null"`
	Username   string    `gorm:"column:username;type:varchar(255);not null"`
	Password   string    `gorm:"column:password;type:varchar(255);null"`
	IdTipeUser string    `gorm:"column:id_tipe_user;type:varchar(2);null"`
	IdProdi    string    `gorm:"column:id_prodi;type:varchar(25);null"`
	NamaProdi  string    `gorm:"column:nama_prodi;type:varchar(255);null"`
}

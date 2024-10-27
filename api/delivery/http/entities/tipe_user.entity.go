package entities

type TipeUser struct {
	Id   string `gorm:"column:id;type:varchar(2);primaryKey;not null"`
	Nama string `gorm:"column:nama;type:varchar(255);not null"`
}

func (t *TipeUser) TableName() string {
	return "cbt_lsp.mst_tipe_user"
}

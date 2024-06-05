package domain

type Role struct {
	RoleID      uint `gorm:"primarykey"`
	RoleName    string
	Description *string
}

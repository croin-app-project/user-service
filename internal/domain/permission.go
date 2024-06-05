package domain

type Permission struct {
	PermissionID   uint `gorm:"primarykey"`
	PermissionName string
	Description    *string
}

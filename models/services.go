package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewServices(connectionInfo string) (*Services, error) {
    _, err := gorm.Open(postgres.Open(connectionInfo))
    if err != nil {
        return nil, err
    }
    return &Services{}, nil
}

type Services struct {
    Gallery GalleryService
    User UserService
}

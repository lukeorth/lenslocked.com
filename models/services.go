package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewServices(connectionInfo string) (*Services, error) {
    db, err := gorm.Open(postgres.Open(connectionInfo))
    if err != nil {
        return nil, err
    }
    return &Services{
        User: NewUserService(db),
        db: db,
    }, nil
}

type Services struct {
    Gallery GalleryService
    User UserService
    db *gorm.DB
}

// DestructiveReset drops all tables and rebuilds them
func (s *Services) DestructiveReset() error {
    if err := s.db.Migrator().DropTable(&User{}, &Gallery{}); err != nil {
        return err
    }
    return s.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate all tables
func (s *Services) AutoMigrate() error {
    return s.db.Migrator().AutoMigrate(&User{}, &Gallery{})
}


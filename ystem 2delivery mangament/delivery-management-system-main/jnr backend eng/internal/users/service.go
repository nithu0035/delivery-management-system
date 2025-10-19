package users

import (
    "context"
    "errors"

    "gorm.io/gorm"
    "golang.org/x/crypto/bcrypt"
)

type Service struct {
    db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
    return &Service{db: db}
}

func (s *Service) Create(ctx context.Context, email, password, role string) (*User, error) {
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    u := &User{Email: email, PasswordHash: string(hash), Role: role}
    if err := s.db.WithContext(ctx).Create(u).Error; err != nil {
        return nil, err
    }
    return u, nil
}

func (s *Service) Authenticate(ctx context.Context, email, password string) (*User, error) {
    var u User
    if err := s.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
        return nil, err
    }
    if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
        return nil, errors.New("invalid credentials")
    }
    return &u, nil
}

func (s *Service) GetByID(ctx context.Context, id uint) (*User, error) {
    var u User
    if err := s.db.WithContext(ctx).First(&u, id).Error; err != nil {
        return nil, err
    }
    return &u, nil
}

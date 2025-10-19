package auth

import (
    "context"
    "errors"    
    "time"

    "delivery/internal/users"

    "github.com/golang-jwt/jwt/v5"
)

type Service struct {
    secret string
    users  *users.Service
}

func NewService(secret string, usersvc *users.Service) *Service {
    if secret == "" { secret = "devsecret" }
    return &Service{secret: secret, users: usersvc}
}

func (s *Service) GenerateToken(ctx context.Context, u *users.User) (string, error) {
    claims := jwt.MapClaims{
        "sub": u.ID,
        "email": u.Email,
        "role": u.Role,
        "exp": time.Now().Add(24 * time.Hour).Unix(),
    }
    t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return t.SignedString([]byte(s.secret))
}

func (s *Service) ParseToken(ctx context.Context, tokenStr string) (map[string]interface{}, error) {
    t, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return []byte(s.secret), nil
    })
    if err != nil || !t.Valid { return nil, errors.New("invalid token") }
    if claims, ok := t.Claims.(jwt.MapClaims); ok {
        return claims, nil
    }
    return nil, errors.New("invalid token claims")
}

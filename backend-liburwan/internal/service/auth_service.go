package service

import (
	"backend-liburwan/internal/model"
	"backend-liburwan/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService struct {
	repo         *repository.KaryawanRepository
	oauthConfig  *oauth2.Config
	jwtSecret    []byte
	jwtExpiry    time.Duration
}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func NewAuthService(repo *repository.KaryawanRepository) *AuthService {
	expiry, _ := time.ParseDuration(os.Getenv("JWT_EXPIRY"))
	if expiry == 0 {
		expiry = 24 * time.Hour
	}

	return &AuthService{
		repo: repo,
		oauthConfig: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Endpoint:     google.Endpoint,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		},
		jwtSecret: []byte(os.Getenv("JWT_SECRET")),
		jwtExpiry: expiry,
	}
}

func (s *AuthService) GetAuthURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state)
}

func (s *AuthService) HandleGoogleCallback(code string) (string, *model.Karyawan, error) {
	token, err := s.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return "", nil, err
	}

	client := s.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	var gUser GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&gUser); err != nil {
		return "", nil, err
	}

	// 1. Check if karyawan exists by email
	karyawan, err := s.repo.GetByEmail(gUser.Email)
	if err != nil {
		return "", nil, errors.New("USER_NOT_REGISTERED")
	}

	// 2. Update google_id if not set
	if karyawan.GoogleID == "" {
		s.repo.UpdateGoogleID(karyawan.ID, gUser.ID)
		karyawan.GoogleID = gUser.ID
	}

	// 3. Generate JWT
	jwtToken, err := s.GenerateJWT(karyawan)
	if err != nil {
		return "", nil, err
	}

	return jwtToken, karyawan, nil
}

func (s *AuthService) GenerateJWT(karyawan *model.Karyawan) (string, error) {
	claims := jwt.MapClaims{
		"karyawan_id": karyawan.ID.String(),
		"role":        karyawan.Role,
		"toko_id":      karyawan.TokoID.String(),
		"exp":          time.Now().Add(s.jwtExpiry).Unix(),
		"iat":          time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

package auth

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type token struct {
	slg *slog.Logger
}

func NewToken(slg *slog.Logger) *token {
	return &token{
		slg: slg,
	}
}

var _ TokenInterface = &token{}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	TokenUuid    string
	RefreshUuid  string
	ATExpiresAt  int64
	RTExpiresAt  int64
}

type AccessDetails struct {
	TokenUuid string
	UID       string
	Role      string
}

type TokenInterface interface {
	CreateToken(uid, role string) (*TokenDetails, error)
	CreateRefreshToken(uid, role string, td *TokenDetails) (string, error)
	ExtractMetadata(*http.Request) (*AccessDetails, error)
}

func (t *token) CreateToken(uid, role string) (*TokenDetails, error) {
	td := &TokenDetails{}

	var err error

	td.ATExpiresAt = time.Now().Add(time.Minute * 15).Unix()
	td.TokenUuid = uuid.New().String()

	authClaims := jwt.MapClaims{}
	authClaims["authorized"] = true
	authClaims["access_uuid"] = td.TokenUuid
	authClaims["uid"] = uid
	authClaims["role"] = role
	authClaims["expires_at"] = td.ATExpiresAt

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims)

	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		t.slg.Error("Error signing access token", "error context", err, "function", "CreateToken")
		return nil, err
	}

	td.RefreshToken, err = t.CreateRefreshToken(uid, role, td)
	if err != nil {
		t.slg.Error("Error creating refresh token", "error context", err, "function", "CreateToken")
		return nil, err
	}

	return td, nil
}

func extractToken(r *http.Request, slg *slog.Logger) string {
	authToken := r.Header.Get("Authorization")

	arr := strings.Split(authToken, " ")
	if len(arr) != 2 {
		slg.Error("header is invalid ", "error context", "Invalid header provided", "function", "extractToken")
		return ""
	}

	return arr[1]
}

func verifyToken(r *http.Request, slg *slog.Logger) (*jwt.Token, error) {
	tokenString := extractToken(r, slg)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			slg.Error("Error parsing token", "error context", "unexpected signing method", "function", "verifyToken")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if err != nil {
		slg.Error("Error parsing token", "error context", err, "function", "verifyToken")
		return nil, err
	}

	return token, nil
}

func ValidateToken(r *http.Request, slg *slog.Logger) error {
	token, err := verifyToken(r, slg)
	if err != nil {
		slg.Error("Error verifying token", "error context", err, "function", "ValidateToken")
		return err
	}

	if !token.Valid {
		slg.Error("Token is invalid", "error context", "Token is invalid", "function", "ValidateToken")
		return fmt.Errorf("token is invalid")
	}

	return nil
}

func (t *token) ExtractMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := verifyToken(r, t.slg)
	if err != nil {
		t.slg.Error("Error verifying token", "error context", err, "function", "ExtractMetadata")
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			t.slg.Error("Error extracting access uuid", "error context", "Cannot get access uuid", "function", "ExtractMetadata")
			return nil, err
		}

		uid, ok := claims["uid"].(string)
		if !ok {
			t.slg.Error("Error extracting uid", "error context", "Cannot get uid", "function", "ExtractMetadata")
			return nil, err
		}

		role, ok := claims["role"].(string)
		if !ok {
			t.slg.Error("Error extracting role", "error context", "Cannot get role", "function", "ExtractMetadata")
			return nil, err
		}

		return &AccessDetails{
			TokenUuid: accessUuid,
			UID:       uid,
			Role:      role,
		}, nil
	}
	return nil, nil
}

func (t *token) CreateRefreshToken(uid, role string, td *TokenDetails) (string, error) {
	td.RTExpiresAt = time.Now().Add(time.Hour * 24 * 3).Unix()
	td.RefreshUuid = td.TokenUuid + "++" + uid

	refreshClaims := jwt.MapClaims{}
	refreshClaims["refresh_uuid"] = td.RefreshUuid
	refreshClaims["uid"] = uid
	refreshClaims["role"] = role
	refreshClaims["expires_at"] = td.RTExpiresAt

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	refreshToken, err := rt.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	if err != nil {
		t.slg.Error("Error signing refresh token", "error context", err, "function", "CreateRefreshToken")
		return "", err
	}

	return refreshToken, nil
}

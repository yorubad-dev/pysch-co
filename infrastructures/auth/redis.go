package auth

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisClient struct {
	c      *redis.Client
	logger *slog.Logger
	ctx    context.Context
}

func NewRedis(logger *slog.Logger) *redisClient {
	c := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return &redisClient{
		c:      c,
		logger: logger,
		ctx:    context.Background(),
	}
}

var _ RedisInterface = &redisClient{}

type RedisInterface interface {
	CreateAuth(uid, role string, td *TokenDetails) error
	FetchAuth(access_uuid string) (string, error)
	FetchRefresh(refresh_uid string) (string, error)
	DeleteRefresh(refresh_uuid string) error
	DeleteTokens(access_details *AccessDetails) error
	CreateOTPReferenceID(reference_id string) error
	FetchOTPReferenceID(reference_id string) (bool, error)
}

func (rc *redisClient) CreateAuth(uid, role string, td *TokenDetails) error {
	atExpiresAt := time.Unix(td.ATExpiresAt, 0)
	rtExpiresAt := time.Unix(td.RTExpiresAt, 0)
	now := time.Now()

	// save the role too
	// data := map[string]string{"uid": uid, "role": role}
	// create a function that extracts the data then unmarshall into a structure
	if err := rc.c.Set(rc.ctx, td.TokenUuid, uid, atExpiresAt.Sub(now)).Err(); err != nil {
		rc.logger.Error("Error setting access token in redis", "error context", err, "function", "CreateAuth")
		return err
	}

	if err := rc.c.Set(rc.ctx, td.RefreshUuid, uid, rtExpiresAt.Sub(now)).Err(); err != nil {
		rc.logger.Error("Error setting refresh token in redis", "error context", err, "function", "CreateAuth")
		return err
	}

	return nil
}

func (rc *redisClient) FetchAuth(access_uuid string) (string, error) {
	uid, err := rc.c.Get(rc.ctx, access_uuid).Result()
	if err != nil {
		rc.logger.Error("Error fetching access token from redis", "error context", err, "function", "FetchAuth")
		return "", err
	}
	return uid, nil
}

func (rc *redisClient) FetchRefresh(refresh_uid string) (string, error) {
	uid, err := rc.c.Get(rc.ctx, refresh_uid).Result()
	if err != nil {
		rc.logger.Error("Error fetching access token from redis", "error context", err, "function", "FetchAuth")
		return "", err
	}
	return uid, nil

}

func (rc *redisClient) DeleteRefresh(refresh_uuid string) error {
	if err := rc.c.Del(rc.ctx, refresh_uuid).Err(); err != nil {
		rc.logger.Error("Error deleting refresh token from redis", "error context", err, "function", "DeleteRefresh")
		return err
	}

	return nil
}

func (rc *redisClient) DeleteTokens(access_details *AccessDetails) error {
	refreshUUID := fmt.Sprintf("%s++%s", access_details.TokenUuid, access_details.UID)

	if err := rc.c.Del(rc.ctx, access_details.TokenUuid).Err(); err != nil {
		rc.logger.Error("Error deleting access token from redis", "error context", err, "function", "DeleteTokens")
		return err
	}

	if err := rc.c.Del(rc.ctx, refreshUUID).Err(); err != nil {
		rc.logger.Error("Error deleting refresh token from redis", "error context", err, "function", "DeleteTokens")
		return err
	}

	return nil
}

func (rc *redisClient) CreateOTPReferenceID(reference_id string) error {
	if err := rc.c.Set(rc.ctx, reference_id, true, 10*time.Minute).Err(); err != nil {
		rc.logger.Error("Error setting OTP reference id in redis", "error context", err, "function", "CreateOTPReferenceID")
		return err
	}
	return nil
}

func (rc *redisClient) FetchOTPReferenceID(reference_id string) (bool, error) {
	status, err := rc.c.Get(rc.ctx, reference_id).Bool()
	if err != nil {
		rc.logger.Error("Error fetching OTP reference id from redis", "error context", err, "function", "FetchOTPReferenceID")
		return false, err
	}
	return status, nil
}

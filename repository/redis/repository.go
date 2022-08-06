package redis

import (
	"fmt"
	"shortner/core"
	"strconv"

	"github.com/pkg/errors"
	errs "github.com/pkg/errors"

	"github.com/go-redis/redis"
)

type redisRepository struct {
	client *redis.Client
}

func newRedisClient(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	// check redi server if it's active
	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

// create new repo for redis
func NewRedisRepository(redisURL string) (core.RedirectRepository, error) {
	repo := &redisRepository{}
	client, err := newRedisClient(redisURL)
	if err != nil {
		return nil, errs.Wrap(err, "repository.NewRedisRepository")
	}
	repo.client = client
	return repo, nil
}

func (r *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (r *redisRepository) Find(code string) (*core.Redirect, error) {
	redirect := &core.Redirect{}
	key := r.generateKey(code)
	data, err := r.client.HGetAll(key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	if len(data) == 0 {
		return nil, errs.Wrap(core.ErrRedirectNotFound, "repository.Redirect.Find")
	}
	// base 10 and parse 64 byte for string
	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	redirect.Code = data["code"]
	redirect.URL = data["url"]
	redirect.CreatedAt = createdAt
	return redirect, nil
}

func (r *redisRepository) Store(redirect *core.Redirect) error {
	key := r.generateKey(redirect.Code)
	data := map[string]interface{}{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	}
	_, err := r.client.HMSet(key, data).Result()
	fmt.Println("bp1", err)
	if err != nil {
		return errs.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}

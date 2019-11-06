package repository

import (
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/gomodule/redigo/redis"
	"net/http"
)

type SessionsRedisRepository struct {
	redisConn redis.Conn
}

func (s *SessionsRedisRepository) GetID(session string) (uint64, error) {
	mkey := "sessions:" + session
	data, err := redis.Uint64(s.redisConn.Do("GET", mkey))
	if err != nil {
		return data, models.NewServerError(err, http.StatusInternalServerError, "can not get session in GetID "+err.Error())
	}
	return data, nil
}

func (s *SessionsRedisRepository) Contains(session string) bool {
	panic("implement me")
}

func (s *SessionsRedisRepository) Put(session string, id uint64) error {
	mkey := "sessions:" + session
	_, err := s.redisConn.Do("SET", mkey, id)
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "can not get session in GetID"+err.Error())
	}
	return nil
}

func (s *SessionsRedisRepository) Remove(session string) error {
	mkey := "sessions:" + session
	_, err := s.redisConn.Do("DEL", mkey)
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "can not get session in GetID"+err.Error())
	}
	return nil
}

func NewSessionRedisStore(redisConn redis.Conn) SessionRepository {
	return &SessionsRedisRepository{redisConn: redisConn}
}

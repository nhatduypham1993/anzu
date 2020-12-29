package deps

import (
	"github.com/mitchellh/goamz/s3"
	"github.com/op/go-logging"
	"github.com/siddontang/ledisdb/ledis"
	"github.com/tryanzu/core/board/legacy/model"
	"gopkg.in/mgo.v2"
	"github.com/go-redis/redis/v8"
)

type Deps struct {
	GamingConfigProvider    *model.GamingRules
	DatabaseSessionProvider *mgo.Session
	DatabaseProvider        *mgo.Database
	LoggerProvider          *logging.Logger
	CacheProvider           *redis.Client
	S3Provider              *s3.Bucket
	LedisProvider           *ledis.DB
}

func (d Deps) GamingConfig() *model.GamingRules {
	return d.GamingConfigProvider
}

func (d Deps) Log() *logging.Logger {
	return d.LoggerProvider

}

func (d Deps) Mgo() *mgo.Database {
	return d.DatabaseProvider
}

func (d Deps) LedisDB() *ledis.DB {
	return d.LedisProvider
}

func (d Deps) MgoSession() *mgo.Session {
	return d.DatabaseSessionProvider
}

func (d Deps) S3() *s3.Bucket {
	return d.S3Provider
}

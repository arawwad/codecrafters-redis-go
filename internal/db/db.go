package db

import (
	"fmt"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/resp/types"
)

type dbValue struct {
	value   types.RespType
	expires time.Time
}

func (val *dbValue) hasExpired() bool {
	return !val.expires.IsZero() && time.Now().After(val.expires)
}

type db map[types.RespType]dbValue

func New() db {
	return db{}
}

func (db db) Set(key types.RespType, value types.RespType, ttl *time.Duration) {
	expires := time.Time{}
	if ttl != nil {
		expires = time.Now().Add(*ttl)
	}
	db[key] = dbValue{
		value:   value,
		expires: expires,
	}
}

func (db db) Get(key types.RespType) (types.RespType, bool) {
	val, ok := db[key]
	if !ok {
		return nil, false
	}

	if val.hasExpired() {
		delete(db, key)
		return nil, false
	}

	return val.value, true
}

func (db db) Incr(key types.RespType) types.RespType {
	val, ok := db.Get(key)
	if !ok {
		db.Set(key, types.BulkString("1"), nil)
		return types.Integer(1)
	}

	num, ok := val.Num()
	if !ok {
		return types.SimpleError("ERR value is not an integer or out of range")
	}

	db.Set(key, types.BulkString(fmt.Sprintf("%d", num+1)), nil)
	return types.Integer(num + 1)
}

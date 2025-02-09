package db

import (
	"fmt"
	"sync"
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

type db struct {
	values map[types.RespType]dbValue
	mu     *sync.RWMutex
}

func New() db {
	return db{values: map[types.RespType]dbValue{}, mu: &sync.RWMutex{}}
}

func (db *db) Set(key types.RespType, value types.RespType, ttl *time.Duration) {
	db.mu.Lock()
	defer db.mu.Unlock()
	expires := time.Time{}
	if ttl != nil {
		expires = time.Now().Add(*ttl)
	}
	db.values[key] = dbValue{
		value:   value,
		expires: expires,
	}
}

func (db *db) Get(key types.RespType) (types.RespType, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	val, ok := db.values[key]
	if !ok {
		return nil, false
	}

	if val.hasExpired() {
		delete(db.values, key)
		return nil, false
	}

	return val.value, true
}

func (db *db) Incr(key types.RespType) types.RespType {
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

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

type DB struct {
	values map[types.RespType]dbValue
	mu     *sync.RWMutex
}

func New() DB {
	return DB{values: map[types.RespType]dbValue{}, mu: &sync.RWMutex{}}
}

func (db *DB) Set(key types.RespType, value types.RespType, ttl *time.Duration) {
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

func (db *DB) Get(key types.RespType) (types.RespType, bool) {
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

func (db *DB) Incr(key types.RespType) types.RespType {
	// TODO keep the expires value when incrementing
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

func (db *DB) AppendToStream(key, id, entryKey, entryValue types.RespType) types.RespType {
	val, ok := db.Get(key)
	if !ok {
		db.Set(key, types.Stream(""), nil)
		return id
	}
	if val.Type() != types.StreamType {
		return types.SimpleError("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	return id
}

package configs

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/VladimirYalumov/logger"
)

const (
	TypeString   = "string"
	TypeInt      = "int"
	TypeFloat    = "float"
	TypeBool     = "bool"
	TypeDuration = "duration"
	TypeIntMap   = "int_map"
)

type WatchFn func(oldValue *Value, newValue *Value) error

type Value struct {
	mu sync.RWMutex

	value any
	key   string
	fns   []WatchFn
	watch atomic.Bool
}

type Connection struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

func New(val any, key string) *Value {
	return &Value{value: val, key: key}
}

func (v *Value) Key() string {
	return v.key
}

func (v *Value) String() string {
	if v == nil {
		return ""
	}

	v.mu.RLock()
	value, ok := v.value.(string)
	v.mu.RUnlock()
	if ok {
		return value
	}
	logErrorType(v, TypeString)
	return ""
}

func (v *Value) Float() float64 {
	if v == nil {
		return 0
	}

	v.mu.RLock()
	value, ok := v.value.(float64)
	v.mu.RUnlock()
	if ok {
		return value
	}
	logErrorType(v, TypeFloat)
	return 0
}

func (v *Value) Bool() bool {
	if v == nil {
		return false
	}

	v.mu.RLock()
	value, ok := v.value.(bool)
	v.mu.RUnlock()
	if ok {
		return value
	}
	logErrorType(v, TypeBool)
	return false
}

func (v *Value) Duration() time.Duration {
	if v == nil {
		return 0
	}

	v.mu.RLock()
	value, ok := v.value.(time.Duration)
	v.mu.RUnlock()
	if ok {
		return value
	}
	logErrorType(v, TypeDuration)
	return 0
}

func (v *Value) Int() int {
	if v == nil {
		return 0
	}

	v.mu.RLock()
	value, ok := v.value.(int)
	v.mu.RUnlock()
	if ok {
		return value
	}
	logErrorType(v, TypeInt)
	return 0
}

func (v *Value) IntMap() map[string]int {
	if v == nil {
		return make(map[string]int)
	}

	v.mu.RLock()
	value, ok := v.value.(map[string]int)
	v.mu.RUnlock()
	if ok {
		return value
	}
	logErrorType(v, TypeDuration)
	return make(map[string]int)
}

func (v *Value) Update(ctx context.Context, val any) {
	if v == nil {
		return
	}

	oldValue := &Value{key: v.key, value: v.value}
	v.mu.Lock()
	v.value = val
	fns := v.fns
	v.mu.Unlock()

	for _, fn := range fns {
		v.mu.RLock()
		if err := fn(oldValue, v); err != nil {
			logger.Error(ctx, err, "update config", "key", v.Key())
		}
		v.mu.RUnlock()
	}
}

func (v *Value) StartWatch() {
	v.watch.Store(true)
}

func (v *Value) IsWatchStarted() bool {
	return v.watch.Load()
}

func (v *Value) AddWatchCallbacks(fns ...WatchFn) {
	v.mu.Lock()
	if v.fns == nil {
		v.fns = make([]WatchFn, 0, len(fns))
	}
	v.fns = append(v.fns, fns...)
	v.mu.Unlock()
}

func logErrorType(v *Value, exceptedType string) {
	logger.Warn(context.Background(),
		fmt.Sprintf("invalid value type: Key: %s, value: %v, Type: %T, Expected: %s",
			v.key, v.value, v.value, exceptedType,
		),
	)
}

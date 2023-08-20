package configs

import "context"

type Provider interface {
	Value(key string) *Value
	Watch(ctx context.Context, key string, fns ...WatchFn) *Value
}

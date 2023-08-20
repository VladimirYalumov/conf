package env

import (
	"context"
	"os"
	"strings"

	"github.com/VladimirYalumov/conf/configs"
	"github.com/VladimirYalumov/conf/keys"
	"github.com/VladimirYalumov/logger"
)

type provider struct {
	values map[string]*configs.Value
}

func Init(ctx context.Context, k map[string]keys.ConfigType) (*provider, error) {
	p := &provider{
		values: make(map[string]*configs.Value),
	}

	for name, key := range k {
		strValue, ok := os.LookupEnv(strings.ToUpper(name))
		if !ok {
			logger.Debug(ctx, "env: value does not exist in env", "key", name)
			if val := configs.Convert(ctx, key.Default, key.Type); val != nil {
				p.values[name] = configs.New(val, name)
			}
			continue
		}
		if val := configs.Convert(ctx, strValue, key.Type); val != nil {
			p.values[name] = configs.New(val, name)
		}
	}

	return p, nil
}

func (y *provider) Value(key string) *configs.Value {
	val, ok := y.values[key]
	if ok {
		return val
	}
	return nil
}

func (y *provider) Watch(_ context.Context, key string, _ ...configs.WatchFn) *configs.Value {
	return y.Value(key)
}

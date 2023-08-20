package conf

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	keysData = `configs:
 key1:
  type: string
  default: "str"
  source: yaml
 key2:
  type: int
  default: 2
  source: yaml
 key3:
  type: float
  default: 1.2
  source: yaml
 key4:
  type: bool
  default: false
  source: yaml
 key5:
  type: duration
  default: 40m
  source: yaml
 key6:
  type: string
  default: "str6"
  source: yaml
 key7:
  type: string
  default: "str3"
  source: yaml
 key10:
  type: int_map
  default: '{"4242424242424242": 3}'
  source: yaml
 key11:
  type: string
  default: "default_env_1"
  source: env
 key12:
  type: string
  default: 'default_env_2'
  source: env`

	configsData = `values:
  key1: "str"
  key2: 1
  key3: 1.1
  key4: true
  key5: 30m
  key6: "str2"
  key10: 2222222222222222:1,5454545454545454:2`
)

func Test_YamlValue(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()

	r := strings.NewReader(keysData)
	r2 := strings.NewReader(configsData)

	_ = os.Setenv("KEY11", "real_env_1")

	_ = New().KeysReader(r).YamlConfigs(r2).Init(ctx)

	t.Run("from yaml", func(t *testing.T) {
		t.Parallel()
		// Act
		stringVal := Value(ctx, "key1").String()
		intVal := Value(ctx, "key2").Int()
		floatVal := Value(ctx, "key3").Float()
		boolVal := Value(ctx, "key4").Bool()
		durationVal := Value(ctx, "key5").Duration()
		exceptDuration, _ := time.ParseDuration("30m")

		exceptIntMap := make(map[string]int)
		exceptIntMap["2222222222222222"] = 1
		exceptIntMap["5454545454545454"] = 2

		stringVal2 := Value(ctx, "key6").String()

		intMap := Value(ctx, "key10").IntMap()

		// Assert
		assert.Equal(t, "str", stringVal)
		assert.Equal(t, "str2", stringVal2)
		assert.Equal(t, 1, intVal)
		assert.Equal(t, 1.1, floatVal)
		assert.Equal(t, true, boolVal)
		assert.Equal(t, exceptDuration, durationVal)
		assert.Equal(t, exceptIntMap, intMap)
	})
	t.Run("from default yaml", func(t *testing.T) {
		t.Parallel()
		// Act
		stringVal := Value(ctx, "key7").String()

		// Assert
		assert.Equal(t, "str3", stringVal)
	})
	t.Run("from default yaml", func(t *testing.T) {
		t.Parallel()
		// Act
		envVal1 := Value(ctx, "key11").String()
		envVal2 := Value(ctx, "key12").String()

		// Assert
		assert.Equal(t, "real_env_1", envVal1)
		assert.Equal(t, "default_env_2", envVal2)

		_ = os.Unsetenv("KEY11")
	})
}

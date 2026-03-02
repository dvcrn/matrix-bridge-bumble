package bumble

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
)

var (
	lastErr error
	env     Env
	once    sync.Once
)

type Env struct {
	DatabasePath string `envconfig:"DATABASE_PATH"`
	Localpart    string `envconfig:"LOCALPART"`
}

func Process() (Env, error) {
	once.Do(func() {
		lastErr = envconfig.Process("", &env)
	})
	return env, lastErr
}

func Reload() (Env, error) {
	lastErr = envconfig.Process("", &env)
	return env, lastErr
}

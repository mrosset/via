package via

import (
	"fmt"
	"sort"
)

// Env provides a hash map for environment variables
type Env map[string]string

// Returns alphabetically sorted keys
func (env Env) keys() (keys []string) {
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return
}

// Value returns key=value for maps key
func (env Env) Value(key string) string {
	return fmt.Sprintf("%s=%s", key, env[key])
}

// KeyValue return Env as key=value sorted alphabetically by key
func (env Env) KeyValue() (kv []string) {
	for _, k := range env.keys() {
		kv = append(kv, env.Value(k))
	}
	return
}

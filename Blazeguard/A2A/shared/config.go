package shared

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if p := strings.TrimSpace(os.Getenv("BLAZEGUARD_ENV_FILE")); p != "" {
		_ = godotenv.Load(p)
		return
	}

	wd, err := os.Getwd()
	if err == nil {
		dir := wd
		for {
			envPath := filepath.Join(dir, ".env")
			if _, e := os.Stat(envPath); e == nil {
				_ = godotenv.Load(envPath)
				return
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	_ = godotenv.Load()
}

func GetEnv(key, fallback string) string {
	v := os.Getenv(key)
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}

func RequireEnv(keys ...string) error {
	var missing []string
	for _, k := range keys {
		if strings.TrimSpace(os.Getenv(k)) == "" {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required env: %s", strings.Join(missing, ", "))
	}
	return nil
}

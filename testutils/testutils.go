package testutils

import (
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

const projectDirName = "toastnotes-server"

// LoadEnv loads env vars from .env.
func LoadEnv() error {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	return godotenv.Load(string(rootPath) + `/.env`)
}

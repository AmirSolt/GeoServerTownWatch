package base

import "github.com/gin-gonic/gin"

type Base struct {
	RootDir string
	*Env
	*Config
	*DB
	*gin.Engine
}

func LoadBase() *Base {

	base := Base{
		RootDir: "./",
	}

	base.loadEnv()
	base.loadConfig()
	base.loadDB()
	base.loadEngine()
	base.loadLogging()

	return &base
}

func (base *Base) Kill() {
	base.killDB()
	base.killLogging()
}

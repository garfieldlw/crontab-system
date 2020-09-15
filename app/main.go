package main

import (
	"github.com/garfieldlw/crontab-system/app/router"
	"github.com/garfieldlw/crontab-system/library/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	engine.Use(middleware.Cors())
	engine.Use(middleware.Logger())
	router.LoadRouter(engine)
	_ = engine.Run(":3005")
}

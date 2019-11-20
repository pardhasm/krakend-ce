package custom

import (
	krakendcors "github.com/devopsfaith/krakend-cors"
	"github.com/devopsfaith/krakend/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	//wrapper "github.com/rs/cors/wrapper/gin"
)

// CorsNew returns a gin.HandlerFunc with the CORS configuration provided in the ExtraConfig
func CorsNew(e config.ExtraConfig) gin.HandlerFunc {
	tmp := krakendcors.ConfigGetter(e)
	if tmp == nil {
		return nil
	}
	cfg, ok := tmp.(krakendcors.Config)
	if !ok {
		return nil
	}

	return NewFunc(cors.Options{
		AllowedOrigins:   cfg.AllowOrigins,
		AllowedMethods:   cfg.AllowMethods,
		AllowedHeaders:   cfg.AllowHeaders,
		ExposedHeaders:   cfg.ExposeHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           int(cfg.MaxAge.Seconds()),
	})
}

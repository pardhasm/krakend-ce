// Package cors/wrapper/gin provides gin.HandlerFunc to handle CORS related
// requests as a wrapper of github.com/rs/cors handler.
package custom

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

// Options is a configuration container to setup the CORS middleware.
type Options = cors.Options

// corsWrapper is a wrapper of cors.Cors handler which preserves information
// about configured 'optionPassthrough' option.
type corsWrapper struct {
	*cors.Cors
	optionPassthrough bool
}

// build transforms wrapped cors.Cors handler into Gin middleware.
func (c corsWrapper) build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c.HandlerFunc(ctx.Writer, ctx.Request)
		if !c.optionPassthrough &&
			ctx.Request.Method == http.MethodOptions &&
			ctx.GetHeader("Access-Control-Request-Method") != "" {
			ctx.Status(http.StatusNoContent)
		}
		ctx.Writer.Header().Del("Access-Control-Allow-Origin")
		if ctx.Request.Method == http.MethodOptions {
			ctx.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		}
		check("http://uat-twl-web-api.incred.com" + ctx.Request.RequestURI)
	}
}

// AllowAll creates a new CORS Gin middleware with permissive configuration
// allowing all origins with all standard methods with any header and
// credentials.
func AllowAll() gin.HandlerFunc {
	return corsWrapper{Cors: cors.AllowAll()}.build()
}

// Default creates a new CORS Gin middleware with default options.
func Default() gin.HandlerFunc {
	return corsWrapper{Cors: cors.Default()}.build()
}

// CorsNew creates a new CORS Gin middleware with the provided options.
func NewFunc(options Options) gin.HandlerFunc {
	return corsWrapper{cors.New(options), options.OptionsPassthrough}.build()
}

func check(url string) {
	var b bytes.Buffer

	r, err := http.NewRequest("OPTIONS", url, &b)
	if err != nil {
		panic(err)
	}
	r.Header.Add("X-Custom", "Copy me!")

	rc, err := http.NewRequest("POST", r.URL.String(), &b)
	if err != nil {
		panic(err)
	}

	rc.Header = r.Header.Clone() // note shallow copy
	fmt.Println("Headers", r.Header, rc.Header)

	// Adjust copy adjusts original
	rc.Header.Add("X-Hello", "World")

	fmt.Println("Headers", r.Header, rc.Header)
}

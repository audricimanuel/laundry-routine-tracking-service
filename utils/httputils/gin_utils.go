package httputils

import "github.com/gin-gonic/gin"

func InvalidateCookie(c *gin.Context, cookieName string) {
	c.SetCookie(
		cookieName,
		"",   // value (empty to delete)
		-1,   // maxAge in seconds (negative = delete immediately)
		"/",  // path
		"",   // domain (empty = default to current)
		true, // secure (set to true if using HTTPS)
		true, // httpOnly
	)
}

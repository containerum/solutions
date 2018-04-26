package middleware

import (
	"context"

	"git.containerum.net/ch/cherry/adaptors/gonic"
	"git.containerum.net/ch/solutions/pkg/sErrors"
	"git.containerum.net/ch/solutions/pkg/server"
	"github.com/gin-gonic/gin"
)

const (
	UserIDHeader   = "X-User-Id"
	UserRoleHeader = "X-User-Role"
)

var hdrToKey = map[string]interface{}{
	UserIDHeader:   server.UserIDContextKey,
	UserRoleHeader: server.UserRoleContextKey,
}

func RequireHeaders(headers ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var notFoundHeaders []string
		for _, v := range headers {
			if ctx.GetHeader(v) == "" {
				notFoundHeaders = append(notFoundHeaders, v)
			}
		}
		if len(notFoundHeaders) > 0 {
			gonic.Gonic(sErrors.ErrRequiredHeadersNotProvided().AddDetails(notFoundHeaders...), ctx)
		}
	}
}

func PrepareContext(ctx *gin.Context) {
	for hn, ck := range hdrToKey {
		if hv := ctx.GetHeader(hn); hv != "" {
			rctx := context.WithValue(ctx.Request.Context(), ck, hv)
			ctx.Request = ctx.Request.WithContext(rctx)
		}
	}
}

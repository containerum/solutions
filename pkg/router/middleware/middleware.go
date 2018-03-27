package middleware

import (
	"context"

	umtypes "git.containerum.net/ch/json-types/user-manager"
	"git.containerum.net/ch/kube-client/pkg/cherry/adaptors/gonic"
	cherryusr "git.containerum.net/ch/kube-client/pkg/cherry/user-manager"
	"git.containerum.net/ch/solutions/pkg/server"
	"github.com/gin-gonic/gin"
)

var hdrToKey = map[string]interface{}{
	umtypes.UserIDHeader:      server.UserIDContextKey,
	umtypes.UserAgentHeader:   server.UserAgentContextKey,
	umtypes.FingerprintHeader: server.FingerPrintContextKey,
	umtypes.SessionIDHeader:   server.SessionIDContextKey,
	umtypes.TokenIDHeader:     server.TokenIDContextKey,
	umtypes.ClientIPHeader:    server.ClientIPContextKey,
	umtypes.PartTokenIDHeader: server.PartTokenIDContextKey,
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
			gonic.Gonic(cherryusr.ErrRequiredHeadersNotProvided().AddDetails(notFoundHeaders...), ctx)
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

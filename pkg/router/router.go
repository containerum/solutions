package router

import (
	"net/http"
	"time"

	//	umtypes "git.containerum.net/ch/json-types/user-manager"
	"git.containerum.net/ch/kube-client/pkg/cherry/adaptors/cherrylog"
	"git.containerum.net/ch/kube-client/pkg/cherry/adaptors/gonic"
	h "git.containerum.net/ch/solutions/pkg/router/handlers"
	m "git.containerum.net/ch/solutions/pkg/router/middleware"

	"git.containerum.net/ch/solutions/pkg/server"

	"git.containerum.net/ch/utils"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	cherryusr "git.containerum.net/ch/kube-client/pkg/cherry/user-manager"
)

//CreateRouter initialises router and middlewares
func CreateRouter(ss *server.SolutionsService) http.Handler {
	e := gin.New()
	initMiddlewares(e, ss)
	initRoutes(e)
	return e
}

func initMiddlewares(e *gin.Engine, ss *server.SolutionsService) {
	/* System */
	e.Use(ginrus.Ginrus(logrus.WithField("component", "gin"), time.RFC3339, true))
	e.Use(gonic.Recovery(cherryusr.ErrInternalError, cherrylog.NewLogrusAdapter(logrus.WithField("component", "gin"))))
	/* Custom */
	e.Use(m.RegisterServices(ss))
	e.Use(m.PrepareContext)
	e.Use(utils.SaveHeaders)
}

// SetupRoutes sets up http router needed to handle requests from clients.
func initRoutes(app *gin.Engine) {
	//	requireIdentityHeaders := m.RequireHeaders(umtypes.UserIDHeader, umtypes.UserRoleHeader)
	//	requireLoginHeaders := m.RequireHeaders(umtypes.UserAgentHeader, umtypes.FingerprintHeader, umtypes.ClientIPHeader)
	//	requireLogoutHeaders := m.RequireHeaders(umtypes.TokenIDHeader, umtypes.SessionIDHeader)

	solutions := app.Group("/solutions")
	{
		solutions.GET("", h.UpdateSolutions, h.GetSolutionsList)
		solutions.GET("/:solution/env", h.UpdateSolutions, h.GetSolutionEnv)
		solutions.GET("/:solution/resources", h.UpdateSolutions, h.GetSolutionResources)
	}
}

package router

import (
	"net/http"
	"time"

	"git.containerum.net/ch/kube-client/pkg/cherry/adaptors/cherrylog"
	"git.containerum.net/ch/kube-client/pkg/cherry/adaptors/gonic"
	h "git.containerum.net/ch/solutions/pkg/router/handlers"
	m "git.containerum.net/ch/solutions/pkg/router/middleware"

	"git.containerum.net/ch/solutions/pkg/server"

	utils "git.containerum.net/ch/utils/httputil"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	cherry "git.containerum.net/ch/kube-client/pkg/cherry/solutions"
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
	e.Use(gonic.Recovery(cherry.ErrInternalError, cherrylog.NewLogrusAdapter(logrus.WithField("component", "gin"))))
	/* Custom */
	e.Use(m.RegisterServices(ss))
	e.Use(m.PrepareContext)
	e.Use(utils.SaveHeaders)
}

// SetupRoutes sets up http router needed to handle requests from clients.
func initRoutes(app *gin.Engine) {
	requireIdentityHeaders := m.RequireHeaders(m.UserIDHeader, m.UserRoleHeader)

	solutions := app.Group("/solutions")
	{
		solutions.GET("", h.UpdateSolutions, h.GetSolutionsList)
		solutions.GET("/:solution/env", h.UpdateSolutions, h.GetSolutionEnv)
		solutions.GET("/:solution/resources", h.UpdateSolutions, h.GetSolutionResources)
	}
	userSolutions := app.Group("/user_solutions")
	{
		userSolutions.GET("", requireIdentityHeaders, h.GetUserSolutionsList)
		userSolutions.GET("/:solution/deployments", requireIdentityHeaders, h.GetUserSolutionsDeployments)
		userSolutions.GET("/:solution/services", requireIdentityHeaders, h.GetUserSolutionsServices)
		userSolutions.POST("", requireIdentityHeaders, h.UpdateSolutions, h.RunSolution)
		userSolutions.DELETE("/:solution", h.DeleteSolution)
	}
}

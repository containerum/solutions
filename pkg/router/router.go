package router

import (
	"net/http"
	"time"

	h "git.containerum.net/ch/solutions/pkg/router/handlers"
	m "git.containerum.net/ch/solutions/pkg/router/middleware"
	"git.containerum.net/ch/solutions/pkg/sErrors"
	"git.containerum.net/ch/solutions/static"
	"github.com/containerum/cherry/adaptors/cherrylog"
	"github.com/containerum/cherry/adaptors/gonic"
	cors "gopkg.in/gin-contrib/cors.v1"

	"git.containerum.net/ch/solutions/pkg/server"

	"github.com/containerum/utils/httputil"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//CreateRouter initialises router and middlewares
func CreateRouter(ss *server.SolutionsService, enableCORS bool) http.Handler {
	e := gin.New()
	initMiddlewares(e, ss, enableCORS)
	initRoutes(e)
	return e
}

func initMiddlewares(e *gin.Engine, ss *server.SolutionsService, enableCORS bool) {
	/* CORS */
	if enableCORS {
		cfg := cors.DefaultConfig()
		cfg.AllowAllOrigins = true
		cfg.AddAllowMethods(http.MethodDelete)
		cfg.AddAllowHeaders(httputil.UserIDXHeader, httputil.UserRoleXHeader)
		e.Use(cors.New(cfg))
	}
	e.Group("/static").
		StaticFS("/", static.HTTP)
	/* System */
	e.Use(ginrus.Ginrus(logrus.WithField("component", "gin"), time.RFC3339, true))
	e.Use(gonic.Recovery(sErrors.ErrInternalError, cherrylog.NewLogrusAdapter(logrus.WithField("component", "gin"))))
	/* Custom */
	e.Use(m.RegisterServices(ss))
	e.Use(httputil.PrepareContext)
	e.Use(httputil.SaveHeaders)
}

// SetupRoutes sets up http router needed to handle requests from clients.
func initRoutes(app *gin.Engine) {
	requireIdentityHeaders := httputil.RequireHeaders(sErrors.ErrRequiredHeadersNotProvided, httputil.UserIDXHeader, httputil.UserRoleXHeader)

	app.Use(requireIdentityHeaders)

	solutions := app.Group("/solutions")
	{
		solutions.Use(h.UpdateSolutions)
		solutions.GET("", h.GetSolutionsList)
		solutions.GET("/:solution/env", h.GetSolutionEnv)
		solutions.GET("/:solution/resources", h.GetSolutionResources)
		solutions.POST("", m.RequireAdminRole, h.AddAvailableSolution)
		solutions.POST("/:solution/activate", m.RequireAdminRole, h.ActivateAvailableSolution)
		solutions.POST("/:solution/deactivate", m.RequireAdminRole, h.DeactivateAvailableSolution)
		solutions.PUT("/:solution", m.RequireAdminRole, h.UpdateAvailableSolution)
		solutions.DELETE("/:solution", m.RequireAdminRole, h.DeleteAvailableSolution)
	}
	userSolutions := app.Group("/user_solutions")
	{
		userSolutions.GET("", h.GetUserSolutionsList)
		userSolutions.GET("/:solution/deployments", h.GetUserSolutionsDeployments)
		userSolutions.GET("/:solution/services", h.GetUserSolutionsServices)
		userSolutions.POST("", h.UpdateSolutions, h.RunSolution)
		userSolutions.DELETE("/:solution", h.DeleteSolution)
	}
}

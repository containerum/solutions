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
func CreateRouter(ss *server.SolutionsService) http.Handler {
	e := gin.New()
	initMiddlewares(e, ss)
	initRoutes(e)
	return e
}

func initMiddlewares(e *gin.Engine, ss *server.SolutionsService) {
	/* CORS */
	cfg := cors.DefaultConfig()
	cfg.AllowAllOrigins = true
	cfg.AddAllowMethods(http.MethodDelete)
	cfg.AddAllowHeaders(httputil.UserIDXHeader, httputil.UserRoleXHeader)
	e.Use(cors.New(cfg))
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
	requireIdentityHeaders := httputil.RequireHeaders(sErrors.ErrInternalError, httputil.UserIDXHeader, httputil.UserRoleXHeader)

	solutions := app.Group("/solutions", requireIdentityHeaders)
	{
		solutions.GET("", h.UpdateSolutions, h.GetSolutionsList)
		solutions.GET("/:solution/env", h.UpdateSolutions, h.GetSolutionEnv)
		solutions.GET("/:solution/resources", h.UpdateSolutions, h.GetSolutionResources)
		solutions.POST("", m.RequireAdminRole, h.UpdateSolutions, h.AddAvailableSolution)
		solutions.PUT("/:solution", m.RequireAdminRole, h.UpdateSolutions, h.UpdateAvailableSolution)
		solutions.PUT("/:solution/activate", m.RequireAdminRole, h.UpdateSolutions, h.ActivateAvailableSolution)
		solutions.PUT("/:solution/deactivate", m.RequireAdminRole, h.UpdateSolutions, h.DeactivateAvailableSolution)
		solutions.DELETE("/:solution", m.RequireAdminRole, h.UpdateSolutions, h.DeleteAvailableSolution)
	}
	userSolutions := app.Group("/user_solutions", requireIdentityHeaders)
	{
		userSolutions.GET("", h.GetUserSolutionsList)
		userSolutions.GET("/:solution/deployments", h.GetUserSolutionsDeployments)
		userSolutions.GET("/:solution/services", h.GetUserSolutionsServices)
		userSolutions.POST("", h.UpdateSolutions, h.RunSolution)
		userSolutions.DELETE("/:solution", h.DeleteSolution)
	}
}

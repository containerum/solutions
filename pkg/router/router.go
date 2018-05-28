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

	templates := app.Group("/templates")
	{
		templates.GET("", h.GetTemplatesList)
		templates.GET("/:solution/env", h.GetTemplatesEnv)
		templates.GET("/:solution/resources", h.GetTemplatesResources)
		templates.POST("", httputil.RequireAdminRole(sErrors.ErrAdminRequired), h.AddTemplate)
		templates.POST("/:solution/activate", httputil.RequireAdminRole(sErrors.ErrAdminRequired), h.ActivateTemplate)
		templates.POST("/:solution/deactivate", httputil.RequireAdminRole(sErrors.ErrAdminRequired), h.DeactivateTemplate)
		templates.PUT("/:solution", httputil.RequireAdminRole(sErrors.ErrAdminRequired), h.UpdateTemplate)
		templates.DELETE("/:solution", httputil.RequireAdminRole(sErrors.ErrAdminRequired), h.DeleteTemplate)
	}
	solutions := app.Group("/solutions")
	{
		solutions.GET("", h.GetSolutionsList)
		solutions.GET("/:solution/deployments", h.GetSolutionsDeployments)
		solutions.GET("/:solution/services", h.GetSolutionsServices)
		solutions.POST("", h.RunSolution)
		solutions.DELETE("/:solution", h.DeleteSolution)
	}
}

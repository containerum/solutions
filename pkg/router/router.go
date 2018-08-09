package router

import (
	"net/http"
	"time"

	h "git.containerum.net/ch/solutions/pkg/router/handlers"
	m "git.containerum.net/ch/solutions/pkg/router/middleware"
	"git.containerum.net/ch/solutions/pkg/solerrors"
	"git.containerum.net/ch/solutions/static"
	"github.com/containerum/cherry/adaptors/cherrylog"
	"github.com/containerum/cherry/adaptors/gonic"
	"gopkg.in/gin-contrib/cors.v1"

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
		cfg.AddAllowHeaders(httputil.UserIDXHeader, httputil.UserRoleXHeader, httputil.UserNamespacesXHeader)
		e.Use(cors.New(cfg))
	}
	e.Group("/static").
		StaticFS("/", static.HTTP)
	/* System */
	e.Use(ginrus.Ginrus(logrus.WithField("component", "gin"), time.RFC3339, true))
	e.Use(gonic.Recovery(solerrors.ErrInternalError, cherrylog.NewLogrusAdapter(logrus.WithField("component", "gin"))))
	/* Custom */
	e.Use(httputil.SaveHeaders)
	e.Use(httputil.PrepareContext)
	e.Use(m.RequiredUserHeaders())
	e.Use(m.RegisterServices(ss))
}

// SetupRoutes sets up http router needed to handle requests from clients.
func initRoutes(app *gin.Engine) {
	requireIdentityHeaders := httputil.RequireHeaders(solerrors.ErrRequiredHeadersNotProvided, httputil.UserIDXHeader, httputil.UserRoleXHeader)

	app.Use(requireIdentityHeaders)

	templates := app.Group("/templates")
	{
		templates.GET("", h.GetTemplatesList)
		templates.GET("/:template/env", h.GetTemplatesEnv)
		templates.GET("/:template/resources", h.GetTemplatesResources)
		templates.POST("", httputil.RequireAdminRole(solerrors.ErrAdminRequired), h.AddTemplate)
		templates.POST("/:template/activate", httputil.RequireAdminRole(solerrors.ErrAdminRequired), h.ActivateTemplate)
		templates.POST("/:template/deactivate", httputil.RequireAdminRole(solerrors.ErrAdminRequired), h.DeactivateTemplate)
		templates.PUT("/:template", httputil.RequireAdminRole(solerrors.ErrAdminRequired), h.UpdateTemplate)
	}
	solutions := app.Group("/solutions")
	{
		solutions.GET("", h.GetSolutionsList)
		solutions.DELETE("", h.DeleteSolutions)
	}
	namespaceSolutions := app.Group("/namespaces/:namespace/solutions")
	{
		namespaceSolutions.GET("", m.ReadAccess, h.GetNamespaceSolutions)
		namespaceSolutions.GET("/:solution", m.ReadAccess, h.GetSolution)
		namespaceSolutions.GET("/:solution/deployments", m.ReadAccess, h.GetSolutionsDeployments)
		namespaceSolutions.GET("/:solution/services", m.ReadAccess, h.GetSolutionsServices)
		namespaceSolutions.POST("", m.WriteAccess, h.RunSolution)
		namespaceSolutions.DELETE("/:solution", m.DeleteAccess, h.DeleteSolution)
		namespaceSolutions.DELETE("", m.DeleteAccess, h.DeleteNamespaceSolutions)
	}
}

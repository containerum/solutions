package router

import (
	"net/http"
	"time"

	"git.containerum.net/ch/auth/static"
	h "git.containerum.net/ch/solutions/pkg/router/handlers"
	m "git.containerum.net/ch/solutions/pkg/router/middleware"
	"git.containerum.net/ch/solutions/pkg/server"
	"git.containerum.net/ch/solutions/pkg/solerrors"
	"github.com/containerum/cherry/adaptors/cherrylog"
	"github.com/containerum/cherry/adaptors/gonic"
	"github.com/containerum/kube-client/pkg/model"

	"github.com/containerum/utils/httputil"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/gin-contrib/cors.v1"
)

//CreateRouter initialises router and middlewares
func CreateRouter(ss *server.SolutionsService, status *model.ServiceStatus, enableCORS bool) http.Handler {
	e := gin.New()
	e.GET("/status", httputil.ServiceStatus(status))
	initMiddlewares(e, ss)
	initRoutes(e, status, enableCORS)
	return e
}

func initMiddlewares(e *gin.Engine, ss *server.SolutionsService) {
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
func initRoutes(app *gin.Engine, status *model.ServiceStatus, enableCORS bool) {
	requireIdentityHeaders := httputil.RequireHeaders(solerrors.ErrRequiredHeadersNotProvided, httputil.UserIDXHeader, httputil.UserRoleXHeader)

	if enableCORS {
		cfg := cors.DefaultConfig()
		cfg.AllowAllOrigins = true
		cfg.AddAllowMethods(http.MethodDelete)
		cfg.AddAllowHeaders(httputil.UserIDXHeader, httputil.UserRoleXHeader, httputil.UserNamespacesXHeader)
		app.Use(cors.New(cfg))
	}
	app.Group("/static").
		StaticFS("/", static.HTTP)

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

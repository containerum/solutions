package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"time"

	"strings"

	m "git.containerum.net/ch/solutions/pkg/router/middleware"
	"git.containerum.net/ch/solutions/pkg/sErrors"
	"git.containerum.net/ch/solutions/pkg/server"
	"github.com/containerum/cherry"
	"github.com/containerum/cherry/adaptors/gonic"
	"github.com/containerum/utils/httputil"
	"github.com/sirupsen/logrus"
)

var lastCheckTime time.Time

const checkInterval = 6 * time.Hour

func UpdateSolutions(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	logrus.Infoln("Last solutions update check:", lastCheckTime.Format(time.RFC1123))
	if lastCheckTime.Add(checkInterval).Before(time.Now()) || (ctx.Query("forceupdate") == "true" && ctx.GetHeader(httputil.UserRoleXHeader) == "admin") {
		logrus.Infoln("Updating solutions")
		err := ss.UpdateAvailableSolutionsList(ctx.Request.Context())
		if err != nil {
			if cherr, ok := err.(*cherry.Err); ok {
				gonic.Gonic(cherr, ctx)
			} else {
				ctx.Error(err)
				gonic.Gonic(sErrors.ErrUnableUpdateSolutionsList(), ctx)
			}
			return
		}
		lastCheckTime = time.Now()
	} else {
		logrus.Infoln("No need to update logs")
	}
	ctx.Status(http.StatusAccepted)
}

func GetSolutionsList(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	resp, err := ss.GetAvailableSolutionsList(ctx.Request.Context())
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableGetSolutionsTemplatesList(), ctx)
		}
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func GetSolutionEnv(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	branch := branchMaster
	if ctx.Query("branch") != "" {
		branch = strings.TrimSpace(ctx.Query("branch"))
	}

	resp, err := ss.GetAvailableSolutionEnvList(ctx.Request.Context(), ctx.Param("solution"), branch)
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableGetSolutionTemplate(), ctx)
		}
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func GetSolutionResources(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	branch := branchMaster
	if ctx.Query("branch") != "" {
		branch = strings.TrimSpace(ctx.Query("branch"))
	}

	resp, err := ss.GetAvailableSolutionResourcesList(ctx.Request.Context(), ctx.Param("solution"), branch)
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableGetSolutionTemplate(), ctx)
		}
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

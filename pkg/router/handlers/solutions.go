package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"time"

	ch "git.containerum.net/ch/kube-client/pkg/cherry"
	"git.containerum.net/ch/kube-client/pkg/cherry/adaptors/gonic"
	cherry "git.containerum.net/ch/kube-client/pkg/cherry/solutions"
	m "git.containerum.net/ch/solutions/pkg/router/middleware"
	"git.containerum.net/ch/solutions/pkg/server"
	"github.com/sirupsen/logrus"
)

var lastCheckTime time.Time

const checkInterval = 6 * time.Hour

func UpdateSolutions(ctx *gin.Context) {
	ssp := ctx.MustGet(m.SolutionsServices).(*server.SolutionsService)
	ss := *ssp
	logrus.Infoln("Last solutions update check:", lastCheckTime.Format(time.RFC1123))
	if lastCheckTime.Add(checkInterval).Before(time.Now()) || (ctx.Query("forceupdate") == "true" && ctx.GetHeader(m.UserRoleHeader) == "admin") {
		logrus.Infoln("Updating solutions")
		err := ss.UpdateAvailableSolutionsList(ctx.Request.Context())
		if err != nil {
			if cherr, ok := err.(*ch.Err); ok {
				gonic.Gonic(cherr, ctx)
			} else {
				ctx.Error(err)
				gonic.Gonic(cherry.ErrUnableUpdateSolutionsList(), ctx)
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
	ssp := ctx.MustGet(m.SolutionsServices).(*server.SolutionsService)
	ss := *ssp
	resp, err := ss.GetAvailableSolutionsList(ctx.Request.Context())
	if err != nil {
		if cherr, ok := err.(*ch.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(cherry.ErrUnableGetSolutionsTemplatesList(), ctx)
		}
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func GetSolutionEnv(ctx *gin.Context) {
	ssp := ctx.MustGet(m.SolutionsServices).(*server.SolutionsService)
	ss := *ssp
	resp, err := ss.GetAvailableSolutionEnvList(ctx.Request.Context(), ctx.Param("solution"), ctx.Query("branch"))
	if err != nil {
		if cherr, ok := err.(*ch.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(cherry.ErrUnableGetSolutionTemplate(), ctx)
		}
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func GetSolutionResources(ctx *gin.Context) {
	ssp := ctx.MustGet(m.SolutionsServices).(*server.SolutionsService)
	ss := *ssp
	resp, err := ss.GetAvailableSolutionResourcesList(ctx.Request.Context(), ctx.Param("solution"), ctx.Query("branch"))
	if err != nil {
		if cherr, ok := err.(*ch.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(cherry.ErrUnableGetSolutionTemplate(), ctx)
		}
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

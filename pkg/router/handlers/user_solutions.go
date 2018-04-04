package handlers

import (
	"net/http"

	stypes "git.containerum.net/ch/json-types/solutions"
	"git.containerum.net/ch/kube-client/pkg/cherry/adaptors/gonic"
	cherry "git.containerum.net/ch/kube-client/pkg/cherry/solutions"
	m "git.containerum.net/ch/solutions/pkg/router/middleware"
	"git.containerum.net/ch/solutions/pkg/server"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

func RunSolution(ctx *gin.Context) {
	ssp := ctx.MustGet(m.SolutionsServices).(*server.SolutionsService)
	ss := *ssp
	logrus.Infoln("Last check time:", lastchecktime)

	var request stypes.UserSolution
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		ctx.Error(err)
		gonic.Gonic(cherry.ErrRequestValidationFailed().AddDetailsErr(err), ctx)
		return
	}

	err := ss.RunSolution(ctx.Request.Context(), request)
	if err != nil {
		ctx.Error(err)
		gonic.Gonic(cherry.ErrUnableCreateSolution(), ctx)
		return
	}

	ctx.Status(http.StatusAccepted)
}

func DeleteSolution(ctx *gin.Context) {
	ssp := ctx.MustGet(m.SolutionsServices).(*server.SolutionsService)
	ss := *ssp

	err := ss.DeleteSolution(ctx.Request.Context(), ctx.Param("solution"))
	if err != nil {
		ctx.Error(err)
		gonic.Gonic(cherry.ErrUnableDeleteSolution(), ctx)
		return
	}

	ctx.Status(http.StatusAccepted)
}

func GetUserSolutionsList(ctx *gin.Context) {
	ssp := ctx.MustGet(m.SolutionsServices).(*server.SolutionsService)
	ss := *ssp
	resp, err := ss.GetUserSolutionsList(ctx.Request.Context())
	if err != nil {
		ctx.Error(err)
		gonic.Gonic(cherry.ErrUnableGetSolutionsList(), ctx)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func GetUserSolutionsDeployments(ctx *gin.Context) {
	ssp := ctx.MustGet(m.SolutionsServices).(*server.SolutionsService)
	ss := *ssp
	resp, err := ss.GetUserSolutionDeployments(ctx.Request.Context(), ctx.Param("solution"))
	if err != nil {
		ctx.Error(err)
		gonic.Gonic(cherry.ErrUnableGetSolution(), ctx)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func GetUserSolutionsServices(ctx *gin.Context) {
	ssp := ctx.MustGet(m.SolutionsServices).(*server.SolutionsService)
	ss := *ssp
	resp, err := ss.GetUserSolutionServices(ctx.Request.Context(), ctx.Param("solution"))
	if err != nil {
		ctx.Error(err)
		gonic.Gonic(cherry.ErrUnableGetSolution(), ctx)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

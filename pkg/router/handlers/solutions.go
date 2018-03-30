package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	m "git.containerum.net/ch/solutions/pkg/router/middleware"
	"git.containerum.net/ch/solutions/pkg/server"
)

func UpdateSolutions(ctx *gin.Context) {
	ssp := ctx.MustGet(m.SolutionsServices).(*server.SolutionsService)
	ss := *ssp
	err := ss.UpdateAvailableSolutionsList(ctx.Request.Context())
	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusAccepted)
}

func GetSolutionsList(ctx *gin.Context) {
	ssp := ctx.MustGet(m.SolutionsServices).(*server.SolutionsService)
	ss := *ssp
	resp, err := ss.GetAvailableSolutionsList(ctx.Request.Context())
	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func GetSolutionEnv(ctx *gin.Context) {
	ssp := ctx.MustGet(m.SolutionsServices).(*server.SolutionsService)
	ss := *ssp
	resp, err := ss.GetAvailableSolutionEnv(ctx.Request.Context(), ctx.Param("solution"), ctx.Query("branch"))
	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func GetSolutionResources(ctx *gin.Context) {
	ssp := ctx.MustGet(m.SolutionsServices).(*server.SolutionsService)
	ss := *ssp
	resp, err := ss.GetAvailableSolutionResources(ctx.Request.Context(), ctx.Param("solution"), ctx.Query("branch"))
	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

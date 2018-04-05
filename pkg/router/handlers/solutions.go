package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"time"

	"fmt"

	"git.containerum.net/ch/kube-client/pkg/cherry/adaptors/gonic"
	cherry "git.containerum.net/ch/kube-client/pkg/cherry/solutions"
	m "git.containerum.net/ch/solutions/pkg/router/middleware"
	"git.containerum.net/ch/solutions/pkg/server"
)

var lastchecktime time.Time

const checkinterval = 6 * time.Hour

func UpdateSolutions(ctx *gin.Context) {
	ssp := ctx.MustGet(m.SolutionsServices).(*server.SolutionsService)
	ss := *ssp
	fmt.Println(lastchecktime)

	if lastchecktime.Add(checkinterval).Before(time.Now()) || ctx.Query("forceupdate") == "true" {
		fmt.Println("Updating solutions")
		err := ss.UpdateAvailableSolutionsList(ctx.Request.Context())
		if err != nil {
			ctx.Error(err)
			gonic.Gonic(cherry.ErrUnableUpdateSolutionsList(), ctx)
			return
		}
		lastchecktime = time.Now()
	} else {
		fmt.Println("Solutions list is still actual")
	}
	ctx.Status(http.StatusAccepted)
}

func GetSolutionsList(ctx *gin.Context) {
	ssp := ctx.MustGet(m.SolutionsServices).(*server.SolutionsService)
	ss := *ssp
	resp, err := ss.GetAvailableSolutionsList(ctx.Request.Context())
	if err != nil {
		ctx.Error(err)
		gonic.Gonic(cherry.ErrUnableGetSolutionsTemplatesList(), ctx)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func GetSolutionEnv(ctx *gin.Context) {
	ssp := ctx.MustGet(m.SolutionsServices).(*server.SolutionsService)
	ss := *ssp
	resp, err := ss.GetAvailableSolutionEnvList(ctx.Request.Context(), ctx.Param("solution"), ctx.Query("branch"))
	if err != nil {
		ctx.Error(err)
		gonic.Gonic(cherry.ErrUnableGetSolutionTemplate(), ctx)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func GetSolutionResources(ctx *gin.Context) {
	ssp := ctx.MustGet(m.SolutionsServices).(*server.SolutionsService)
	ss := *ssp
	resp, err := ss.GetAvailableSolutionResourcesList(ctx.Request.Context(), ctx.Param("solution"), ctx.Query("branch"))
	if err != nil {
		ctx.Error(err)
		gonic.Gonic(cherry.ErrUnableGetSolutionTemplate(), ctx)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

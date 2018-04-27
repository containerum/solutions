package handlers

import (
	"net/http"

	"fmt"

	"strings"

	stypes "git.containerum.net/ch/solutions/pkg/models"
	m "git.containerum.net/ch/solutions/pkg/router/middleware"
	"git.containerum.net/ch/solutions/pkg/sErrors"
	"git.containerum.net/ch/solutions/pkg/server"
	"github.com/containerum/cherry"
	"github.com/containerum/cherry/adaptors/gonic"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

const (
	branchMaster = "master"
)

func RunSolution(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	logrus.Infoln("Last check time:", lastCheckTime)

	var request stypes.UserSolution
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		gonic.Gonic(sErrors.ErrRequestValidationFailed().AddDetailsErr(err), ctx)
		return
	}

	valerrs := []error{}
	if request.Template == "" {
		valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "Template"))
	}
	if request.Name == "" {
		valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "Name"))
	}
	if request.Namespace == "" {
		valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "Namespace"))
	}
	if len(valerrs) > 0 {
		gonic.Gonic(sErrors.ErrRequestValidationFailed().AddDetailsErr(valerrs...), ctx)
		return
	}
	if request.Branch != "" {
		request.Branch = strings.TrimSpace(request.Branch)
	} else {
		request.Branch = branchMaster
	}

	solutionFile, solutionName, err := ss.DownloadSolutionConfig(ctx.Request.Context(), request)
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableCreateSolution(), ctx)
		}
		return
	}

	solutionConfig, solutionUUID, err := ss.ParseSolutionConfig(ctx.Request.Context(), solutionFile, request)
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableCreateSolution(), ctx)
		}
		return
	}

	ret, err := ss.CreateSolutionResources(ctx.Request.Context(), *solutionConfig, request, *solutionName, *solutionUUID)
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableCreateSolution(), ctx)
		}
		return
	}

	ctx.JSON(http.StatusAccepted, ret)
}

func DeleteSolution(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	err := ss.DeleteSolution(ctx.Request.Context(), ctx.Param("solution"))
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableDeleteSolution(), ctx)
		}
		return
	}

	ctx.Status(http.StatusAccepted)
}

func GetUserSolutionsList(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	resp, err := ss.GetUserSolutionsList(ctx.Request.Context())
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableGetSolutionsList(), ctx)
		}
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func GetUserSolutionsDeployments(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	resp, err := ss.GetUserSolutionDeployments(ctx.Request.Context(), ctx.Param("solution"))
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableGetSolution(), ctx)
		}
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func GetUserSolutionsServices(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	resp, err := ss.GetUserSolutionServices(ctx.Request.Context(), ctx.Param("solution"))
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableGetSolution(), ctx)
		}
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

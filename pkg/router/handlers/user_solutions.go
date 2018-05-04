package handlers

import (
	"net/http"

	"fmt"

	"strings"

	m "git.containerum.net/ch/solutions/pkg/router/middleware"
	"git.containerum.net/ch/solutions/pkg/sErrors"
	"git.containerum.net/ch/solutions/pkg/server"
	"github.com/containerum/cherry"
	"github.com/containerum/cherry/adaptors/gonic"
	stypes "github.com/containerum/kube-client/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const (
	branchMaster = "master"
)

// swagger:operation POST /user_solutions UserSolutions RunSolution
// Run solution.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserIDHeader'
//  - name: body
//    in: body
//    schema:
//      $ref: '#/definitions/UserSolution'
// responses:
//  '202':
//    description: solution created
//    schema:
//      $ref: '#/definitions/RunSolutionResponce'
//  default:
//    $ref: '#/responses/error'
func RunSolution(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)

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

// swagger:operation DELETE /user_solutions/{solution} UserSolutions DeleteSolution
// Delete solution.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserIDHeader'
//  - name: solution
//    in: path
//    type: string
//    required: true
// responses:
//  '202':
//    description: solution deleted
//  default:
//    $ref: '#/responses/error'
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

// swagger:operation GET /user_solutions UserSolutions GetUserSolutionsList
// Get running solutions list.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserIDHeader'
// responses:
//  '200':
//    description: running solutions list
//    schema:
//      $ref: '#/definitions/UserSolutionsList'
//  default:
//    $ref: '#/responses/error'
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

// swagger:operation GET /user_solutions/{solution}/deployments UserSolutions GetUserSolutionsDeployments
// Get solution deployments.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserIDHeader'
//  - name: solution
//    in: path
//    type: string
//    required: true
// responses:
//  '200':
//    description: solution deployments
//    schema:
//      $ref: '#/definitions/DeploymentsList'
//  default:
//    $ref: '#/responses/error'
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

// swagger:operation GET /user_solutions/{solution}/services UserSolutions GetUserSolutionsServices
// Get solution services.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserIDHeader'
//  - name: solution
//    in: path
//    type: string
//    required: true
// responses:
//  '200':
//    description: solutions services
//    schema:
//      $ref: '#/definitions/ServicesList'
//  default:
//    $ref: '#/responses/error'
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

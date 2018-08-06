package handlers

import (
	"net/http"

	m "git.containerum.net/ch/solutions/pkg/router/middleware"
	"git.containerum.net/ch/solutions/pkg/sErrors"
	"git.containerum.net/ch/solutions/pkg/server"
	"git.containerum.net/ch/solutions/pkg/validation"
	"github.com/containerum/cherry"
	"github.com/containerum/cherry/adaptors/gonic"
	kube_types "github.com/containerum/kube-client/pkg/model"
	"github.com/containerum/utils/httputil"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const (
	branchMaster = "master"
)

// swagger:operation GET /solutions Solutions GetSolutionsList
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
//      $ref: '#/definitions/SolutionsList'
//  default:
//    $ref: '#/responses/error'
func GetSolutionsList(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	resp, err := ss.GetSolutionsList(ctx.Request.Context(), ctx.GetHeader(httputil.UserRoleXHeader) == "admin")
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

// swagger:operation GET /namespaces/{namespace}/solutions Solutions GetNamespaceSolutions
// Get running namespace solutions list.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserIDHeader'
//  - name: namespace
//    in: path
//    type: string
//    required: true
// responses:
//  '200':
//    description: running namespace solutions list
//    schema:
//      $ref: '#/definitions/SolutionsList'
//  default:
//    $ref: '#/responses/error'
func GetNamespaceSolutions(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	resp, err := ss.GetNamespaceSolutionsList(ctx.Request.Context(), ctx.Param("namespace"), ctx.GetHeader(httputil.UserRoleXHeader) == "admin")
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

// swagger:operation GET /namespaces/{namespace}/solutions/{solution} Solutions GetSolution
// Get running solution.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserIDHeader'
//  - name: namespace
//    in: path
//    type: string
//    required: true
//  - name: solution
//    in: path
//    type: string
//    required: true
// responses:
//  '200':
//    description: running solution
//    schema:
//      $ref: '#/definitions/SolutionsList'
//  default:
//    $ref: '#/responses/error'
func GetSolution(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	resp, err := ss.GetSolution(ctx.Request.Context(), ctx.Param("namespace"), ctx.Param("solution"), ctx.GetHeader(httputil.UserRoleXHeader) == "admin")
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

// swagger:operation GET /namespaces/{namespace}/solutions/{solution}/deployments Solutions GetSolutionsDeployments
// Get solution deployments.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserIDHeader'
//  - name: namespace
//    in: path
//    type: string
//    required: true
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
func GetSolutionsDeployments(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	resp, err := ss.GetSolutionDeployments(ctx.Request.Context(), ctx.Param("namespace"), ctx.Param("solution"))
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

// swagger:operation GET /namespaces/{namespace}/solutions/{solution}/services Solutions GetSolutionsServices
// Get solution services.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserIDHeader'
//  - name: namespace
//    in: path
//    type: string
//    required: true
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
func GetSolutionsServices(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	resp, err := ss.GetSolutionServices(ctx.Request.Context(), ctx.Param("namespace"), ctx.Param("solution"))
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

// swagger:operation POST /namespaces/{namespace}/solutions Solutions RunSolution
// Run solution.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserIDHeader'
//  - name: namespace
//    in: path
//    type: string
//    required: true
//  - name: body
//    in: body
//    schema:
//      $ref: '#/definitions/Solution'
// responses:
//  '202':
//    description: solution created
//    schema:
//      $ref: '#/definitions/RunSolutionResponse'
//  default:
//    $ref: '#/responses/error'
func RunSolution(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)

	var request kube_types.Solution
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		gonic.Gonic(sErrors.ErrRequestValidationFailed().AddDetailsErr(err), ctx)
		return
	}

	request.Namespace = ctx.Param("namespace")

	if err := validation.ValidateSolution(request); err != nil {
		gonic.Gonic(err, ctx)
		return
	}

	if request.Branch == "" {
		request.Branch = branchMaster
	}

	ret, err := ss.RunSolution(ctx.Request.Context(), request)
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

// swagger:operation DELETE /namespaces/{namespace}/solutions/{solution} Solutions DeleteSolution
// Delete solution.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserIDHeader'
//  - name: namespace
//    in: path
//    type: string
//    required: true
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
	if err := ss.DeleteSolution(ctx.Request.Context(), ctx.Param("namespace"), ctx.Param("solution")); err != nil {
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

// swagger:operation DELETE /solutions Solutions DeleteSolutions
// Delete user solution.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserIDHeader'
// responses:
//  '202':
//    description: user solutions deleted
//  default:
//    $ref: '#/responses/error'
func DeleteSolutions(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	if err := ss.DeleteSolutions(ctx.Request.Context()); err != nil {
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

// swagger:operation DELETE /namespaces/{namespace}/solutions Solutions DeleteNamespaceSolutions
// Delete all namespace solutions.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserIDHeader'
//  - name: namespace
//    in: path
//    type: string
//    required: true
// responses:
//  '202':
//    description: solutions deleted
//  default:
//    $ref: '#/responses/error'
func DeleteNamespaceSolutions(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	if err := ss.DeleteNamespaceSolutions(ctx.Request.Context(), ctx.Param("namespace")); err != nil {
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

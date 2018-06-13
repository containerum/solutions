package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	m "git.containerum.net/ch/solutions/pkg/router/middleware"
	"git.containerum.net/ch/solutions/pkg/sErrors"
	"git.containerum.net/ch/solutions/pkg/server"
	"git.containerum.net/ch/solutions/pkg/validation"
	"github.com/containerum/cherry"
	"github.com/containerum/cherry/adaptors/gonic"
	kube_types "github.com/containerum/kube-client/pkg/model"
	"github.com/containerum/utils/httputil"
	"github.com/gin-gonic/gin/binding"
)

// swagger:operation GET /templates Templates GetTemplatesList
// Get solutions templates list.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserIDHeader'
// responses:
//  '200':
//    description: available solutions
//    schema:
//      $ref: '#/definitions/AvailableSolutionsList'
//  default:
//    $ref: '#/responses/error'
func GetTemplatesList(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	resp, err := ss.GetTemplatesList(ctx.Request.Context(), ctx.GetHeader(httputil.UserRoleXHeader) == "admin")
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableGetTemplatesList(), ctx)
		}
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// swagger:operation GET /templates/{solution}/env Templates GetTemplatesEnv
// Get solution templates environment variables.
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
//    description: available solution envs
//    schema:
//      $ref: '#/definitions/SolutionEnv'
//  default:
//    $ref: '#/responses/error'
func GetTemplatesEnv(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	branch := branchMaster
	if ctx.Query("branch") != "" {
		branch = ctx.Query("branch")
	}

	resp, err := ss.GetTemplatesEnvList(ctx.Request.Context(), ctx.Param("solution"), branch)
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableGetTemplate(), ctx)
		}
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// swagger:operation GET /templates/{solution}/resources Templates GetTemplatesResources
// Get solution templates resources.
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
//    description: available solution resources
//    schema:
//      $ref: '#/definitions/SolutionResources'
//  default:
//    $ref: '#/responses/error'
func GetTemplatesResources(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	branch := branchMaster
	if ctx.Query("branch") != "" {
		branch = ctx.Query("branch")
	}

	resp, err := ss.GetTemplatesResourcesList(ctx.Request.Context(), ctx.Param("solution"), branch)
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableGetTemplate(), ctx)
		}
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// swagger:operation POST /templates Templates AddTemplate
// Add template.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserIDHeader'
//  - name: body
//    in: body
//    schema:
//      $ref: '#/definitions/AvailableSolution'
// responses:
//  '201':
//    description: solution added
//  default:
//    $ref: '#/responses/error'
func AddTemplate(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	var request kube_types.AvailableSolution
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		gonic.Gonic(sErrors.ErrRequestValidationFailed().AddDetailsErr(err), ctx)
		return
	}

	if err := validation.ValidateTemplate(request); err != nil {
		gonic.Gonic(err, ctx)
		return
	}

	if err := ss.AddTemplate(ctx.Request.Context(), request); err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableAddTemplate(), ctx)
		}
		return
	}
	ctx.Status(http.StatusCreated)
}

// swagger:operation PUT /templates/{solution} Templates UpdateTemplate
// Update template.
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
//  - name: body
//    in: body
//    schema:
//      $ref: '#/definitions/AvailableSolution'
// responses:
//  '202':
//    description: solution updated
//  default:
//    $ref: '#/responses/error'
func UpdateTemplate(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)

	var request kube_types.AvailableSolution
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		gonic.Gonic(sErrors.ErrRequestValidationFailed().AddDetailsErr(err), ctx)
		return
	}

	request.Name = ctx.Param("solution")

	if err := validation.ValidateTemplate(request); err != nil {
		gonic.Gonic(err, ctx)
		return
	}

	if err := ss.UpdateTemplate(ctx.Request.Context(), request); err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableUpdateTemplate(), ctx)
		}
		return
	}
	ctx.Status(http.StatusAccepted)
}

// swagger:operation POST /templates/{solution}/activate Templates ActivateTemplate
// Activate template.
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
//    description: solution activated
//  default:
//    $ref: '#/responses/error'
func ActivateTemplate(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)

	if err := ss.ActivateTemplate(ctx.Request.Context(), ctx.Param("solution")); err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableActivateTemplate(), ctx)
		}
		return
	}

	ctx.Status(http.StatusAccepted)
}

// swagger:operation POST /templates/{solution}/deactivate Templates DeactivateTemplate
// Deactivate template.
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
//    description: solution deactivated
//  default:
//    $ref: '#/responses/error'
func DeactivateTemplate(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)

	if err := ss.DeactivateTemplate(ctx.Request.Context(), ctx.Param("solution")); err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableDeactivateTemplate(), ctx)
		}
		return
	}

	ctx.Status(http.StatusAccepted)
}

// swagger:operation DELETE /templates/{solution} Templates DeleteTemplate
// Delete template.
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
func DeleteTemplate(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	if err := ss.DeleteTemplate(ctx.Request.Context(), ctx.Param("solution")); err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableDeleteTemplate(), ctx)
		}
		return
	}

	ctx.Status(http.StatusAccepted)
}

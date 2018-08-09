package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	m "git.containerum.net/ch/solutions/pkg/router/middleware"
	"git.containerum.net/ch/solutions/pkg/server"
	"git.containerum.net/ch/solutions/pkg/solerrors"
	"git.containerum.net/ch/solutions/pkg/validation"
	"github.com/containerum/cherry"
	"github.com/containerum/cherry/adaptors/gonic"
	kubeTypes "github.com/containerum/kube-client/pkg/model"
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
//      $ref: '#/definitions/SolutionsTemplatesList'
//  default:
//    $ref: '#/responses/error'
func GetTemplatesList(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	resp, err := ss.GetTemplatesList(ctx.Request.Context(), ctx.GetHeader(httputil.UserRoleXHeader) == m.RoleAdmin)
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(solerrors.ErrUnableGetTemplatesList(), ctx)
		}
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// swagger:operation GET /templates/{template}/env Templates GetTemplatesEnv
// Get solution templates environment variables.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserIDHeader'
//  - name: template
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

	resp, err := ss.GetTemplatesEnvList(ctx.Request.Context(), ctx.Param("template"), branch)
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(solerrors.ErrUnableGetTemplate(), ctx)
		}
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// swagger:operation GET /templates/{template}/resources Templates GetTemplatesResources
// Get solution templates resources.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserIDHeader'
//  - name: template
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

	resp, err := ss.GetTemplatesResourcesList(ctx.Request.Context(), ctx.Param("template"), branch)
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(solerrors.ErrUnableGetTemplate(), ctx)
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
//      $ref: '#/definitions/SolutionTemplate'
// responses:
//  '201':
//    description: solution added
//  default:
//    $ref: '#/responses/error'
func AddTemplate(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	var request kubeTypes.SolutionTemplate
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		gonic.Gonic(solerrors.ErrRequestValidationFailed().AddDetailsErr(err), ctx)
		return
	}

	if err := validation.ValidateTemplate(request); err != nil {
		gonic.Gonic(err, ctx)
		return
	}

	if err := ss.ValidateTemplate(ctx.Request.Context(), request); err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(solerrors.ErrTemplateValidationFailed().AddDetailsErr(err), ctx)
		}
		return
	}

	if err := ss.AddTemplate(ctx.Request.Context(), request); err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(solerrors.ErrUnableAddTemplate(), ctx)
		}
		return
	}
	ctx.Status(http.StatusCreated)
}

// swagger:operation PUT /templates/{template} Templates UpdateTemplate
// Update template.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserIDHeader'
//  - name: template
//    in: path
//    type: string
//    required: true
//  - name: body
//    in: body
//    schema:
//      $ref: '#/definitions/SolutionTemplate'
// responses:
//  '202':
//    description: solution updated
//  default:
//    $ref: '#/responses/error'
func UpdateTemplate(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)

	var request kubeTypes.SolutionTemplate
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		gonic.Gonic(solerrors.ErrRequestValidationFailed().AddDetailsErr(err), ctx)
		return
	}

	request.Name = ctx.Param("template")

	if err := validation.ValidateTemplate(request); err != nil {
		gonic.Gonic(err, ctx)
		return
	}

	if err := ss.ValidateTemplate(ctx.Request.Context(), request); err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(solerrors.ErrTemplateValidationFailed().AddDetailsErr(err), ctx)
		}
		return
	}

	if err := ss.UpdateTemplate(ctx.Request.Context(), request); err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(solerrors.ErrUnableUpdateTemplate(), ctx)
		}
		return
	}
	ctx.Status(http.StatusAccepted)
}

// swagger:operation POST /templates/{template}/activate Templates ActivateTemplate
// Activate template.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserIDHeader'
//  - name: template
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

	if err := ss.ActivateTemplate(ctx.Request.Context(), ctx.Param("template")); err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(solerrors.ErrUnableActivateTemplate(), ctx)
		}
		return
	}

	ctx.Status(http.StatusAccepted)
}

// swagger:operation POST /templates/{template}/deactivate Templates DeactivateTemplate
// Deactivate template.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserIDHeader'
//  - name: template
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

	if err := ss.DeactivateTemplate(ctx.Request.Context(), ctx.Param("template")); err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(solerrors.ErrUnableDeactivateTemplate(), ctx)
		}
		return
	}

	ctx.Status(http.StatusAccepted)
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"strings"

	"fmt"

	m "git.containerum.net/ch/solutions/pkg/router/middleware"
	"git.containerum.net/ch/solutions/pkg/sErrors"
	"git.containerum.net/ch/solutions/pkg/server"
	"github.com/containerum/cherry"
	"github.com/containerum/cherry/adaptors/gonic"
	stypes "github.com/containerum/kube-client/pkg/model"
	"github.com/containerum/utils/httputil"
	"github.com/gin-gonic/gin/binding"
)

/*
var lastCheckTime time.Time

const checkInterval = 6 * time.Hour

func UpdateSolutions(ctx *gin.Context) {
	ss := ctx.MustGet(m.SolutionsServices).(server.SolutionsService)
	logrus.Infoln("Last solutions update check:", lastCheckTime.Format(time.RFC1123))
	if lastCheckTime.Add(checkInterval).Before(time.Now()) || (ctx.Query("forceupdate") == "true" && ctx.GetHeader(httputil.UserRoleXHeader) == "admin") {
		logrus.Infoln("Updating solutions")
		err := ss.UpdateTemplatesList(ctx.Request.Context())
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
}*/

// swagger:operation GET /solutions AvailableSolutions GetSolutionsList
// Get available solutions list.
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
			gonic.Gonic(sErrors.ErrUnableGetSolutionsTemplatesList(), ctx)
		}
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// swagger:operation GET /solutions/{solution}/env AvailableSolutions GetTemplatesEnv
// Get available solution environment variables.
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
		branch = strings.TrimSpace(ctx.Query("branch"))
	}

	resp, err := ss.GetTemplatesEnvList(ctx.Request.Context(), ctx.Param("solution"), branch)
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

// swagger:operation GET /solutions/{solution}/resources AvailableSolutions GetTemplatesResources
// Get available solution resources.
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
		branch = strings.TrimSpace(ctx.Query("branch"))
	}

	resp, err := ss.GetTemplatesResourcesList(ctx.Request.Context(), ctx.Param("solution"), branch)
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

// swagger:operation POST /solutions AvailableSolutions AddTemplate
// Add available solution.
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

	var request stypes.AvailableSolution
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		gonic.Gonic(sErrors.ErrRequestValidationFailed().AddDetailsErr(err), ctx)
		return
	}

	valerrs := []error{}
	if request.Name == "" {
		valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "Name"))
	}
	if request.Limits == nil {
		valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "Limits"))
	} else {
		if request.Limits.RAM == "" {
			valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "RAM"))
		}
		if request.Limits.CPU == "" {
			valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "CPU"))
		}
	}
	if len(request.Images) == 0 {
		valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "Images"))
	}
	if len(request.URL) == 0 {
		valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "URL"))
	}
	if len(valerrs) > 0 {
		gonic.Gonic(sErrors.ErrRequestValidationFailed().AddDetailsErr(valerrs...), ctx)
		return
	}

	err := ss.AddTemplate(ctx.Request.Context(), request)
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableGetSolutionsTemplatesList(), ctx)
		}
		return
	}

	ctx.Status(http.StatusCreated)
}

// swagger:operation PUT /solutions/{solution} AvailableSolutions UpdateTemplate
// Update available solution.
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

	var request stypes.AvailableSolution
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		gonic.Gonic(sErrors.ErrRequestValidationFailed().AddDetailsErr(err), ctx)
		return
	}

	request.Name = ctx.Param("solution")

	valerrs := []error{}
	if request.Limits == nil {
		valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "Limits"))
	} else {
		if request.Limits.RAM == "" {
			valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "RAM"))
		}
		if request.Limits.CPU == "" {
			valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "CPU"))
		}
	}
	if len(request.Images) == 0 {
		valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "Images"))
	}
	if len(request.URL) == 0 {
		valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "URL"))
	}
	if len(valerrs) > 0 {
		gonic.Gonic(sErrors.ErrRequestValidationFailed().AddDetailsErr(valerrs...), ctx)
		return
	}

	err := ss.UpdateTemplate(ctx.Request.Context(), request)
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableUpdateSolutionsList(), ctx)
		}
		return
	}

	ctx.Status(http.StatusAccepted)
}

// swagger:operation POST /solutions/{solution}/activate AvailableSolutions ActivateTemplate
// Activate available solution.
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

	err := ss.ActivateTemplate(ctx.Request.Context(), ctx.Param("solution"))
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableUpdateSolutionsList(), ctx)
		}
		return
	}

	ctx.Status(http.StatusAccepted)
}

// swagger:operation POST /solutions/{solution}/deactivate AvailableSolutions DeactivateTemplate
// Deactivate available solution.
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

	err := ss.DeactivateTemplate(ctx.Request.Context(), ctx.Param("solution"))
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableUpdateSolutionsList(), ctx)
		}
		return
	}

	ctx.Status(http.StatusAccepted)
}

// swagger:operation DELETE /solutions/{solution} AvailableSolutions DeleteTemplate
// Delete available solution.
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
	err := ss.DeleteTemplate(ctx.Request.Context(), ctx.Param("solution"))
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(sErrors.ErrUnableUpdateSolutionsList(), ctx)
		}
		return
	}

	ctx.Status(http.StatusAccepted)
}

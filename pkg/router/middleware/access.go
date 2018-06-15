package middleware

import (
	"git.containerum.net/ch/solutions/pkg/model"
	"git.containerum.net/ch/solutions/pkg/sErrors"
	"github.com/containerum/cherry/adaptors/gonic"
	kubeModel "github.com/containerum/kube-client/pkg/model"
	headers "github.com/containerum/utils/httputil"
	"github.com/gin-gonic/gin"
)

var (
	readLevels = []kubeModel.AccessLevel{
		kubeModel.Owner,
		kubeModel.Write,
		kubeModel.ReadDelete,
		kubeModel.Read,
	}

	deleteLevels = []kubeModel.AccessLevel{
		kubeModel.Owner,
		kubeModel.Write,
		kubeModel.ReadDelete,
	}

	writeLevels = []kubeModel.AccessLevel{
		kubeModel.Owner,
		kubeModel.Write,
	}
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

func IsAdmin(ctx *gin.Context) {
	if role := GetHeader(ctx, headers.UserRoleXHeader); role != RoleAdmin {
		gonic.Gonic(sErrors.ErrAdminRequired(), ctx)
		return
	}
	return
}

func ReadAccess(ctx *gin.Context) {
	CheckAccess(ctx, readLevels)
}

func DeleteAccess(ctx *gin.Context) {
	CheckAccess(ctx, deleteLevels)
}

func WriteAccess(ctx *gin.Context) {
	CheckAccess(ctx, writeLevels)
}

func CheckAccess(ctx *gin.Context, level []kubeModel.AccessLevel) {
	ns := ctx.Param("namespace")
	if GetHeader(ctx, headers.UserRoleXHeader) == RoleUser {
		var userNsData *kubeModel.UserHeaderData
		nsList := ctx.MustGet(UserNamespaces).(*model.UserHeaderDataMap)
		for _, n := range *nsList {
			if ns == n.ID {
				userNsData = &n
				break
			}
		}
		if userNsData != nil {
			if ok := containsAccess(userNsData.Access, level...); ok {
				return
			}
			gonic.Gonic(sErrors.ErrAccessError(), ctx)
			return
		}
		gonic.Gonic(sErrors.ErrSolutionNotExist().AddDetails("project is not found"), ctx)
		return
	}
	return
}

func containsAccess(access kubeModel.AccessLevel, in ...kubeModel.AccessLevel) bool {
	contains := false
	userAccess := kubeModel.AccessLevel(access)
	for _, acc := range in {
		if acc == userAccess {
			return true
		}
	}
	return contains
}

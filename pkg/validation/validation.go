package validation

import (
	"fmt"

	"git.containerum.net/ch/solutions/pkg/solerrors"
	"github.com/containerum/cherry"
	kube_types "github.com/containerum/kube-client/pkg/model"
)

const (
	fieldShouldExist = "field %v should be provided"
)

func ValidateTemplate(template kube_types.SolutionTemplate) *cherry.Err {
	valerrs := []error{}
	if template.Name == "" {
		valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "Name"))
	}
	if len(template.URL) == 0 {
		valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "URL"))
	}
	if len(valerrs) > 0 {
		return solerrors.ErrRequestValidationFailed().AddDetailsErr(valerrs...)
	}
	return nil
}

func ValidateSolution(solution kube_types.Solution) *cherry.Err {
	valerrs := []error{}
	if solution.Template == "" {
		valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "Template"))
	}
	if solution.Name == "" {
		valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "Name"))
	}
	if solution.Namespace == "" {
		valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "Namespace"))
	}
	if len(valerrs) > 0 {
		return solerrors.ErrRequestValidationFailed().AddDetailsErr(valerrs...)
	}
	return nil
}

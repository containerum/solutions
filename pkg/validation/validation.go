package validation

import (
	"fmt"

	"git.containerum.net/ch/solutions/pkg/sErrors"
	"github.com/containerum/cherry"
	stypes "github.com/containerum/kube-client/pkg/model"
)

const (
	fieldShouldExist = "field %v should be provided"
)

func ValidateTemplate(template stypes.AvailableSolution) *cherry.Err {
	valerrs := []error{}
	if template.Name == "" {
		valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "Name"))
	}
	if template.Limits == nil {
		valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "Limits"))
	} else {
		if template.Limits.RAM == "" {
			valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "RAM"))
		}
		if template.Limits.CPU == "" {
			valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "CPU"))
		}
	}
	if len(template.Images) == 0 {
		valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "Images"))
	}
	if len(template.URL) == 0 {
		valerrs = append(valerrs, fmt.Errorf(fieldShouldExist, "URL"))
	}
	if len(valerrs) > 0 {
		return sErrors.ErrRequestValidationFailed().AddDetailsErr(valerrs...)
	}
	return nil
}

func ValidateSolution(solution stypes.UserSolution) *cherry.Err {
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
		return sErrors.ErrRequestValidationFailed().AddDetailsErr(valerrs...)
	}
	return nil
}
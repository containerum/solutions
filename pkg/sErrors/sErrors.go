// Code generated by noice. DO NOT EDIT.
package sErrors

import (
	bytes "bytes"
	cherry "github.com/containerum/cherry"
	template "text/template"
)

const ()

func ErrAdminRequired(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Admin access required", StatusHTTP: 403, ID: cherry.ErrID{SID: "Solutions", Kind: 0x1}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrRequiredHeadersNotProvided(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Required headers not provided", StatusHTTP: 400, ID: cherry.ErrID{SID: "Solutions", Kind: 0x2}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrRequestValidationFailed(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Request validation failed", StatusHTTP: 400, ID: cherry.ErrID{SID: "Solutions", Kind: 0x3}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableUpdateTemplate(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to update template", StatusHTTP: 500, ID: cherry.ErrID{SID: "Solutions", Kind: 0x4}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableGetTemplatesList(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to get solutions templates list", StatusHTTP: 500, ID: cherry.ErrID{SID: "Solutions", Kind: 0x5}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableGetTemplate(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to get solutions template", StatusHTTP: 500, ID: cherry.ErrID{SID: "Solutions", Kind: 0x6}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableGetSolutionsList(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to get user solutions list", StatusHTTP: 500, ID: cherry.ErrID{SID: "Solutions", Kind: 0x7}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableGetSolution(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to get user solution", StatusHTTP: 500, ID: cherry.ErrID{SID: "Solutions", Kind: 0x8}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableCreateSolution(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to create solution", StatusHTTP: 500, ID: cherry.ErrID{SID: "Solutions", Kind: 0x9}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableDeleteSolution(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to delete solution", StatusHTTP: 500, ID: cherry.ErrID{SID: "Solutions", Kind: 0xa}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrSolutionNotExist(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Solution with this name doesn't exist", StatusHTTP: 404, ID: cherry.ErrID{SID: "Solutions", Kind: 0xc}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrInternalError(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Internal error", StatusHTTP: 500, ID: cherry.ErrID{SID: "Solutions", Kind: 0xd}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableAddTemplate(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to add solution template", StatusHTTP: 500, ID: cherry.ErrID{SID: "Solutions", Kind: 0xe}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrResourceAlreadyExists(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Resource with this name already exists", StatusHTTP: 409, ID: cherry.ErrID{SID: "Solutions", Kind: 0xf}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrTemplateNotExist(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Template with this name doesn't exist", StatusHTTP: 404, ID: cherry.ErrID{SID: "Solutions", Kind: 0x10}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableActivateTemplate(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to activate template", StatusHTTP: 500, ID: cherry.ErrID{SID: "Solutions", Kind: 0x11}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableDeactivateTemplate(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to activate template", StatusHTTP: 500, ID: cherry.ErrID{SID: "Solutions", Kind: 0x12}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrUnableDeleteTemplate(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Unable to activate template", StatusHTTP: 500, ID: cherry.ErrID{SID: "Solutions", Kind: 0x13}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrAccessError(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Access denied", StatusHTTP: 403, ID: cherry.ErrID{SID: "Solutions", Kind: 0x14}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrInvalidRole(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Invalid user role", StatusHTTP: 403, ID: cherry.ErrID{SID: "Solutions", Kind: 0x15}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrTemplateValidationFailed(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Template validation failed", StatusHTTP: 400, ID: cherry.ErrID{SID: "Solutions", Kind: 0x16}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}
func renderTemplate(templText string) string {
	buf := &bytes.Buffer{}
	templ, err := template.New("").Parse(templText)
	if err != nil {
		return err.Error()
	}
	err = templ.Execute(buf, map[string]string{})
	if err != nil {
		return err.Error()
	}
	return buf.String()
}

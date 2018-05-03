package models

// UserSolutionsList -- list of running solution
//
// swagger:model
type UserSolutionsList struct {
	Solutions []UserSolution `json:"solutions"`
}

// UserSolution -- running solution
//
// swagger:model
type UserSolution struct {
	Branch string            `json:"branch"`
	Env    map[string]string `json:"env"`
	// required: true
	Template string `json:"template"`
	// required: true
	Name string `json:"name"`
	// required: true
	Namespace string `json:"namespace"`
}

// DeploymentsList -- list of deployments
//
// swagger:model
type DeploymentsList struct {
	Deployments []*interface{} `json:"deployments"`
}

// ServicesList -- list of services
//
// swagger:model
type ServicesList struct {
	Services []*interface{} `json:"services"`
}

// RunSolutionResponce -- responce to run solution request
//
// swagger:model
type RunSolutionResponce struct {
	Created    int      `json:"created"`
	NotCreated int      `json:"not_created"`
	Errors     []string `json:"errors,omitempty"`
}

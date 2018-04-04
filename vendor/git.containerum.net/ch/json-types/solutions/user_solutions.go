package solutions

type UserSolutionsList struct {
	Solutions []UserSolution `json:"solutions"`
}

type UserSolution struct {
	Branch    string            `json:"branch"`
	Env       map[string]string `json:"env"`
	Template  string            `json:"template"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
}

type DeploymentsList struct {
	Deployments []*interface{} `json:"deployments"`
}

type ServicesList struct {
	Services []*interface{} `json:"services"`
}

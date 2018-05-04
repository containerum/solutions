package models

// AvailableSolutionsList -- list of available solutions
//
// swagger:model
type AvailableSolutionsList struct {
	Solutions []AvailableSolution `json:"solutions"`
}

// AvailableSolution -- solution which user can run
//
// swagger:model
type AvailableSolution struct {
	Name   string   `json:"name"`
	Limits *Limits  `json:"limits"`
	Images []string `json:"images"`
	URL    string   `json:"url"`
	Active bool
}

// Limits -- solution resources limits
//
// swagger:model
type Limits struct {
	CPU string `json:"cpu"`
	RAM string `json:"ram"`
}

// SolutionEnv -- solution environment variables
//
// swagger:model
type SolutionEnv struct {
	Env map[string]string `json:"env"`
}

// SolutionResources -- list of solution resources
//
// swagger:model
type SolutionResources struct {
	Resources map[string]int `json:"resources"`
}

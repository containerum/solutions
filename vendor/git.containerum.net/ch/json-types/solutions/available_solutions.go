package solutions

type AvailableSolutionsList struct {
	Solutions []AvailableSolution `json:"solutions"`
}

type AvailableSolution struct {
	Name   string   `json:"name"`
	Limits *Limits  `json:"limits"`
	Images []string `json:"images"`
	URL    string   `json:"url"`
}

type Limits struct {
	CPU string `json:"cpu"`
	RAM string `json:"ram"`
}

type SolutionEnv struct {
	Env []string `json:"env"`
}

type SolutionResources struct {
	Resources map[string]int `json:"resources"`
}

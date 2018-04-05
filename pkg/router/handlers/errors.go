package handlers

const (
	noContainer         = "Container %v is not found in deployment"
	fieldShouldExist    = "Field %v should be provided"
	invalidReplicas     = "Invalid replicas number: %v. It must be between 1 and %v"
	invalidPort         = "Invalid port: %v. It must be between %v and %v"
	invalidProtocol     = "Invalid protocol: %v. It must be TCP or UDP"
	invalidOwner        = "Owner should be UUID"
	invalidName         = "Invalid name: %v. %v"
	invalidIP           = "Invalid IP: %v. It must be a valid IP address, (e.g. 10.9.8.7)"
	invalidCPUQuota     = "Invalid CPU quota: %v. It must be between %v and %v"
	invalidMemoryQuota  = "Invalid memory quota: %v. It must be between %v and %v"
	subPathRelative     = "Invalid Sub Path: %v. It must be relative path"
	invalidResourceKind = "Invalid resource kind: %v. Shoud be %v"
	invalidApiVersion   = "Invalid API Version: %v. Shoud be %v"
)

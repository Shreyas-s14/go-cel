package yamlbuilder


// generated as per the number of ValidationInfo structs provided in the input array(?)
// NOT USED
type VAPBinding struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind string `yaml:"kind"`
	Metadata Meta `yaml:"metadata"`
	Spec VAPBindingSpec `yaml:"spec"`
}

// skipped: paramRef
type VAPBindingSpec struct {
	PolicyName string `yaml:"policyName"`
	ValidationActions []string `yaml:"validationActions"`
	MatchResources MatchResourcesInfo `yaml:"matchResources"`
}

// skipped: namespaceSelector, objectSelector
type MatchResourcesInfo struct {
	MatchPolicy string `yaml:"matchPolicy,omitempty"`  // allowed :("Exact", "Equivalent"), default: "Equivalent"
	ResourceRules []Rules `yaml:"resourceRules"`
	ExcludeResourceRules []Rules `yaml:"excludeResourceRules,omitempty"`
}
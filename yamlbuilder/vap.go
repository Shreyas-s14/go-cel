package yamlbuilder
import (
	"k8s.io/api/admissionregistration/v1"
)

// NOT USED

// add field Reason (string)?? -> "Reason represents a machine-readable description of why this validation failed" - manifests.io/kubernetes/1.30/io.k8s.api.admissionregistration.v1.Validation?linked=ValidatingAdmissionPolicy.spec.validations description

type ValidationInfo struct {
	Expression string   // Must evaluate to bool
	Message string `yaml:"message,omitempty"`
	MessageExpression string `yaml:"messageExpression,omitempty"`// not used in current VAP, can be added in the future for custom messages. Validation-> must evaluate to string
}


// ValidatingAdmissionPolicy struct . fields as per expected in the yaml
type VAP struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind string `yaml:"kind"`
	Metadata Meta `yaml:"metadata"`
	Spec VAPSpec `yaml:"spec"`
	test v1.ValidatingAdmissionPolicy
}

// add other fields (annotations, labels, etc)
type Meta struct {
	Name string `yaml:"name"`
	Namespace string `yaml:"namespace,omitempty"`
}

type VAPSpec struct {
	FailurePolicy string `yaml:"failurePolicy"`
	MatchConstraints MatchConstraintsInfo `yaml:"matchConstraints"`
	Validations []ValidationInfo `yaml:"validations"`
}


type MatchConstraintsInfo struct {
	ResourceRules []Rules `yaml:"resourceRules"`
	
}

type Rules struct {
	Operations []v1.OperationType `yaml:"operations"`
	APIGroups []string `yaml:"apiGroups"`
	APIVersions []string `yaml:"apiVersions"`
	Resources []string `yaml:"resources"`
}
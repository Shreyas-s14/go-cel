package expression

/*
	MatchConditions -> further filtering. filtering for requests. eg: check annotations and request.user etc.
	MatchCOnstraints -> apply the valdiation rules on certain resources in certain api groups only.
	1) structs for: a) variableinfo(if in case variables are used.)  b) celValidationinfo(mainly used)  c) rest of the non singular fields.

*/

import (
	"errors"
	"fmt"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/api/admissionregistration/v1"
)

type CelVariableInfo struct {
	Name string
	Expression string
}

type CelValidationInfo struct {
	Expression string
	Message string
	MessageExpression string // not used in the current VAP.
}
// in case action in policybinding is set as audit instead of warn/deny
type CelAuditAnnotationsInfo struct {
	Key        string
	Expression string
}
type CelMatchConditionsInfo struct {
	Name string
	Expression string
}
// fields: a) excluderesiurcerules b) matchpolicy  c) namespaceSelector d) objectSelector e) resourcerules.
// Only taking : resourceRules namespaceSelector objectselector.      matchpolicy default -> equivalent. // TODO: check with exact.
// TODO: add checks for the resource/apigroup validity.
type CelMatchConstraintsInfo struct {
	ResourceRules []Rules
	
}

type Rules struct {
	Operations []v1.OperationType
	APIGroups []string
	APIVersions []string
	Resources []string
	Scope []string // not used in current VAP.
}


// general CEL struct for VAP:
// not namspaced for now.
type CelInformation struct {
	Name string
	Variables              []CelVariableInfo
	Validations            []CelValidationInfo
	AuditAnnotations       []CelAuditAnnotationsInfo
	MatchConditions        []CelMatchConditionsInfo
	MatchConstraints 	   CelMatchConstraintsInfo
}

// Function to parse the byte array and give out the CelInformation struct.
func ExtractCelInfoFromFile(input []byte) (*CelInformation, error) {
	decoder := scheme.Codecs.UniversalDeserializer()

	runtimeObject, _, err := decoder.Decode(input, nil, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to decode the yaml, %w", err))
	} 
	policy, ok := runtimeObject.(*v1.ValidatingAdmissionPolicy)
	if !ok {
		return nil, fmt.Errorf("Type error for %T", runtimeObject)
	}
	return ExtractVAPInfo(policy)
}

// take in the runtime object and 
func ExtractVAPInfo(policy *v1.ValidatingAdmissionPolicy) (*CelInformation, error) {
	name := policy.ObjectMeta.GetName()

	variables := []CelVariableInfo{} // for variables if any
	for  _, variable := range policy.Spec.Variables {
		variables = append(variables, CelVariableInfo{
			Name:       variable.Name,
			Expression: variable.Expression,
		})
	}

	// for validations.
	validations := []CelValidationInfo{}
	for _, validation := range policy.Spec.Validations {
		validations = append(validations, CelValidationInfo{
			Expression:        validation.Expression,
			Message:           validation.Message,
			MessageExpression: validation.MessageExpression,
		})
	}

	// not used: audit annotations
	auditAnnotations := []CelAuditAnnotationsInfo{}
	for _, auditAnnotation := range policy.Spec.AuditAnnotations {
		auditAnnotations = append(auditAnnotations, CelAuditAnnotationsInfo{
			Key:        auditAnnotation.Key,
			Expression: auditAnnotation.ValueExpression,
		})
	}
	// not used in the current VAP spec.
	matchConditions := []CelMatchConditionsInfo{}
	for _, matchCondition := range policy.Spec.MatchConditions {
		matchConditions = append(matchConditions, CelMatchConditionsInfo{
			Name:       matchCondition.Name,
			Expression: matchCondition.Expression,
		})
	}

	matchConstraints := CelMatchConstraintsInfo{}
	if policy.Spec.MatchConstraints != nil {
		rules := make([]Rules, 0)
    	for _, rule := range policy.Spec.MatchConstraints.ResourceRules {
    	    rules = append(rules, Rules{
    	        Operations:  rule.Operations,
    	        APIGroups:  rule.APIGroups,
    	        APIVersions: rule.APIVersions,
    	        Resources:  rule.Resources,
    	    })
    	}
		matchConstraints = CelMatchConstraintsInfo{ // skipped Namespace and ObjectSelectors for now.
			ResourceRules: rules,
		}

	}

	return &CelInformation{name, variables, validations, auditAnnotations, matchConditions, matchConstraints}, nil
	
}


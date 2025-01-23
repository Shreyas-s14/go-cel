package yamlbuilder

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/google/cel-go/cel"
	// "gopkg.in/yaml.v3"
	"sigs.k8s.io/yaml"
	"k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	
)

/*	TODO:
	1) have a struct with fields expression and message. An array of such structs -> one VAP.
	2) Validate the expressions in each of the struct given the array of such structs
	3) Generate a struct for the VAP and VAP Binding.
	4) to have 2 such bindings and VAPs (deny / allow)
	5) Optimize?

*/

/*
	Input:

	struct{
		name string
		validations []validationInfo
	}


	helper: return (VAP, VAPbinding, error)
	main : []helper() -> yaml (vap), yaml(vapbinding), err
*/

// each FuncInput Struct -> 1x Vap + 1x VapBinding
type FuncInput struct {
	Name string
	Validations []ValidationInfo
}

// name: Policy name -> FuncInput.Name
func ValidateCELExpressions(env *cel.Env, expression string, index int, name string) (error) {
	if expression == "" {
		return fmt.Errorf("Empty expression in index %d for Policy %s", index, name)
	}
	// compile = parse + check(type check)
	ast, issues := env.Compile(expression)
	if issues != nil && issues.Err() != nil {
		return fmt.Errorf("Error while compiling expression %d for policy %s : %s", index, name, issues.Err())
	}

	// returns prg, err. prg: can be further evaluated using Eval, ContextEval. -> unit tests
	_, err := env.Program(ast)
	if err != nil {
		return fmt.Errorf("Program Construction Error: index: %d, Policy: %s, err: %s", index,  )
	}
	return nil
	
}

// Check the VAP for identical names. TODO: add further checks. (no expression check?)
func PreCheckVAP(inputs []FuncInput) error {
	if len(inputs) == 0 {
		return nil
	}

	nameMap := make(map[string]bool)

	for i, input := range inputs {
		if nameMap[input.Name] {
			return fmt.Errorf("duplicate name found in index %d: %s", i, input.Name)
		}
		nameMap[input.Name] = true
	}
	return nil
}

func EvaluateValidations(env *cel.Env, validations []ValidationInfo, policyName string) error {
	for i, validation := range validations {
		if err := ValidateCELExpressions(env, validation.Expression, i, policyName); err != nil {
			return err
		}
	}
	return nil
	
	
}

// TODO: Take in the array of inputs and generate structs for the 
func GenerateVAPYamlStructs(env *cel.Env, policy []FuncInput)  ([]v1.ValidatingAdmissionPolicy, []v1.ValidatingAdmissionPolicyBinding, error) {
	if err := PreCheckVAP(policy); err != nil {
		return nil, nil, err
	}
	// return values. 
	var vaps []v1.ValidatingAdmissionPolicy
	var bindings []v1.ValidatingAdmissionPolicyBinding

	for _, input := range(policy) {
		if err := EvaluateValidations(env, input.Validations, input.Name); err != nil {
			return nil, nil, err
		}
		
		validations := make([]v1.Validation, len(input.Validations))
		for i, validation := range input.Validations {
			validations[i] = v1.Validation{
				Expression: validation.Expression,
				Message:    validation.Message,
			}
		}
		// struct for VAP:
		failurePolicy := v1.Fail
		vap := v1.ValidatingAdmissionPolicy{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "admissionregistration.k8s.io/v1",
				Kind:       "ValidatingAdmissionPolicy",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: input.Name,
			},
			Spec: v1.ValidatingAdmissionPolicySpec{
				FailurePolicy: &failurePolicy,
				MatchConstraints: &v1.MatchResources{
					ResourceRules: []v1.NamedRuleWithOperations{
						{
							RuleWithOperations: v1.RuleWithOperations{
								Operations: []v1.OperationType{
									v1.Create,
									v1.Update,
								},
								Rule: v1.Rule{
									APIGroups:   []string{"druid.gardener.cloud"},
									APIVersions: []string{"v1alpha1"},
									Resources:   []string{"etcds"},
								},
							},
						},
					},
				},
				Validations: validations,
			},
		}
		vaps = append(vaps, vap)

		binding := v1.ValidatingAdmissionPolicyBinding{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "admissionregistration.k8s.io/v1",
				Kind:       "ValidatingAdmissionPolicyBinding",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("%s", input.Name),
			},
			Spec: v1.ValidatingAdmissionPolicyBindingSpec{
				PolicyName: input.Name,
				ValidationActions: []v1.ValidationAction{
					v1.Deny,
				},
				MatchResources: &v1.MatchResources{
					ResourceRules: []v1.NamedRuleWithOperations{
						{
							RuleWithOperations: v1.RuleWithOperations{
								Operations: []v1.OperationType{
									v1.Create,
									v1.Update,
								},
								Rule: v1.Rule{
									APIGroups:   []string{"druid.gardener.cloud"},
									APIVersions: []string{"v1alpha1"},
									Resources:   []string{"etcds"},
								},
							},
						},
					},
				},
			},
		}
		bindings = append(bindings, binding)
	}

	return vaps, bindings, nil



}

// take input as the array of validations FuncIputs. Return type: error. Generate yaml in the path.
func GenerateYAMLFromStruct(env *cel.Env, policy []FuncInput) error {
	vaps, bindings, err := GenerateVAPYamlStructs(env, policy)
	if err != nil {
		return fmt.Errorf("failed to generate VAP and VAPBinding structs: %v", err)
	}
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Errorf("Failed to get directory")
	}

	outputPath := filepath.Join(pwd, "generatedyamls")
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	for _, vap := range vaps {
		filename := filepath.Join(outputPath, fmt.Sprintf("%s-vap.yaml", vap.Name))
		err = YamlWriter(vap, filename)
		if err != nil {
			return fmt.Errorf("helper error")
		}
	}
	// binding vapbinding
	for _, binding := range bindings {
		filename := filepath.Join(outputPath, fmt.Sprintf("%s-binding.yaml", binding.Name))
		err = YamlWriter(binding, filename)
		if err != nil {
			return fmt.Errorf("helper error")
		}
	}

	return nil


}
// any-> alias for interface{}
func YamlWriter(data any, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %s", filename, err)
	}
	defer file.Close()

	yamlData, err := yaml.Marshal(data)

	if err != nil {
		return fmt.Errorf("failed to marshal to YAML for %s: %s", err, filename)
	}

	if _, err := file.Write(yamlData); err != nil {
		return fmt.Errorf("failed to write YAML to file %s: %s", filename, err)
	}
	return nil
}


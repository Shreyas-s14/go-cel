package yamlbuilder

import (
	"k8s.io/api/admissionregistration/v1"
	"github.com/google/cel-go/cel"
	"fmt"
)
func generateVAPYamlStructs(env *cel.Env, policy []FuncInput)  ([]VAP, []VAPBinding, error) {
	if err := PreCheckVAP(policy); err != nil {
		return nil, nil, err
	}
	// return values. 
	var vaps []VAP
	var bindings []VAPBinding

	for _, input := range(policy) {
		if err := EvaluateValidations(env, input.Validations, input.Name); err != nil {
			return nil, nil, err
		}
		
		vap := VAP{
			ApiVersion: "admissionregistration.k8s.io/v1",
			Kind:       "ValidatingAdmissionPolicy",
			Metadata: Meta{
				Name: input.Name,
			},
			Spec: VAPSpec{
				FailurePolicy: "Fail",
				MatchConstraints: MatchConstraintsInfo{
					ResourceRules: []Rules{
						{
							APIGroups: []string{
								"druid.gardener.cloud",
							},
							APIVersions: []string{
								"v1alpha1",
							},
							Operations: []v1.OperationType{
								v1.Create,
								v1.Update,
							},
							Resources: []string{
								"etcds",
							},
						},
					},
				},
				Validations: input.Validations,
			},
		}
		vaps = append(vaps, vap)


		binding := VAPBinding{
			ApiVersion: "admissionregistration.k8s.io/v1",
			Kind:       "ValidatingAdmissionPolicyBinding",
			Metadata: Meta{
				Name: fmt.Sprintf("%s-binding", input.Name),
			},
			Spec: VAPBindingSpec{
				PolicyName:        input.Name,
				ValidationActions: []string{"Deny"},
				MatchResources: MatchResourcesInfo{
					ResourceRules: []Rules{
						{
							APIGroups: []string{
								"druid.gardener.cloud",
							},
							APIVersions: []string{
								"v1alpha1",
							},
							Operations: []v1.OperationType{
								v1.Create,
								v1.Update,
							},
							Resources: []string{
								"etcds",
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
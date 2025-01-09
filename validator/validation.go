package validator

/*
	TODO:
	function to validate the expression after taking input as cel environment and 

*/

import (
	"github.com/google/cel-go/cel"
	"gocel/expression"
	"fmt"
)

// main -> create cel environment
func EvaluateVAPCel(env *cel.Env,policy *expression.CelInformation) error {
	fmt.Println(env)
	if env == nil {
		return fmt.Errorf("CEL environment cannot be nil")
	}

	if policy == nil {
		return fmt.Errorf("CelInformation cannot be nil")
	}
	// for validations:
	// fmt.Println(env.CELTypeProvider())
	for i, validations := range policy.Validations {
		ast, issues := env.Compile(validations.Expression)
		fmt.Println(issues.Errors())
		if issues != nil && issues.Err() != nil {
			return fmt.Errorf("Error occured while compiling expression %d, : %s", i, issues.Err())
		} 
		fmt.Println(ast.OutputType())
		fmt.Println(cel.AstToParsedExpr(ast)) // parsed tree
	}
	return nil
	
}
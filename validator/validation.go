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


func EvaluateCelExpression(env *cel.Env, expression string, index int, field string ) error {
	if expression == "" {
		return nil // no error if no expression?
	}
	// compile = parse + check.
	ast, issues := env.Compile(expression)
		// fmt.Println(issues.Errors())
	if issues != nil && issues.Err() != nil {
		return fmt.Errorf("Error occured while compiling expression %d for field %s : %s", index, field, issues.Err())
	} 
	fmt.Println(ast.IsChecked())
	// fmt.Println(cel.AstToParsedExpr(ast)) // parsed tree
	return nil
}
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
		if validations.Expression == "" && validations.Message != "" {
			return fmt.Errorf("Validation message is present, but Expression is missing")
		}
		if err := EvaluateCelExpression(env, validations.Expression,i, "validations"); err != nil {
			return err
		}
		
	}

	// variables: add the if condition later.
	for i, variables := range policy.Variables {
		if err := EvaluateCelExpression(env, variables.Expression, i, "variables"); err != nil {
			return err
		}
		
	}

	//audit annotations -> not used as of now:
	// same for match conditions. -> add later
	return nil
	
}
package builder

import (
	"errors"
	"fmt"
)

// generate expressions based on keyword??
// fields: validation []string ??,  message: string
// 0) first keyword: optional / required
// 1) object path -> TODO Add validation for path based on yaml provided
// 2) object type (needed for casting) : string/ quantity/ integer/float/duration
// 3) Comparison type(in case casting is to be done): array(no casting)/ string(casting based on prev type)/ quantity(casting..prev)/ integer/ float
// 4) validation keyword : matches(only string or if casted to string)/(less, equal, greater..)(quantity/ int /float)/ in (only for array)
// 5) object/expression to be validated against -> 1) 
type Expression struct{
	validation []string
	message string
}

// generated from the Expression struct. Array of such structs accepted by validation field in the resource- ValidatingAdmissionPolicy spec.
type GeneratedExp struct {
	expression string
	message string
}

func validationChecker(keys []string, expression string) (string, error) {
	// 1) & 2)
	if len(keys) == 0 || len(keys) < 6 {
		return "",errors.New("No Validation expression provided / incomplete validation provided")
	}
	// 2) can be done with routines
	if keys[0] == "optional" {
		expression = fmt.Sprintf("!has(%s) || ",keys[1])
	}
	// 3)
	switch keys[2] {
	case "string":
		if keys[3] != "string" && keys[3] != "array" {
			return "", errors.New("incorrect types for comparisons: object type 'string' must match 'string' or 'array' comparison type")
		}
		// TODO: Include other 
		if keys[4] == "matches" {
			expression += fmt.Sprintf("matches(%s, \"%s\")", keys[2], keys[5])
		} else {
			return "", errors.New("invalid validation keyword for string type")
		}
	// change the expression type from exp() to <>.exp()   + Add casting to string.
	case "quantity":
		if keys[3] != "quantity" {
			return "", errors.New("incorrect types for comparisons: object type 'quantity' must match 'quantity' comparison type")
		}
		if keys[4] == "less" || keys[4] == "greater" || keys[4] == "equal" {
			expression += fmt.Sprintf("%s(%s, %s)", keys[4], keys[2], keys[5])
		} else if keys[4] == "compareTo" {
			expression += fmt.Sprintf("quantity(%s).compareTo(quantity('%s'))", keys[2], keys[5])
			if keys[6] == "<=" {
				expression += " <= 0"
			} else if keys[6] == ">=" {
				expression += " >= 0"
			} else {
				return "", errors.New("invalid comparison operator for quantity type")
			} 
		} else {
			return "", errors.New("invalid validation keyword for quantity type")
		}

	case "integer":
		if keys[3] != "integer" && keys[3] != "array" {
			return "", errors.New("incorrect types for comparisons: object type 'integer' must match 'integer' or 'array' comparison type")
		}
		if keys[4] == "less" || keys[4] == "greater" || keys[4] == "equal" {
			expression += fmt.Sprintf("%s(%s, %s)", keys[4], keys[2], keys[5])
		} else {
			return "", errors.New("invalid validation keyword for integer type")
		}

	case "float":
		if keys[3] != "float" && keys[3] != "array" {
			return "", errors.New("incorrect types for comparisons: object type 'float' must match 'float' or 'array' comparison type")
		}
		if keys[4] == "less" || keys[4] == "greater" || keys[4] == "equal" {
			expression += fmt.Sprintf("%s(%s, %s)", keys[4], keys[2], keys[5])
		} else {
			return "", errors.New("invalid validation keyword for float type")
		}
	default:
		return "", errors.New("unsupported object type")
	}
	return "", nil
}

// Goes through the Expression struct and generates validation, message pair
/* For Expression.validation: 
	1) check if the validation is empty -> err : empty validation 
	2) check for all the fields-> err: please provide <field> (use routines???)
	3) optional -> add "!has(object path) || " to generated.validation, required -> continue
	4) comparison type, obj type: {array, <anything}} -> no casting, {string(regex, url), !string}-> add "string(" to expression. 
	5) keyword: {array, in}, TODO: handle edge cases
   For Expression.message:
	no validation
*/
func (exp *Expression) GenerateExp() (GeneratedExp, error) {
	p := new(GeneratedExp)
	p.expression = ""


	return *p, nil	
}

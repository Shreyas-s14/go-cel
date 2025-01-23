package yamlbuilder

// testing for valid/ invalid objects
// for each of the validation, 


// testfiles n

// Name-> field to be validated-> take from expression
type TestCase struct {
	Name string 
	Type string // positive/ negative. To be parsed from the filename 
	Data map[string]interface{} // yaml data
}
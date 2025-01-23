package main

import (
	// "context"
	// "fmt"
	"fmt"
	"os"
	"path/filepath"

	"gocel/expression"
	"gocel/validator"
	"gocel/yamlbuilder"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/ext"
	"k8s.io/apiserver/pkg/cel/library"
	// "k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/tools/clientcmd"
	"github.com/google/cel-go/checker/decls"
)

/*
	1) builder-> an idea... can be removed. (expression generator). FUTURE????
	2) expression-> where the structs will be defined. function takes in the byte array and makes it into a struct with each field as the yaml field. return this struct.
	3) validator: a) setup environment b) generate ast and validate the expressions in the struct byb calling the map key. return type err??
	4) main: calls these after reading the yaml file. returns the value from the validator.
	5) In case of the druid based implementation, create an object for the VAP.

*/
func main() {
	
	// fmt.Println(cel.DefaultMaxRequestSizeBytes)
	// read yaml file from location. 
	// TODO: CHeck on whether it is to be read as yaml only??
	yamlFile := filepath.Join("resources", "validating-policy.yaml")
	res, err := os.ReadFile(yamlFile)
	if err!=nil {
		fmt.Println(err)
		return 
	}

	celInfo, err := expression.ExtractCelInfoFromFile(res)
	if err != nil {
		fmt.Println(err)
	}
	// to include regex, URL, Lists and Quantity checks. -> k8s library predefined functions
	var CelEnv = []cel.EnvOption{
		cel.EagerlyValidateDeclarations(true),
		cel.DefaultUTCTimeZone(true),
		ext.Strings(ext.StringsVersion(2)),
		cel.CrossTypeNumericComparisons(true),
		cel.OptionalTypes(),
		library.URLs(),
		library.Regex(),
		library.Lists(),
		library.Quantity(),

		// cel.Variable("object", 
        // cel.ObjectType("object"),),
		// decls to include . Got error while using above-> cannot resolve Object type
		// using map[string]any -> resolves the above
		cel.Declarations(
            decls.NewVar("object", decls.NewMapType(decls.String, decls.Any)),
        ),

		// add functions?? -> decls.NewFunction() // TODO: research on args and return types.
	}

	env, err := cel.NewEnv(CelEnv...)
	if err != nil {
		fmt.Printf("Error: %s\n",err)
		return
	}
	err = validator.EvaluateVAPCel(env,celInfo)
	if err != nil {
		fmt.Printf("Error : %s\n", err)
		return
	} else {
		fmt.Println("worked.")
	}

	err = yamlbuilder.GenerateYAMLFromStruct(env, yamlbuilder.ValidatingAdmissionPolicies)
	if err != nil {
		fmt.Println(err)
	}
	



	


}
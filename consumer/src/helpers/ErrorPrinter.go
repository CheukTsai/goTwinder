package helpers

import(
	"log"
)

func FailOnFatalError(err error, msg string, funcName string) {
	if err != nil {
		log.Fatalf("Fatal error executing function %s;\n Error message: %s - %s", funcName, msg, err);
	} 
}

func FailOnError(err error, msg string, funcName string) {
	if err != nil {
		log.Panicf("Error executing function %s;\n Error message: %s - %s", funcName, msg, err);
	} 
}

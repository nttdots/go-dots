package dots_common

import (
	"errors"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
)

const SCHEMA_DIR string = "schemas"

/*
 * get the path for the schema files.
 */
func getSchemaPath(schemaDir string, schemaName string) string {
	return fmt.Sprintf("%s/%s.json", schemaDir, schemaName)
}

/*
 * validate the Json strings.
 */
func ValidateJson(schemaName string, request string) error {
	schemaPath := getSchemaPath(SCHEMA_DIR, schemaName)

	schemaFile, err := Asset(schemaPath)
	if err != nil {
		return err
	}
	schemaLoader := gojsonschema.NewBytesLoader(schemaFile)
	loader := gojsonschema.NewStringLoader(request)

	result, err := gojsonschema.Validate(schemaLoader, loader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		error_msg := "The document is not valid. see errors :"
		for _, desc := range result.Errors() {
			error_msg += fmt.Sprintf("- %s\n", desc)
		}
		return errors.New(error_msg)
	}
	return nil
}

// DigiStratum GoLib - JSON
package json

/*

Dealing with JSON at a level of abstraction above encoding/json.

*/

import(
	"os"
	"fmt"
	"encoding/json"
)

type JsonIfc interface {
	Load(target interface{}) error
}

type Json struct {
	source	string
	path	string
	json	*string
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewJson(jsonString *string) *Json {
	return &Json{ json: jsonString, source: "string" }
}

// Make a new one of these (from file)!
func NewJsonFromFile(path string) *Json {
	return &Json{ path: path, source: "file" }
}

// -------------------------------------------------------------------------------------------------
// JsonIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Generic JSON load (into ANY interface)
// The provided target should be a pointer to where we will dump the decoded JSON result
func (r *Json) Load(target interface{}) error {
	switch (r.source) {
		case "string":
			if (nil == r.json) || ("" == *r.json) {
				return fmt.Errorf(
					"Json.Load(): We were given nil or empty string for the JSON (string)",
				)
			}
			if err := json.Unmarshal([]byte(*r.json), &target); err != nil {
				return fmt.Errorf(
					"Json.Load(): Failed to unmarshall JSON (string): %s",
					err.Error(),
				)
			}
			return nil

		case "file":
			file, err := os.Open(r.path)
			defer file.Close()
			if nil == err {
				decoder := json.NewDecoder(file)
				err = decoder.Decode(target)
				if nil == err { return nil }
			}
			// Decorate the errror with a little more context
			return fmt.Errorf(
				"Json.Load(): (file='%s'): '%s'", r.path, err.Error(),
			)
	}

	return fmt.Errorf("Json.Load(): Unsupported JSON source (%s)", r.source)
}

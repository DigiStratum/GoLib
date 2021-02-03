// DigiStratum GoLib - JSON
package golib

/*

Dealing with JSON at a level of abstraction above encoding/json.

*/

import(
	"os"
	"fmt"
	gojson "encoding/json"
	"errors"
)

type json struct {
	source	string
	path	string
	json	*string
}

type JSON struct {
	Json	json
}

// Make a new one of these (from string)!
func NewJson(jsonString *string) *json {
	return &json{ json: jsonString, source: "string" }
}

// Make a new one of these (from file)!
func NewJsonFromFile(path string) *json {
	return &json{ path: path, source: "file" }
}

// Make a new one of these (from string)!
func NewJSON(jsonString *string) *JSON {
	return &JSON{
		Json: json{
			json: jsonString,
			source: "string",
		},
	}
}

// Make a new one of these (from file)!
func NewJSONFromFile(path string) *JSON {
	return &JSON{
		Json: json{
			path: path,
			source: "file",
		},
	}
}

// Generic JSON load (into ANY interface)
// The provided target should be a pointer to where we will dump the decoded JSON result
func (j *json) Load(target interface{}) error {
	switch (j.source) {
		case "string":
			if (nil == j.json) || ("" == *(j.json)) {
				return errors.New(
					"Json.Load(): We were given nil or empty string for the JSON (string)",
				)
			}
			if err := gojson.Unmarshal([]byte(*(j.json)), &target); err != nil {
				return errors.New(fmt.Sprintf(
					"Json.Load(): Failed to unmarshall JSON (string): %s",
					err.Error(),
				))
			}
			return nil

		case "file":
			file, err := os.Open(j.path)
			defer file.Close()
			if nil == err {
				decoder := gojson.NewDecoder(file)
				err = decoder.Decode(target)
				if nil == err { return nil }
			}
			// Decorate the errror with a little more context
			return errors.New(fmt.Sprintf(
				"Json.Load(): (file='%s'): '%s'", j.path, err.Error(),
			))
	}

	return errors.New(fmt.Sprintf("Json.Load(): Unsupported JSON source (%s)", j.source))
}

// Generic JSON load (or panic)
// The provided target should be a pointer to where we will dump the decoded JSON result
func (j *json) LoadOrPanic(target interface{}) {
	if err := j.Load(target); nil != err { panic(err.Error()) }
}


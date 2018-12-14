// DigiStratum GoLib - JSON
package golib

/*

Dealing with JSON at a level of abstraction above encoding/json.

*/

import(
	"os"
	"fmt"
	"encoding/json"
	"errors"
)

type Json struct {
	source	string
	path	string
	json	*string
}

// Make a new one of these (from string)!
func NewJson(json *string) *Json {
	return &Json{ json: json, source: "string" }
}

// Make a new one of these (from file)!
func NewJsonFromFile(path string) *Json {
	return &Json{ path: path, source: "file" }
}

// Generic JSON load (into ANY interface)
// The provided target should be a pointer to where we will dump the decoded JSON result
func (j *Json) Load(target interface{}) error {
	switch (j.source) {
		case "string":
			if nil == j.json {
				return errors.New(
					"Json.Load(): We were given nil for the JSON (string)",
				)
			}
			if err := json.Unmarshal([]byte(*(j.json)), &target); err != nil {
				return errors.New(fmt.Sprintf(
					"Json.Load(): Failed to unmarshall JSON (string): %s",
					err.Error(),
				))
			}
			return nil

		case "file":
			file, err := os.Open(j.path)
			defer file.Close()
			if nil != err {
				decoder := json.NewDecoder(file)
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
func (j *Json) LoadOrPanic(target interface{}) {
	if err := j.Load(target); nil != err { panic(err.Error()) }
}


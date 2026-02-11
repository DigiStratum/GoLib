package serializable

/*
General purpose de|serializer

TODO:
 * Allow transcoder to be nil and auto-detect/instatiate a Transcoder based on detected encoding/method
 * Add support for []byte variants of the string based methods
 * Add support for file path variants
*/

import (
	"fmt"
	"regexp"

	xc "github.com/DigiStratum/GoLib/Data/transcoder"
)

type SerializerIfc interface {
	Serialize(data *string, typeName string) (*string, error)
	Deserialize(data *string) (*string, error)
}

type Serializer struct {
	transcoder	xc.TranscoderIfc
}

func NewSerializer(transcoder xc.TranscoderIfc) *Serializer {
	return &Serializer{
		transcoder:	transcoder,
	}
}

func (r Serializer) Serialize(data *string, typeName string) (*string, error) {

	// Encoding base64, JSON data
	if (nil == r.transcoder) || ("base64" != r.transcoder.GetEncoderSchemeName()) {
		return nil, fmt.Errorf("Serialization requires Transcoder with EncodingSchemeBase64 Encoder")
	}

	// Overall format: "ser[{Method}:{Type}:{Data}]"
	method := "j64"
	edata, err := r.transcoder.Encode(data)
	if (nil != err) || (nil == edata) {
		return nil, fmt.Errorf("Error encoding serialized data")
	}

	etype, err := r.transcoder.Encode(&typeName)
	if (nil != err) || (nil == etype) {
		return nil, fmt.Errorf("Error encoding serialized data type")
	}

	serialized := fmt.Sprintf("ser[%s:%s:%s]", method, *etype, *edata)
	return &serialized, nil

}

func (r Serializer) Deserialize(data *string) (*string, error) {
	// TODO: Allow transcoder to be nil and auto-detect/instatiate a Transcoder based on detected encoding/method
	if nil == data { return nil, fmt.Errorf("Cannot deserialize nil data") }

	// Overall format: "ser[{Method}:{Type}:{Data}]"
	matchNameMethod := "method"     // Method of encoding (determines which Transcoder to expect)
	matchNameEType := "etype"
	matchNameEData := "edata"
	re := regexp.MustCompile(
		fmt.Sprintf(
			"^ser\\[(?P<%s>\\w+):(?P<%s>[A-Za-z0-9+/=]*):(?P<%s>[A-Za-z0-9+/=]*)\\]$",
			matchNameMethod, matchNameEType, matchNameEData,
		),
	)
	if matched := re.MatchString(*data); !matched {
		return nil, fmt.Errorf("Data does not match expected format")
	}

	matches := re.FindStringSubmatch(*data)
	if (nil == matches) || (len(matches) < 3) {
		return nil, fmt.Errorf("Unexpected mismatch on parameters for serialized value")
	}

	method := matches[re.SubexpIndex(matchNameMethod)]
	etype := matches[re.SubexpIndex(matchNameEType)]
	edata := matches[re.SubexpIndex(matchNameEData)]

	// We support multiple methods of object deserialization
	switch method {
		case "j64":
			// Encoding base64, JSON data
			if (nil == r.transcoder) || ("base64" != r.transcoder.GetDecoderSchemeName()) {
				return nil, fmt.Errorf("Deserialization requires Transcoder with EncodingSchemeBase64 Decoder")
			}

			utype, err := r.transcoder.Decode(&etype)
			if (nil != err) || (nil == utype) || ("Object" != *utype) {
				return nil, fmt.Errorf("Error decoding serialized data type")
			}

			return r.transcoder.Decode(&edata)
	}

	return nil, fmt.Errorf("Unsupported serialization method '%s'", method)
}


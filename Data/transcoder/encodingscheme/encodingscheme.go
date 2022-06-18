package encodingscheme

type EncodingSchemeIfc interface {
	GetName() string
	SetEncodedValue(source *string) error
	GetEncodedValue() (*string, error)
	SetDecodedValue(source *string) error
	GetDecodedValue() (*string, error)
}


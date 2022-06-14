package encodingscheme

type EncodingSchemeIfc interface {
	SetEncodedValue(source *string) error
	GetEncodedValue() (*string, error)
	SetDecodedValue(source *string) error
	GetDecodedValue() (*string, error)
}


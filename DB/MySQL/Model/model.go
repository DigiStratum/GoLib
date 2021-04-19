package mmodel;

type FieldType int

const (
        INT FieldType = iota
	TINYINT
	SMALLINT
	MEDIUMINT
	BIGINT
	DECIMAL
	NUMERIC
	FLOAT
	DOUBLE
	BIT
        DATE
	TIME
	YEAR
	DATETIME
	TIMESTAMP
	CHAR
	VARCHAR
	BINARY
	VARBINARY
	BLOB
	TINYBLOB
	MEDIUMBLOB
	LONGBLOB
	TEXT
	TINYTEXT
	MEDIUMTEXT
	LONGTEXT
	ENUM
	SET
)


type FieldTypeConstraints struct {
	Type		FieldType
	Min		int
	Max		int
}

type Field struct {
	Name		string
	Type		FieldType
}


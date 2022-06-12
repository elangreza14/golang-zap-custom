package logger

type (
	EncodingType string

	option struct {
		EncodingType     EncodingType
		NameService      string
		EnableStackTrace bool
	}
	Option option
)

const (
	EncodingTypeJson    EncodingType = "json"
	EncodingTypeConsole EncodingType = "console"
)

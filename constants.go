package restful

const (
	MIME_JSON = "application/json"
	MIME_XML  = "application/xml"
	MIME_TEXT = "text/plain"
	MIME_HTML = "text/html"
	MIME_FORM = "application/x-www-form-urlencoded"
)

type FilterLevel int

const (
	FilterLevelContainer FilterLevel = iota
	FilterLevelWebService
	FilterLevelRoute
)

const (
	DefaultSuccessStatus = 200
	DefaultCreateStatus  = 201
	DefaultDeleteStatus  = 204
	DefaultErrorStatus   = 500
)

var (
	defaultProduces = []string{}
	defaultConsumes = []string{}
)

func SetDefaultProduces(m ...string) {
	defaultProduces = append([]string(nil), m...)
}

func SetDefaultConsumes(m ...string) {
	defaultConsumes = append([]string(nil), m...)
}

func DefaultProduces() []string {
	return append([]string(nil), defaultProduces...)
}

func DefaultConsumes() []string {
	return append([]string(nil), defaultConsumes...)
}

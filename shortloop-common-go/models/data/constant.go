package data

type DataSizeUnit string
type DataCategory string
type HTTPRequestMethod string
type ContentType string

const (
	DataCategoryMessageAttributeKey = "DATA_CATEGORY"
	ContentTypeHeaderLowerCase      = "content-type"
	ContentTypeHeaderJSONValue      = "application/json"
	MaskingRulesConfigurationName   = "PIIMaskingRules"
	DefaultMaskReplacement          = "MASKED"
	MaskedValueSuffix               = "SL_MASK"

	KB DataSizeUnit = "KB"
	MB DataSizeUnit = "MB"
	GB DataSizeUnit = "GB"

	SampleData DataCategory = "SAMPLE_DATA"

	GET     HTTPRequestMethod = "GET"
	HEAD    HTTPRequestMethod = "HEAD"
	POST    HTTPRequestMethod = "POST"
	PUT     HTTPRequestMethod = "PUT"
	PATCH   HTTPRequestMethod = "PATCH"
	DELETE  HTTPRequestMethod = "DELETE"
	OPTIONS HTTPRequestMethod = "OPTIONS"
	TRACE   HTTPRequestMethod = "TRACE"
)

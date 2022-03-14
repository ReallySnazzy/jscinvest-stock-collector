package iexcloud

const (
	iexcloudBaseUrl = "https://cloud.iexapis.com"
)

type Iexcloud struct {
	apiKey    string
	dateRange string
}

func CreateIexcloudContext(apiKey string, dateRange string) Iexcloud {
	context := Iexcloud{}
	context.apiKey = apiKey
	context.dateRange = dateRange
	return context
}

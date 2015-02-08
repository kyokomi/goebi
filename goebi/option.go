package goebi

import "fmt"

// Options goebi Client options
type Options struct {
	// Host (required) hostName example `http://localhost:3000`
	Host string
	// ProjectID is found AirBrake.io in the URL.
	ProjectID string
	// ApiPath apiPath example `/api/v3/projects`
	APIPath string
	// ApiKey (required) apiKey is issued when you register the app to errbit.
	APIKey string
}

// v3 API
func (opt Options) createNoticeBaseURL() string {

	host := opt.Host
	if host[0] == '/' {
		host = host[1:]
	}

	if host[len(host)-1] == '/' {
		host = host[:len(host)-1]
	}

	apiPath := opt.APIPath
	if apiPath == "" {
		apiPath = "/api/v3/projects"
	}

	if apiPath[0] == '/' {
		apiPath = apiPath[1:]
	}

	if apiPath[len(apiPath)-1] == '/' {
		apiPath = apiPath[:len(apiPath)-1]
	}

	// ProjectIDなしならApiKeyを使う
	projectID := opt.ProjectID
	if projectID == "" {
		projectID = opt.APIKey
	}

	// http://localhost:8000/api/v3/projects/xxxxx/notices?key=xxxxx
	return fmt.Sprintf("%s/%s/%s/%s", host, apiPath, projectID, "notices?key="+opt.APIKey)
}

package gobit

import "fmt"

// Options gobit Client options
type Options struct {
	// Host hostName example `http://localhost:3000`
	Host string
	// ApiPath apiPath example `/api/v3/projects`
	ApiPath string
	// ApiKey apiKey is issued when you register the app to errbit.
	ApiKey string
}

func (opt Options) createNoticeBaseURL() string {

	if opt.Host[len(opt.Host)-1] == '/' {
		opt.Host = opt.Host[:len(opt.Host)-1]
	}

	if opt.ApiPath[len(opt.ApiPath)-1] == '/' {
		opt.ApiPath = opt.ApiPath[:len(opt.ApiPath)-1]
	}
	
	// http://localhost:8000/api/v3/projects/xxxxx/notices?key=xxxxx
	return fmt.Sprintf("%s/%s/%s/%s", opt.Host, opt.ApiPath, opt.ApiKey, "notices?key=" + opt.ApiKey)
}

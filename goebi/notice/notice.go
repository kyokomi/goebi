package notice

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
)

// Notice  error内容
type Notice struct {
	Notifier Notifier               `json:"notifier"`
	Errors   []ErrorReport          `json:"errors"`
	Env      map[string]interface{} `json:"environment"`
	Session  map[string]interface{} `json:"session"`
	Params   map[string]interface{} `json:"params"`
	Context  Context                `json:"context`
}

// Notifier error送信者
type Notifier struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	URL     string `json:"url"`
}

// ErrorReport エラー情報
type ErrorReport struct {
	ErrorType string      `json:"type"`
	Message   string      `json:"message"`
	Backtrace []BackTrace `json:"backtrace"`
}

// BackTrace stackTrace
type BackTrace struct {
	File string `json:"file"`
	Line int    `json:"line"`
	Func string `json:"function"`
}

// Context context
type Context struct {
	URL           string `json:"url"`
	OS            string `json:"is"`
	Language      string `json:"language"`
	Environment   string `json:"environment"`
	RootDirectory string `json:"rootDirectory"`
	Version       string `json:"version"`
}

// NewNotice エラー通知を作成
func NewNotice(notifier Notifier, err interface{}, stack []BackTrace) *Notice {

	n := &Notice{}

	n.Notifier = notifier

	n.Errors = []ErrorReport{
		ErrorReport{
			ErrorType: fmt.Sprintf("%T", err),
			Message:   fmt.Sprint(err),
			Backtrace: stack,
		},
	}

	n.Context = Context{}
	n.Env = make(map[string]interface{})
	n.Session = make(map[string]interface{})
	n.Params = make(map[string]interface{})

	return n
}

// SetHTTPRequest http.Requestの内容を通知内容に設定します
func (n *Notice) SetHTTPRequest(req *http.Request) {

	n.Context.URL = req.URL.String()

	if ua := req.Header.Get("User-Agent"); ua != "" {
		n.Env["userAgent"] = ua
	}

	for k, v := range req.Header {
		if len(v) == 1 {
			n.Env[k] = v[0]
		} else {
			n.Env[k] = v
		}
	}

	// TODO: jsonのParamsがとれない　いずれ対応する...
	if err := req.ParseForm(); err != nil {
		return
	}

	for k, v := range req.Form {
		if len(v) == 1 {
			n.Params[k] = v[0]
		} else {
			n.Params[k] = v
		}
	}
}


}

// SetRuntime setup context default runtime.
func (n *Notice) SetRuntime() {
	n.Context.Language = runtime.GOOS
	n.Context.Version = runtime.Version()

	if hostname, err := os.Hostname(); err == nil {
		n.Context.URL = hostname
	}
	if wd, err := os.Getwd(); err == nil {
		n.Context.RootDirectory = wd
	}
}

// SetEnvRuntime setup context and env default runtime.
func (n *Notice) SetEnvRuntime() {
	n.SetRuntime()

	n.Env["language"] = n.Context.Language
	n.Env["version"] = n.Context.Version

	n.Env["architecture"] = runtime.GOARCH
}

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
	Context  Context                `json:"context"`
	Errors   []ErrorReport          `json:"errors"`

	// optional
	Env      map[string]interface{} `json:"environment"`
	Params   map[string]interface{} `json:"params"`
	Session  map[string]interface{} `json:"session"`
}

// Notifier error送信者
type Notifier struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	URL     string `json:"url"`
}

// Context context
type Context struct {
	// エラーになったURL等
	URL              string `json:"url"`

	// TODO: 未使用？
	SourceMapEnabled bool   `json:"sourceMapEnabled"`

	// Where
	// Controllerなどを指定
	Component        string `json:"component"`
	// Controllerのメソッド等を指定（Handler）
	Action           string `json:"action"`

	// AppServerの情報
	Language         string `json:"language"`
	Version          string `json:"version"`

	// User情報
	User

	RootDirectory    string `json:"rootDirectory"`
}

type User struct {
	UserID           int    `json:"userId"`
	UserName         string `json:"userName"`
	UserUsername     string `json:"userUsername"`
	UserEmail        string `json:"userEmail"`
	UserAgent        string `json:"userAgent"`
}

// ErrorReport エラー情報
type ErrorReport struct {
	ErrorType string      `json:"type"`
	Message   string      `json:"message"`
	Backtrace []BackTrace `json:"backtrace"`
}

// BackTrace stackTrace
type BackTrace struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Func     string `json:"function"`
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
		n.Context.UserAgent = ua
	}

	for k, v := range req.Header {
		if len(v) == 1 {
			n.Env["HTTP_" + k] = v[0]
		} else {
			n.Env["HTTP_" + k] = v
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

// SetUserInfo setup context.user
func (n *Notice) SetUserInfo(user User) {
	n.Context.User = user
}

// SetWhere setup context.where
func (n *Notice) SetWhere(packageName string, methodName string) {
	n.Context.Component = packageName
	n.Context.Action = methodName
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

package gobit

import "github.com/kyokomi/gobit/gobit/notice"

var appNotifier = notice.Notifier{
	Version: "0.1.0",
	Name:    "gobit",
	URL:     "https://github.com/kyokomi/gobit",
}

// NewNotice エラー通知内容を作成します
func NewNotice(err interface{}) *notice.Notice {
	return notice.NewNotice(appNotifier, err, defaultStackTrace())
}

// NewNoticeWithFilter StackFilterを指定してエラー通知内容を作ります
func NewNoticeWithFilter(err interface{}, filter StackFilterFunc) *notice.Notice {
	return notice.NewNotice(appNotifier, err, stackTrace(filter))
}

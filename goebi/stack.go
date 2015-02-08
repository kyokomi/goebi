package goebi

import (
	"runtime"
	"strings"

	"github.com/kyokomi/goebi/goebi/notice"
)

// StackFilterFunc stackTraceのFilter
type StackFilterFunc func(file string, line int, packageName, funcName string) bool

func defaultStackTrace() []notice.BackTrace {
	return stackTrace(func(_ string, _ int, packageName, funcName string) bool {
		return packageName == "runtime" && funcName == "panic"
	})
}

func stackTrace(filter StackFilterFunc) []notice.BackTrace {
	// stackTrace -> newNotice -> NewNoticeWithFilter or NewNotice
	startFrame := 3

	stack := []notice.BackTrace{}
	for i := startFrame; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		packageName, funcName := packageFuncName(pc)
		if filter(file, line, packageName, funcName) {
			stack = stack[:0]
			continue
		}
		stack = append(stack, notice.BackTrace{
			File: file,
			Line: line,
			Func: funcName,
		})
	}

	return stack
}

func packageFuncName(pc uintptr) (string, string) {
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "", ""
	}

	packageName := ""
	funcName := f.Name()

	if ind := strings.LastIndex(funcName, "/"); ind > 0 {
		packageName += funcName[:ind+1]
		funcName = funcName[ind+1:]
	}
	if ind := strings.Index(funcName, "."); ind > 0 {
		packageName += funcName[:ind]
		funcName = funcName[ind+1:]
	}

	return packageName, funcName
}

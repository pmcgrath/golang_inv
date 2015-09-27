package main

import (
	"fmt"
	"log"
	"runtime"
)

var isVerbose = false

func logDebugf(format string, v ...interface{}) {
	if isVerbose {
		log.Printf(getCallerFuncName()+": "+format, v...)
	}
}

func logDebugln(v ...interface{}) {
	if isVerbose {
		if len(v) > 0 {
			v[0] = fmt.Sprintf("%s: %v", getCallerFuncName(), v[0])
		}

		log.Println(v...)
	}
}

func getCallerFuncName() string {
	// See 	http://stackoverflow.com/questions/10742749/get-name-of-function-using-google-gos-reflection
	//	http://play.golang.org/p/teu5CnHoek
	pc, _, _, ok := runtime.Caller(2) // 2 as we have the function, the log function and this function
	if !ok {
		return "unknown"
	}
	return runtime.FuncForPC(pc).Name()
}

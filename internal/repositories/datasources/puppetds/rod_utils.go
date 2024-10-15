package puppetds

import (
	"time"

	"github.com/go-rod/rod"
	"github.com/ysmood/gson"
)

func obtainJSValues[T any](
	page *rod.Page, jsCode string, dst *T,
) (gson.JSON, error) {
	runtimeObj, err := page.Eval(jsCode)
	if runtimeObj == nil || err != nil {
		return gson.JSON{}, err
	}

	resultVal := runtimeObj.Value
	switch destinationPointer := any(dst).(type) {
	case *int:
		*destinationPointer = resultVal.Int()
	case *string:
		*destinationPointer = resultVal.String()
	case *float64:
		*destinationPointer = resultVal.Num()
	}
	return runtimeObj.Value, nil
}

func scrollPageDown(page *rod.Page) (err error) {
	var (
		pageHeight   = 500
		windowHeight = 100
	)
	_, err = obtainJSValues(page, "() => document.body.scrollHeight", &pageHeight)
	_, err = obtainJSValues(page, `() => window.innerHeight`, &windowHeight)

	totalChunks := pageHeight / (windowHeight * 2)
	for chunk := 0; chunk < totalChunks; chunk += 1 {
		page.Mouse.MustScroll(0, float64(windowHeight*chunk))
		time.Sleep(time.Millisecond << 8)
	}

	return
}

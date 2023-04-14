package httphelper

import (
	"bytes"
	"fmt"
	"github.com/artisancloud/httphelper/dataflow"
	"log"
	"net/http"
	"net/http/httputil"
)

func HttpDebugMiddleware(debug bool) dataflow.RequestMiddleware {
	return func(handle dataflow.RequestHandle) dataflow.RequestHandle {
		return func(request *http.Request, response *http.Response) (err error) {
			if debug {
				// Print request
				dumpReq, _ := httputil.DumpRequest(request, true)
				formattedReq := bytes.ReplaceAll(dumpReq, []byte("\r\n"), []byte("\n"))
				log.Print(fmt.Sprintf("[HTTP DEBUG] Request:\n%s\n", formattedReq))

				// Handle the request
				err = handle(request, response)
				if err != nil {
					return err
				}

				// Print response
				dumpRes, _ := httputil.DumpResponse(response, true)
				formattedRes := bytes.ReplaceAll(dumpRes, []byte("\r\n"), []byte("\n"))

				log.Print(fmt.Sprintf("------------------\n[HTTP DEBUG] Response:\n%s\n", formattedRes))
			} else {
				err = handle(request, response)
				if err != nil {
					return err
				}
			}

			return
		}
	}
}

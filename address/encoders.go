package address

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/getsafepay/communist/payments/common"
	"io/ioutil"
	"net/http"
)

func encodeValidateRequest(c context.Context, r *http.Request, request interface{}) error {
	common.SetHeaders(r, request)

	if req, ok := request.(*ValidateRequest); ok {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(req); err != nil {
			return err
		}
		r.Body = ioutil.NopCloser(&buf)
		r.ContentLength = int64(buf.Len())
	}

	return nil
}

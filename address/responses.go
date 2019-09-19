package address

import "github.com/ziyadparekh/communist/payments/common"

type ValidateResponse struct {
  Data interface{} `json:"data"`

  Status common.Status `json:"status"`
}

func (vr *ValidateResponse) DoesError() []string {
  return vr.Status.Errors
}

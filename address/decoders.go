package address

import (
  "net/http"
  "encoding/json"
  "context"
)

func decodeValidateReponse(c context.Context, r *http.Response) (interface{}, error) {
  var response ValidateResponse
  if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
    return nil, err
  }
  return response, nil
}

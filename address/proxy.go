package address

import (
  "github.com/go-kit/kit/endpoint"
  httptransport "github.com/go-kit/kit/transport/http"
  "strings"
  "net/url"
)

const (
  ADDRESS_VALIDATE_ENDPOINT string = "/meta/v1/"
)

func makeValidateProxy(instance string) endpoint.Endpoint {
  if !strings.HasPrefix(instance, "http") {
    instance = "http://" + instance
  }
  u, err := url.Parse(instance)
  if err != nil {
    panic(err)
  }
  if u.Path == "" {
    u.Path = ADDRESS_VALIDATE_ENDPOINT
  }

  return httptransport.NewClient(
    "POST",
    u,
    encodeValidateRequest,
    decodeValidateReponse,
  ).Endpoint()
}

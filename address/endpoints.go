package address

import (
  "github.com/go-kit/kit/endpoint"
  "context"

  "errors"
  "strings"
  "github.com/go-kit/kit/log"
  "github.com/go-kit/kit/log/level"
  "fmt"
)

type Addressor interface {
  Validate(ctx context.Context, req *ValidateRequest) (interface{}, error)
}

type address struct {
  validate endpoint.Endpoint
  logger log.Logger
}

func NewAddress(url string, logger log.Logger) Addressor {
  var v endpoint.Endpoint
  v = makeValidateProxy(url)

  return &address{
    validate: v,
    logger:logger,
  }
}

func (a *address) Validate(ctx context.Context, req *ValidateRequest) (interface{}, error) {
  res, err := a.validate(ctx, req)
  if err != nil {
    level.Error(a.logger).Log("message", "error contacting address service", "err", err.Error())
    return nil, err
  }

  result, ok := res.(ValidateResponse)
  if !ok {
    level.Error(a.logger).Log("message", "unexpected response from address service", "res", fmt.Sprintf("%v", result))
    return nil, errors.New("unexpected response from address")
  }

  if len(result.DoesError()) > 0 {
    level.Error(a.logger).Log("message", "address service returned errors", "err", strings.Join(result.DoesError(), ","))
    return nil, errors.New(strings.Join(result.DoesError(), ","))
  }

  return result.Data, nil
}

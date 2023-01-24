package main

import (
	"context"
	nativeerr "errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/pkg/errors"
)

var ErrNoValueFound error = nativeerr.New("no value found")

// LocalSSM
// implements ssmiface.SSMAPI.
type LocalSSM struct {
	ssmiface.SSMAPI
	store map[string]string
}

func NewLocalSSM(store map[string]string) *LocalSSM {
	if store == nil {
		store = make(map[string]string)
	}
	return &LocalSSM{
		store: store,
	}
}

// override.
func (s *LocalSSM) GetParameterWithContext(
	ctx context.Context, input *ssm.GetParameterInput, opts ...request.Option,
) (*ssm.GetParameterOutput, error) {
	val, ok := s.store[*input.Name]
	if !ok {
		return nil, errors.WithStack(ErrNoValueFound)
	}
	return &ssm.GetParameterOutput{
		Parameter: &ssm.Parameter{
			Value: aws.String(val),
		},
	}, nil
}

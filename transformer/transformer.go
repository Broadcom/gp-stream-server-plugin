package transformer

import (
	"errors"
	"github.com/sirupsen/logrus"
)

type TransformStatus int

var (
	ErrorTransformerLoad = errors.New("error load transformer")
	ErrorTransformerInit = errors.New("error init transformer")
)

const (
	// TransformStatusAccept means transform success
	TransformStatusAccept = iota + 1

	// TransformStatusReject means transform failed, data dropped but will not return error to gpss
	TransformStatusReject

	// TransformStatusError means transform failed, error
	TransformStatusError
)

type TransformContext interface {
	GetLogger() logrus.FieldLogger
	GetProperties() map[string]string
	GetInput() []byte

	SetTransformStatus(TransformStatus)
	SetOutput([]byte)
	SetError(error)
}

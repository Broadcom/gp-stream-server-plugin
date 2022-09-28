package transformer

import (
	"errors"
	common "github.com/greenplum-db/gp-stream-server-plugin"
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
	common.BaseContext

	// GetInput returns data to be transformed by the transformer
	GetInput() []byte

	// SetTransformStatus sets the transform result,
	// it should be TransformStatusAccept/TransformStatusReject/TransformStatusError
	SetTransformStatus(TransformStatus)

	// SetOutput sets the transformed data
	SetOutput([]byte)

	// SetError sets the error during transforming
	SetError(error)
}

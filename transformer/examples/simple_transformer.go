package main

/**
This is a simple transformer plugin implementation, it's used mainly for example and e2e tests.
The transformer accepts data like {"A":2,"B":1} as input, and takes operator from properties,
if operator=="+", it puts C=A+B back.
if operator=='-', it puts C=A-B back.
otherwise, an error occurs
*/
import (
	"encoding/json"
	"errors"
	common "github.com/greenplum-db/gp-stream-server-plugin"
	"github.com/greenplum-db/gp-stream-server-plugin/transformer"
)

var UnsupportedOpError = errors.New("simple transformer error: unsupported op")

type record struct {
	A int
	B int
	C int
}

var name string

func SimpleTransformOnInit(ctx common.BaseContext) error {
	properties := ctx.GetProperties()
	name = properties["name"]
	simulateInitError := properties["simulate-init-error"]
	ctx.GetLogger().Infof("plugin:%s init start...", name)
	// do some init work...
	if simulateInitError == "true" {
		// simulate an error
		ctx.GetLogger().Errorf("plugin:%s init error", name)
		return transformer.ErrorTransformerInit
	} else {
		ctx.GetLogger().Infof("plugin:%s init finished", name)
		return nil
	}
}

func SimpleTransform(ctx transformer.TransformContext) {
	logger := ctx.GetLogger()
	logger.Infof("start to transform using %s", name)

	properties := ctx.GetProperties()
	op := properties["op"]
	errorIf := properties["error"]

	input := ctx.GetInput()

	r := &record{}
	json.Unmarshal(input, r)

	if op == "+" || op == "-" {
		if op == "+" {
			r.C = r.A + r.B
		} else {
			r.C = r.A - r.B
		}
		output, _ := json.Marshal(r)
		ctx.SetOutput(output)
		ctx.SetTransformStatus(transformer.TransformStatusAccept)
	} else {
		ctx.SetError(UnsupportedOpError)
		if errorIf == "reject" {
			ctx.SetTransformStatus(transformer.TransformStatusReject)
		} else {
			ctx.SetTransformStatus(transformer.TransformStatusError)
		}
	}
	logger.Infof("finished to transform using %s", name)
}

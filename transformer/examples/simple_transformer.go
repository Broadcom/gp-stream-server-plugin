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
	"github.com/greenplum-db/gp-stream-server-plugin/transformer"
)

var UnsupportedOpError = errors.New("simple transformer error: unsupported op")

type record struct {
	A int
	B int
	C int
}

func SimpleTransformOnInit() error {
	return nil
}

func SimpleTransform(ctx transformer.TransformContext) {
	logger := ctx.GetLogger()
	logger.Infof("start to transform using adder.")

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
	logger.Infof("finished to transform using adder")
	return
}

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

// SimpleTransformOnInit is the init function of this plugin, it is only invoked once when plugin is first loaded.
func SimpleTransformOnInit(ctx common.BaseContext) error {
	// get property specified in job config file using ctx.GetProperties()
	properties := ctx.GetProperties()
	name = properties["name"]
	// use ctx.GetLogger() to print your log
	ctx.GetLogger().Infof("plugin:%s init start...", name)

	simulateInitError := properties["simulate-init-error"]
	if simulateInitError == "true" {
		// simulate an error
		ctx.GetLogger().Errorf("plugin:%s init error", name)
		// return an error if plugin init failed
		return transformer.ErrorTransformerInit
	} else {
		// do some init work...
		ctx.GetLogger().Infof("plugin:%s init finished", name)
		return nil
	}
}

// SimpleTransform is the transform function of this plugin, it is invoked after each message is consumed.
func SimpleTransform(ctx transformer.TransformContext) {
	logger := ctx.GetLogger()
	logger.Infof("start to transform using %s", name)

	// get property specified in job config file using ctx.GetProperties()
	properties := ctx.GetProperties()
	op := properties["op"]

	// get input data from ctx.GetInput(), it's a raw kafka/rabbitmq message body
	input := ctx.GetInput()

	r := &record{}
	if err := json.Unmarshal(input, r); err != nil {
		ctx.SetError(err)
		errorPolicy := properties["error-policy"]
		if errorPolicy == "reject" {
			// set transform status as rejected
			ctx.SetTransformStatus(transformer.TransformStatusReject)
		} else {
			// set transform status as error
			ctx.SetTransformStatus(transformer.TransformStatusError)
		}
		return
	}

	// let's perform some transformation logic according to properties["op"] config...
	if op == "+" || op == "-" {
		if op == "+" {
			r.C = r.A + r.B
		} else {
			r.C = r.A - r.B
		}
		output, _ := json.Marshal(r)
		// set transform result back
		ctx.SetOutput(output)
		// set transform status as accepted
		ctx.SetTransformStatus(transformer.TransformStatusAccept)
	} else {
		ctx.SetError(UnsupportedOpError)
		// set transform status as error
		ctx.SetTransformStatus(transformer.TransformStatusError)
	}
	logger.Infof("finished to transform using %s", name)
}

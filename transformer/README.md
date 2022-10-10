# Introduction

GPSS 1.9.0 introduces a "Transformer" type plugin for Kafka and RabbitMQ job. Transformer gives customer a chance to inject some custom logic right after messages were read from Kafka/RabbitMQ and before they were sent to downstream. We'll give an example to show how to write your own transformer plugin in the following sections. The source code of example is in [examples](./examples) directory.

# How to write and use your own transformer

In this section, we'll use an example to show how to write your own transformer plugin. Our example transformer is [simple_transformer.go](./examples/simple_transformer.go). The implementation is really simple, it takes json data like `{"A":2,"B":1}` as input and takes an operator from properties.If operator=="+", it puts C=A+B back, if operator=='-', it puts C=A-B back.

## 1 Implement `init` and `transform` function

There are two main functions to be implemented. We call them `init` function and `transform` function,  this doesn't mean the function must have name "init" or "transform", you can name your function according to your needs.

###  1.1 init function

Init function will be invoked when the transformer plugin is first loaded. It's guaranteed to be invoked only once even if the plugin is used in multi jobs.

Init function should have type `func(ctx common.BaseContext) error`, for example, in [simple_transformer.go](./examples/simple_transformer.go) init function is `func SimpleTransformOnInit(ctx common.BaseContext) error`.

Init function has one parameter with type `common.BaseContext`,you can get logger or property from it. It's defined in [common.go](../common.go).

If any error is occurred during init, return the error to the function, otherwise just return nil.

### 1.2 transform function

Transform function will be invoked after each message is read from Kafka or RabbitMQ and before the message is sent to downstream in GPSS job process pipeline.

Transform function should have type `func(ctx transformer.TransformContext)`, for example, in [simple_transformer.go](./examples/simple_transformer.go) transform function is `func SimpleTransform(ctx transformer.TransformContext)`.

Transform function has one parameter with type `transformer.TransformContext`.It's defined in [transformer.go](./transformer.go).

Here is an explanation for functions of `transformer.TransformContext`:

| Function                            | Description                                                                             |
|-------------------------------------|-----------------------------------------------------------------------------------------|
| GetLogger() logrus.FieldLogger      | you should use this logger to print your logs and you can find logs in GPSS log files   |
| GetProperties() map[string]string   | return plugin property defined in job yml config file with this function                |
| GetInput() []byte                   | return data to be transformed (Kafka or RabbitMQ message body)by the transformer        |
| SetTransformStatus(TransformStatus) | set the transform status                                                                |
| SetOutput([]byte)                   | set the transform result                                                                |
| SetError(error)                     | if any error is occurred during transform,use this function to report the error to GPSS |

Here is an explanation for `TransformStatus`

| Value                 | Description                                                                                                                                    |
|-----------------------|------------------------------------------------------------------------------------------------------------------------------------------------|
| TransformStatusAccept | TransformStatusAccept means transform success                                                                                                  |
| TransformStatusReject | TransformStatusReject means transform failed, message will be dropped, but GPSS will ignore the error and continue to process the next message |
| TransformStatusError  | TransformStatusError means transform failed, GPSS will stop the job because of the error                                                       |

## 2 Build plugin

We should use buildmode=plugin to build a plugin, for example:

```
go build -buildmode=plugin -o simple-transformer.so simple_transformer.go
```

## 3 Add plugin config to Kafka/RabbitMQ job

We should add config to yml job config file, for example:

```
  KAFKA:
    INPUT:
      SOURCE:
        BROKERS: localhost:9092
        TOPIC: test
      TRANSFORMER:
        PATH: /tmp/simple-transformer.so
        ON_INIT: SimpleTransformOnInit
        TRANSFORM: SimpleTransform
        PROPERTIES:
          name: my-simple-transformer
          simulate-init-error: false
          op: +
          error-policy: reject
```

we add a `TRANSFORMER` section, it has 4 fields:

| Field      | Description                                                                                                                                          |
|------------|------------------------------------------------------------------------------------------------------------------------------------------------------|
| PATH       | the path to plugin dynamically linked file, if the path is not existed, an error will be reported in job preparation period                          |
| ON_INIT    | init function name                                                                                                                                   |
| TRANSFORM  | transform function name                                                                                                                              |
| PROPERTIES | some custom key-value properties, when `ON_INIT` func or `TRANSFORM` func is invoked, `PROPERTIES` are passed to the function together transparently |


# Introduction

The Greenplum Stream Server (GPSS) is an ETL (extract, transform, load) tool for [greenplum-db](https://github.com/greenplum-db/gpdb).

GPSS 1.9.0 introduces a "kernel-plugin" mode based on [go native plugin](https://pkg.go.dev/plugin) mechanism to help customer extend their custom logic when loading data. GPSS acts as the kernel,it is responsible for loading plugin and invoking plugin function at specific stage during the job lifecycle. Customers write their own plugins and build them with `go build -buildmode=plugin` to a dynamically linked file, then specify plugin config in job config file.

This repository helps you implement your own GPSS plugin. It contains some common interfaces and constants used to implement a plugin, and also some simple examples.

Currently the following plugin types are supported:

- Transformer

# Known issues and Limitations

## Go version

The plugin should be built with same go version as greenplum stream server, we use go1.17.6 now.

## Dependency version

If a dependency that your plugin uses is also used by the greenplum stream server, their versions should be same.

This is the dependencies of greenplum stream server:

```
require (
   github.com/blang/semver v3.5.1+incompatible
   github.com/confluentinc/confluent-kafka-go v0.11.5-0.20180718091254-d2ab44a05b46
   github.com/elodina/go-avro v0.0.0-20160406082632-0c8185d9a3ba
   github.com/elodina/go-kafka-avro v0.0.0-20160422130714-ab6b1d1c9a23
   github.com/golang/mock v1.3.1
   github.com/golang/protobuf v1.5.2
   github.com/greenplum-db/gp-common-go-libs v0.0.0-20190801065458-13e3e4bfe0e4
   github.com/greenplum-db/pq v1.1.2-0.20200902100711-efba7a31999d
   github.com/jmoiron/sqlx v1.2.1-0.20190826204134-d7d95172beb5
   github.com/onsi/ginkgo v1.16.4
   github.com/onsi/gomega v1.20.0
   github.com/pkg/errors v0.9.1
   github.com/prometheus/client_golang v1.5.1
   github.com/sirupsen/logrus v1.6.0
   github.com/spf13/cobra v1.5.0
   github.com/spf13/pflag v1.0.5
   github.com/spf13/viper v1.7.0
   github.com/stretchr/testify v1.5.1
   github.com/valyala/fasthttp v1.2.0
   google.golang.org/grpc v1.28.0
   gopkg.in/yaml.v2 v2.4.0
)

require (
   github.com/beorn7/perks v1.0.1 // indirect
   github.com/cespare/xxhash/v2 v2.1.1 // indirect
   github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
   github.com/davecgh/go-spew v1.1.1 // indirect
   github.com/fsnotify/fsnotify v1.4.9 // indirect
   github.com/golang/snappy v0.0.4 // indirect
   github.com/google/go-cmp v0.5.8 // indirect
   github.com/greenplum-db/gp-stream-server-plugin v0.0.0-20220928072037-c4ed44e1ed5e // indirect
   github.com/hashicorp/hcl v1.0.0 // indirect
   github.com/inconshreveable/mousetrap v1.0.0 // indirect
   github.com/klauspost/compress v1.15.9 // indirect
   github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
   github.com/kr/pretty v0.3.0 // indirect
   github.com/magiconair/properties v1.8.1 // indirect
   github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
   github.com/mitchellh/mapstructure v1.1.2 // indirect
   github.com/nxadm/tail v1.4.8 // indirect
   github.com/pelletier/go-toml v1.2.0 // indirect
   github.com/pierrec/lz4 v2.6.1+incompatible // indirect
   github.com/pmezard/go-difflib v1.0.0 // indirect
   github.com/prometheus/client_model v0.2.0 // indirect
   github.com/russross/blackfriday/v2 v2.1.0 // indirect
   github.com/spf13/afero v1.1.2 // indirect
   github.com/spf13/cast v1.3.0 // indirect
   github.com/spf13/jwalterweatherman v1.0.0 // indirect
   github.com/subosito/gotenv v1.2.0 // indirect
   github.com/valyala/bytebufferpool v1.0.0 // indirect
   golang.org/x/net v0.0.0-20220526153639-5463443f8c37 // indirect
   golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
   golang.org/x/text v0.3.7 // indirect
   google.golang.org/genproto v0.0.0-20191108220845-16a3f7862a1a // indirect
   google.golang.org/protobuf v1.28.0 // indirect
   gopkg.in/ini.v1 v1.51.0 // indirect
   gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
   gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
   github.com/greenplum-db/gssapi v0.0.0-20190813095309-58507edf9144 // indirect
   github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
   github.com/prometheus/common v0.10.0 // indirect
   github.com/prometheus/procfs v0.1.3 // indirect
   github.com/rabbitmq/amqp091-go v1.3.4
   github.com/rabbitmq/rabbitmq-stream-go-client v1.0.1-rc.2
   gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0 // indirect
   gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)
```
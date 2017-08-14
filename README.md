# schemarshal [![wercker status](https://app.wercker.com/status/ebb1f8ec249177acd0d47bd8a6a59dd2/s/master "wercker status")](https://app.wercker.com/project/byKey/ebb1f8ec249177acd0d47bd8a6a59dd2)

[![Coverage Status](https://coveralls.io/repos/github/aaharu/schemarshal/badge.svg)](https://coveralls.io/github/aaharu/schemarshal)
[![Go Report Card](https://goreportcard.com/badge/github.com/aaharu/schemarshal)](https://goreportcard.com/report/github.com/aaharu/schemarshal)

Generates Go struct types from a [JSON Schema](http://json-schema.org/).

## Installation

```bash
go get -u github.com/aaharu/schemarshal
```

## Usage

```
SYNOPSIS
  schemarshal [options] [<json_schema_file>]
OPTIONS
  -h, -help
           Show this help message.
  -f <file>, -file <file>
           Input file name.
  -o <file>, -output <file>
           Write output to file instead of stdout.
  -p <package>, -package <package>
           Package name for output. (default `main`)
  -t <package>, -type <package>
           Set default Type name.
  -v, -version
           Show version.
```

```bash
# with args
schemarshal -o api/schema/gen.go -p gen -f schema.json

# pipe
curl -s "https://raw.githubusercontent.com/aaharu/schemarshal/master/test_data/disk.json" | schemarshal
```

## TODO

- [ ] use go/ast

## [Examples](examples.md)

- [example a.json](examples.md#a.json)
- [example qiita schema](examples.md#qiita-v2-schema)

## Dependencies

* https://github.com/lestrrat/go-jsschema
  - JSON Schema parser
  - MIT License
* golang.org/x/crypto/ssh/terminal
  - https://godoc.org/golang.org/x/crypto/ssh/terminal
  - https://golang.org/LICENSE

## Similar Projects

* https://github.com/idubinskiy/schematyper
* https://github.com/dameleon/structr
* https://github.com/interagent/schematic

## License

BSD-2-Clause

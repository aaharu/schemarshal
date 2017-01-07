# schemarshal [![wercker status](https://app.wercker.com/status/ebb1f8ec249177acd0d47bd8a6a59dd2/s/master "wercker status")](https://app.wercker.com/project/byKey/ebb1f8ec249177acd0d47bd8a6a59dd2)

[![Coverage Status](https://coveralls.io/repos/github/aaharu/schemarshal/badge.svg)](https://coveralls.io/github/aaharu/schemarshal)

Generates Go struct types from a [JSON Schema](http://json-schema.org/).

## Installation

```bash
go get github.com/aaharu/schemarshal
```

## Usage

```
SYNOPSIS
  schemarshal [options] [<json_schema_file>]
OPTIONS
  -h, -help
           Show this help message.
  -o <file>, -output <file>
           Write output to file instead of stdout.
  -p <package>, -package <package>
           Package name for output. (default `main`)
  -t <package>, -type <package>
           Set default Type name.
  -v, -version
           Show version.
```

## TODO

- [ ] unit tests
- [ ] continuous integration , wercker?
- [ ] refactoring
- [ ] write doc

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

# schemarshal

Generates Go struct types from a [JSON Schema](http://json-schema.org/).

## Installation

```bash
go get github.com/aaharu/schemarshal/cmd/schemarshal
```

## Usage

```
SYNOPSIS
  schemarshal [options] <json_shcema_file>
OPTIONS
  -h, -help
           Show help message.
  -o <file>, -output <file>
           Write output to <file> instead of stdout.
  -p <package>, -package <package>
           Package name for output. (default `main`)
  -version
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

## Similar Projects

* https://github.com/idubinskiy/schematyper
* https://github.com/dameleon/structr
* https://github.com/interagent/schematic

## License

BSD-2-Clause

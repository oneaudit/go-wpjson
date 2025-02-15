<h1 align="center">
  GO WPJSON
</h1>
<h4 align="center">A toolkit to parse WordPress Rest API specification.</h4>

## WPJSON As Library 📚

```
```

## WPJSON CLI Usage 🤖

WPJson requires **Go 1.22+** to install successfully.

```console
CGO_ENABLED=1 go install github.com/oneaudit/go-wpjson/cmd/go-wpjson@latest
wpjson-go -h
```

This will display help for the tool. Here are all the switches it supports.

```
A toolkit to parse WordPress Rest API specification.

Usage:
  wpjson_go [flags]

Flags:
TARGET:
   -t, -target string  target input file or URL to parse

CONFIGURATION:
   -config string  path to the katana-ng configuration file

OUTPUT:
   -o, -output string  output file to save the results
   -silent             display output only
   -v, -verbose        display verbose output
   -debug              display debug output
   -version            display project version
```
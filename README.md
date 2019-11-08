# glearn-cli

## Installation
Install glearn-cli using `go install`:

```
go install github.com/Galvanize-IT/glearn-cli
```

- We will soon be adding other installation methods that don't require
golang to be installed

## Development
Build
```
go build -o glearn-cli main.go
```

Run
```
./glearn-cli [command...] [flag...]
```

Or for quicker iterations:
```
go run main.go [command...] [flag...]
```

## Examples

Setting your API token:
```
glearn-cli settoken my_neat_token_123_456
```

Creating new:
```
glearn-cli new
```

Building:
```
glearn-cli build
```

Publishing:
```
glearn-cli publish
```
# glearn-cli

## Installation
Install glearn-cli using `go install`:

```
go install github.com/Galvanize-IT/glearn-cli
```

- We will soon be adding other installation methods that don't require
golang to be installed

Make sure to set your config variables in `~/.glearn-config.yaml`. You can do this by either manually
editing the file:
```
api_token: YOUR_API_TOKEN
aws_access_key_id: S3_USER_ACCESS_KEY_ID
aws_secret_access_key: S3_USER_SECRET_ACCESS_KEY
```

Or by using the set commands:
```
glearn-cli setapitoken [token]
```

```
glearn-cli setawsaccesskeyid [access_key_id]
```

```
glearn-cli setawssecretaccesskey [secret_access_key]
```

## Development
Add a `.env` with the variables set from the `.env.example`

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

# Releases

Create/add a github token with `repo` access. This gives you the ability to push releases and their binaries.

For test release run:
```
goreleaser --snapshot --skip-publish --rm-dist
```

To release run:
```
GITHUB_TOKEN=your_githhub_token \
    AWS_S3_BUCKET=bucket_from_env \
    AWS_KEY_PREFIX=key_prefix_from_env \
    goreleaser release
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
# glearn-cli

## Installation

Make sure to set your config variables in `~/.glearn-config.yaml`. You can do this by either manually
editing the file:
```
api_token: YOUR_API_TOKEN
aws_access_key_id: S3_USER_ACCESS_KEY_ID
aws_secret_access_key: S3_USER_SECRET_ACCESS_KEY
aws_s3_bucket: S3_BUCKET_NAME
aws_s3_key_prefix: S3_BUCKET_KEY_PREFIX
```

Or by using the set commands:
```
glearn-cli set [...flags]
```

_**Option A:**_

```
brew tap Galvanize-IT/glearn-cli
brew install Galvanize-IT/glearn-cli/glearn-cli
```

_**Option B:**_

If you mosey on over to [releases](https://github.com/Galvanize-IT/glearn-cli/releases), you'll find binaries for darwin, linux, windows, and amd64. You can download directly from there.

_**Option C:**_

If you have Go installed on your machine, use `go install`:

```
go install github.com/bradford-hamilton/monkey-lang
```

This will place the binary in your `go/bin` and is ready to use.

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

Create a new semantic version tag (ex. 0.1.0)
```
git tag -a v{semantic_version} -m "Some new release commit"
```

Push new tag
```
git push origin v{semantic_version}
```

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
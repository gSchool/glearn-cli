# glearn-cli

## Installation

Make sure to set your `api_token` in `~/.glearn-config.yaml`. You can do this by either manually editing the file:
```
api_token: YOUR_API_TOKEN
```

Or by using the set command:
```
glearn set --api_token=neat_token_123
```

_**Option A:**_

Create/add a github token with full `repo` access. If you are logged into github, you can go here: [your tokens](https://github.com/settings/tokens) to add one. This gives you the ability to use brew to manage glearn.

```
HOMEBREW_GITHUB_API_TOKEN=YOUR_TOKEN brew tap Galvanize-IT/glearn
HOMEBREW_GITHUB_API_TOKEN=YOUR_TOKEN brew install Galvanize-IT/glearn/glearn
```

_**Option B:**_

If you mosey on over to [releases](https://github.com/Galvanize-IT/glearn-cli/releases), you'll find binaries for darwin, linux, windows, and amd64. You can download directly from there.

_**Option C:**_

If you have Go installed on your machine, use `go install`:

```
go install github.com/Galvanize-IT/glearn-cli
```

This will place the binary in your `go/bin` and is ready to use.

## Development
Build
```
go build -o glearn-cli main.go
```

Run
```
./glearn [commands...] [flags...]
```

Or for quicker iterations:
```
go run main.go [commands...] [flags...]
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
GITHUB_TOKEN=your_githhub_token goreleaser release
```

## Examples

Setting your API token:
```
glearn set --api_token=neat_token_123
```

Preview a `test_curriculum` directory:
```
glearn preview test_curriculum
```

Creating new:
```
glearn new
```

Building:
```
glearn build
```

Publishing:
```
glearn publish
```

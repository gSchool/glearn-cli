# glearn-cli

# Installation

Make sure to set your `api_token` in `~/.glearn-config.yaml`. You can do this by either manually editing the file:
```
api_token: YOUR_API_TOKEN
```

Or by using the set command:
```
learn set --api_token=my_neat_token_123
```

_**Option A:**_

Homebrew:
```
brew tap gSchool/learn
brew install gSchool/learn/learn
```

Uninstalling:
```
brew uninstall gSchool/learn/learn
```

_**Option B:**_

If you mosey on over to [releases](https://github.com/gSchool/glearn-cli/releases), you'll find binaries for darwin, linux, windows, and amd64. You can download directly from there.

_**Option C:**_

If you have Go installed on your machine, use `go install`:

```
go install github.com/gSchool/glearn-cli
```

This will place the binary in your `go/bin` and is ready to use.

# Development
Build
```
go build -o glearn-cli main.go
```

Run
```
./learn [commands...] [flags...]
```

Or for quicker iterations:
```
go run main.go [commands...] [flags...]
```

### Specifying Learn App URL

By default, the CLI tool will use Learn's base url `https://learn-2.galvanize.com`. This value can be changed by exporting the environment variable `LEARN_BASE_URL` to specify the desired address. This is convenient for testing stage/PR environments.

# Releases

Create a github token with `repo` access. This gives you the ability to push releases and their binaries and allows glearn-cli write commits when necessary.

Create a new semantic version tag (ex. 0.1.0)
```
git tag -a v0.1.0 -m "Some new release commit"
```

Push new tag
```
git push origin v0.1.0
```

To release run:
```
GITHUB_TOKEN=your_githhub_token goreleaser release --rm-dist
```

## Examples

Setting your API token:
```
learn set --api_token=neat_token_123
```

Preview a `test_curriculum` directory:
```
learn preview test_curriculum
```

Publishing:
```
learn publish test_curriculum
```

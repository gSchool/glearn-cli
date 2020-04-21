# learn-cli

This is the command line interface for developing, previewing, and publishing curriculum on Learn.

## Installation with Homebrew

### Install
```
brew tap gSchool/learn
```
```
brew install learn
```
```
learn set --api_token=YOUR_LEARN_API_TOKEN
```

You can get your Learn API token from https://learn-2.galvanize.com/api_token

### Get Started

Run

```
learn help
```

### Update
```
brew upgrade learn
```

### Uninstall
```
brew uninstall learn
```

## Alternatives to Homebrew

### Install

Option: Use curl on Mac
```
curl -L $(curl -s https://api.github.com/repos/gSchool/glearn-cli/releases/latest | grep -o "http.*Darwin_x86_64.tar.gz") | tar -xzf - -C /usr/local/bin
```

Option: Use curl on Linux
```
curl -L $(curl -s https://api.github.com/repos/gSchool/glearn-cli/releases/latest | grep -o "http.*Linux_x86_64.tar.gz") | tar -xzf - -C /usr/local/bin
```

Option: Download binaries for all platforms directly from
https://github.com/gSchool/glearn-cli/releases

After using any of these options, set your API token with
```
learn set --api_token=YOUR_LEARN_API_TOKEN
```

You can get your Learn API token from https://learn-2.galvanize.com/api_token

### Get Started

Create a temp directory somewhere
```
mkdir test-content && cd test-content
```

Then run
```
learn new
```

## Example Usage

See a list of commands
```
learn help
```

Preview a single file
```
learn preview my_file.md
```

Preview an entire directory:
```
learn preview my_curriculum_directory
```

Publishing an entire repo
* add/commit/push to github
* if block doesn't exist, create and publish new block
* if block exists, update existing block and create new release
```
learn publish
```

## Development
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

## Releases

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

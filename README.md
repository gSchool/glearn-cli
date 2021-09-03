# learn-cli

This is the command line interface for developing, previewing, and publishing curriculum on Learn.

## Option 1: Installation with Homebrew (using WSL2 on Windows)

### Install
```
brew tap gSchool/learn
```
```
brew install learn
```

### Set your API token
```
learn set --api_token=[Your API Token from https://learn-2.galvanize.com/api_token]
```

## Option 2: Install binaries directly from Github

### Download

Download binaries for all platforms directly from
https://github.com/gSchool/glearn-cli/releases

Place in an appropriate location and include in your system Path so that commands can be called from any directory.

### Set your API token
```
learn set --api_token=[Your API Token from https://learn-2.galvanize.com/api_token]
```

## Option 3: Install with curl (using WSL2 on Windows)

### Curl commands

Mac users
```
curl -L $(curl -s https://api.github.com/repos/gSchool/glearn-cli/releases/latest | grep -o "http.*Darwin_x86_64.tar.gz") | tar -xzf - -C /usr/local/bin
```

Linux & WSL2 users
```
curl -L $(curl -s https://api.github.com/repos/gSchool/glearn-cli/releases/latest | grep -o "http.*Linux_x86_64.tar.gz") | tar -xzf - -C /usr/local/bin
```

Both of these commands should place the binary in a location that is already covered by your system Path.

### Set your API token
```
learn set --api_token=[Your API Token from https://learn-2.galvanize.com/api_token]
```

## Get Started: Walkthrough

Visit https://galvanize-learn.zendesk.com/hc/en-us/articles/1500000930401-Introduction for a short walkthrough to publish your first curriculum.

## Get Started: Quick Edits to Existing Curriculum

1. Clone and edit curriculum
2. Preview your changes. Run:
    `learn preview -o <directory|file>`
3. Git add / commit / push changes to the master branch
4. Publish changes for any cohort in Learn. Run:
    `learn publish`

## Help with other commands

```
learn help
```

## Update
Depending on how you installed above--

Homebrew: `brew upgrade learn`

Binary download: Download new binaries from https://github.com/gSchool/glearn-cli/releases

Curl: Run curl commands listed under "install" again.

## Uninstall

Homebrew: `brew uninstall learn`

Other users: delete binary

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

Create a new semantic version tag (ex. v0.1.0)
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

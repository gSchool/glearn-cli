# Developing the `learn` CLI

Build
```
go build -o glearn main.go
```

Run
```
./glearn [commands...] [flags...]
```

Or for quicker iterations:
```
go run main.go [commands...] [flags...]
```

# Specifying Learn App URL

By default, the CLI tool will use Learn's base url `https://learn-2.galvanize.com`. This value can be changed by exporting the environment variable `LEARN_BASE_URL` to specify the desired address. This is convenient for testing stage/PR environments.

# Releases

Create a github token with `repo` access. This gives you the ability to push releases and their binaries and allows `glearn-cli` write commits when necessary.

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
GITHUB_TOKEN=<your_githhub_token> ./release-new-version
```

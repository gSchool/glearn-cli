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

## Adding a new command Markdown subcommand

The Markdown subcommands are broken out into groups. Identify which group your
subcommand should belong in:

* [files](./app/cmd/markdown/files/) - Commands that generate files
* [questions](./app/cmd/markdown/questions/) - Commands that generate the markup
  for questions/challenges
* [yaml](./app/cmd/markdown/yaml/) - Commands that generate YAML configuration
  files
* [others](./app/cmd/markdown/others/) - The place where other commands go

Each of those directories have a `root.go` file in it that provides the factory
method for the associated subcommands. The [markdown](./app/cmd/markdown/)
directory has its own `root.go` file to orchestrate the creation of each of the
commands with the main `learn` command.

To create a new subcommand, follow the pattern that you find in another file in
the same directory. The pattern is declarative with the majority of the code
existing in the associated `root.go` file.

Then, add a new test case in the `XXX_test.go` file in the same directory.

Finally, register the new subcommand in the appropriate `addXXXCommands`
function in [markdown/root.go](./app/cmd/markdown/root.go).

## Specifying Learn App URL

By default, the CLI tool will use Learn's base url `https://learn-2.galvanize.com`. This value can be changed by exporting the environment variable `LEARN_BASE_URL` to specify the desired address. This is convenient for testing stage/PR environments.

## Releases

Create a github token with `repo` access. This gives you the ability to push releases and their binaries and allows `glearn-cli` write commits when necessary.

Create a new semantic version tag (ex. v0.1.0)

```
git tag -a v0.10.11 -m "Some new release commit"
```

Push new tag
```
git push origin v0.10.11
```

To release run:
```
GITHUB_TOKEN=<your_githhub_token> ./release-new-version
```

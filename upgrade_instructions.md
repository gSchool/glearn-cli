# Update Learn CLI

Depending on how you installed your cli tool, please use the appropriate method to upgrade below:

## Homebrew

```
brew tap gSchool/learn
brew update
brew upgrade learn
```

## Curl

```
curl -sSL $(curl -sSL https://api.github.com/repos/gSchool/glearn-cli/releases/latest | grep -o "http.*$(uname -sm | sed 's/ /_/').tar.gz") | sudo tar -C /usr/local/bin -xzf - learn
```

## Binary download

Download a new version from https://github.com/gSchool/glearn-cli/releases/latest

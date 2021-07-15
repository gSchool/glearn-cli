# Update Learn CLI
Depending on how you installed your cli tool, please use the appropriate method to upgrade below:

## Homebrew
Homebrew: `brew tap gSchool/learn && brew upgrade learn`

## Binary download
Binary download: Download new binaries from https://github.com/gSchool/glearn-cli/releases

## Curl 
Mac users
```
curl -L $(curl -s https://api.github.com/repos/gSchool/glearn-cli/releases/latest | grep -o "http.*Darwin_x86_64.tar.gz") | tar -xzf - -C /usr/local/bin
```

Linux users
```
curl -L $(curl -s https://api.github.com/repos/gSchool/glearn-cli/releases/latest | grep -o "http.*Linux_x86_64.tar.gz") | tar -xzf - -C /usr/local/bin

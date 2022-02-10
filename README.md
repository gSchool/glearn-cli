# Learn CLI

The Learn CLI is the command line interface for developing, previewing, and publishing curriculum on Learn.

## Requirements

The learn command is supported on Mac, Linux, and Windows.

## Installation

### Option 1: Install with Homebrew (use WSL2 on Windows)

The easiest installation is with Homebrew, which is for MacOS and Linux. Follow these directions to install Homebrew. Once installed, run these commands on your command line.

```
brew tap gSchool/learn
brew install learn
```

### Option 2: Install with curl (use WSL2 on Windows)

Use the command line utility `curl` to download and install the latest version. The learn command will be placed in the `/usr/local/bin` directory.

```
curl -sSL $(curl -sSL https://api.github.com/repos/gSchool/glearn-cli/releases/latest | grep -o "http.*$(uname -sm | sed 's/ /_/').tar.gz") | tar -C /usr/local/bin -xzf - learn
```

### Option 3: Install binaries directly from GitHub

Download binaries for all platforms directly from https://github.com/gSchool/glearn-cli/releases

Place the `learn` executable in a location included in your `PATH` so that it can be called from any directory.

## Set API Token

After installation, you must set your API token. Copy your token from https://learn-2.galvanize.com/api_token and run this command, replacing YOUR_LEARN_API_TOKEN with your token.

```
learn set --api_token=YOUR_LEARN_API_TOKEN
```

## Confirm Installation
Run the command

```
learn version
```

You should get a response like

```
v0.10.0
```

## Help with other commands

```
learn help
```

## Get Started: Walkthrough

Visit [Guru](https://app.getguru.com/boards/iEdB57dT/Creating-Content-in-Learn) for a short walkthrough to publish your first curriculum.

### Get Started: Quick Edits to Existing Curriculum

1. Clone and edit curriculum
2. Preview your changes. Run:
    `learn preview -o <directory|file>`
3. Git add / commit / push changes to any branch
4. Publish changes as a new Block revision. Run:
    `learn publish`

## Update

Follow the instructions in the [upgrade](./upgrade_instructions.md) document.

## Uninstall

Homebrew: `brew uninstall learn`

Other installations: delete `learn` executable

## Development

To contribute, look at the [development](./development_instructions.md) instructions.

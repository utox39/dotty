# dotty

- [Description](#description)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
    - [Configuration](#configuration)
- [Roadmap](#roadmap)
- [Contributing](#contributing)

## Description

Backup your dotfiles of choice in a folder

## Requirements

- [Go](https://go.dev/)

## Installation

```bash
# Clone the repo
$ git clone ...

# cd to the repo
$ cd path/to/dotty

# Build dotty
$ go build dotty.go

# Then move it somewhere in your $PATH. Here is an example:
$ mv dotty ~/bin

# Create the config folder
$ mkdir -p ~/.config/dotty

# Copy the config file
$ cp path/to/dotty/config.json ~/.config/dotty/config.json
```

## Usage

```bash
$ dotty
```
### Configuration

To add the paths to your dotfiles go to the configuration file
in ~/.config/dotty/config.json

> N.B Always remember to use the absolute path to the file

```json
{
  "dotfiles" : [
    "~/.example", "~/Documents/.example2"
  ]
}
```

## Roadmap

- Improve checks
- Publish it as a package
- Automatic creation of the configuration file
- Ability to add a file via command line command

## Contributing

If you would like to contribute to this project just create a pull request which I will try to review as soon as
possible.

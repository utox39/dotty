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
$ git clone https://github.com/utox39/dotty.git

# cd to the repo
$ cd path/to/dotty

# Build dotty
$ go build ./...

# Then move it somewhere in your $PATH. Here is an example:
$ mv dotty ~/bin

# Create the config folder
$ mkdir -p ~/.config/dotty

# Copy the config file
$ cp path/to/dotty/config.json ~/.config/dotty/config.json
```

## Usage

#### Backup the dotfiles

```bash
$ dotty backup
```

#### Add new dotfile

```bash
$ dotty add .foo
```

### Configuration

To add the dotfile paths and destination path go to the configuration file
in ~/.config/dotty/config.json

> N.B Always remember to use the absolute path

```json
{
  "dotfiles" : [
    "~/.example", "~/Documents/.example2"
  ],
  "destination-path" : "~/dotfiles/"
}
```

## Roadmap

- Improve customization
- Improve checks
- Publish it as a package
- Automatic creation of the configuration file

## Contributing

If you would like to contribute to this project just create a pull request which I will try to review as soon as
possible.

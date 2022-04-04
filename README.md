# GOG - Opinionated Git

GOG is an Opinionated Git workflow CLI tool written in Golang to help developers enforce consistency and strong project/version management within their projects without the added complexity of manual process documentation. Rather GOG allows you to focus on developing features and bringing value to stakeholders rather than managing versions, dependencies and CHANGELOGs.

## Installation

### Linux / MAC OSX

```bash

curl -sSL https://raw.githubusercontent.com/SystemFiles/GOG/master/INSTALL.sh | bash /dev/stdin

```

This will install the latest version of GOG on a Linux or OSX target.

### Windows

TBD

## Basic Usage

```bash

gog (feature | push | finish) [options ...] [-h] [-help]

```

## Testing

Currently GOG has been tested on the following deployment targets:

- Linux (AMD64)
- Darwin (AMD64)
- Windows (AMD64) - Some functionality limited (#12)

## Credits

Credit for original concept goes to Daniel Waespe (STATCAN)

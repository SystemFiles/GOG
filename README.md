# GOG - Opinionated Git

GOG is an Opinionated Git workflow CLI tool written in Golang to help developers enforce consistency and strong project/version management within their projects without the added complexity of manual process documentation. Rather GOG allows you to focus on developing features and bringing value to stakeholders rather than managing versions, dependencies and CHANGELOGs.

![Demo](./.github/GOG-Demo.gif)

## Installation

### Linux / MAC OSX (darwin)

```bash

curl -sSL https://raw.githubusercontent.com/SystemFiles/GOG/master/INSTALL.sh | bash /dev/stdin

```

This will install the latest version of GOG on a Linux or OSX target.

### Windows

TBD

## Basic Usage

```bash

gog (feature(feat) | push(p) | finish(fin)) [options ...] [-h] [-help]

```

### Feature Creation

```bash

Usage: gog (feature | feat) <jira> <comment> [-from-feature] [-h] [-help]

-------====== Feature Arguments ======-------

jira
      specifies the JIRA issue we are working under
comment
      specifies a human-readable comment describing the issue/feature

------================================------

  -from-feature
      specifies if this feature will be based on the a current feature branch
  -prefix string
      optionally specifies a version prefix to use for this feature which will override existing prefix in global GOG config

-------================================-------

```

### Intermediate Changes

```bash

Usage: gog (push | p) [message] [-h] [-help]

-------====== Push Arguments ======-------

message
      specifies a commit message for this feature push

-------================================-------

```

### Feature Release

```bash

Usage: gog (finish | fin) (-major | -minor | -patch) [-h] [-help]

-------====== Finish Arguments ======-------

  -major
      specifies that in this freature you make incompatible API changes (breaking changes)
  -minor
      specifies that in this feature you add functionality in a backwards compatible manner (non-breaking)
  -patch
      specifies that in this feature you make backwards compatible bug fixes small backwards compatible updates

-------================================-------

```

### Simple Push (no feature attached)

While this does not fit into the opinionated workflow defined by the commands above, it is sometimes necessary to perform a simple push when collaborating on projects that do not exactly follow the workflow.

```bash

Usage: ../GOG/dist/gog (simple-push | sp) [message] [-h] [-help]

Simple-Push is a utility to allow non-feature related code pushes directly to the current remote branch. If used without a message one will be generated for you.

-------====== Simple-Push Arguments ======-------

message
    (optionally) specifies a commit message for this simple push operation

-------================================-------

```

## Updating GOG

Updating GOG (if on Darwin or Linux) can be done in-place using the `gog update` command.

```bash

Usage: ../GOG/dist/gog update [-tag TAG] [-h] [-help]

--------======= Tag Arguments =======--------

  -tag string
        specifies a specific version tag to use for update

-------================================-------

```

## Testing

Currently GOG has been tested on the following deployment targets:

- Linux (AMD64)
- Darwin (AMD64)
- Windows (AMD64) - Some functionality limited (#12)

## Credits

Credit for original concept goes to Daniel Waespe (STATCAN)

<p align="center">
  <img src="./static/images/gopher_full.png" alt="Gitrob" width="200" />
</p>

# Wraith: Putting the Open Source in OSINT
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mattyjones/gitrob)![GitHub release (latest by date)](https://img.shields.io/github/v/release/mattyjones/gitrob)![GitHub](https://img.shields.io/github/license/mattyjones/gitrob)

![Code Climate maintainability](https://img.shields.io/codeclimate/maintainability/mattyjones/gitrob)![Code Climate technical debt](https://img.shields.io/codeclimate/tech-debt/mattyjones/gitrob)![Code Climate issues](https://img.shields.io/codeclimate/issues/mattyjones/gitrob)

[![Build Status](https://travis-ci.org/mattyjones/gitrob.svg?branch=master)](https://travis-ci.org/mattyjones/gitrob)


Wraith is a tool to help find potentially sensitive information pushed to repositories on GitLab or Github. Wraith will clone repositories belonging to a user or group/organization down to a configurable depth and iterate through the commit history and flag files and/or commit content that match signatures for potentially sensitive information. The findings will be presented through a web interface for easy browsing and analysis.

## Features

- Scan the following sources:
  - Gitlab repositories
  - Github.com repositories
  - Local git repos
- Exclude files, paths, and extensions
- Web interface for real-time results
- Configurable commit depth
- Use environment variables, a config file, or flags
- Uses sub-commands for easier, more modular, functionality
- Clone a repo to memory instead of disk

This currently in beta, check the [roadmap][1] for planned functionality

## Usage

For a full list of use cases and configuration options use the included help functionality.

`gitrob --help`


## Configuration

**IMPORTANT** If you are targeting a GitLab group, please give the **group ID** as the target argument.  You can find the group ID just below the group name in the GitLab UI.  Otherwise, names with suffice for the target arguments. This id can be found on the group homepage.

There are multiple was to configure the tool for a scan. The easiest way is via commandline flags. To get a full list of available flags and their purpose use `gitrob <subcommand> --help`. This will pring out a list of flags and how they interact with the base scan. You can also set all flags as environment variables or use a static config file in YAML format. This config file can be used to store targets for multiple scan targets.

The order of precendence with each item taking precedence over the item below it is:

- explicit call to Set
- commandline flag
- environment variable
- configuration file
- key/value store
- default value

The various values are configured independently of each other so if you set all values in a config file, you can then override just the ones you want on the commandline. A sample config file looks like:

```yaml
---
commit-depth: 0
gitlab-targets:
    - codeemitter
    - mattyjones1
    - 8692959
silent: false
debug: true
gitlab-api-token: <token>
github-api-token: <token>
github-targets:
    - mattyjones
    - phantomSecrets
ignore-path:
    - cmd/
    - docs/
ignore-extension:
    - .go
    - .log
in-mem-clone: true
repo-dirs:
    - ../../../mattyjones/telegraf
```

## Examples

Scan a GitLab group assuming your access token has been added to the environment variable or a config file.  Look for file signature matches only:

    gitrob scanGitlab <gitlab_group_id>

Scan a multiple GitLab groups assuming your access token has been added to the environment variable or a config file.  Clone repositories into memory for faster analysis.  Set the scan mode to 2 to scan each file match for a content match before creating a result.:

    gitrob scanGitlab -in-mem-clone -mode 2  "<gitlab_group_id_1> <gitlab_group_id_2>"

Scan a GitLab groups assuming your access token has been added to the environment variable or a config file. Clone repositories into memory for faster analysis.  Set the scan mode to 3 to scan each commit for content matches only.:

    gitrob scanGitlab -in-mem-clone -mode 3 "<gitlab_group_id>"

Scan a Github user setting your Github access token as a parameter.  Clone repositories into memory for faster analysis.

    gitrob scangithub -github-access-token <token> -in-mem-clone "<github_user_name>"

### Editing File and Content Regular Expressions

Regular expressions are included in the [filesignatures.json](./rules/filesignatures.json) and [contentsignatures.json](./rules/contentsignatures.json) files respectively.  Edit these files to adjust your scope and fine-tune your results.

Gitrob will start its web interface and serve the results for analysis.

## Installation

At this stage the only option is to build from source from this repository.

To install from source, make sure you have a correctly configured **Go >= 1.14** environment and that `$GOPATH/bin` is in your `$PATH`.

    $ git clone git@gitlab.com:mattyjones1/gitrob.git
    $ cd ~/go/src/gitrob
    $ make build
    $ ./bin/gitrob-<ARCH> <sub-command>
    
In the future there will be binary releases of the code

## Access Tokens

Gitrob will need either a GitLab or Github access token in order to interact with the appropriate API.  You can create a [GitLab personal access token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html), or [a Github personal access token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/) and save it in an environment variable in your `.bashrc` or similar shell configuration file:

    export GITROB_GITLAB_ACCESS_TOKEN=deadbeefdeadbeefdeadbeefdeadbeefdeadbeef
    export GITROB_GITHUB_ACCESS_TOKEN=deadbeefdeadbeefdeadbeefdeadbeefdeadbeef

Alternatively you can specify the access token with the `-gitlab-access-token` or `-github-access-token` option on the command line, but watch out for your command history! A configuration file can also be used, an example is provided above.

[1]: docs/development/roadmap.md

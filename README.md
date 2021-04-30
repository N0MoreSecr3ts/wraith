<h1 align="center">
  <br>
    <img src="./static/images/gopher_full.png" alt="wraith" width="200"/>
  <br>
  Wraith
  <br>
</h1>

<h4 align="center">Finding digital secrets that were never meant to be found in all corners of the net.</h4>

<p align="center">
  <img alt="GitHub go.mod Go version (branch)" src="https://img.shields.io/github/go-mod/go-version/n0moresecr3ts/wraith/master?style=for-the-badge">
  <img alt="GitHub go.mod Go version (branch)" src="https://img.shields.io/github/go-mod/go-version/n0moresecr3ts/wraith/develop?style=for-the-badge">
  <img alt="GitHub release (latest SemVer)" src="https://img.shields.io/github/v/release/n0moresecr3ts/wraith?style=for-the-badge">
  <img alt="GitHub commits since latest release (by SemVer)" src="https://img.shields.io/github/commits-since/n0moresecr3ts/wraith/latest?style=for-the-badge">
<br>
  <img alt="GitHub issues by-label" src="https://img.shields.io/github/issues-raw/n0moresecr3ts/wraith/Bug?color=RED&label=BUGS&style=for-the-badge">
  <img alt="GitHub issues by-label" src="https://img.shields.io/github/issues-raw/n0moresecr3ts/wraith/Feature%20Request?color=38BED3&label=FEATURE%20REQUESTS&style=for-the-badge">
  <img alt="Travis (.org) branch" src="https://img.shields.io/travis/n0moresecr3ts/wraith/master?label=BUILD%20MASTER&style=for-the-badge">
  <img alt="Travis (.org) branch" src="https://img.shields.io/travis/n0moresecr3ts/wraith/develop?label=BUILD%20DEVELOP&style=for-the-badge">
<br>
  <img alt="Code Climate maintainability" src="https://img.shields.io/codeclimate/maintainability/N0MoreSecr3ts/wraith?style=for-the-badge">
  <img alt="Code Climate technical debt" src="https://img.shields.io/codeclimate/tech-debt/N0MoreSecr3ts/wraith?style=for-the-badge">
  <img alt="Code Climate issues" src="https://img.shields.io/codeclimate/issues/N0MoreSecr3ts/wraith?style=for-the-badge">
<br>
  <img alt="GitHub" src="https://img.shields.io/github/license/n0moresecr3ts/wraith?color=blue&style=for-the-badge">
  <img alt="GitHub All Releases" src="https://img.shields.io/github/downloads/n0moresecr3ts/wraith/total?style=for-the-badge">

</p>

<p align="center">
  <a href="#capabilities">Capabilities</a> •
  <a href="#screenshots">Screenshots</a> •
  <a href="#quickstart">Quickstart</a> •
  <a href="#documentation">Documentation</a> •
  <a href="#contributing">Contributing</a> •
  <a href="#credits">Credits</a> •
  <a href="#faq">FAQ</a> •
  <a href="#related">Related</a>
</p>

Wraith uncovers forgotten secrets and brings them back to life, haunting security and operations teams. It can be used to scan hosted and local git repos as well as local filesystems.

## Capabilities

### Targets
- Gitlab.com repositories and projects
- Github.com repositories and organizations
- Local git repositories
- Local filesystem

### Major Features

- Exclude files, paths, and extensions
- Web and terminal interfaces for real-time results (very much alpha)
- Configurable commit depth
- Built with [Viper][1] to manage environment variables, config files, or flags
- Uses [Cobra][2] sub-commands for easier, more modular, functionality

## Screenshots
<p>
  <img width="537" alt="Screen Shot 2020-08-16 at 11 23 25 PM" src="https://user-images.githubusercontent.com/672940/90354541-9f515a80-e017-11ea-8669-97a2d7823cbb.png">
  <img width="365" alt="Screen Shot 2020-08-16 at 11 23 43 PM" src="https://user-images.githubusercontent.com/672940/90354550-a11b1e00-e017-11ea-9bb6-5f7c6209f7b0.png">
</p>
<br>

## Quickstart

1. Download the latest [release][3] and either build it yourself with `make build` or you can use a prebuilt binary, currently they only exist for OSX. This project uses a branching git flow. Details can be found in the developer doc, suffice it to say **Master** is stable **develop** shoud be considered beta.
2. Download or clone the latest set of [signatures][4] and either copy *signatures/default.yaml* to *~/.wraith/signatures/* or adjust the location in your configuration file. A sample is shown below
3. Copy *config/sample.yaml* or the below configuration to *~/.wraith/config.yaml*. This will allow you to get up and running for basic scans without having to figure out the flags. Any of these values can be overwritten on the commnd line as well. You will need to generate your own api tokens for github and gitlab if you are scanning against them.
4. Once you have this done, just run a scan command.
- `wraith scanGithub`
- `wraith scanGitlab`
- `wraith scanLocalGitRepo`
- `wraith scanLocalPath`

```yaml
---
commit-depth: -1
debug: false
github-api-token: <token>
github-orgs:
    - <org 1>
    - <org 2>
github-repos:
    - <repo 1>
github-users:
    - <user 1>
gitlab-api-token: <token>
gitlab-targets:
    - <repo 1>
    - <project 1>
    - <user 1>
scan-forks: false
scan-tests: false
ignore-extension:
    - .html
    - .css
    - .log
ignore-path:
    - static/
    - docs/
    - .idea/
in-mem-clone: false
local-paths:
    - ../relative/path/to/file.md
    - $HOME/path/to/file.pl
    - /absolute/path/to/file.rb
confidence-level: 3
num-threads: 1
local-repos:
    - ../wraith
    - /home/bob/Go/src/foo
silent: false
web-server: false
json: false
csv: false
test-signatures: false
signatures-path: "~/.wraith/signatures"
signature-file:
    - ../wraith-signatures/signatures/default.yaml
signatures-url: "https://github.com/N0MoreSecr3ts/wraith-signatures"
signatures-version: "0.2.0"
```

## Documentation

### Build from source
At this stage the best option is to build from source from this repository.

To install from source, make sure you have a correctly configured **Go >= 1.14** environment and that `$GOPATH/bin` is in your `$PATH`.
```shell
    $ cd $GOPATH/src
    $ git clone git@github.com:N0MoreSecr3ts/wraith.git
    $ cd wraith
    $ make build
    $ ./bin/wraith-<ARCH> <sub-command>
```

### Signatures
Signatures are the current method used to detect secrets within the a target source. They are broken out into the [wraith-signatures][4] repo for extensability purposes. This allows them to be independently versioned and developed without having to recompile the code. To makes changes just edit an existing signature or create a new one. Check the [README][5] in that repo for additional details.

### Authencation
Wraith will need either a GitLab or Github access token in order to interact with their appropriate API's.  You can create a [GitLab personal access token][6], or [a Github personal access token][7] and save it in an environment variable in your **bashrc**, add it to a wraith config file, or pass it in on the command line. Passing it in on the commandline should be avoided if possible for security reasons. Of course if you want to eat your own dog food, go ahead and do it that way, then point wraith at your command history file. :smiling_imp:

### Additional Documentation
Additional documentation is forthcoming

## Contributing

[Contributing.md][14]

There is a [roadmap][13] as well, but at this point it's little more than a glorified TODO list and personal braindump. I am using that instead of issues, due to my velocity and general laziness towards process at this point. When the project becomes stable, most likely after Milestone 1, the roadmap will probably fall away and be captured in [Issues][15].

## Credits
- [@michenriksen][8] for writing [gitrob][9] which serves as the foundation for wraith
- [@codeemitter][11] for contributing several major features including in memory clones and gitlab support. His version is the immediate parent to wraith.
- [@mattyjones][10] (Maintainer)

## Related
There are several other projects that wraith owes some lineage to including:
- [Trufflehog][12]
- all the many recon and OSINT tools already existing


[1]: https://github.com/spf13/viper
[2]: https://github.com/spf13/cobra
[3]: https://github.com/N0MoreSecr3ts/wraith/releases
[4]: https://github.com/N0MoreSecr3ts/wraith-signatures
[5]: https://github.com/N0MoreSecr3ts/wraith-signatures/blob/master/README.md
[6]: https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html
[7]: https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
[8]: https://github.com/michenriksen
[9]: https://github.com/michenriksen/gitrob
[10]: https://github.com/mattyjones
[11]: https://github.com/codeEmitter/
[12]: https://github.com/dxa4481/truffleHog
[13]: https://github.com/N0MoreSecr3ts/wraith/blob/develop/docs/development/roadmap.md
[14]: https://github.com/N0MoreSecr3ts/wraith/blob/master/Contributing.md
[15]: https://github.com/N0MoreSecr3ts/wraith/issues

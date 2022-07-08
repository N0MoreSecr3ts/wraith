
# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.0.9] - 2022-07-08
### Changed

### Fixed
- Support for the new GHE api token format [#104][104]
- version/version.go had the wrong version number defined [#106][106]
- reset the page count when getting a list of user repos [#107][107]
- findings now work in realtime in the web ui [#126][126]
- add help to makefile [#125][125]

### Added
- Add persistent flags to the root to streamline the adding of additional functionality [#95][95]
- Add ability to expand org members in github [#115][115] [@circleous][circleous]

### Removed

## [0.0.8] - 2021-04-05
### Changed
- Adjust flag descriptions
- Enhanced debug output and header
- Update sample.yaml and quickstart documentation
- Condensed the increment stats to provide a more streamlined approach
- Added the concept of persistent flags at the command root, to make sub-command flags more streamlined

### Fixed
- Flag names were pointing to the wrong variables
- Change yml -> yaml in the documentation across the board
- The web interface now has the correct stats, and the progress bar works again, threading is still an issue though

### Added
- Flag to enable/disable the webserver, default is disable
- Add csv output format
- Add json output format
- Ability to update signatures to a given release

### Removed
- Several dead functions

## [0.0.6] - 2021-01-22
### Changed
- rule -> signature throughout the code

### Added
- support for enterprise github with personal token using basic auth
- ascii banner for the project

### Removed
- remove "los" throughout the code
- remove signatures from wraith by default

### Fixed
- local path scans where not working due to bug in trying to get file changes
- performance issues with scanning local git repos and github repos

## [0.0.4] - 2020-08-10
### Changed
- change the internal name from gitrob to wraith
- condense number of packages to remove cyclic dependencys. Better code organization is still needed.
- rules are known as signatures
- all signatures are in a single yaml file
- how targets are calculated and their count is displayed and referenced

### Added
- rule metadata
- a confidence level for each signature (match level)
- a bit flip to enable/disable a given signature
- ability for signatures to be versioned
- ability to specify one or more signature files
- ability to pull signatures from a default location automatically
- enhanced filtering of match relative to a signature
- additional performance and metrics to the summary output
- additional metrics to the real-time output

### Fixed
- bug in scanLocalGitRepo configuration flags
- gitlab and gihub scans were not working due to the wrong clone function being called.

## [0.0.3] - 2020-08-06
### Changed
- massive rewrite and uplift of code from @codeemitter and @mattyjones

## [0.0.2] - 2020-05-27
### Fixed
- fix issue where if a token was not found or invalid it would panic
- fix issue where is a user or org was not found it would panic

## [0.0.1] - 2020-05-20
### Added
- initial release for the new project

[Unreleased]: https://github.com/mattyjones/wraith/compare/0.0.8...HEAD
[0.0.8]: https://github.com/mattyjones/wraith/releases/tag/0.0.8
[0.0.6]: https://github.com/mattyjones/wraith/releases/tag/0.0.6
[0.0.4]: https://github.com/mattyjones/wraith/releases/tag/0.0.4
[0.0.3]: https://github.com/mattyjones/wraith/releases/tag/0.0.3
[0.0.2]: https://github.com/mattyjones/gitrob/releases/tag/0.0.2
[0.0.1]: https://github.com/mattyjones/gitrob/releases/tag/0.0.1

[95]: https://github.com/N0MoreSecr3ts/wraith/issues/95
[104]: https://github.com/N0MoreSecr3ts/wraith/pull/104
[106]: https://github.com/N0MoreSecr3ts/wraith/issues/106
[107]: https://github.com/N0MoreSecr3ts/wraith/issues/107
[115]: https://github.com/N0MoreSecr3ts/wraith/pull/115

[circleous]: https://github.com/circleous

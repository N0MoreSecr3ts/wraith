# Roadmap

## Scanning

### Targets
- [ ] Scan Github Enterprise Org Repos (needs an org)
- [ ] Scan Github Enterprise User Repos (needs a user)
- [ ] Scan Github Enterprise Single Repo (needs a user/org and a repo)
- [X] ~~Scan Github.com Org Repos (needs an org)~~
- [X] ~~Scan Github.com User Repos (needs a user)~~
- [ ] Scan Github.com Single Repo (needs a repo and a user/org)
- [ ] Scan Github [Gists][1]


- [X] ~~Scan Gitlab.com Group Repos (needs a group ID)~~
- [X] ~~Scan Gitlab.com User Repos (needs a user name)~~
- [ ] Scan Gitlab Snippets
- [ ] Scan Gitlab On-Prem Org Repos
- [ ] Scan Gitlab On-Prem User Repos


- [ ] Scan Bitbucket.com Project Repos
- [ ] Scan Bitbucket.com User Repos
- [ ] Scan Bitbucket On-Prem Org Repos
- [ ] Scan Bitbucket On-Prem User Repos


- [ ] Scan Local Files
- [X] ~~Scan Local Git Repos~~


- [ ] Scan AWS Code Commit
- [ ] Scan Azure DevOps


- [ ] Scan Wiki's
- [ ] Scan Pastebin
- [ ] Scan Confluence


- [ ] Scan OneDrive
- [ ] Scan Dropbox
- [ ] Scan GoogleDrive
- [ ] Scan iCloud


- [ ] Scan MS Office Docs
- [ ] Scan OpenOffice/LibreOffice Docs
- [ ] Scan Evernote
- [ ] Scan GSuite 
- [ ] Scan Quip

### Scaning Features
**default is to include everything**
- [ ] Entrophy Checks
- [ ] Scan specific branches
- [ ] Scan since a given commit
- [X] ~~Exclude files based on extension~~
- [ ] Exclude or include files based on mime type
- [ ] Exclude Test Files
- [ ] Exclude Forks
- [ ] Exclude a default path/extension default exclusion
- [ ] Exclude specific branches or tags
- [ ] Only scan selected branches or tags
- [ ] Exclude public or private repos
- [X] ~~Ability to set commit depth of scan~~
- [ ] Confidence level for regexes
- [X] ~~Should clone to memory, not disk~~
- [X] ~~Exclude  paths~~
- [ ] Exclude Users or Repos in an org scan
- [X] ~~Status output of a session~~
- [X] ~~Ability to silence the output~~
- [ ] JSON or CSV Output
- [ ] Max file size to scan

### UX Features
- [ ] Database Backend
- [ ] Web Frontend For Configuration
- [X] ~~Web Frontend For Output~~
- [ ] Specific YAML Configuration File
- [ ] Signatures in a seperate repo (after on new sigs)
- [ ] Signatures in either yaml or json format (after new sigs)
- [ ] Update Signatures command (after they are moved to a new repo)
- [ ] Fully Instrumented with Performance Stats
- [ ] Finding should have an ID (Hash)
- [X] ~~Ability to configure via environment variables~~
- [ ] Ability to version signatures (after they are moved to a new repo)
- [X] ~~When running silent, no webserver is started~~


## Bugs
- [ ] DB Code is smelly
- [ ] Regex's are not performant (after they are moved to a new repo)
- [ ] Code organization is horrible
- [ ] Consistent search on all platforms
- [X] ~~Web interface line in stdio is borked~~
- [X] ~~Can only find a single target~~
- [X] ~~Silent still displays the gitlab logo~~
- [X] ~~Slient does not print you need to hit Ctrl-C to stop the webserver~~
- [ ] need to update the go [git library][2] used
- [ ] web interface cannot handle local files
- [ ] web interface is gitlab specific by default

## TODO
- [ ] Thread the scanning of commits
- [ ] Make sure we clean up the temp directories
- [ ] Pre-compiled binaries
- [X] ~~Use YAML arrays~~
- [ ] Implement MJ Stats (waiting on new matching)
- [ ] Break out global vs command specific variables
- [ ] Combine all shell scripts into Makefile
- [ ] Split rules into a seperate repo (after new sigs)
- [ ] Combine the rules and sigs into a single yaml file
- [X] ~~Plug into gitlab ci pipeline~~
- [ ] Remove the common package and integrate it with core
- [ ] Add copyright notices
- [X] ~~Remove github traces~~
- [ ] Test all regexes (after new repo)
- [ ] Alpha sort structs, functions, flags
- [ ] Unit tests for all code
    - [ ] common
    - [ ] config
    - [ ] core
    - [ ] github
    - [ ] gitlab
    - [ ] matching
    - [ ] version
    - [ ] rules
- [ ] Debug Info
- [X] ~~How do we want to handle authN~~
- [ ] Better Logging using Logrus
- [ ] Error Handling
    - [ ] common
    - [ ] config
    - [ ] core
    - [ ] github
    - [ ] gitlab
    - [ ] matching
    - [ ] version
    - [ ] rules
- [ ] Code Test Coverage
- [ ] 3PP Scans
- [X] ~~Security.txt~~
- [X] ~~Update Readme.md~~
- [ ] Security Scans
- [ ] Sanitize user inputs
- [ ] Contributing.md
- [ ] Makefile
- [X] ~~.editorconfig~~
- [X] ~~.gitignore~~
- [ ] Mascot (waiting on Mandy)
- [X] Name (wraith)
- [X] Gitlab group
- [ ] go doc strings
    - [ ] common
    - [ ] config
    - [ ] core
    - [ ] github
    - [ ] gitlab
    - [ ] matching
    - [ ] version
    - [ ] rules
- [X] ~~Implement cobra/viper~~
- [X] ~~License.txt~~
- [ ] Table driven tests
- [ ] Need to make a flag in the Makefile to update the dependencies
- [ ] Split out the web go code into a specific package
- [ ] Swap to libgit2 where it makes sense for scaling
- [ ] Golint needs to pass

## Notes
- [ ] Can we Go for the web front-end
- [ ] Language Parsers

[1]: https://github.com/eth0izzle/shhgit/blob/master/core/github.go#L91
[2]: https://pkg.go.dev/github.com/go-git/go-git/v5?tab=doc#example-Clone

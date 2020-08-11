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


- [ ] Scan Local Files **Next**
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
- [ ] Confidence level for regexes (signature uplift)
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
- [ ] Signatures in a seperate repo (signature uplift)
- [ ] Signatures in either yaml or json format (signature uplift)
- [ ] Update Signatures command (signature uplift)
- [ ] Fully Instrumented with Performance Stats
- [X] ~~Finding should have an ID (Hash)~~
- [X] ~~Ability to configure via environment variables~~
- [X] ~~Ability to version signatures (signature uplift)~~
- [X] ~~When running silent, no webserver is started~~


## Bugs
- [X] ~~DB Code is smelly (reomved feature)~~
- [ ] Regex's are not performant (signature uplift)
- [ ] Code organization is horrible
- [ ] Consistent search on all platforms
- [X] ~~Web interface line in stdio is borked~~
- [X] ~~Can only find a single target~~
- [X] ~~Silent still displays the gitlab logo~~
- [X] ~~Slient does not print you need to hit Ctrl-C to stop the webserver~~
- [ ] need to update the go [git library][2] used
- [ ] web interface cannot handle local files (requires mucking with bindata.go)
- [ ] web interface is gitlab specific by default (requires mucking with bindata.go)
- [ ] why is the web interface using to old index.html (requires mucking with bindata.go)
- [X] ~~gitlab scans are failing~~
- [ ] not consistently finding secrets for all sources
- [ ] web interface progress bar not working
- [X] ~~web interface links are broken for local files~~
- [ ] web interface links to the file should be more detailed and point to the commit/line in the code
- [ ] web interface is not dynamic, I need to refresh it manually
- [ ] no error when rules file is not found
- [ ] in-mem-clone is not working properly
- [ ] working with local repos is not threaded
- [ ] there are no findings in a local search
- [ ] there are no findings in a gitlab search
- [ ] there are multiple generateid functions
- [ ] need to reorg the code again



## TODO
- [ ] Remove the repo count per target during a run
- [X] ~~Repositiores -> reposScanned~~
- [X] ~~commits -> commitsScanned~~
- [X] ~~findings -> findings total~~
- [X] ~~files -> files scanned~~
- [ ] Add config details to summary output
- [ ] Add content to summary
- [ ] cleanup issues in summary output
- [X] ~~Implement match level for sigs~~
- [ ] Create stats for signatures
- [ ] Move sigs to a different repo
- [ ] Implement rules in either json or yaml
- [ ] call all rules sigs
- [ ] implement comand to update sigs from repo
- [X] ~~port all grover stats~~
- [ ] tests for all regex's
- [ ] remove all debugging statements
- [ ] remove all dead code
- [ ] add more debuging info
- [ ] add additional stats to web interface
- [ ] what errors should stop the run
- [ ] add a flag to point to a custom rules file
- [X] ~~add flag for setting the match level~~
- [ ] document the match level
- [ ] document all stats
- [ ] implement flag for setting the thread count manually
- [ ] document the differance between targets and repos
- [ ] document all code completely
- [ ] create a developer doc with the design and code execution flow
- [ ] Look at the clone configs
- [ ] Create new ascii art
- [ ] Rebuild the web interface
- [ ] Copy existing tests to the new codebase (need to reorg the existing codebase first)
- [ ] Make a single function to create a temp dir
- [ ] Need to list the flag defaults on the help screen
- [ ] If no arg's are given for a command, then list the help screen
- [ ] Thread the scanning of commits
- [ ] Check all urls point to the right repos (requires mucking with bindata.go)
- [ ] Update Code Climate for Wraith
- [ ] Write a new README
- [ ] Make sure we clean up the temp directories
- [X] ~~Pre-compiled binaries~~
- [X] ~~Use YAML arrays~~
- [X] ~~Implement MJ Stats (waiting on new matching)~~
- [ ] Break out global vs command specific variables
- [ ] Combine all shell scripts into Makefile
- [ ] Split rules into a seperate repo (signature uplift)
- [X] ~~Combine the rules and sigs into a single yaml file (signature uplift)~~
- [X] ~~Plug into gitlab ci pipeline~~
- [X] ~~Remove the common package and integrate it with core~~
- [X] ~~Add copyright notices~~
- [X] ~~Remove github traces~~
- [ ] Test all regexes (signature uplift)
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
- [X] ~~Name (wraith)~~
- [X] ~~Gitlab group~~
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
- [ ] Convert Repository type, etc into github ones if needed
- [ ] How do we deal with abuse msgs, and rate limiting

[1]: https://github.com/eth0izzle/shhgit/blob/master/core/github.go#L91
[2]: https://pkg.go.dev/github.com/go-git/go-git/v5?tab=doc#example-Clone

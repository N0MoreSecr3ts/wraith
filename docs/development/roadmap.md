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


- [X] ~~Scan Local Files~~
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

## Milestone 1

### Bugs
- [X] ~~temp directories are not getting deleted~~
- [X] ~~in-mem-clone is not working properly~~
- [X] ~~there are no findings in a gitlab search~~
- [X] ~~there are multiple generateID functions~~
- [X] ~~Web Frontend For Output~~
- [X] ~~Specific YAML Configuration File~~
- [X] ~~Signatures in a seperate repo (signature uplift)~~
- [X] ~~Signatures in either yaml or json format (signature uplift)~~
- [X] ~~Finding should have an ID (Hash)~~
- [X] ~~Ability to configure via environment variables~~
- [X] ~~Ability to version signatures (signature uplift)~~
- [X] ~~When running silent, no webserver is started~~
- [X] ~~DB Code is smelly (reomved feature)~~
- [X] ~~Web interface line in stdio is borked~~
- [X] ~~Can only find a single target~~
- [X] ~~Silent still displays the gitlab logo~~
- [X] ~~Slient does not print you need to hit Ctrl-C to stop the webserver~~
- [X] ~~gitlab scans are failing~~
- [X] ~~not consistently finding secrets for all sources~~
- [X] ~~web interface links are broken for local files~~
- [X] ~~no error when rules file is not found~~
- [X] ~~there are no findings in a local git search~~
- [ ] Gitlab client does not follow redirects
- [ ] Need to confirm if github client follows redirects
- [ ] In gitlab you can have a project w/ no repo, this will error out
- [ ] In github you can have a project w/ no repo, does this error out
- [X] ~~Searching through a commit history is present but not effective~~
- [X] ~~Secret ID's are possibly not unique~~
- [ ] Duplicate findings are being displayed
- [X] ~~Findings with a line number of 0 are being displayed~~
- [ ] Review all flags to ensure they are needed
- [ ] Expanding orgs is not working
- [ ] Setting the number of commits to 1 does not scan anything
- [X] ~~Number of total commits was wrong~~
- [X] ~~Line numbers are wrong in the patches being scanned~~
- [ ] Number of dirty commits is wrong, it should be more than is showing
- [X] ~~**Redo github enterprise bits to clean them up**~~
- [X] ~~**Refactor the client to follow G's method**~~
- [ ] **Set the debug like G**
- [ ] Set the csv and json like G
- [X] ~~**Set orgs and repos like G**~~
- [ ] **gitlab api endpoint**
- [ ] **Port gitlab to match G**
- [ ] ~/ does not work in the config file due to a missing '/'
- [ ] yaml list does not work for github enterprise 
- [ ] need to be able to scan all the repos for a specific user
- [ ] need to be able to scan all the orgs a user is a member of
- [X] ~~remove all lazy, one-off data structures and use the official ones~~
- [ ] refactor how we do stats on commits
- [ ] Validate the gitlab api token
- [ ] Uplift gitlab functionality to match github

### Documentation
- [ ] **Document all flags**
- [ ] Document how to add a new command or source
- [ ] Document the tech debt using colors and a shell script
- [ ] Document all stats
- [X] ~~**Document the differance between targets and repos**~~
- [ ] **Document all code completely**
- [ ] Create a developer doc with the design and code execution flow
- [ ] Contributing.md
    - [ ] wraith
    - [ ] wraith-tests
    - [ ] wraith-signatures
- [ ] README.md
    - [X] ~~wraith~~
    - [ ] wraith-tests
    - [ ] wraith-signatures
- [ ] Security.txt
    - [X] ~~wraith~~
    - [ ] wraith-tests
    - [ ] wraith-signatures
- [ ] Initial blog post
- [ ] Detailed documentation published on the net and with source control
- [X] ~~Cleanup issues~~
- [ ] Changelog.md
    - [ ] wraith
    - [ ] wraith-tests
    - [ ] wraith-signatures
- [X] ~~Label issues for begineer and hacktoberfest~~
- [ ] **Go doc strings**
    - [ ] common
    - [ ] config
    - [ ] core
    - [ ] version
- [ ] Issue template
- [ ] PR template
- [ ] Submit story to hackernews
- [ ] Submit story to changelog.com
- [ ] Add a built w/ section
- [ ] Call out individual contributers after N merges
    
### Testing
- [ ] Copy existing tests to the new codebase
- [X] ~~Confirm hide secrets~~
- [X] ~~Update Code Climate for Wraith~~
- [ ] Update CodeCov for Wraith
- [ ] Golint needs to pass
- [ ] Convert tests to testify
- [ ] **Code review and remove debug statements**
- [X] ~~Set unique secrets in the test~~

### Features
- [ ] Change name from threads to go routines or make that clear
- [ ] Refactor how threads are calculated
- [X] ~~**Configure repos~~
- [X] ~~**Configure Org~~
- [ ] Created a dedicated GPG key
- [ ] Enforce https for all connections
- [ ] Enforce https for the site
- [ ] Fully Instrumented with Performance Stats
- [ ] JSON or CSV Output
- [ ] **Exclude/Include Forks**
- [ ] Entrophy Checks
- [ ] If we find a .git directory in a localPath scan just ignore it and process the dir as localPath
- [X] ~~Change empty string defaults to nil~~
- [X] ~~Add content to summary~~
- [ ] Cleanup issues in summary output
- [ ] **Remove all print debugging statements**
- [ ] **Remove all dead code**
- [ ] Add more debuging info
- [X] ~~Implement flag for setting the thread count manually~~
- [X] ~~**Look at the clone configs**~~
- [X] ~~Make a single function to create a temp dir~~
- [ ] Need to list the flag defaults on the help screen
- [ ] If no arg's are given for a command, then list the help screen
- [X] ~~Make sure we clean up the temp directories~~
- [ ] Alpha sort structs, functions, flags
- [X] ~~Exclude files based on extension~~
- [X] ~~Exclude Test Files~~
- [X] ~~Ability to set commit depth of scan~~
- [X] ~~Confidence level for regexes (signature uplift)~~
- [X] ~~Should clone to memory, not disk~~
- [X] ~~Exclude  paths~~
- [X] ~~Status output of a session~~
- [X] ~~Ability to silence the output~~
- [X] ~~Max file size to scan~~
- [ ] Only export the functions and variables necessary
- [ ] Capture the error if no sig file is given and the default does not exist
- [ ] Break out checking if a file is to be scanned into a single function
- [ ] Add a flag to de-dupe findings


### Milestone 2

### Bugs
- [ ] web interface cannot handle local files (requires mucking with bindata.go)
- [ ] web interface is gitlab specific by default (requires mucking with bindata.go)
- [ ] why is the web interface using to old index.html (requires mucking with bindata.go)
- [ ] web interface progress bar not working
- [ ] web interface links to the file should be more detailed and point to the commit/line in the code
- [ ] web interface is not dynamic, I need to refresh it manually
- [ ] validate all user input

### Documentation

### Features
- [ ] Scan specific branches
- [ ] Scan since a given commit
- [ ] Update Signatures command
- [ ] Implement threading for local path scans
- [ ] Ability to use the .gitignoe when scanning for ingoring paths and files

### Testing
- [ ] Make tech debt fail build process
- [ ] Add config details to debug statement at the start of a run
- [ ] Create stats for signatures
- [ ] Structured Logging
- [ ] Create a standard set of error codes
- [ ] Error Handling
    - [ ] config
    - [ ] core
    - [ ] version
- [ ] Security Scans
- [ ] Sanitize user inputs
- [ ] Code should be optimized into multiple packages

## Milestone 3

### Bugs

### Features
- [ ] working with local repos is not threaded
- [ ] Exclude or include files based on mime type
- [ ] Exclude a default path/extension default exclusion
- [ ] Exclude specific branches or tags
- [ ] Only scan selected branches or tags
- [ ] Exclude public or private repos
- [ ] Exclude Users or Repos in an org scan
- [ ] Database Backend
- [ ] Web Frontend For Configuration

### Testing

### Documentation

- [ ] Consistent search on all platforms
- [ ] need to update the go [git library][2] used
- [ ] Regex's are not performant
- [ ] Move the repo count per target during a run to a debug statement
- [ ] Create stats for signatures
- [ ] tests for all regex's
- [ ] add additional stats to web interface
- [ ] Create new ascii art
- [ ] Rebuild the web interface
- [ ] Thread the scanning of commits
- [ ] Check all urls point to the right repos (requires mucking with bindata.go)
- [ ] Break out global vs command specific variables
- [ ] Combine all shell scripts into Makefile
- [ ] Test all regexes
- [ ] Unit tests for all code
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
- [ ] Table driven tests
- [ ] Need to make a flag in the Makefile to update the dependencies
- [ ] Split out the web go code into a specific package
- [ ] Swap to libgit2 where it makes sense for scaling
    
## Research
- [ ] Do we want to add files,dirs,repos,etc to an ignore list when they are not found or they error out
- [ ] Look at using the gitignore when scanning repos
- [ ] what errors should stop the run
- [ ] Mascot

## Notes
- [ ] Can we Go for the web front-end
- [ ] Language Parsers
- [ ] Convert Repository type, etc into github ones if needed
- [ ] How do we deal with abuse msgs, and rate limiting

[1]: https://github.com/eth0izzle/shhgit/blob/master/core/github.go#L91
[2]: https://pkg.go.dev/github.com/go-git/go-git/v5?tab=doc#example-Clone

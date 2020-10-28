# Roadmap

## Scanning

### Targets
- [X] ~~Scan Github Enterprise Org Repos (needs an org)~~
- [X] ~~Scan Github Enterprise User Repos (needs a user)~~
- [X] ~~Scan Github Enterprise Single Repo (needs a user/org and a repo)~~
- [X] ~~Scan Github.com Org Repos (needs an org)~~
- [X] ~~Scan Github.com User Repos (needs a user)~~
- [X] ~~Scan Github.com Single Repo (needs a repo and a user/org)~~
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
- [ ] Scan S3 Buckets


- [ ] Scan MS Office Docs
- [ ] Scan OpenOffice/LibreOffice Docs
- [ ] Scan Evernote
- [ ] Scan GSuite 
- [ ] Scan Quip

## Milestone 1

### Bugs
- [X] ~~Temp directories are not getting deleted~~
- [X] ~~In-mem-clone is not working properly~~
- [X] ~~There are no findings in a gitlab search~~
- [X] ~~There are multiple generateID functions~~
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
- [X] ~~Gitlab scans are failing~~
- [X] ~~Not consistently finding secrets for all sources~~
- [X] ~~Web interface links are broken for local files~~
- [X] ~~No error when rules file is not found~~
- [X] ~~There are no findings in a local git search~~
- [X] ~~Searching through a commit history is present but not effective~~
- [X] ~~Secret ID's are possibly not unique~~
- [X] ~~Findings with a line number of 0 are being displayed~~
- [X] ~~Number of total commits was wrong~~
- [X] ~~Line numbers are wrong in the patches being scanned~~
- [X] ~~Path ~/ does not work in the config file due to a missing '/'~~
- [X] ~~Yaml list does not work for github enterprise~~
- [X] ~~Yaml list does not work for github~~
- [X] ~~Yaml list does not work for local git~~
- [X] ~~Not all user repos get pulled~~
- [X] ~~In github you can have an org w/ no repo, this will error out~~
- [X] ~~Scaning forks is not working~~
- [X] ~~Change default commit depth to -1~~
- [ ] Still havea lot of missing files when scanning (turn on debug)
- [X] ~~Change github orgs, repos, users to use a slice~~
- [X] ~~Change github enterprise orgs, repos, user, to use a slice~~
- [X] ~~Change ignore path to use a slice~~
- [X] ~~Change ignore extension to use a slice~~
- [X] ~~Change match-level to confidence level~~
- [X] ~~Change default thread count to -1~~
- [ ] In memory clone returns no findings (I think this has something to do with the path not being found)
- [ ] Fix how all repos are gathered (org repos is threaded and general, user repos is not threaded and github specific)
- [X] ~~Repo totals are getting counted twice~~
- [ ] Gitlab client does not follow redirects
- [ ] Github does not follow redirects
- [ ] In gitlab you can have a project w/ no repo, this will error out
- [ ] Expanding orgs is not working
- [X] Not scanning tests is not working
- [X] Max file size is not working
- [X] Setting the number of commits to 1 does not scan anything
- [ ] Number of dirty commits is wrong, it should be more than is showing


### Documentation
- [X] ~~Document the differance between targets and repos~~
- [X] ~~Document the tech debt using colors and a shell script~~
- [0] Document all flags
- [0] Document all code completely
- [ ] Document how to add a new command or source
- [ ] Document all stats
- [ ] Create a developer doc with the design and code execution flow
- [ ] Contributing.md
    - [X] ~~wraith~~
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
- [0] Go doc strings
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
- [X] ~~Confirm hide secrets~~
- [X] ~~Update Code Climate for Wraith~~
- [X] ~~Set unique secrets in the test~~
- [0] Code review and remove debug statements
- [0] Ensure that an error status will exit the program, if not swap to a warning status so it does
- [ ] Sanity check testing plan
- [ ] Copy existing tests to the new codebase
- [ ] Make sure we use https so keys are not necessary
- [ ] Update CodeCov for Wraith
- [ ] Golint needs to pass
- [ ] Convert tests to testify
- [ ] Review all flags to ensure they are needed

### Features
- [X] ~~Redo github enterprise bits to clean them up~~
- [X] ~~Refactor the client to follow G's method~~
- [X] ~~Set orgs and repos like G~~
- [X] ~~Configure repos~~
- [X] ~~Configure Org~~
- [X] ~~Change empty string defaults to nil~~
- [X] ~~Add content to summary~~
- [X] ~~Implement flag for setting the thread count manually~~
- [X] ~~Look at the clone configs~~
- [X] ~~Make a single function to create a temp dir~~
- [X] ~~Make sure we clean up the temp directories~~
- [X] ~~Exclude files based on extension~~
- [X] ~~Exclude Test Files~~
- [X] ~~Ability to set commit depth of scan~~
- [X] ~~Confidence level for regexes (signature uplift)~~
- [X] ~~Should clone to memory, not disk~~
- [X] ~~Exclude  paths~~
- [X] ~~Status output of a session~~
- [X] ~~Ability to silence the output~~
- [X] ~~Max file size to scan~~
- [X] ~~Need to be able to scan all the repos for a specific user~~
- [X] ~~Need to be able to scan a single repo for a specific user~~
- [ ] -1 Confidence level loads all signatures
- [ ] Signature file flag should be a slice
- [X] Exclude/Include Forks 
- [0] Need to drop in the org in the realtime output
- [0] Add the status to all functions for use in the web interface*__*
- [0] Need to find gitlab api endpoint
- [0] Port gitlab to match G 
- [0] Set the debug like G 
- [0] Remove all print debugging statements
- [0] Remove all dead code
- [ ] Change id -> ID
- [ ] Refactor how we do stats on commits
- [ ] Need to be able to scan all the orgs a user is a member of
- [ ] Need to be able to scan all the forks of a given repo that we can reach
- [ ] Set the csv and json like G
- [ ] Change name from threads to go routines or make that clear
- [ ] Refactor how threads are calculated
- [ ] Created a dedicated GPG key
- [ ] Enforce https for all connections
- [ ] Enforce https for the site
- [ ] Fully Instrumented with Performance Stats
- [ ] Entrophy Checks
- [0] If we find a .git directory in a localPath scan just ignore it and process the dir as localPath
- [ ] Cleanup issues in summary output
- [ ] Add more debuging info
- [X] ~~Alpha sort structs, functions, flags~~
- [ ] Need to list the flag defaults on the help screen
- [ ] If no arg's are given for a command, then list the help screen
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
- [ ] Ability to use the .gitignore when scanning for ignoring paths and files

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

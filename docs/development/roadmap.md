# Roadmap

## Scanning

### Targets

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

- [X] ~~Still have a lot of missing files when scanning (turn on debug)~~
- [0] In memory clone returns no findings (I think this has something to do with the path not being found)
- [0] Fix how all repos are gathered (org repos is threaded and general, user repos is not threaded and github specific)
- [0] Gitlab client does not follow redirects
- [0] Github does not follow redirects
- [0] In gitlab you can have a project w/ no repo, this will error out
- [0] Expanding orgs is not working
- [0] Number of dirty commits is wrong, it should be more than is showing

### Documentation
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
- [ ] Changelog.md
    - [ ] wraith
    - [ ] wraith-tests
    - [ ] wraith-signatures
- [0] Go doc strings
    - [0] core
    - [0] version
- [ ] Issue template
- [ ] PR template
- [ ] Submit story to hackernews
- [ ] Submit story to changelog.com
- [ ] Add a built w/ section
- [ ] Call out individual contributers after N merges
    
### Testing
- [0] Code review and remove debug statements
- [0] Ensure that an error status will exit the program, if not swap to a warning status so it does
- [0] Sanity check testing plan
- [0] Copy existing tests to the new codebase
- [ ] Make sure we use https so keys are not necessary
- [ ] Update CodeCov for Wraith
- [0] Golint needs to pass
- [ ] Convert tests to testify
- [0] Review all flags to ensure they are needed

### Features
- [ ] -1 Confidence level loads all signatures
- [ ] Tor capability
- [ ] IP switching to hide itself
- [ ] Ability to use multiple tokens/keys for a single service (several GH keys)
- [0] Need to drop in the org in the realtime output
- [0] Add the status to all functions for use in the web interface*__*
- [0] Need to find gitlab api endpoint
- [0] Port gitlab to match G 
- [0] Set the debug like G 
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

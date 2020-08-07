## Release Process
Currently releases are managed by hand but strictly follow semver2 and the changelog should be considered authoritative.

When a Pull Request is submitted, it will be reviewed and any feedback will be given within the PR. When a committer wishes to deploy a new release the following procedure should be followed:

1. Update CHANGELOG to reflect all of the changes that has happened between last release and now. The Unreleased link in the CHANGELOG gives you a nice diff.
1. Make sure the README is updated as necessary.
1. Update the version using semver2.
1. Push the commit that bumps the changelog.
1. Make a git release using the changelog commit.
1. Add any necessary binaries to the release they were built with.

## Issue and Pull Request Submissions
If you see something wrong or come across a bug please open up an issue, try to include as much data in the issue as possible. If you feel the issue is critical than tag a team member and we will respond as soon as is feasible.

Pull requests need to follow the guidelines below for the quickest possible merge. These not only make our lives easier, but also keep the repo and commit history as clean as possible.

- Please do a git pull --rebase both before you start working on the repo and then before you commit. This will help ensure the most up to date codebase. It will also go along way towards cutting down or eliminating(hopefully) annoying merge commits.
- The CHANGELOG follows the standard conventions laid out [here](https://keepachangelog.com/en/1.0.0/). Every PR has to include an updated CHANGELOG and README (if needed), this makes our lives easier, increases the accuracy of the codebase, and gets your PR deployed much faster.
- When suggesting a new version please keep the following in mind
    - The patch version is for any non-breaking changes to existing code or the addition of minor functionality to existing code
    - The minor version is for the addition of any new functionality. Even though this is generally non-breaking, it is a major change and should be indicated as such
    - The major version should only be bumped by an admin/owner. This is for major breaking or non-breaking changes that affect widespread functionality. Examples of this would be a wholesale refactor of the repo or a switch away from an established method such as going from SOAP to REST.

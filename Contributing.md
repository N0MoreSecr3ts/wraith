# Contributing

Pull requests are welcome from everyone.

Keep an open mind! Improving documentation, bug triaging, or writing blog posts are all examples of helpful contributions that mean less work for the maintainers.

All new functionality should have an issue openned prior to submitting a pull request. This will speed up to merging of your pull request by a significant amount. Being able to discuss it ahead of time, provide feedback on the code, and any additional details will help you in the long run.

This guide for contributions will make things much easier for everyone involved. You will know what is needed to get an issue triaged or a pull requested merged with as little non-development effort as possible.

The maintainers will be able to review and merge code faster if all tests are passing and the proper formating and style have already been implemented.

Bug fixes and simple non-functional changes such as fixing typo's does not need an issue openned prior to making a pull request.


# Ground Rules

The issue tracker is the preferred channel for bug reports, features requests and submitting pull requests, but please respect the following restrictions.

Please do not derail or troll issues. Keep the discussion on topic and respect the opinions of others.

This includes not just how to communicate with others (being respectful, considerate, etc) but also technical responsibilities (importance of testing, project dependencies, etc). Deatils can be found in the [Code of Conduct](wraith_coc.html).


## Responsibilities
* Ensure cross-platform compatibility for every change that's accepted. 
    * OSX
    * Windows
    * Linux
        * CentOS
        * Debian
        * Ubuntu
        * Arch
* Ensure that code merged into wraith meets all the requirements [here](https://github.com/golang/go/wiki/CodeReviewComments)
* Create issues for any major changes and enhancements that you wish to make. Discuss things transparently and get community feedback.
* Keep feature versions as small as possible, preferably one new feature per version. This is directly related to semver.
* Be welcoming to newcomers and encourage diverse new contributors from all backgrounds.

# First Contributions

Unsure where to begin contributing to Wraith? You can start by looking through issues labeled with [Beginner](https://github.com/N0MoreSecr3ts/wraith/issues?q=is%3Aopen+is%3Aissue+label%3ABeginner) or [Hacktoberfest](https://github.com/N0MoreSecr3ts/wraith/labels/Hacktoberfest):

**Beginner** - Issues which should only require a few lines of code, and a test or two.

**Hacktoberfest** - issues which are more broad and good for new developers or those that want to contribute in other ways

 [How to Contribute to an Open Source Project on GitHub](https://egghead.io/series/how-to-contribute-to-an-open-source-project-on-github)


# Getting started

## General Contributer Guidelines
* All merge conflicts must be resolved before the pull request can be reviewed
* The changelog must be updated following these [conventions](https://keepachangelog.com/en/1.0.0/). Your entrys should go in the unreleased section.
* Commit messages should follow these [guidelines](https://chris.beams.io/posts/git-commit/)
* This project follows [git flow](https://guides.github.com/introduction/flow/index.html)
* When you submit code changes, your submissions are understood to be under the same [MIT License](https://choosealicense.com/licenses/mit/) that covers the project. Feel free to contact the maintainers if that's a concern.
* This project adheres to [SemVer 2.0](https://semver.org).

## New Features, larger fixes, contributions
For larger changes including new features this is the fastest way to get it merged:

1. Open an issue to describe your proposed improvement or feature
1. Ensure you have a proper Golang development environment
1. Fork the repo and create your branch from **develop**
1. Create your feature branch (`git checkout -b my-new-feature`)
1. Add a new changelog entry matching the previous entries and using this as a [guide](https://keepachangelog.com/en/1.0.0/).
1. If you've added code that should be tested, add tests
1. Update any documentation
4. Ensure the test suite passes
1. Push your feature branch (`git push origin my-new-feature`)
1. Create a Pull Request as appropriate based on the issue discussion
    * Please add **WIP** to the pull request title if it is not ready to be merged. This will allow us to prioritize code reviews.
    * [WIP] My New Feature
1. Respond to any comments by the maintainers within two weeks

## Simple bug fixes, typo's and documentation
These are one or two line fixes, spelling issues or non-functional changes such as adding tests, logging, debuging information. Documentation updates or expansions, including godoc strings also qualify.

1. Ensure you have a proper Golang development environment
1. Fork the repo and create your branch from `develop`
1. Create your feature branch (`git checkout -b my-new-feature`)
1. Add a new changelog entry matching the previous entries and using this as a [guide](https://keepachangelog.com/en/1.0.0/).
1. If you've added code that should be tested, add tests
1. Update any documentation
4. Ensure the test suite passes
1. Push your feature branch (`git push origin my-new-feature`)
1. Create a Pull Request as appropriate based on the issue discussion
    * Please add **WIP** to the pull request title if it is not ready to be merged. This will allow us to prioritize code reviews.
    * [WIP] My New Feature
1. Respond to any comments by the maintainers within two weeks

# Bug Reports
If you find a security issue please do not open a public issue, guideline for this can be found in *security.txt* within the repository root.

When filing an issue, please make sure to answer these questions and fill in any other details. The more details, screenshots, configuration infomation you can provide, the faster we can react to the bug:

1. What version of Wraith and Wraith Signatures are you using?
1. Did you use a pre-compiled release ot did you build it yourself?
1. What operating system and processor architecture are you using?
1. What did you do?
    * Detailed steps to reproduce the bug
    * Sample code or screenshots
1. What did you expect to see?
1. What did you see instead?
1. Have you tried anything to fix the issue yourself or do you know what may be the cause.

# Feature Suggestions
There is currently a running roadmap in the docs directory of the repo. These are sorted by non-specific milestones with no due date. If you would like to pick something out of this list please open an issue or pull request with **[WIP]** at the beginning of the title.

This information will give contributors context before they make suggestions that may not align with the project’s needs.

If you find yourself wishing for or needing a feature that doesn't exist for an engagement, you are probably not alone. Many of the features, indeed this entire project, were born from a need to scan something. Open an [issue](https://github.com/N0MoreSecr3ts/wraith/issues) on GitHub which describes the feature you would like to see, why you need it, and how it should work.

# Code review process
After a pull request has been submited, or the **WIP** tage has been removed from the title it will be reviewd within two weeks. This is an open source project and everyone has day jobs and other commitments. It will most likely happen before then, but it will happen by then.

When looking at a pull request, it will be given a code review for obvious security or functional issues, style concerns, or other issues that may present themselves further down the line. All libraries will be checked against automated services for known vulnerabilities.

# Community
For updated release information and bug information please follow the project on [twitter](https://twitter.com/N0MoreSecr3ts).

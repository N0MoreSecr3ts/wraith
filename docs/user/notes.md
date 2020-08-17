https://michenriksen.com/blog/gitrob-now-in-go/
https://michenriksen.com/blog/gitrob-putting-the-open-source-in-osint/


- when a file cannot be found/read we print the error and add it to the ignore list. We should find a way to print out the ignore list. The error printing out is good for now though.

- in memory cloning should be turned on with care. If the repo or the targets are big enough memory exhaustion issues may present themselves.

- if you dont want to scan any commit depth (1) you should use the local file scan

- a content scan will find all secrets in a file using FindAll. This will lead to duplicate findings such that a given string may be either an artifactory password due to surronding context or simply a generic password.

--- title
cloudbuild github app
--- description
there's a what now?
--- main


[Google Cloud Build](https://cloud.google.com/cloud-build/)
recently got a new
[Github App](https://github.com/marketplace/google-cloud-build)

This is apparently different from the old
oauth2 connect repository workflow

### _new_ things

- directly select a repository, no need to connect
- pull request trigger

### things that _broke_

- source provenance
  - the source now appears to be a cloud storage submitted build
  - you also lose the repository name
  - and the specific trigger
  - workaround is to **manually** add back the info as tags

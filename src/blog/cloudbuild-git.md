title = cloudbuild git
date = 2019-07-08
desc = git stuff done on cloudbuild

---

[Google Cloud Build](https://cloud.google.com/cloud-build/),
the minimalist CI system that is kinda close to
[Knative Build](https://knative.dev/docs/build/)

Anyways,
here's to working with git in the cloud build environment

# Environment

you are _root_,
your repo is shallow cloned into `/workspace`

## _Github_ Permissions

goals: make a mountable volume with github permissions

use by mounting the volume into other steps,
`git push`... should just work

1. create key `ssh-keygen -t ed25519 -C some_descriptive_comment -f ssh_keyfile`
2. add key to Github repo > Settings > Deploy Keys > Add deploy key > copy and paste `ssh_keyfile.pub`
3. create GCP KMS keyring `gcloud kms keyrings create keyring_name --location global`
4. create GCP KMS key `gcloud kms keys create key_name --location=global --keyring=keyring_name --purpose=encryption`
5. encrypt ssh private key `gcloud kms encrypt --plaintext-file=ssh_keyfile --ciphertext-file=ssh_keyfile.kms --location=global --keyring=keyring_name --key=key_name`
6. add cloudbuild step

cloudbuild.yaml:

```
secrets:
  - kmsKeyName: projects/project_id/locations/global/keyRings/keyring_name/cryptoKeys/key_name
    secretEnv:
      GH_KEY: ouput of $(base64 ssh_keyfile.kms)
steps:
  - id: setup git permissions
    name: gcr.io/cloud-builders/gcloud
    # this is Github's public key
    # get by running: ssh-keyscan -t rsa github.com
    env:
      - GH_KNOWN=github.com ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8V
    secretEnv:
      - GH_KEY
    entrypoint: bash
    volumes:
      - name: config
        path: /root
    args:
      - -c
      - >-
        set -euxo pipefail;
        mkdir -p /root/.ssh &&
        printenv GH_KNOWN > /root/.ssh/known_hosts &&
        printenv GH_KEY > /root/.ssh/id_ed25519 && chmod 0600 /root/.ssh/id_ed25519 &&
        echo Hostname github.com >> /root/.ssh/config && echo IdentityFile /root/.ssh/id_ed25519 >> /root/.ssh/config &&
        git config user.email $(gcloud auth list --filter=status:ACTIVE --format='value(account)')


  - id: test git push
    name: gcr.io/cloud-builders/git
    volumes:
      - name: config
        path: /root
    entrypoint: bash
    args:
      - -c
      - >-
        set -euxo pipefail &&
        touch test_file &&
        git add test_file &&
        git commit -m "add test file" &&
        git push git@github.com:user/repo $BRANCH_NAME

```

## Getting More _History_

goals: restore history and branch names

cloudbuild shallow clones the current commit to save time,
it is also set as `master`

```
- id: restore history
  name: gcr.io/cloud-builders/git
  entrypoint: bash
  args:
  - -c
  - >-
    set -euxo pipefail &&
    git fetch --unshallow &&              # gets all history use --depth=n for more limited history
    git branch -m $BRANCH_NAME &&         # rename current branch to its original name
    git checkout -b master origin/master  # set master to upstream master

```

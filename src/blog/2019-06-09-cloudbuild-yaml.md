--- title
cloudbuild yaml
--- description
current preferred setup of Google Cloud Build
--- main


# Google Cloud Build

my current use of `cloudbuild.yaml`

## cloudbuild.yaml

```
steps:
  # build steps to start in parallel
  - name: gcr.io/$PROJECT_ID/parcel:latest
    id: js build 1
    waitFor: [-]
    args:
      - build
      - ...
  - name: gcr.io/$PROJECT_ID/parcel:latest
    id: js build 2
    waitFor: [-]
    args:
      - build
      - ...

  # docker builds will respect subdirs
  # remember to push
  - name: gcr.io/cloud-builders/docker:latest
    dir: subdir
    args:
      - build
      - -t
      - gcr.io/$PROJECT_ID/subdir
      - .

  # kaniko needs context to be explicitly set if not in root dir (/workspace)
  # default cache time is 2 weeks
  # pushes layers as they are built
  - name: gcr.io/kaniko-project/executor:latest
    dir: an-img
    args:
      - --cache
      - --context=/workspace/cf-dns-update
      - --destination=gcr.io/$PROJECT_ID/an-img:latest

  # final build waits for everything to complete
  - name: gcr.io/$PROJECT_ID/site-builder:latest
    id: static site gen
    args:
      - build
      # env vars need to be double escaped
      - $$TOKEN
    env:
      - E1=env1
    secretEnv:
      - TOKEN



substitutions:
  # must start with _underscore
  _KEYRING: cloudbuilder
  _KEYNAME: akey


secrets:
  # create key ring:
  #   gcloud kms keyrings create $KEYRING --location global
  # create key
  #   gcloud kms keys create $KEYNAME --keyring $KEYRING --location global --purpose encryption
  # encrypt file
  #   gcloud kms encrypt --plaintext-file $INPUTFILE --ciphertext-file $INTERMEDIATEFILE --location global --keyring $KEYRING --key $KEYNAME
  #   base64 $INTERMEDIATEFILE > $OUTPUTFILE
  # or encrypt STDIN
  #   echo $SECRET | gcloud kms encrypt --plaintext-file - --ciphertext-file - --location global --keyring $KEYRING --key $KEY | base64
  - kmsKeyName: projects/$PROJECT_ID/locations/global/keyRings/$_KEYRING/cryptoKeys/$_KEYNAME
    secretEnv:
      TOKEN: base64encodedtoken

images:
  # only needed with docker builds
  - image:tag
```

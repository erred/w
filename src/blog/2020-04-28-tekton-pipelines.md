---
description: using tekton pipelines
title: tekton pipelines
---

### _tekton_ pipelines

Fancy ci/cd system / task runner with a hard dependency on k8s.
Why would you want to run it?
You want to declaritvely manage your entire setup so
`kustomize build | kubectl apply -f -` will have your entire build system setup.

#### _tasks_ and pipelines

- _step:_ smallest unit of execution, a container image
- _task:_ multiple steps run sequentially
- _pipeline:_ multiple tasks, can be in parallel or any DAG
- _taskrun:_ an execution of a task
- _pipelinerun:_ an execution of a pipeline
- _workspace:_ directories you can pass between tasks, mounted in all steps in task
- _params:_ strings (or arrays) you can declare and pass one level down

taskruns and pipelineruns are implemented as pods
which are left in _exited_ state after completion.
keep them to keep the logs, but they clutter up your namespace

array params are mostly useless,
a lot of the time you need to pass them into a bash script / inline json
and it's easier to just use a string

secrets can only be mounted as files so you have to invoke a shell and do `$(cat path/to/secret)`

persistent volumes (for volumes) are finicky,
they keep state between executions
and is shared between all of them, no clean slate.
not even sure is readwriteonce is respected

#### _triggers_

- _eventlistener:_ a deployment that recieves requests and triggers pipelineruns
- _triggerbinding:_ translate incoming request json to params for pipelineruns
- _triggertemplate:_ template for pipelineruns to be executed
- _interceptor:_ filter and mangle the incoming request / json on eventlisteners

the interceptors (including cel) are super limited
and you will almost immediately want to run your own deployment
for a webhook interceptor. (ex you can't do replaceall, split-join...)

it's supposed to be reuseable, but i seriously doubt it

#### example

##### _event_ listener

- the deployment is created with `el-<name of event listener>`
- mangledtag: task/pipeline/trigger names need to be dns names so no dots, ex in versions

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: simple-el
spec:
  entryPoints:
    - https
  routes:
    - kind: Rule
      match: Host(`build.seankhliao.com`)
      services:
        - kind: Service
          name: el-simple-el
          namespace: build
          port: 8080
  tls: {}
---
apiVersion: triggers.tekton.dev/v1alpha1
kind: EventListener
metadata:
  name: simple-el
spec:
  serviceAccountName: tekton-triggers-admin
  triggers:
    - name: simple-container
      interceptors:
        - github:
            secretRef:
              secretName: github-webhook-token
              secretKey: shared
            eventTypes:
        - cel:
            filter: "header.match('X-GitHub-Event', 'push') && (split(body.ref, '/')[1] == 'tags') && (body.repository.name in ['calproxy', 'goproxy', 'http-server', 'statslogger', 'vanity', 'webstyle'])"
            overlays:
              - key: extensions.tag_name
                expression: "split(body.ref, '/')[2]"
              - key: extensions.mangledtag
                expression: "split(split(body.ref, '/')[2], '.')[0]+'-'+split(split(body.ref, '/')[2], '.')[1]+'-'+split(split(body.ref, '/')[2], '.')[2]"
      bindings:
        - name: simple-container
      template:
        name: simple-container
```

##### _trigger_ binding

map of json to params

```yaml
apiVersion: triggers.tekton.dev/v1alpha1
kind: TriggerBinding
metadata:
  name: simple-container
spec:
  params:
    # https://developer.github.com/v3/activity/events/types/#pushevent
    - name: url
      value: $(body.repository.clone_url)
    - name: revision
      value: $(body.extensions.tag_name)
    - name: image
      value: $(body.repository.name)
    - name: mangledtag
      value: $(body.extensions.mangledtag)
```

##### _trigger_ template

- inline the pipeline declaration in the pipelinerun
- inline the task definition in the pipeline

```yaml
apiVersion: triggers.tekton.dev/v1alpha1
kind: TriggerTemplate
metadata:
  name: simple-container
spec:
  params:
    - name: url
      description: The git repository url
    - name: revision
      description: The git revision
    - name: image
      description: container image name
    - name: mangledtag
      description: used in naming
      default: $(uid)
  resourcetemplates:
    - apiVersion: tekton.dev/v1beta1
      kind: PipelineRun
      metadata:
        name: z-$(params.image)-$(params.mangledtag)
      spec:
        serviceAccountName: build-bot
        pipelineSpec:
          workspaces:
            - name: src
          tasks:
            - name: clone
              taskRef:
                name: git-clone
              params:
                - name: url
                  value: $(params.url)
                - name: revision
                  value: $(params.revision)
              workspaces:
                - name: src
                  workspace: src
            - name: build
              taskRef:
                name: kaniko
              runAfter:
                - clone
              params:
                - name: image
                  value: seankhliao/$(params.image)
                - name: tag
                  value: $(params.revision)
              workspaces:
                - name: src
                  workspace: src
        workspaces:
          - name: src
            persistentVolumeClaim:
              claimName: simple-container
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: simple-container
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
```

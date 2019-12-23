--- title
kustomize all the things
--- description
managing k8s configs with kustomize
--- main


[kustomize](https://github.com/kubernetes-sigs/kustomize)
is a cli built in (-ish) to `kubectl` that allows for more complicated workflows

As luck would have it,
I decided to start using it one day after the v0.2.1 release,
_bleeding edge_,
yay no docs

# magic and faerie dust

`kustomize` understands k8s config files,
at least the non CRD kinds.

so it:

- combines all files listed in `resources`
- if it's a dir, build using the `kustomization.yaml` inside
- apply patches (native, json, and shortcuts)
- output a keyname sorted canonical yaml (amazing for diffs)

# dir layout

```
service-root
├── base
│  ├── config.yaml
│  ├── deployment.yaml
│  ├── ingress.yaml
│  ├── kustomization.yaml
│  └── service.yaml
├── config.yaml
├── ingress.yaml
├── kustomization.yaml
├── preemptible.yaml
└── resources.yaml

```

or if you have multiple deployments

```
service-root
├── base
│  ├── config.yaml
│  ├── deployment.yaml
│  ├── ingress.yaml
│  ├── kustomization.yaml
│  └── service.yaml
└── overlays
   ├── deployment-1
   │  ├── config.yaml
   │  ├── ingress.yaml
   │  ├── kustomization.yaml
   │  ├── preemptible.yaml
   │  └── resources.yaml
   └── deployment-2
      ├── config.yaml
      ├── ingress.yaml
      ├── kustomization.yaml
      ├── preemptible.yaml
      └── resources.yaml
```

# _good_ use cases

dump upstream reference configs into the `base` dir,
apply all customizations in overlays

Or just to splut out your config into easily manageable parts,
a file for each bit you _care about_

# _problems_

- figure out CRDs
- make it more clear exactly how much is required to match for patching (all parents + metadata name?)

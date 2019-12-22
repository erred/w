--- title
cube uhhm nates.md
--- description
k8s concepts
--- main


# kubternetes

the thing that eats up all your computers

## cluster

collection of resources (computers) managed by kubernetes

## node

a single compute resource (vm or hardware computer)
to schedule tasks on

## containers

smallest packaging unit,
everything is built on top

## pod

smallest schedulable unit,
collection of _containers_ that run together,
shares a filesystem and network stack

## deployment

collection of identical _pods_

## service

_unified interface_ to the multiple pods in a deployment

## ingress

_exposed_ interface to the outside world,
read and fulfilled by ingress controllers

## ingressroute

fancier version of _ingress_

## ingress controller

edge router that sits between cluster and outside world

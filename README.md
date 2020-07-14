# Kleaner

![](https://img.shields.io/github/go-mod/go-version/zerodayyy/kleaner?style=flat-square)
![](https://img.shields.io/docker/cloud/build/zerodayyy/kleaner?style=flat-square)

A Kubernetes operator that deletes evicted pods

##### Why?

Sometimes an application can have a severe memory leak, which will cause it to consume more and more RAM over time, and eventually get evicted. Evicted pods are not cleaned up automatically, which creates visual noise in CLI output and tools like ArgoCD. This tool is a quick fix for such problem.

##### How?

The heart of this tool is a small piece of Go code, which scans all namespaces for pods with `Evicted` in their status, and deletes them.
This code is run as a Kubernetes CronJob.

### Installation

The recommended way of using Kleaner is installing it using Helm:

```
helm install kleaner helm/kleaner
```

You can also use tools like ArgoCD or Flux.


### Configuration

There are a few variables exposed through `values.yaml` file:

- **`interval`** (*valid values: 1–59*) — how many minutes to wait between runs
- **`namespaces`** — list of namespaces to run in, defaults to all namespaces if not specified
- **`debug`** — if something works not as you expected, this can give you more insight (Go knowledge advised)


### Changelog

##### v1.0.0 — *July 14, 2020*
Initial release

Demo
----

Run our custom Flux K8s controller which connects to the Flux source
controller to monitor this repo. Adapted from

- https://fluxcd.io/docs/gitops-toolkit/source-watcher/


### Prerequisites

Install

* Nix - https://nixos.org/guides/install-nix.html
* docker >= 19.03


### Setting up the local K8s cluster

Open a terminal in the repo root dir, then

```bash
$ nix-shell
$ flux check --pre
$ kind create cluster --name dev
$ flux install \
    --namespace=flux-system \
    --network-policy=false \
    --components=source-controller
$ kubectl -n flux-system port-forward svc/source-controller 8181:80
```

### Running our custom controller

We'll run it outside the cluster so we can debug it too easily.
Open another terminal in the repo root dir and run

```bash
$ nix-shell
$ export SOURCE_HOST=localhost:8181
$ make run
```


### Simulating commits to GitHub

Open yet another terminal in the repo root dir, then create a test GitHub
source within Flux to monitor our own repo

```bash
$ nix-shell
$ flux create source git test \
    --url=https://github.com/c0c0n3/source-watcher \
    --tag=test.0.0.1
```

Now if you switch back to the terminal running our custom controller, you
should be able to see it processing the files in the `_deployment_` dir as
it was at tag `test.0.0.1`. It only logs the OSM command it would call without
actually running it since we don't have a working OSM cluster yet.

When it's done processing, we can simulate a new commit by switching to
tag `test.0.0.2`

```bash
$ flux create source git test \
    --url=https://github.com/c0c0n3/source-watcher \
    --tag=test.0.0.2
```

Again you should see OSM command lines being logged by our custom controller.


### Clean up

Kill all the processes running in your terminals, then get rid of the
cluster too

```bash
$ kind delete cluster --name dev
```

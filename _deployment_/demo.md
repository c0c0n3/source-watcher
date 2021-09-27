Demo
----

Run our custom Flux K8s controller which connects to the Flux source
controller to monitor this repo. Adapted from

- https://fluxcd.io/docs/gitops-toolkit/source-watcher/


### Prerequisites

Install

* Nix - https://nixos.org/guides/install-nix.html
* docker >= 19.03
* multipass >= 1.6.2

Keep in mind you're going to need a beefy box to run this demo smoothly.
With lots of effort and patience, I've managed to run it on my 4 core, 8
GB RAM laptop but my guess is that you'd need at box with at least double
that horse power.


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


### Setting up the local OSM cluster

Spin up an Ubuntu 18.04 VM with Multipass and install OSM in it by following
the steps in [multipass.install.sh][osm-install]. Notice I had to patch the
OSM install scripts to make it work.

When done, upload the OSM OpenLDAP packages we're going to use to create
NS instances:

```bash
$ cd _tmp/osm-pkgs
$ multipass mount ./ osm:/mnt/osm-pkgs
$ multipass shell osm
% cd /mnt/osm-pkgs
% osm nfpkg-create openldap_knf.tar.gz
% osm nspkg-create openldap_ns.tar.gz
% exit
```

Note down the VM's IPv4 address where the OSM NBI can be accessed:

```bash
$ multipass info osm
```

It should be the first one on the list the `192.168.*` one.


### Running our custom controller

We'll run it outside the cluster so we can debug it easily and we won't
have to set up a bridge between the K8s Kind cluster and the Multipass
OSM VM. Open another terminal in the repo root dir and run

```bash
$ nix-shell
$ export SOURCE_HOST=localhost:8181
$ make run
```

We're going to have our controller create and update an NS instance
by looking at the OSM GitOps files in this repo's `_deployment_` dir
at tags `test.0.0.3` and `test.0.0.4`. The `osm_ops_config.yaml` file
is the same for both tags and points to an NBI connection file sitting
on the same box where the controller runs: `/tmp/osm_ops_secret.yaml`.
You need to create this file with the following content:

```yaml
hostname: 192.168.64.19:80
project: admin
user: admin
password: admin
```

but replace `192.168.64.19` with the OSM IP address you noted down
earlier.


### Simulating commits to GitHub

Open yet another terminal in the repo root dir, then create a test GitHub
source within Flux to monitor our own repo

```bash
$ nix-shell
$ flux create source git test \
    --url=https://github.com/c0c0n3/source-watcher \
    --tag=test.0.0.3
```

Now if you switch back to the terminal running our custom controller, you
should be able to see it processing the files in the `_deployment_` dir as
it was at tag `test.0.0.3`. It should call OSM NBI to create an NS instance
using the OSM OpenLDAP package we uploaded earlier with two replicas as
specified in the `ldap.ops.yaml` in `_deployment_/kdu`. It's going to take
a while for the deployment state to reflect in the OSM Web UI, but you
can check what's going on under the bonnet by shelling into the OSM VM

```bash
$ multipass shell osm
% kubectl get ns
#              ^ pick the one that looks like an UUID
% kubectl -n fada443a-905c-4241-8a33-4dcdbdac55e7 get pods
# ... you should see two pods being created for the OpenLDAP service
```

Wait until the two K8s pods are up and running and the deployment state
got updated in the OSM Web UI. Then, we can simulate a new commit by
switching to tag `test.0.0.4`

```bash
$ flux create source git test \
    --url=https://github.com/c0c0n3/source-watcher \
    --tag=test.0.0.4
```

The content of `ldap.ops.yaml` at tag `test.0.0.4` is the same as that
of tag `test.0.0.3` except for the replica count which is `1`. So you
should see that eventually your NS instance for the OpenLDAP service
gets scaled down to one K8s pod. Be patient, unless you've got a beefy
box, this too will take a while.


### Clean up

Kill all the processes running in your terminals, then get rid of the
K8s cluster

```bash
$ kind delete cluster --name dev
```

Finally, zap the Multipass VM with the clean up commands you'll find
in [multipass.install.sh][osm-install].




[osm-install]: ../_tmp/osm-install/multipass.install.sh

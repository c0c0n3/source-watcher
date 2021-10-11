Introduction
------------
> The why, the what and the how.

This introductory section first touches on project motivation and goals,
then goes on to sketching the architecture conceptual model and how
it has been implemented through a Kubernetes operator plugged into
the FluxCD framework.


### Motivation and goals

The [Affordable5G][a5g] project adopts Open Source MANO (OSM) to
virtualise and orchestrate network functions, simplify infrastructure
operation, and achieve faster service deployment. One of the Affordable5G
objectives is to explore the continuous delivery of services through
GitOps workflows whereby the state of an OSM Kubernetes deployment
is described by version-controlled text files which a tool then interprets
to achieve the desired deployment state in the live OSM cluster.
Although OSM features a sophisticated toolset for the packaging,
deployment and operation of services, GitOps workflows for Kubernetes
network functions (KNFs) are not fully supported yet. Hence the need,
within Affordable5G, of a software to complement OSM’s capabilities
with GitOps.

Automated, version-controlled service delivery has several benefits.
Automation shortens deployment time and ensures reproducibility of
deployment states. In turn, reproducibility dramatically reduces the
time needed to recover from severe production incidents caused by
faulty deployments as the OSM cluster can swiftly be reverted to a
previous, known-to-be-working deployment state stored in the Git
repository. Thus, overall cluster stability and service availability
are enhanced. Moreover, the Git repository stores information about
who modified the OSM cluster state when, thus furnishing an audit
trail that may help to detect security breaches and failure to comply
with regulations such as GDPR.


### Conceptual model

OSM Ops is a cloud-native micro-service to implement GitOps workflows
within OSM. The basic idea is to describe the state of an OSM deployment
through version-controlled text files hosted in an online Git repository.
Each file declares a desired instantiation and runtime configuration
for some of the services in a specified OSM cluster. Collectively,
the files at a given Git revision describe the deployment state of
the these services at a certain point in time. OSM Ops monitors the
Git repository in order to automatically reconcile the desired deployment
state with the actual live state of the OSM cluster. OSM Ops is implemented
as a [Kubernetes][k8s] operator that plugs into the [FluxCD][flux]
framework in order to leverage the rich Kubernetes/FluxCD GitOps
ecosystem. The following visual illustrates the context in which
OSM Ops operates and exemplifies the GitOps workflow resulting in
the creation and update of KNFs from version-controlled deployment
declarations.

![Architecture context diagram.][dia.ctx]

**TODO**: narrative


### Implementation overview

Having defined the abstract ideas, we are now ready to explain how
they have been realised. In a nutshell, OSM Ops is a Kubernetes operator
that gets notified of any changes to an online Git repository monitored
by FluxCD and then uses OSM’s north-bound interface (NBI) to realise
the KNF deployment configurations found in that repository.

These deployment configurations are declared through OSM Ops YAML
files. Each file specifies a desired instantiation and runtime configuration
(e.g. number of replicas) of a KNF previously defined within OSM by
installing suitable OSM descriptor packages, Helm charts, etc. For
example, the following YAML file demands that the live OSM cluster
run a 2-replica NS instance called `ldap` within the VIM identified
by the given VIM account and that the service be configured according
to the definitions found in the named OSM descriptors—the referenced
NSD, VNFD and KDU are actually defined in the OpenLDAP OSM packages
published by Telefonica.

```yaml
kind: NsInstance
name: ldap
description: Demo LDAP NS instance
nsdName: openldap_ns
vnfName: openldap
vimAccountName: mylocation1
kdu:
    name: ldap
    params:
        replicaCount: "2"
```

Source Controller is a FluxCD service that, among other things, manages
interactions with online Git repositories—e.g. repositories hosted
on GitHub or GitLab. OSM Ops depends on it both for monitoring repositories
and for fetching the repository content at a given revision. Source
Controller arranges a Kubernetes custom resource for each repository
that it monitors and then polls each repository to detect new revisions.
As soon as a new revision becomes available, Source Controller updates
the corresponding Git repository custom resource in Kubernetes.

OSM Ops implements the Kubernetes Operator interface to get notified
of any changes to Git repository custom resources. Thus, soon after
Source Controller updates a Git repository custom resource, Kubernetes
dispatches an update event to OSM Ops. This arrangement is akin to
the publish-subscribe pattern often found in messaging systems: Source
Controller, the publisher, sends a message to Kubernetes, the broker,
which results in the broker notifying OSM Ops, the subscriber. The
publisher and the subscriber have no knowledge of each other (no space
coupling) and communication is asynchronous (no time coupling).

At this point, OSM Ops enters the reconcile phase in which it tries
to align the deployment state declared in the OSM Ops YAML files with
that of the live OSM cluster. It fetches the content of the notified
Git revision from Source Controller as a tarball and then uses OSM's
NBI to transition the OSM cluster to the deployment state declared
in the OSM Ops YAML files found in the tarball. For each file, OSM
Ops determines whether to create or update the KNF specified in the
file and then configures it according to the KNF parameters given
in that file.

The UML communication diagram below summarises the typical workflow
through which OSM Ops turns the deployment state declared in a Git
repository into actual NS instances. The workflow begins with a system
administrator pushing a new revision, `v6`, to the online Git repository.
It then continues as just explained, with Source Controller updating
the Git custom resource, Kubernetes notifying OSM Ops and OSM Ops
calling the NBI to achieve the deployment state declared in `v6`.

![Implementation overview.][dia.impl]


### Rationale

What is the rationale behind our design decisions? A few explanatory
words are in order.

**TODO**
- evaluated two leading GitOps solutions: ArgoCD & FluxCD
- similar capabilities but ArgoCD comes with powerful UI
- convergence: the two projects will likely be merged in the
  future—ref merger plans
- FluxCD has better docs about extending it with custom functionality
  which is what in the end tipped the balance in its favour
- Go was a natural PL choice b/c of FluxCD and K8s libs are both
  written in Go




[a5g]: https://www.affordable5g.eu/
    "Affordable5G"
[dia.ctx]: ./arch.context.png
[dia.impl]: ./arch.impl-overview.png
[flux]: https://fluxcd.io/
    "Flux - the GitOps family of projects"
[k8s]: https://en.wikipedia.org/wiki/Kubernetes
    "Kubernetes"
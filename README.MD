## Distributed Cowboys

Deployment manifests are in this repository. Algorithm implementation is at https://github.com/josefkarasek/distributed-cowboys


### Goal
A little fun project to synchronize simultaneous processes.
N number of _cowboys_ shoot at each other, last man standing reports his survival at the end.

### Design
The given problem is a Consensus problem - solution which would most likely work would import
existing Consensus implementation, such as etcd (raft). But that would be boring.

On about 350 lines of code I tried to solve this task with only:
* golang channels
* mutexes

This, of course, cannot satisfy the criteria of Safety and Liveness, but it was fun working on!
I also tried to implement this within a time limit of 4 hours.

State:
Initial state is stored in etcd via Custom Resource `Shootout`. This state is distributed to cowboys
in raw JSON form via a Config Map, that load it in memory and maintain it in a later formed ring.

Node roles:
* leader: one node is designated as leader, this is static/hard-coded designation
* follower: all other nodes are followers

The leader node:
* maintains timing (1 second beats)
* initiates state synchronization

The nodes form a ring, in which the leader node takes the head/tail role.
While cowboys can shoot at any _live_ cowboy (anycast), state is synchronized in the formed ring
and only after the end of a beat. This can lead to a cowboy being struck and killed more
than once during a beat.

### Implementation

Each cowboy runs as a separate process, in a kubernetes Pod, managed by a Job and is reachable through a Service.
Jobs are spawned by a Operator, based on configuration.

Example config:
```yaml
apiVersion: example.com/v1alpha1
kind: Shootout
metadata:
  name: full
spec:
  cowboys: >
    [
      {
      "name": "John",
      "health": 10,
      "damage": 1
      },
      {
      "name": "Bill",
      "health": 8,
      "damage": 2
      },
      {
      "name": "Sam",
      "health": 10,
      "damage": 1
      },
      {
      "name": "Peter",
      "health": 5,
      "damage": 3
      },
      {
      "name": "Philip",
      "health": 15,
      "damage": 1
      }
    ]

```
Cowboys use HTTP to communicate with each other.

### Instructions for testing

```
$ kind create cluster --name=devenv
$ kubectl create namespace cowboys
$ kubectl config set-context --current --namespace=cowboys
$ kubectl kustomize build config/crd/ | kubectl apply -f -
$ go run main.go
$ k apply -f config/samples/_v1alpha1_shootout.yaml
$ stern -n cowboys full
```
As there's multiple pods running at the same time, I recommend using [stern](https://github.com/wercker/stern) to tail their logs simultaneously
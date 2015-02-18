Disco
====

_A Thing that links Docker containers together._

* * *

![disco arch](https://dl.dropboxusercontent.com/u/42154947/blog%20pics/disco.png)

Disco is an experiment in service discovery, cross-host connectivity, and eventually scheduling.  Strictly, and officially, just an experiment.

## Types

### `disco.Container`

Disco provides an abstract Container type which can be used across providers.  Here you will find only Docker and the native Disco containers conforming to this type, but it is designed in such a way that any provider could be swapped in given some valid marshaling rules.  Aside from the core information Disco needs to keep track of a container, `disco.Container` also holds a collection of the Link IDs which pertain to it.  These Links are separated in a map where the key is either source or target, and the value is an array of Link ids.

### `disco.Node`

Represents a node in the Disco cluster.  It keeps track of an ephemeral node id, and the Node's IP addresses

### `disco.Link`

Link represents the connection between two containers.

## Modules

Disco has primarily 2 parts: The daemon, and the CLI.

### Daemon

The daemon is made up of an API server, and a poller which polls the Docker socket on the local node.

#### Poller

The poller takes a config, and starts a loop that polls and collects data about Docker _and_ Disco containers, their exposed ports, and any links the containers may have.  The data from each is marshalled into a common type, `disco.Container`, and the differences between the two APIs are reconciled.  Docker is taken as the authority in conflicts, as it is the SoT on container state.  The Poller takes the resultant data and adds any Links that may be relevant to the containers.  It then writes this data to the Disco API.

#### API

The Disco API is a namespace-based Unix Socket server.  The Socket takes requests via a namespace, say `/disco/api/get_containers`, reads from etcd or the Docker daemon, and returns a JSON byte array.  In the case that a Payload is given, it is separated by a `\n` newline character, where a request would look like this:

```
/disco/api/add_container/{a_uuid}
{"Name: "test", "NodeId": a_different_uuid}
```


### CLI

The CLI interacts, through a client, with the Disco API.  It takes user provided data and links containers.  It also lists containers with their Link state, and finally, it handles other management-based queries such as getting the local Node ID.

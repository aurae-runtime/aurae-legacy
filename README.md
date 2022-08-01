# Aurae

Experimental peer-to-peer application runtime. 

 - [X] NAT Traversal using [libp2p](https://github.com/libp2p/specs/tree/master/identify#identify-v100) identify v.1.0.0 
 - [X] MicroVM management using [Firecracker](https://github.com/firecracker-microvm/firecracker)
 - [X] Service Discovery using IPFS Distributed Hash Table (DHT) 
 - [X] Local socket primitives
 - [X] Fuse filesystem mount over the database

## Empty Loop Architecture

> Everything and nothing all at once.

Minimal scope for `auraed` which runs as PID 1 on a Linux kernel and will ultimately replace `systemd`. 

For reliability purposes `auraed` will *always* start. However without TLS certificate material in place and registered with the daemon it will be an empty loop. 

Additionally, after `auraed` has registered certificate material it will, by default, do nothing other than expose its registered endpoints over gRPC on a Unix Domain Socket. 

We only have 2 primary flags, that have reasonable defaults.

``` 
 --sock (default: "/run/aurae.sock")
 --key  (default: "home/user/.ssh/id_aura")
```

This empty loop status of the daemon is it's default and most stable state. Once the loop is `healthy` work and capabilities can be registered to the daemon and the system can be composed at runtime.

## Spaces 

Aurae maintains a strict isolation boundary between two spaces that is seperated by the aurae socket API.

 - Daemon space
 - Aurae space 

Deamon space is most closely related to what Linux refers to as "kernel space". 
The daemon is responsible for executing various tasks on behalf of clients of the system.
Unlike Linux, the daemon space is not limited to the some constraints of a hardware kernel. 
Memory management, process management, IPC, and basic operating system functionality is assumed.

Aurae space is most closely related to what Linux refers to as "User space".
This is the part of the system that is destined for users of the system. 
Here is where applications exist, and where multitenancy is enforced.

First Principal: Everything in Aurae space exists in an isolation zone and MUST access the system via the local unix domain socket.

## Capabilities API

Capabilities are similar to subsystems and will come with various sub-resources that can be implemented and registered against the empty loop.

### Compute 

 - Runtime
 - Schedule
 - Observe

### Storage

 - Mount
 - Network
 - Object
 - File
 - Secrets

### Network

 - Ingress
 - Egress 
 - Peer
 - Firewall
 
 ## Aurae connection Syntax
 
 Imagine an `nginx` server running on a node in the network.
 
 ```
 <entity>@<peer>@<domain>
 ```
 
 Where domain is a registered domain that allows peers to identify each other in the public internet.

 
 ## Relationship with POSIX
 
The project will be responsible for maintaing various POSIX compliant CLI tools.
 
These will leverage the `auare.sock` as a UNIX pipe which can be managed in a traditional POSIX environment.

We will house a series of improved Core utils that can be leveraged with Aurae nodes

```
auraefs # Used to mount filesystems that are registered with the `Storage` capability
ascp    # Traditional scp but leverages the aurae syntax
assh    # Traditional ssh but leverages the aurae socket for connections to nodes
```

 
 

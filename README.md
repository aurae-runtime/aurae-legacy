# Aurae

Experimental peer-to-peer application runtime. 

 - [X] NAT Traversal using [libp2p](https://github.com/libp2p/specs/tree/master/identify#identify-v100) identify v.1.0.0 
 - [X] MicroVM management using [Firecracker](https://github.com/firecracker-microvm/firecracker)
 - [X] Service Discovery using IPFS Distributed Hash Table (DHT) 
 - [X] Local socket primitives
 - [X] Fuse filesystem mount over the database

## Empty Loop Architecture

Minimal scope for `auraed` which runs as PID 1 on a Linux kernel and will ultimately replace `systemd`. 

For reliability purposes `auraed` will *always* start. However without TLS certificate material in place and registered with the daemon it will be an empty loop. 

Additionally after `auraed` has registered certificate material it will, by default, do nothing other than expose its registered endpoints over gRPC on a Unix Domain Socket. 

``` 
/run/aurae.sock
/pki/aurae.pem
/pki/aurae.pub
```

This empty loop status of the daemon is it's default and most stable state. Once the loop is `healthy` work and capabilities can be registered to the daemon and the system can be composed at runtime.

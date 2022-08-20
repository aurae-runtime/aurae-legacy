# Stories 

### Scheduling a Pod 

We have a use case is to manage a [mastodon](https://joinmastodon.org/) instance with `aurae`.

 - Need to decide if "pod" is even what we want to call it or do
 - Need to manage postgresql (and other databases)
 - Need to declare (or at least encapsulate) our app
 - Make it possible to replace systemd
 - Connect to firecracker and manage workloads

User experience for scheduling a pod

 - Get aurae status (all systems go)
 - Run pod
 - Get pod definition 

How do we want to programmatically do this?

 - Using Go
 - Using Cue
 - Using *New DSL*
 - Using YAML
 - Using JSON

Requirements for user scheduling:

 - Turing complete 
 - Statically linked (discreet executable)
 - We need a runtime and a daemon 

Conclusion from this thought exercise:

 - We need a "runtime" that will process turing complete code (script, procedural, similar to python, ruby, etc)
 - We need a "daemon" that will be the new systemcalls of the distributed system
 - The gRPC /rpc directory are the core APIs, the runtime is just fluff

Decisions:

 - Right now we are considering using a Rust DSL for the end-user logic
 - Because the Rust DSL will just wrap the aurae API, Cue and other languages are also relevant at this point


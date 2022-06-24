# aurae

A Simple Distributed System. Built for App teams.

Reverse your relationship with the cloud. Shift your application posture from reliant upon cloud dependencies, to harvesting cloud services on your terms.



### Constraints 
 
 - aurae must run in local mode only
 - aurae must run in remote mode only
 - aurae must be simple
 - All applications MUST translate to capabilities
 - aurae status should always work
 - aurae init should mostly work

### aurae Run

aurae uses a 3 part naming convention to express `image@node@domain` where `node` and `domain` are assumed if not provided.

``` 
aurae run nginx                    # Run nginx locally. No different than docker run nginx.
aurae run nginx@localhost          # Same as above.
aurae run nginx@emily@nivenly.com  # Pin nginx to the "emily" node in the "nivenly.com" domain.
aurae run nginx@nivenly.com        # Run nginx anywhere that is available in the "nivenly.com" domain.
```

# Quickstart 

Start aurae on a Linux system.

``` 
aurae status
aurae init 
aurae status
```

# Capabilities

aurae is built on capabilities. 

Capabilities are hierarchical tasks you and your team might need to perform.

# Features

 - Application abstraction
 - Run containers
 - Run micro VMs
 - Manage ingress 
 - Mange DNS

# Development Notes

 - [ ] Structured output. Everything should be structured and encodedable. EG: json
 - [ ] No locks. We hold opinion on filesystem locks. It is an anti pattern. 
 - [ ] 
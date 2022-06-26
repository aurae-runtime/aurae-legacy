# Memory Filesystem 

This is a stateless in-memory filesystem that (by design) will never persist past the current execution context.

This filesystem should be extremely simplified compared to the POSIX filesystem implementations. The only caveat
to this filesystem is that a directory can also have content. 

For example imagine a structure 

``` 
tree
├── beeps
│     └── meeps
└── boops
```

All nodes can have file content.

The "node" `beeps` can contain file content even though it is a directory, as well as the "nodes" `meeps` and `boops`.

### Package State 

The `memfs` package has package level state. 

This is because we know there will only be a single root node in our tree at one time. 

We will need to be able to set values in the tree arbitrarily.
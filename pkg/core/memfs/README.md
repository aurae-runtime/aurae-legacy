# Memory Filesystem 

This is a stateless in-memory filesystem that (by design) will never persist past the current execution context.

The filesystem is simple, is not opinionated and does not store POSIX filesystem meta values such as file attributes, inode details, permissions, users, or groups.

``` 
root
├── dir1
│     └── file1
└── file0
```

### Package State 

The `memfs` package has package level state. 

This is because we know there will only be a single root node in our tree at one time. 

We will need to be able to set values in the tree arbitrarily.
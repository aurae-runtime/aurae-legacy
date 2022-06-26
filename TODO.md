# (aurae) Client 

 - [ ] Build tree 

# (auraefs) Filesystem

 - [ ] The shutdown mechanism can hang if somebody is inside one of the mounted dirs
 - [ ] We need to manage the `client.Client` better
 - [ ] Refactor the filesystem now that we understand what we are doing! The filesystem should "watch" the client!

# (auraed) Daemon

- [ ] Fix the watcher, it needs to be modular. The watcher is how the filesystem will register files!
- [ ] Fix user/group settings (cross-platform). We need to domain socket to be clean. The client should work as your user.
- [ ] List can be cleaner
- [ ] In memfs we need to manage the `rootNode` better 

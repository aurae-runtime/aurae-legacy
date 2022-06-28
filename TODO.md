# (aurae) Client 


# (auraefs) Filesystem

 - [ ] We need to manage the `client.Client` better. We need an audit system, and to manage nil clients/reconnects.

# (auraed) Daemon

- [ ] Fix the watcher, it needs to be modular. The watcher is how the filesystem will register files!
- [ ] Fix user/group settings (cross-platform). We need to domain socket to be clean. The client should work as your user.
- [ ] List can be cleaner
- [ ] In memfs we need to manage the `rootNode` better 


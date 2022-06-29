# aurae 

Aurae is a decentralized peer-to-peer runtime mesh that works well with container runtimes, and micro VMs.

The client command line tool for the `aurae` suite.

### Run a Container 

Download `aurae` and install the binaries.

```bash
$ git clone git@github.com:kris-nova/aurae.git
$ cd aurae
$ make
$ sudo make install
```

Initialize the local aurae node with a hostname and optional domain name.

The hostname can be anything, and does not need to match the hostname of your computer.

The domain name should be a domain you own and plan on exposing on the internet. If you are unsure, leave it blank for now.

```bash
$ aurae up <hostname>[@<domain>]
$ aurae up alice 
$ aurae up alice@nivenly.com
```

Aurae is now running on the system and ready to serve requests. 

```bash
$ aurae run <image>@<hostname>@<domain> [options]
$ aurae run nginx
```

Your container is now running on your local `aurae` node.

Identify the peering token on a `aurae` node (Emily).

```bash
# Will return your peering token if no token is passed
# Hostname: Emily
$ aurae proxy
<token>
```

You can peer your node to another node using `aurae peer`.

```bash
# When a token is passed, aurae attempts to join Alice -> Emily
# Hostname: Alice
$ aurae proxy <token>
```

Peering creates bidirectional executive trust.
Either peer can now control each other.

You can now run a container on any node in your domain using the same parlance as before.

```bash
$ aurae run <image>@<hostname>@<domain> [options]

# Alice can schedule on Emily
$ aurae run nginx@emily 

# Emily can schedule on Alice
$ aurae run nginx@alice
```

### Explore your virtual filesystem 

All data in `aurae` is simplified and will be familiar to anyone comfortable on a command line.

Each node stores its own state. In order to mutate state on a node the client must route directly to the node with the correct permissions.

There is an intuitive way of exploring the filesystem and potentially making changes to a node's state.

Mounting `auraefs` on the host wil use the same cert material the `aurae` client uses. You can mount the filesystem directly to the host using `libfuse`.

``` 
auraefs mount <path>
auraefs mount /aurae
```

You can now interact directly with your mTLS encrypted `aurae` state using any of the methods available on your computer.

### Exploring Auraespace with Ash (Aurae SHell)

Additionally, it is possible to run a UNIX style shell to navigate the entire `aurae` mesh.

We call the mesh space "Auraespace" as it resembles userspace in Linux.

The `ash` shell or `aurae shell` runs as a container and be scheduled in `aurae`.

```bash
ash <hostname>
```

Once inside an `ash` shell you will be inside a container.

Peers are listed in the `/peers` directory. 

The local `auraefs` filesystem is located in the `/aurae` directory. 

You can enter another peer's space by running a new instance of `ash` and pointing directly at the `peer/hostname` file. 

```bash
# Running on localhost
$ ash # No arguments assumes localhost (Alice)

# Now running ash on Alice
$ ls /aurae # Show Alice state
$ ash /proxy/emily

# Now running ash on Emily
$ ls /aurae # Show Emily state
$ exit

# Now back on Alice
$ exit

# Now back on localhost
$ 
```
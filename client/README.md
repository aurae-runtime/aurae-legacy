# Client 

The `aurae` client is a client that is composed of the various RPC services known to the `aurae` server, as well as routing and service discovery logic. 

The (outward) goal of the client is to be as user-friendly as possible, while adhering to best security practices.

The (inward) goal of the client is to make mTLS management easier, while serving as a proxy for client agents to other nodes in the mesh.


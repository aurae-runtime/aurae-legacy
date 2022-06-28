# Core Service

The `CoreService` of `Aurae` is a simple key/value storage mechanism that assumes a filesystem like experience based on Linux parlance.

The `CoreService` is liable to mutate data, which is by design.

For example if a record `key` is saved with content `value` the name of the record is mutated to `/key` and is saved in the top level of an in-memory tree.

The `CoreService` values simplicity over expression.

The `CoreService` does not attempt to implement POSIX, but rather play nicely with POSIX.

### Database implementations

 - [SpiceDB](https://github.com/authzed/spicedb/blob/main/proto/internal/core/v1/core.proto)
 - [etcd](https://github.com/etcd-io/etcd/blob/main/api/etcdserverpb/rpc.proto)
# auraeFS 

The auraeFS is a pseudo virtual filesystem.

### Notes

The filesystem should be transparent, however mutations should be locked via mTLS.

### /app

The `/app` directory is where application owners define applications. 

### /cap (read-only)

The `/cap` directory is a virtual directory that will be empty by default.

As the `aurae` daemon runs, it will identify any configured infrastructure, and reflect the capabilities in this directory.

Applications require capabilities in order to run. 

### /infra 

The `/infra` directory is where infrastructure owners define infrastructure.



module github.com/kris-nova/aurae

go 1.18

// Taken from https://github.com/ory/kratos/blob/master/go.mod
replace (
	//github.com/bradleyjkemp/cupaloy/v2 => github.com/aeneasr/cupaloy/v2 v2.6.1-0.20210924214125-3dfdd01210a3
	github.com/gorilla/sessions => github.com/ory/sessions v1.2.2-0.20220110165800-b09c17334dc2
	github.com/jackc/pgconn => github.com/jackc/pgconn v1.10.1-0.20211002123621-290ee79d1e8d
	github.com/knadh/koanf => github.com/aeneasr/koanf v0.14.1-0.20211230115640-aa3902b3267a
	// github.com/luna-duclos/instrumentedsql => github.com/ory/instrumentedsql v1.2.0
	// github.com/luna-duclos/instrumentedsql/opentracing => github.com/ory/instrumentedsql/opentracing v0.0.0-20210903114257-c8963b546c5c
	github.com/mattn/go-sqlite3 => github.com/mattn/go-sqlite3 v1.14.7-0.20210414154423-1157a4212dcb
	//github.com/oleiade/reflections => github.com/oleiade/reflections v1.0.1
	// Use the internal httpclient which can be generated in this codebase but mark it as the
	// official SDK, allowing for the Ory CLI to consume Ory Kratos' CLI commands.
	//github.com/ory/kratos-client-go => ./internal/httpclient

	//go.mongodb.org/mongo-driver => go.mongodb.org/mongo-driver v1.4.6
	golang.org/x/sys => golang.org/x/sys v0.0.0-20220319134239-a9b59b0215f8
//gopkg.in/DataDog/dd-trace-go.v1 => gopkg.in/DataDog/dd-trace-go.v1 v1.27.1-0.20201005154917-54b73b3e126a
)

require (
	github.com/fsnotify/fsnotify v1.5.4
	github.com/hanwen/go-fuse/v2 v2.1.0
	github.com/sirupsen/logrus v1.8.1
	github.com/urfave/cli/v2 v2.4.1
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/stretchr/testify v1.7.1 // indirect
	golang.org/x/net v0.0.0-20211020060615-d418f374d309 // indirect
	golang.org/x/sys v0.0.0-20220517195934-5e4e11fc645e // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20220517211312-f3a8303e98df // indirect
	google.golang.org/genproto v0.0.0-20211208223120-3a66f561d7aa // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

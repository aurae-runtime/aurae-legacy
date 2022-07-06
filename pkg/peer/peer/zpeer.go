package zpeer

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	golog "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
)

// Notes
//
// Listen Port: 1235
// Client Port: 1234
//
// TODO WE need to manage golog for the sub-libraries here 	//golog.SetAllLoggers(golog.LevelInfo) // Change to INFO for extra info

const (
	DefaultGenerateKeyPairBits int = 2048
	DefaultListenPort          int = 8709
	DefaultPeerPort            int = 8708
)

var emptyKey crypto.PrivKey = &crypto.Ed25519PrivateKey{}

type Peer struct {
	uniqKey     crypto.PrivKey
	routedHost  rhost.RoutedHost
	host        host.Host
	established bool
}

func NewPeer() *Peer {
	golog.SetAllLoggers(golog.LevelPanic)
	golog.SetupLogging(golog.Config{
		Stdout: false,
		Stderr: false,
	})
	randSeeder := rand.Reader
	key, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, DefaultGenerateKeyPairBits, randSeeder)
	if err != nil {
		logrus.Errorf("unable to GenerateKeyPair for new peer: %v", err)
		key = emptyKey
	}
	return &Peer{
		uniqKey:     key,
		established: false,
	}
}

func (p *Peer) Establish(ctx context.Context, offset int) error {
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", DefaultListenPort+offset)),
		libp2p.Identity(p.uniqKey),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
	}

	basicHost, err := libp2p.New(opts...)
	if err != nil {
		return err
	}
	p.host = basicHost

	dstore := dsync.MutexWrap(ds.NewMapDatastore())

	// Make the DHT
	dht := dht.NewDHT(ctx, basicHost, dstore)

	// Make the routed host
	routedHost := rhost.Wrap(basicHost, dht)
	p.routedHost = *routedHost

	// connect to the chosen ipfs nodes
	err = bootstrapConnect(ctx, routedHost, IPFS_PEERS)
	if err != nil {
		return err
	}

	// Bootstrap the host
	err = dht.Bootstrap(ctx)
	if err != nil {
		return err
	}

	// Build host multiaddress
	//hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", routedHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	// addr := routedHost.Addrs()[0]
	//addrs := routedHost.Addrs()
	////log.Println("I can be reached at:")
	//for _, addr := range addrs {
	//	//log.Println(addr.Encapsulate(hostAddr))
	//}

	// targetF = Pretty()
	//log.Printf("Now run \"aurae -d %s%s\" on a different terminal\n",  routedHost.ID().Pretty())
	logrus.Infof("ID: %s", routedHost.ID().Pretty())
	p.established = true
	return nil
}

func (p *Peer) To(peerID string) error {
	if !p.established {
		return fmt.Errorf("unable to stream, first establish in the mesh")
	}
	p.routedHost.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
		//log.Println("Got a new stream!")
		if err := doEcho(s); err != nil {
			//log.Println(err)
			s.Reset()
		} else {
			s.Close()
		}
	})

	id, err := peer.Decode(peerID)
	if err != nil {
		return err
	}

	s, err := p.routedHost.NewStream(context.Background(), id, "/echo/1.0.0")
	if err != nil {
		return err
	}

	// Echo hello world
	_, err = s.Write([]byte("Hello, world!\n"))
	if err != nil {
		return err
	}

	out, err := ioutil.ReadAll(s)
	if err != nil {
		return err
	}

	//log.Printf("read reply: %q\n", out)
	logrus.Infof("%q", out)
	return nil
}

func (p *Peer) Stream() error {
	if !p.established {
		return fmt.Errorf("unable to stream, first establish in the mesh")
	}
	p.routedHost.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
		log.Println("Got a new stream!")
		if err := doEcho(s); err != nil {
			log.Println(err)
			s.Reset()
		} else {
			s.Close()
		}
	})
	select {} // hang forever
}

// doEcho reads a line of data from a stream and writes it back
func doEcho(s network.Stream) error {
	buf := bufio.NewReader(s)
	str, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	log.Printf("read: %s\n", str)
	_, err = s.Write([]byte(str))
	return err
}

package zpeer

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
)

// Notes
//
// Listen Port: 1235
// Client Port: 1234

const (
	DefaultGenerateKeyPairBits int = 2048
	DefaultListenPort          int = 8709
	DefaultPeerPort            int = 8708
)

var emptyKey crypto.PrivKey = &crypto.Ed25519PrivateKey{}

type Peer struct {
	uniqKey    crypto.PrivKey
	routedHost rhost.RoutedHost
	host       host.Host
}

func NewPeer() *Peer {
	randSeeder := rand.Reader
	key, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, DefaultGenerateKeyPairBits, randSeeder)
	if err != nil {
		logrus.Errorf("unable to GenerateKeyPair for new peer: %v", err)
		key = emptyKey
	}
	return &Peer{
		uniqKey: key,
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
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", routedHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	// addr := routedHost.Addrs()[0]
	addrs := routedHost.Addrs()
	log.Println("I can be reached at:")
	for _, addr := range addrs {
		log.Println(addr.Encapsulate(hostAddr))
	}

	// targetF = Pretty()
	//log.Printf("Now run \"aurae -d %s%s\" on a different terminal\n",  routedHost.ID().Pretty())
	log.Printf("ID: %s", routedHost.ID().Pretty())
	return nil
}

func RunClient(input string) {
	//golog.SetAllLoggers(golog.LevelInfo) // Change to INFO for extra info

	// Parse options from the command line
	//listenF := flag.Int("l", 0, "wait for incoming connections")
	//target := flag.String("d", "", "target peer to dial")
	//seed := flag.Int64("seed", 0, "set random seed for id generation")
	//global := flag.Bool("global", false, "use global ipfs peers for bootstrapping")
	//flag.Parse()
	//
	//if *listenF == 0 {
	//	log.Fatal("Please provide a port to bind on with -l")
	//}

	p := NewPeer()
	err := p.Establish(context.Background(), 0)
	if err != nil {
		log.Fatal(err)
	}

	// Make a host that listens on the given multiaddress
	//var bootstrapPeers []peer.AddrInfo
	//bootstrapPeers := IPFS_PEERS
	//ha, err := makeRoutedHost(1234, 0, bootstrapPeers, "global")
	//if err != nil {
	//	log.Fatal(err)
	//}

	// Set a stream handler on host A. /echo/1.0.0 is
	// a user-defined protocol name.
	ha := p.routedHost
	ha.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
		log.Println("Got a new stream!")
		if err := doEcho(s); err != nil {
			log.Println(err)
			s.Reset()
		} else {
			s.Close()
		}
	})

	peerid, err := peer.Decode(input)
	if err != nil {
		log.Fatalln(err)
	}

	// peerinfo := peer.AddrInfo{ID: peerid}
	log.Println("opening stream")
	// make a new stream from host B to host A
	// it should be handled on host A by the handler we set above because
	// we use the same /echo/1.0.0 protocol
	s, err := ha.NewStream(context.Background(), peerid, "/echo/1.0.0")

	if err != nil {
		log.Fatalln(err)
	}

	_, err = s.Write([]byte("Hello, world!\n"))
	if err != nil {
		log.Fatalln(err)
	}

	out, err := ioutil.ReadAll(s)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("read reply: %q\n", out)
}

func RunServer() {
	p := NewPeer()
	err := p.Establish(context.Background(), 1)
	if err != nil {
		log.Fatal(err)
	}
	ha := p.routedHost
	ha.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
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

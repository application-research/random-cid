package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/application-research/random-cid/ipfslite"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multicodec"
	"github.com/multiformats/go-multihash"
	"github.com/urfave/cli/v2"
)

type newCid struct {
	Reader  io.Reader
	Version int
}

var newFiles = make(chan newCid)
var cidout = make(chan string)
var peerID peer.ID

func newIpfsDaemon() *ipfslite.Peer {
	ctx := context.Background()

	ds := ipfslite.NewInMemoryDatastore()
	priv, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	listen, _ := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/4005")

	h, dht, err := ipfslite.SetupLibp2p(
		ctx,
		priv,
		nil,
		[]multiaddr.Multiaddr{listen},
		ds,
		ipfslite.Libp2pOptionsExtra...,
	)

	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	lite, err := ipfslite.New(ctx, ds, h, dht, nil)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	fmt.Println("IPFS peer address: ", h.ID())
	peerID = h.ID()
	return lite
}

func ipfsDaemon() {
	fmt.Println("Starting IPFS daemon")
	ctx := context.Background()

	lite := newIpfsDaemon()
	lite.Bootstrap(ipfslite.DefaultBootstrapPeers())
	fmt.Println("IPFS daemon started!")

	for {
		select {
		case newCid := <-newFiles:
			pref := cid.Prefix{
				Version:  uint64(newCid.Version),
				Codec:    uint64(multicodec.Raw),
				MhType:   multihash.SHA2_256,
				MhLength: -1, // default length
			}

			fmt.Println("Adding new cid to blockstore")
			node, err := lite.AddFile(ctx, newCid.Reader, &ipfslite.AddParams{Prefix: &pref}) // add file
			if err != nil {
				fmt.Println("Could not add node to blockstore: ", err)

			}
			fmt.Println("Added CID to blockstore: ", node.Cid())
			cidout <- node.Cid().String()
		}
	}
}

func NewCid(version int) (string, error) {
	// And then feed it some data
	data := []byte(gofakeit.HackerPhrase()) // random hacker phrase

	fmt.Println("Adding cid to local blockstore...")
	newFiles <- newCid{Reader: bytes.NewReader(data), Version: version}
	addedCid := <-cidout
	fmt.Println("Added cid: ", addedCid)

	return addedCid, nil
}

func getPeerID(ctx *gin.Context) {
	ctx.String(http.StatusOK, peerID.String()+"\n")
}

func getCidV1(ctx *gin.Context) {
	version := 1
	c, err := NewCid(version)
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}
	ctx.String(http.StatusOK, c+"\n")
}

func getCidV0(ctx *gin.Context) {
	version := 0
	c, err := NewCid(version)
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}
	ctx.String(http.StatusOK, c+"\n")
}

func main() {
	app := &cli.App{
		Name:  "random-cid",
		Usage: "generate random CIDs",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "cversion",
				Aliases:     []string{"c"},
				Value:       1,
				Usage:       "specify CID version to generate",
				DefaultText: "1",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "api",
				Usage: "start the random CID HTTP API",
				Action: func(c *cli.Context) error {
					router := gin.Default()
					router.GET("/", getCidV1)
					router.GET("/v1", getCidV1)
					router.GET("/v0", getCidV0)
					router.GET("/peer", getPeerID)

					port := os.Getenv("PORT")
					if port == "" {
						port = "8081"
					}
					if err := router.Run(":" + port); err != nil {
						fmt.Printf("error: %s\n", err)
					}
					return nil
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			// Create a cid manually by specifying the 'prefix' parameters
			cidVersion := ctx.Int("cversion")
			if cidVersion != 0 && cidVersion != 1 {
				return fmt.Errorf("Invalid CID version. Got %d, should be 1 or 0", cidVersion)
			}

			c, err := NewCid(cidVersion)
			if err != nil {
				return err
			}

			fmt.Println(c)
			return nil
		},
	}

	go ipfsDaemon()

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("error: ", err)
	}
}

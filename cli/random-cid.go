package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/brianvoe/gofakeit/v6"
	cid "github.com/ipfs/go-cid"
	mc "github.com/multiformats/go-multicodec"
	mh "github.com/multiformats/go-multihash"
	"github.com/urfave/cli/v2"
)

func NewCid(version int) (cid.Cid, error) {
	pref := cid.Prefix{
		Version:  uint64(version),
		Codec:    uint64(mc.Raw),
		MhType:   mh.SHA2_256,
		MhLength: -1, // default length
	}

	// And then feed it some data
	c, err := pref.Sum([]byte(gofakeit.HackerPhrase())) // random hacker phrase
	if err != nil {
		return cid.Cid{}, err
	}
	return c, nil
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

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

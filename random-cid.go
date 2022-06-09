package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
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

func getCidV1(ctx *gin.Context) {
	version := 1
	c, err := NewCid(version)
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}
	ctx.String(http.StatusOK, c.String()+"\n")
}

func getCidV0(ctx *gin.Context) {
	version := 0
	c, err := NewCid(version)
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}
	ctx.String(http.StatusOK, c.String()+"\n")
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

					port := os.Getenv("PORT")
					if port == "" {
						port = "8080"
					}
					if err := router.Run(":" + port); err != nil {
						log.Panicf("error: %s", err)
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

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"bufio"
	"flag"
	"github.com/nyudlts/go-aspace"
	"log"
	"os"
	"strings"
)

type AspaceDOID struct {
	RepoID int
	DOID   int
}

var (
	inFile      string
	config      string
	environment string
	uris        []AspaceDOID
	client      *aspace.ASClient
)

func init() {
	flag.StringVar(&inFile, "input-file", "", "the list of do uris to delete")
	flag.StringVar(&config, "config", "", "the location of a go-aspace config file")
	flag.StringVar(&inFile, "environment", "", "the environment to delete files from")
}

func main() {
	//parse flags
	flag.Parse()

	//parse uri list
	parselist()

	//create a client
	var err error
	client, err = aspace.NewClient(config, environment, 20)
	if err != nil {
		panic(err)
	}

	//delete aspace dos
	for _, uri := range uris {
		msg, err := client.DeleteDigitalObject(uri.RepoID, uri.DOID)
		if err != nil {
			log.Printf("[ERROR] %s", strings.ReplaceAll("\n", "", err.Error()))
		} else {
			log.Println("[INFO] %s", strings.ReplaceAll("\n", "", msg))
		}
	}
}

func parselist() {
	list, err := os.Open(inFile)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(list)
	for scanner.Scan() {
		repoId, doID, err := aspace.URISplit(scanner.Text())
		if err != nil {
			panic(err)
		}

		if repoId > 0 && doID > 0 {
			uris = append(uris, AspaceDOID{repoId, doID})
		} else {
			panic("INVALID URI: " + scanner.Text())
		}
	}
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nyudlts/go-aspace"
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
	test        bool
)

const placeholder = "https://hdl.handle.net/2333.1/material-request-placeholder"

func init() {
	flag.StringVar(&inFile, "input-file", "", "the list of do uris to delete")
	flag.StringVar(&config, "config", "", "the location of a go-aspace config file")
	flag.StringVar(&environment, "environment", "", "the environment to delete files from")
	flag.BoolVar(&test, "test", false, "")
}

func main() {
	//parse flags
	flag.Parse()

	//setup the log
	logFile, err := os.Create(fmt.Sprintf("aspace-do-delete-%s.log", environment))
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)

	//create a client
	fmt.Println(test)
	client, err = aspace.NewClient(config, environment, 20)
	if err != nil {
		panic(err)
	}

	//parse uri list
	parselist()

	//delete aspace dos
	for _, uri := range uris {

		//get the DO metadata
		domd, err := client.GetDigitalObject(uri.RepoID, uri.DOID)
		if err != nil {
			doErrMsg := fmt.Sprintf("[ERROR] %d %d %s", uri.RepoID, uri.DOID, strings.ReplaceAll(err.Error(), "\n", " "))
			fmt.Println(doErrMsg)
			log.Println(doErrMsg)
			continue
		}

		//get the uris from the file version
		fileversionUris := []string{}
		for _, fv := range domd.FileVersions {
			fileversionUris = append(fileversionUris, fv.FileURI)
		}

		infoMsg1 := fmt.Sprintf("[INFO] DO-URI: %s, TITLE: %s, FILE-URIS: %s", domd.URI, domd.Title, fileversionUris)
		fmt.Println(infoMsg1)
		log.Println(infoMsg1)

		if len(fileversionUris) == 1 && fileversionUris[0] == placeholder {
			//delete the do
			if !test {
				msg, err := client.DeleteDigitalObject(uri.RepoID, uri.DOID)
				if err != nil {
					errMsg := fmt.Sprintf("[ERROR] %s", strings.ReplaceAll(err.Error(), "\n", " "))
					fmt.Println(errMsg)
					log.Println(errMsg)
					continue
				} else {
					infoMsg2 := fmt.Sprintf("[INFO] DELETED %s %s", domd.URI, strings.ReplaceAll(msg, "\n", " "))
					fmt.Println(infoMsg2)
					log.Println(infoMsg2)
				}
			} else {
				msg := fmt.Sprintf("[INFO] skipping %s - in test mode", domd.URI)
				fmt.Println(msg)
				log.Println(msg)
			}
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

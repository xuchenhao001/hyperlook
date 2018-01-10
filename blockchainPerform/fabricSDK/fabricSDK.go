package fabricSDK

import (
	"github.com/hyperledger/fabric-sdk-go/def/fabapi"

	"log"
	"os"
)


// Initialize reads the configuration file and sets up the client, chain and event hub
func New() {

	log.Print(os.Getwd())

	setup := fabapi.Options{
		ConfigFile: "./blockchainPerform/fabricSDK/config.yaml",
	}

	_,err := fabapi.NewSDK(setup)
	if err != nil {
		log.Printf("Error initializing SDK: %s", err)
	}

	// Default channel client (uses organisation from client configuration)
	//_, err = sdk.NewChannelClient("mychannel", "admin")
	//if err != nil {
	//	log.Printf("Failed to create new channel client: %s", err)
	//}

	os.Exit(0)
}
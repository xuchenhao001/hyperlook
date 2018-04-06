package main

import (
	"flag"
	"fmt"
	"os"
)

var(
	// Definition of the Fabric SDK properties
	fabricSetup = FabricSetup{
		OrgAdmin:        "Admin",
		OrgName1:        "Org1",
		OrgName2:        "Org2",
		ConfigFile:      "config.yaml",

		UserName:        "User1",

		// Chaincode parameters
		ChainCodeID:     "examplecc",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "chaincode/",

		// Channel parameters
		ChannelID:       "mychannel",
		ChannelConfig:   "/var/fabric-net/channel-artifacts/channel.tx",
	}
)

// Create channel
func createChannel() {
	err := fabricSetup.CreateChannel()
	if err != nil {
		fmt.Printf("Unable to create channel: %v\n", err)
	}
}

func joinChannel() {
	// 2 organizations join to the channel
	err := fabricSetup.JoinChannel()
	if err != nil {
		fmt.Printf("Unable to join channel: %v\n", err)
	}
}

func installChaincode() {
	// 2 organizations install the chaincode
	err := fabricSetup.InstallCC()
	if err != nil {
		fmt.Printf("Unable to install chaincode: %v\n", err)
	}
}

func instantiateChaincode() {
	// 2 organizations instantiate the chaincode
	err := fabricSetup.InstantiateCC()
	if err != nil {
		fmt.Printf("Unable to instantiate chaincode: %v\n", err)
	}
}

func invokeChaincode() {
	_, err := fabricSetup.Invoke()
	if err != nil {
		fmt.Printf("Unable to invoke chaincode: %v\n", err)
	}
}

func queryChaincode() {
	_, err := fabricSetup.Query()
	if err != nil {
		fmt.Printf("Unable to query chaincode: %v\n", err)
	}
}

func main() {
	flag.Parse()

	switch flag.Arg(0) {
		case "createchannel":
			createChannel()
		case "joinchannel":
			joinChannel()
		case "installCC":
			installChaincode()
		case "instantiateCC":
			instantiateChaincode()
		case "invokeCC":
			invokeChaincode()
		case "queryCC":
			queryChaincode()
		default:
			fmt.Println("Please input an operation: " +
				"createchannel joinchannel installCC instantiateCC invokeCC queryCC")
	}
}

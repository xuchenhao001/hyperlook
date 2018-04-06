package main

import (
	"flag"
	"fmt"
)

var(
	// Definition of the Fabric SDK properties
	fabricSetupOrg1 = FabricSetup{
		OrgAdmin:        "Admin",
		OrgName:         "Org1",
		ConfigFile:      "org1-client-config.yaml",

		// Channel parameters
		ChannelID:       "mychannel",
		ChannelConfig:   "/var/fabric-net/channel-artifacts/channel.tx",
	}
	fabricSetupOrg2 = FabricSetup{
		OrgAdmin:        "Admin",
		OrgName:         "Org2",
		ConfigFile:      "org2-client-config.yaml",

		// Channel parameters
		ChannelID:       "mychannel",
		ChannelConfig:   "/var/fabric-net/channel-artifacts/channel.tx",
	}
)

// Create channel
func createChannel() {
	err := fabricSetupOrg1.CreateChannel()
	if err != nil {
		fmt.Printf("Unable to create channel: %v\n", err)
	}
}

func joinChannel() {
	// 2 organizations join to the channel
	err := fabricSetupOrg1.JoinChannel()
	if err != nil {
		fmt.Printf("Unable to join channel: %v\n", err)
	}
	err = fabricSetupOrg2.JoinChannel()
	if err != nil {
		fmt.Printf("Unable to join channel: %v\n", err)
	}
}

func installChaincode() {

}

func main() {
	flag.Parse()
	fmt.Println("Your operate: ", flag.Arg(0))

	switch flag.Arg(0) {
		case "createchannel":
			createChannel()
		case "joinchannel":
			joinChannel()
		case "instalcc":
			installChaincode()
		case "":
			fmt.Println("Please input an operation: " +
				"createchannel, joinchannel, instalcc, instancc, invokecc, or querycc")
	}
}
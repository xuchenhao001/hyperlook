package main

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
	chmgmt "github.com/hyperledger/fabric-sdk-go/api/apitxn/chmgmtclient"
	resmgmt "github.com/hyperledger/fabric-sdk-go/api/apitxn/resmgmtclient"
	"github.com/hyperledger/fabric-sdk-go/pkg/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
)

// FabricSetup implementation
type FabricSetup struct {
	ConfigFile      string
	OrgID           string
	ChannelID       string
	ChainCodeID     string
	ChannelConfig   string
	ChaincodeGoPath string
	ChaincodePath   string
	OrgAdmin        string
	OrgName1        string
	OrgName2        string
	UserName        string
	client          chclient.ChannelClient
	admin           resmgmt.ResourceMgmtClient
	sdk             *fabsdk.FabricSDK
}

func (setup *FabricSetup) localConfig() error {
	// Initialize the SDK with the configuration file
	sdk, err := fabsdk.New(config.FromFile(setup.ConfigFile))
	if err != nil {
		return fmt.Errorf("failed to create sdk: %v", err)
	}
	setup.sdk = sdk
	return nil
}

// CreateChannel reads the configuration file and sets up the client, chain and event hub
func (setup *FabricSetup) CreateChannel() error {
	// Initialize the SDK with the configuration file
	err := setup.localConfig()
	if err != nil {
		return fmt.Errorf("failed to initialize sdk: %v", err)
	}

	// Channel management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	chMgmtClient, err := setup.sdk.NewClient(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName1)).ChannelMgmt()
	if err != nil {
		return fmt.Errorf("failed to add Admin user to sdk: %v", err)
	}

	// Org admin user is signing user for creating channel.
	// The session method is the only way for now to get the user identity.
	orgAdminUser, err := setup.sdk.NewClient(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName1)).Session()
	if err != nil {
		return fmt.Errorf("failed to get session for %s, %s: %s", setup.OrgName1, setup.OrgAdmin, err)
	}

	// Creation of the channel. A channel can be understood as a private network inside the main network between two or more specific network Organizations
	// The channel is defined by its : Organizations, anchor peer (A peer node that all other peers can discover and communicate with. Every Organizations have one), the shared ledger, chaincode application(s) and the ordering service node(s)
	// Each transaction on the network is executed on a channel.
	req := chmgmt.SaveChannelRequest{ChannelID: setup.ChannelID, ChannelConfig: setup.ChannelConfig, SigningIdentity: orgAdminUser}
	if err = chMgmtClient.SaveChannel(req); err != nil {
		return fmt.Errorf("failed to create channel: %v", err)
	}

	fmt.Println("\n===== Create Channel Success =====")
	return nil
}

// Join channel
func (setup *FabricSetup) JoinChannel() error {
	// Initialize the SDK with the configuration file
	err := setup.localConfig()
	if err != nil {
		return fmt.Errorf("failed to initialize sdk: %v", err)
	}
	// Org1 admin
	setup.admin, err = setup.sdk.NewClient(fabsdk.WithUser(setup.OrgAdmin)).ResourceMgmt()
	if err != nil {
		return fmt.Errorf("failed to create new resource management client: %v", err)
	}
	// peers join channel
	if err = setup.admin.JoinChannel(setup.ChannelID); err != nil {
		return fmt.Errorf("org peers failed to join the channel: %v", err)
	}

	// Change to org2
	setup.admin, err = setup.sdk.NewClient(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName2)).ResourceMgmt()
	if err != nil {
		return fmt.Errorf("failed to create new resource management client: %v", err)
	}
	// peers join channel
	if err = setup.admin.JoinChannel(setup.ChannelID); err != nil {
		return fmt.Errorf("org peers failed to join the channel: %v", err)
	}

	fmt.Println("\n===== Join Channel Success =====")
	return nil
}

// Install chaincode, package chaincode to tarGZ first
func (setup *FabricSetup) InstallCC() error {
	// Initialize the SDK with the configuration file
	err := setup.localConfig()
	if err != nil {
		return fmt.Errorf("failed to initialize sdk: %v", err)
	}
	// Org1 admin
	setup.admin, err = setup.sdk.NewClient(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName1)).ResourceMgmt()
	if err != nil {
		return fmt.Errorf("failed to create new resource management client: %v", err)
	}

	// Create a new go lang chaincode package and initializing it with our chaincode
	ccPkg, err := packager.NewCCPackage(setup.ChaincodePath, setup.ChaincodeGoPath)
	if err != nil {
		return fmt.Errorf("failed to create chaincode package: %v", err)
	}

	// Install our chaincode on org peers
	// The resource management client send the chaincode to all peers in its channel in order for them to store it and interact with it later
	installCCReq := resmgmt.InstallCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodePath, Version: "1.0", Package: ccPkg}
	_, err = setup.admin.InstallCC(installCCReq)
	if err != nil {
		return fmt.Errorf("failed to install cc to org peers %v", err)
	}

	// Org2 admin
	setup.admin, err = setup.sdk.NewClient(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName2)).ResourceMgmt()
	if err != nil {
		return fmt.Errorf("failed to create new resource management client: %v", err)
	}
	_, err = setup.admin.InstallCC(installCCReq)
	if err != nil {
		return fmt.Errorf("failed to install cc to org peers %v", err)
	}

	fmt.Println("\n===== Chaincode Install Success ====")
	return nil
}

func (setup *FabricSetup) InstantiateCC() error {
	// Initialize the SDK with the configuration file
	err := setup.localConfig()
	if err != nil {
		return fmt.Errorf("failed to initialize sdk: %v", err)
	}
	// Org1 Admin
	setup.admin, err = setup.sdk.NewClient(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName1)).ResourceMgmt()
	if err != nil {
		return fmt.Errorf("failed to create new resource management client: %v", err)
	}

	// Set up chaincode policy
	// The chaincode policy is required if your transactions must follow some specific rules
	// If you don't provide any policy every transaction will be endorsed, and it's probably not what you want
	// In this case, we set the rule to : Endorse the transaction if the transaction have been signed by a member from any org
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP","Org2MSP"})

	// Instantiate our chaincode on org peers
	// The resource management client tells to all peers in its channel to instantiate the chaincode previously installed
	initArgs := [][]byte{[]byte("init"), []byte("a"), []byte("1000"), []byte("b"), []byte("2000")}
	initRequest := resmgmt.InstantiateCCRequest{
		Name:    setup.ChainCodeID,
		Path:    setup.ChaincodePath,
		Version: "1.0",
		Args:    initArgs,
		Policy:  ccPolicy,
	}
	err = setup.admin.InstantiateCC(setup.ChannelID, initRequest)
	if err != nil {
		return fmt.Errorf("failed to instantiate the chaincode: %v", err)
	}

	// Change to Org2 Admin
	//setup.admin, err = setup.sdk.NewClient(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName2)).ResourceMgmt()
	//if err != nil {
	//	return fmt.Errorf("failed to create new resource management client: %v", err)
	//}
	//err = setup.admin.InstantiateCC(setup.ChannelID, initRequest)
	//if err != nil {
	//	return fmt.Errorf("failed to instantiate the chaincode: %v", err)
	//}

	fmt.Println("\n===== Chaincode Instantiate Success ====")
	return nil
}

func (setup *FabricSetup) Invoke() (string, error) {
	// Initialize the SDK with the configuration file
	err := setup.localConfig()
	if err != nil {
		return "", fmt.Errorf("failed to initialize sdk: %v", err)
	}
	// Channel client is used to query and execute transactions
	setup.client, err = setup.sdk.NewClient(fabsdk.WithUser(setup.UserName)).Channel(setup.ChannelID)
	if err != nil {
		return "", fmt.Errorf("failed to create new channel client: %v", err)
	}

	eventID := "eventInvoke"

	// Register a notification handler on the client
	notifier := make(chan *chclient.CCEvent)
	rce, err := setup.client.RegisterChaincodeEvent(notifier, setup.ChainCodeID, eventID)
	if err != nil {
		return "", fmt.Errorf("failed to register chaincode evet: %v", err)
	}

	// Create a request (proposal) and send it
	invokeArgs := [][]byte{[]byte("a"), []byte("b"), []byte("10")}
	invokeRequest := chclient.Request{
		ChaincodeID: setup.ChainCodeID,
		Fcn: "invoke",
		Args: invokeArgs,
	}
	response, err := setup.client.Execute(invokeRequest)
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}

	// Wait for the result of the submission
	select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %s\n", ccEvent)
	case <-time.After(time.Second * 120):
		return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	}

	// Unregister the notification handler previously created on the client
	err = setup.client.UnregisterChaincodeEvent(rce)

	fmt.Println("\n===== Chaincode Invoke Success ====")
	return response.TransactionID.ID, nil
}

func (setup *FabricSetup) Query() (string, error) {
	// Initialize the SDK with the configuration file
	err := setup.localConfig()
	if err != nil {
		return "", fmt.Errorf("failed to initialize sdk: %v", err)
	}
	// Channel client is used to query and execute transactions
	setup.client, err = setup.sdk.NewClient(fabsdk.WithUser(setup.UserName)).Channel(setup.ChannelID)
	if err != nil {
		return "", fmt.Errorf("failed to create new channel client: %v", err)
	}

	// Prepare arguments
	queryArgs := [][]byte{[]byte("invoke"), []byte("query"), []byte("a")}
	queryRequest := chclient.Request{
		ChaincodeID: setup.ChainCodeID,
		Fcn: "invoke",
		Args: queryArgs,
	}
	response, err := setup.client.Query(queryRequest)
	if err != nil {
		return "", fmt.Errorf("failed to query: %v", err)
	}

	return string(response.Payload), nil
}

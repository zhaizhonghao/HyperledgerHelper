package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/zhaizhonghao/configtxTool/services/configtx"
	"github.com/zhaizhonghao/configtxTool/services/crypto"
)

type Success struct {
	Payload string `json:"Payload"`
	Message string `json:"Message"`
}

type ChannelID struct {
	Name string `json:"Name"`
}

type PeerInfo struct {
	Org     string `json:"Org"`
	Port    string `json:"Port"`
	Channel string `json:"Channel"`
}

var tpl *template.Template

//Pipeline
func GetIdOfCouchDB(a int, b int) int {
	return 2*a + b - 2
}

func GetPortOfCouchDB(a int, b int) int {
	return ((2*a+b-2)+5)*1000 + 984
}

func GetGeneralPortOfPeer(a int, b int) int {
	return ((2*a+b-2)+7)*1000 + 51
}

func GetChaincodePortOfPeer(a int, b int) int {
	return ((2*a+b-2)+7)*1000 + 52
}

func GetIdOfBootstrapNode(a int) int {
	if a == 0 {
		return 1
	} else {
		return 0
	}
}

func GetPortOfBootstrapNode(a int) int {
	if a == 0 {
		return 8051
	} else {
		return 7051
	}
}

var fm = template.FuncMap{
	"GetIdOfCouchDB":         GetIdOfCouchDB,
	"GetPortOfCouchDB":       GetPortOfCouchDB,
	"GetGeneralPortOfPeer":   GetGeneralPortOfPeer,
	"GetChaincodePortOfPeer": GetChaincodePortOfPeer,
	"GetIdOfBootstrapNode":   GetIdOfBootstrapNode,
	"GetPortOfBootstrapNode": GetPortOfBootstrapNode,
}

func main() {
	router := mux.NewRouter()
	//Section configtx
	router.HandleFunc("/configtx", requestConfigtx).Methods("POST", http.MethodOptions)

	router.HandleFunc("/configtx", revokeConfigtx).Methods("DELETE")

	router.HandleFunc("/configtx", updateConfigtx).Methods("PUT")

	router.HandleFunc("/configtx", patchConfigtx).Methods("PATCH")

	//Section Crypto
	router.HandleFunc("/crypto", requestCrypto).Methods("POST", http.MethodOptions)

	//Section Node Deployment
	router.HandleFunc("/node", nodeDeploy).Methods("GET", http.MethodOptions)

	//Section Channel Management
	router.HandleFunc("/channel", createChannel).Methods("POST", http.MethodOptions)

	router.HandleFunc("/channel/join", joinChannel).Methods("POST", http.MethodOptions)

	//Section Smart Contract
	router.HandleFunc("/contract/package", packageChaincode).Methods("POST", http.MethodOptions)

	router.HandleFunc("/contract/install", installChaincode).Methods("POST", http.MethodOptions)

	router.HandleFunc("/contract/approve", approveChaincode).Methods("POST", http.MethodOptions)

	router.HandleFunc("/contract/checkApprove", checkApprove).Methods("POST", http.MethodOptions)

	router.HandleFunc("/contract/instantiate", instantiateChaincode).Methods("POST", http.MethodOptions)

	router.HandleFunc("/contract/initialize", initializeChaincode).Methods("POST", http.MethodOptions)

	router.Use(mux.CORSMethodMiddleware(router))

	fmt.Println("Server is listenning on localhost:8080")

	http.ListenAndServe(":8080", router)
}

func requestConfigtx(w http.ResponseWriter, r *http.Request) {
	//set the header
	setHeader(w)

	var cfgtx = configtx.Configtx{}
	err := json.NewDecoder(r.Body).Decode(&cfgtx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Get configtx", cfgtx)
	//Generate configtx.yaml
	tpl = template.Must(template.ParseGlob("templates/configtx/*.yaml"))
	file, err := os.Create("channel/configtx.yaml")
	if err != nil {
		fmt.Println("Fail to create file!")
	}
	err = configtx.GenerateConfigTxTemplate(cfgtx, tpl, file)
	if err != nil {
		fmt.Println("Fail to generate configtx.yaml", err)
	}
	//Generate the genesis.block for the sys channel and the [channelName].tx for specific channel
	cmd := exec.Command("configtxgen", "-profile", "OrdererGenesis", "-configPath", "channel", "-channelID", "sys-channel", "-outputBlock", "channel/genesis.block")

	err = cmd.Run()
	if err != nil {
		fmt.Println("Execute Command for generating genesis.block failed:" + err.Error())
		return
	}
	//Generate the [channelName].tx for specific channel
	fmt.Println("channel name", cfgtx.Channel.Name)
	cmd = exec.Command("configtxgen", "-profile", cfgtx.Channel.Name, "-configPath", "channel", "-outputCreateChannelTx", "channel/"+strings.ToLower(cfgtx.Channel.Name)+".tx", "-channelID", strings.ToLower(cfgtx.Channel.Name))
	err = cmd.Run()
	if err != nil {
		fmt.Println("Execute Command for generating channel.tx failed:" + err.Error())
		return
	}
	json.NewEncoder(w).Encode(cfgtx)
}

func revokeConfigtx(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func updateConfigtx(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func patchConfigtx(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func requestCrypto(w http.ResponseWriter, r *http.Request) {
	setHeader(w)

	var configCp = crypto.ConfigCp{}
	err := json.NewDecoder(r.Body).Decode(&configCp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Get configCp", configCp)
	//Generate crypto-config.yaml
	tpl = template.Must(template.ParseGlob("templates/crypto/*.yaml"))
	file, err := os.Create("channel/crypto-config.yaml")
	if err != nil {
		fmt.Println("Fail to create file!")
	}
	err = crypto.GenerateCryptoTemplate(configCp, tpl, file)
	if err != nil {
		fmt.Println("Fail to generate crypto-config", err)
	}
	//cryptogen.GenerateCryptoConfig()
	//Generate docker-compose.yaml
	file1, err1 := os.Create("docker-compose.yaml")
	if err1 != nil {
		fmt.Println("Fail to create file!")
	}
	defer file1.Close()
	tpl = template.Must(template.New("").Funcs(fm).ParseFiles("templates/docker/dockerComposeTemplate.yaml"))
	err1 = tpl.ExecuteTemplate(file1, "dockerComposeTemplate.yaml", configCp)
	if err1 != nil {
		fmt.Println(err)
	}
	//Generate the crypto-config file
	cmd := exec.Command("cryptogen", "generate", "--config=channel/crypto-config.yaml", "--output=channel/crypto-config")

	err = cmd.Run()
	if err != nil {
		fmt.Println("Execute Command failed:" + err.Error())
		return
	}
	json.NewEncoder(w).Encode(configCp)
}

func nodeDeploy(w http.ResponseWriter, r *http.Request) {
	setHeader(w)

	fmt.Println("Deploying node")

	//Creating or starting the docker containers
	cmd := exec.Command("docker-compose", "up", "-d")

	err := cmd.Run()
	if err != nil {
		fmt.Println("Execute docker-compose Command failed:" + err.Error())
		return
	}

	success := Success{
		Payload: "nodes are deployed successfully",
		Message: "200 OK",
	}
	json.NewEncoder(w).Encode(success)
}

func createChannel(w http.ResponseWriter, r *http.Request) {
	setHeader(w)

	var channelID = ChannelID{}
	err := json.NewDecoder(r.Body).Decode(&channelID)
	if err != nil {
		fmt.Println(err)
	}
	if channelID.Name == "" {
		return
	}
	channelName := strings.ToLower(channelID.Name)
	fmt.Println("Creating channel", channelName)

	setEnvironmentForPeer("org1", "7051")

	//create the channel
	out, err1 := exec.Command(
		"peer", "channel", "create", "-o", "localhost:7050", "-c", channelName,
		"--ordererTLSHostnameOverride", "orderer1.example.com",
		"-f", "./channel/"+channelName+".tx",
		"--outputBlock", "./channel-artifacts/"+channelName+".block",
		"--tls", os.Getenv("CORE_PEER_TLS_ENABLED"),
		"--cafile", os.Getenv("ORDERER_CA")).Output()
	if err1 != nil {
		fmt.Println("Create channel " + channelName + " failed:" + err1.Error())
		fmt.Println(string(out))
		return
	}

	fmt.Println(string(out))

	success := Success{
		Payload: "Channel " + channelName + " is created successfully",
		Message: "200 OK",
	}
	json.NewEncoder(w).Encode(success)
}

func joinChannel(w http.ResponseWriter, r *http.Request) {
	setHeader(w)

	var peerInfo = PeerInfo{}
	err := json.NewDecoder(r.Body).Decode(&peerInfo)
	if err != nil {
		fmt.Println(err)
	}
	if peerInfo.Org == "" {
		return
	}

	var channel = strings.ToLower(peerInfo.Channel)

	setEnvironmentForPeer(peerInfo.Org, peerInfo.Port)
	fmt.Println(peerInfo.Org, peerInfo.Port, "is joining the channel...")
	//join the channel
	out, err1 := exec.Command("peer", "channel", "join", "-b", "./channel-artifacts/"+channel+".block").Output()
	if err1 != nil {
		fmt.Println("Join the channel " + channel + " failed:" + err1.Error())
		fmt.Println(string(out))
		return
	}

	fmt.Println(string(out))

	success := Success{
		Payload: "Join Channel " + channel + " successfully",
		Message: "200 OK",
	}
	json.NewEncoder(w).Encode(success)
}

type ContractInfo struct {
	PeerInfo     PeerInfo `json:"PeerInfo"`
	Language     string   `json:"Language"`
	Version      string   `json:"Version"`
	ContractName string   `json:"ContractName"`
}

func packageChaincode(w http.ResponseWriter, r *http.Request) {
	setHeader(w)

	var contractInfo = ContractInfo{}
	err := json.NewDecoder(r.Body).Decode(&contractInfo)
	if err != nil {
		fmt.Println(err)
	}
	if contractInfo.ContractName == "" {
		return
	}

	setEnvironmentForPeer(contractInfo.PeerInfo.Org, contractInfo.PeerInfo.Port)
	fmt.Println(contractInfo.PeerInfo.Org, contractInfo.PeerInfo.Port, "packaging...")
	//package the contract
	fmt.Println("Packing the contract", contractInfo.ContractName)
	fmt.Println("contract label:" + contractInfo.ContractName + "_" + contractInfo.Version)
	out, err1 := exec.Command(
		"peer", "lifecycle", "chaincode", "package", contractInfo.ContractName+".tar.gz",
		"--path", "./src/github.com/"+contractInfo.ContractName+"/go",
		"--lang", contractInfo.Language,
		"--label", contractInfo.ContractName+"_"+contractInfo.Version).Output()
	if err1 != nil {
		fmt.Println("package the contract " + contractInfo.ContractName + " failed:" + err1.Error())
		fmt.Println(string(out))
		return
	}

	fmt.Println(string(out))

	success := Success{
		Payload: "Package the contract " + contractInfo.ContractName + " successfully",
		Message: "200 OK",
	}
	json.NewEncoder(w).Encode(success)

}

func installChaincode(w http.ResponseWriter, r *http.Request) {
	setHeader(w)

	var contractInfo = ContractInfo{}
	err := json.NewDecoder(r.Body).Decode(&contractInfo)
	if err != nil {
		fmt.Println(err)
	}
	if contractInfo.ContractName == "" {
		return
	}
	setEnvironmentForPeer(contractInfo.PeerInfo.Org, contractInfo.PeerInfo.Port)
	fmt.Println(contractInfo.PeerInfo.Org, contractInfo.PeerInfo.Port, "installing...")
	//Install the contract
	fmt.Println("Install the contract", contractInfo.ContractName)

	out, err1 := exec.Command(
		"peer", "lifecycle", "chaincode", "install", contractInfo.ContractName+".tar.gz").Output()
	if err1 != nil {
		fmt.Println("install the contract " + contractInfo.ContractName + " failed:" + err1.Error())
		fmt.Println(string(out))
		return
	}

	fmt.Println(string(out))

	//query the contract
	fmt.Println("query the contract")

	out, err1 = exec.Command(
		"peer", "lifecycle", "chaincode", "queryinstalled").Output()
	if err1 != nil {
		fmt.Println("query the contract failed:" + err1.Error())
		fmt.Println(string(out))
		return
	}

	fmt.Println(string(out))
	parts := strings.Split(string(out), "Package ID: ")
	parts = strings.Split(parts[1], ",")
	packageID := parts[0]

	success := Success{
		Payload: packageID,
		Message: "200 OK",
	}
	json.NewEncoder(w).Encode(success)
}

type ApproveInfo struct {
	ContractInfo ContractInfo `json:"ContractInfo"`
	PackageID    string       `json:"PackageID"`
}

func approveChaincode(w http.ResponseWriter, r *http.Request) {
	setHeader(w)
	//query the package ID of installed contract
	var approveInfo = ApproveInfo{}
	err := json.NewDecoder(r.Body).Decode(&approveInfo)
	if err != nil {
		fmt.Println(err)
	}
	if approveInfo.ContractInfo.ContractName == "" {
		return
	}
	setEnvironmentForPeer(approveInfo.ContractInfo.PeerInfo.Org, approveInfo.ContractInfo.PeerInfo.Port)
	fmt.Println(approveInfo.ContractInfo.PeerInfo.Org, approveInfo.ContractInfo.PeerInfo.Port, "aproving...")
	//export PRIVATE_DATA_CONFIG=${PWD}/artifacts/private-data/collections_config.json
	//approve the contract
	var channel = strings.ToLower(approveInfo.ContractInfo.PeerInfo.Channel)
	out, err1 := exec.Command(
		"peer", "lifecycle", "chaincode", "approveformyorg", "-o", "localhost:7050",
		"--ordererTLSHostnameOverride", "orderer1.example.com", "--tls",
		"--cafile", os.Getenv("ORDERER_CA"),
		"--channelID", channel,
		"--name", approveInfo.ContractInfo.ContractName,
		"--version", approveInfo.ContractInfo.Version,
		"--init-required",
		"--package-id", approveInfo.PackageID,
		"--sequence", approveInfo.ContractInfo.Version).Output()
	if err1 != nil {
		fmt.Println("approve the contract failed:" + err1.Error())
		fmt.Println(string(out))
		return
	}

	success := Success{
		Payload: "approve the contract " + approveInfo.PackageID + " successfully!",
		Message: "200 OK",
	}
	json.NewEncoder(w).Encode(success)
}

func checkApprove(w http.ResponseWriter, r *http.Request) {
	setHeader(w)

	var approveInfo = ApproveInfo{}
	err := json.NewDecoder(r.Body).Decode(&approveInfo)
	if err != nil {
		fmt.Println(err)
	}
	if approveInfo.ContractInfo.ContractName == "" {
		return
	}

	setEnvironmentForPeer(approveInfo.ContractInfo.PeerInfo.Org, approveInfo.ContractInfo.PeerInfo.Port)

	//check the approval
	var channel = strings.ToLower(approveInfo.ContractInfo.PeerInfo.Channel)
	out, err1 := exec.Command(
		"peer", "lifecycle", "chaincode", "checkcommitreadiness",
		"--channelID", channel,
		"--name", approveInfo.ContractInfo.ContractName,
		"--version", approveInfo.ContractInfo.Version,
		"--sequence", approveInfo.ContractInfo.Version,
		"--output", "json",
		"--init-required").Output()
	if err1 != nil {
		fmt.Println("check the approvals failed:" + err1.Error())
		fmt.Println(string(out))
		return
	}

	fmt.Println(string(out))

	success := Success{
		Payload: string(out),
		Message: "200 OK",
	}
	json.NewEncoder(w).Encode(success)

}

type InstantiateInfo struct {
	Language     string     `json:"Language"`
	Version      string     `json:"Version"`
	ContractName string     `json:"ContractName"`
	Approvers    []PeerInfo `json:"Approvers"`
}

func instantiateChaincode(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start instantiating the contract")
	setHeader(w)

	var instantiateInfo = InstantiateInfo{}
	err := json.NewDecoder(r.Body).Decode(&instantiateInfo)
	if err != nil {

		fmt.Println("parse the instantiateInfo failed", err)
	}
	if instantiateInfo.ContractName == "" {
		return
	}

	setEnvironmentForPeer(instantiateInfo.Approvers[0].Org, instantiateInfo.Approvers[0].Port)
	fmt.Println(instantiateInfo.Approvers[0].Org, instantiateInfo.Approvers[0].Port, "instantiating")
	//instantiat the contract
	var channel = strings.ToLower(instantiateInfo.Approvers[0].Channel)
	args := []string{
		"lifecycle", "chaincode", "commit",
		"-o", "localhost:7050",
		"--ordererTLSHostnameOverride", "orderer1.example.com",
		"--tls", os.Getenv("CORE_PEER_TLS_ENABLED"),
		"--cafile", os.Getenv("ORDERER_CA"),
		"--channelID", channel,
		"--name", instantiateInfo.ContractName,
		"--version", instantiateInfo.Version,
		"--sequence", instantiateInfo.Version,
		"--init-required",
	}
	for i := 0; i < len(instantiateInfo.Approvers); i++ {
		arg1 := "--peerAddresses"
		arg2 := "localhost:" + instantiateInfo.Approvers[i].Port
		arg3 := "--tlsRootCertFiles"
		arg4 := os.Getenv("PWD") + "/channel/crypto-config/peerOrganizations/" + instantiateInfo.Approvers[i].Org + ".example.com/peers/peer0." + instantiateInfo.Approvers[i].Org + ".example.com/tls/ca.crt"
		args = append(args, arg1, arg2, arg3, arg4)
	}
	fmt.Println("instantiating the contract:", args)
	fmt.Println("args:", args)
	out, err1 := exec.Command("peer", args...).Output()
	if err1 != nil {
		fmt.Println("instantiate the contract failed:" + err1.Error())
		fmt.Println(string(out))
		return
	}

	fmt.Println(string(out))

	success := Success{
		Payload: "instantiate the contract " + instantiateInfo.ContractName + " successfully!",
		Message: "200 OK",
	}
	json.NewEncoder(w).Encode(success)
}

type InitializeInfo struct {
	InstantiateInfo InstantiateInfo `json:"InstantiateInfo"`
	ArgsJSONString  string          `json:"ArgsJSONString"`
}

func initializeChaincode(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Start initializing the contract")
	setHeader(w)

	var initializeInfo = InitializeInfo{}
	err := json.NewDecoder(r.Body).Decode(&initializeInfo)
	if err != nil {

		fmt.Println("parse the initializeInfo failed", err)
	}
	if initializeInfo.InstantiateInfo.ContractName == "" {
		return
	}

	setEnvironmentForPeer(initializeInfo.InstantiateInfo.Approvers[0].Org, initializeInfo.InstantiateInfo.Approvers[0].Port)
	//instantiat the contract
	var channel = strings.ToLower(initializeInfo.InstantiateInfo.Approvers[0].Channel)
	fmt.Println("argsJSON", initializeInfo.ArgsJSONString)
	args := []string{
		"chaincode", "invoke",
		"-o", "localhost:7050",
		"--ordererTLSHostnameOverride", "orderer1.example.com",
		"--tls", os.Getenv("CORE_PEER_TLS_ENABLED"),
		"--cafile", os.Getenv("ORDERER_CA"),
		"-C", channel,
		"-n", initializeInfo.InstantiateInfo.ContractName,
		"--isInit",
		"-c", "'" + initializeInfo.ArgsJSONString + "'",
	}
	for i := 0; i < len(initializeInfo.InstantiateInfo.Approvers); i++ {
		arg1 := "--peerAddresses"
		arg2 := "localhost:" + initializeInfo.InstantiateInfo.Approvers[i].Port
		arg3 := "--tlsRootCertFiles"
		arg4 := os.Getenv("PWD") + "/channel/crypto-config/peerOrganizations/" + initializeInfo.InstantiateInfo.Approvers[i].Org + ".example.com/peers/peer0." + initializeInfo.InstantiateInfo.Approvers[i].Org + ".example.com/tls/ca.crt"
		args = append(args, arg1, arg2, arg3, arg4)
	}
	fmt.Println("iniliazing the contract:", args)
	fmt.Println("args:", args)
	out, err1 := exec.Command("peer", args...).Output()
	if err1 != nil {
		fmt.Println("iniliazing the contract failed:" + err1.Error())
		fmt.Println(string(out))
		return
	}

	fmt.Println(string(out))

	success := Success{
		Payload: "iniliazing the contract " + initializeInfo.InstantiateInfo.ContractName + " successfully!",
		Message: "200 OK",
	}
	json.NewEncoder(w).Encode(success)
}

func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

func setHeader(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")
	w.Header().Set("X-Powered-By", "3.2.1")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
}

func setEnvironmentForPeer(org string, port string) {
	os.Setenv("CORE_PEER_TLS_ENABLED", "true")
	os.Setenv("ORDERER_CA", os.Getenv("PWD")+"/channel/crypto-config/ordererOrganizations/example.com/orderers/orderer1.example.com/msp/tlscacerts/tlsca.example.com-cert.pem")
	os.Setenv("FABRIC_CFG_PATH", os.Getenv("PWD")+"/channel/config/")
	//set global variables for peer
	os.Setenv("CORE_PEER_LOCALMSPID", Capitalize(org)+"MSP")
	os.Setenv("CORE_PEER_TLS_ROOTCERT_FILE", os.Getenv("PWD")+"/channel/crypto-config/peerOrganizations/"+org+".example.com/peers/peer0."+org+".example.com/tls/ca.crt")
	os.Setenv("CORE_PEER_MSPCONFIGPATH", os.Getenv("PWD")+"/channel/crypto-config/peerOrganizations/"+org+".example.com/users/Admin@"+org+".example.com/msp")
	os.Setenv("CORE_PEER_ADDRESS", "localhost:"+port)
}

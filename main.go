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

	router.Use(mux.CORSMethodMiddleware(router))

	fmt.Println("Server is listenning on localhost:8080")

	http.ListenAndServe(":8080", router)
}

func requestConfigtx(w http.ResponseWriter, r *http.Request) {
	//set the header
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")
	w.Header().Set("X-Powered-By", "3.2.1")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")

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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")
	w.Header().Set("X-Powered-By", "3.2.1")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")

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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")
	w.Header().Set("X-Powered-By", "3.2.1")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")

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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")
	w.Header().Set("X-Powered-By", "3.2.1")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")

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

	os.Setenv("CORE_PEER_TLS_ENABLED", "true")
	os.Setenv("ORDERER_CA", os.Getenv("PWD")+"/channel/crypto-config/ordererOrganizations/example.com/orderers/orderer1.example.com/msp/tlscacerts/tlsca.example.com-cert.pem")
	os.Setenv("FABRIC_CFG_PATH", os.Getenv("PWD")+"/channel/config/")
	//set global variables for peer0 of org1
	os.Setenv("CORE_PEER_LOCALMSPID", "Org1MSP")
	os.Setenv("CORE_PEER_TLS_ROOTCERT_FILE", os.Getenv("PWD")+"/channel/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt")
	os.Setenv("CORE_PEER_MSPCONFIGPATH", os.Getenv("PWD")+"/channel/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp")
	os.Setenv("CORE_PEER_ADDRESS", "localhost:7051")

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
		fmt.Println(out)
		return
	}

	fmt.Println(out)

	success := Success{
		Payload: "Channel " + channelName + " is created successfully",
		Message: "200 OK",
	}
	json.NewEncoder(w).Encode(success)
}

func joinChannel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")
	w.Header().Set("X-Powered-By", "3.2.1")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	var peerInfo = PeerInfo{}
	err := json.NewDecoder(r.Body).Decode(&peerInfo)
	if err != nil {
		fmt.Println(err)
	}
	if peerInfo.Org == "" {
		return
	}
	fmt.Println(Capitalize(peerInfo.Org) + "MSP")
	fmt.Println(peerInfo.Port)
	fmt.Println(peerInfo.Channel)
	var channel = strings.ToLower(peerInfo.Channel)
	os.Setenv("CORE_PEER_TLS_ENABLED", "true")
	os.Setenv("ORDERER_CA", os.Getenv("PWD")+"/channel/crypto-config/ordererOrganizations/example.com/orderers/orderer1.example.com/msp/tlscacerts/tlsca.example.com-cert.pem")
	os.Setenv("FABRIC_CFG_PATH", os.Getenv("PWD")+"/channel/config/")
	//set global variables for peer
	os.Setenv("CORE_PEER_LOCALMSPID", Capitalize(peerInfo.Org)+"MSP")
	os.Setenv("CORE_PEER_TLS_ROOTCERT_FILE", os.Getenv("PWD")+"/channel/crypto-config/peerOrganizations/"+peerInfo.Org+".example.com/peers/peer0."+peerInfo.Org+".example.com/tls/ca.crt")
	os.Setenv("CORE_PEER_MSPCONFIGPATH", os.Getenv("PWD")+"/channel/crypto-config/peerOrganizations/"+peerInfo.Org+".example.com/users/Admin@"+peerInfo.Org+".example.com/msp")
	os.Setenv("CORE_PEER_ADDRESS", "localhost:"+peerInfo.Port)

	//join the channel
	out, err1 := exec.Command("peer", "channel", "join", "-b", "./channel-artifacts/"+channel+".block").Output()
	if err1 != nil {
		fmt.Println("Join the channel " + channel + " failed:" + err1.Error())
		fmt.Println(out)
		return
	}

	fmt.Println(out)

	success := Success{
		Payload: "Join Channel " + channel + " successfully",
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

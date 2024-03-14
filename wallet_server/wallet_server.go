package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/erdincmutlu/goblockchain/block"
	"github.com/erdincmutlu/goblockchain/utils"
	"github.com/erdincmutlu/goblockchain/wallet"
)

const tempDir = "templates"

type WalletServer struct {
	port    uint16
	gateway string
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port: port, gateway: gateway}
}

func (ws *WalletServer) Port() uint16 {
	return ws.port
}

func (ws *WalletServer) Gateway() string {
	return ws.gateway
}

func (ws *WalletServer) Index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t, _ := template.ParseFiles(path.Join(tempDir, "index.html"))
		t.Execute(w, "")
	default:
		log.Printf("ERROR: Invalid http method")
	}
}

func (ws *WalletServer) Wallet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")
		myWallet := wallet.NewWallet()
		m, _ := myWallet.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid http method")
	}
}

func (ws *WalletServer) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var t wallet.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("ERROR: %s\n", err.Error())
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		if !t.Validate() {
			log.Printf("ERROR: missing fields\n")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		privateKey := utils.PrivateKeyFromString(*t.SenderPrivateKey, publicKey)
		value, err := strconv.ParseFloat(*t.Value, 32)
		if err != nil {
			log.Printf("ERROR: parse error\n")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		value32 := float32(value)
		w.Header().Add("Content-Type", "application/json")

		transaction := wallet.NewTransaction(privateKey, publicKey,
			*t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, value32)
		signature := transaction.GenerateSignature()
		signatureStr := signature.String()

		bt := &block.TransactionRequest{
			SenderBlockchainAddress:    t.SenderBlockchainAddress,
			RecipientBlockchainAddress: t.RecipientBlockchainAddress,
			SenderPublicKey:            t.SenderPublicKey,
			Value:                      &value32,
			Signature:                  &signatureStr,
		}
		m, _ := json.Marshal(bt)
		buf := bytes.NewBuffer(m)
		resp, _ := http.Post(ws.Gateway()+"/transactions", "application/json", buf)
		if resp.StatusCode == http.StatusCreated {
			io.WriteString(w, string(utils.JsonStatus("success")))
			return
		}
		io.WriteString(w, string(utils.JsonStatus("fail")))

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid http method")
	}
}

func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/transaction", ws.CreateTransaction)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.Port())), nil))
}

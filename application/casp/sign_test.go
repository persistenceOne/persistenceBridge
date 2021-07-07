package casp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

const API_TOKEN = "Bearer cHVuZWV0TmV3QXBpa2V5MTI6OWM1NDBhMzAtNTQ5NC00ZDdhLTljODktODA3MDZiNWNhYzQ1"

type request struct {
	DataToSign   []string `json:"dataToSign"`
	Description  string   `json:"description"`
	ProviderData string   `json:"providerData"`
	Details      string   `json:"details"`
	PublicKeys   []string `json:"publicKeys"`
}

func TestSignTx(t *testing.T) {
	//Encode the data
	postBody, _ := json.Marshal(request{
		DataToSign:   []string{"55C53F5D490297900CEFA825D0C8E8E9532EE8A118ABE7D8570762CD38BE9818"},
		Description:  "Test",
		ProviderData: "",
		Details:      "",
		PublicKeys:   []string{"3056301006072A8648CE3D020106052B8104000A03420004A910F239D0799A0DB8B67B16E0CD9902299E6AA9D941A27209F217895367A7D4715840C016750A5D395E1886DBA4D0ED7312E284D74A30A0A5DC6CA80F44425B"},
	})
	responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}}
	request, err := http.NewRequest("POST", "https://65.2.149.241:443/casp/api/v1.0/mng/vaults/509fd89a-762a-40ec-bd4b-0745b06e2d3d/sign", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	request.Header.Set("authorization", API_TOKEN)
	request.Header.Set("Content-Type", "application/json")
	fmt.Println(request.Header)
	//Read the response body
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)
}

func TestGet(t *testing.T) {

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}}
	//request, err := http.NewRequest("GET", "https://65.2.149.241:443/casp/api/v1.0/mng/auth/users", nil)
	//request, err := http.NewRequest("GET", "https://65.2.149.241:443/casp/api/v1.0/mng/vaults/509fd89a-762a-40ec-bd4b-0745b06e2d3d/coins/118/accounts/0/chains/all/addresses?encoding=uncompressed", nil)
	request, err := http.NewRequest("GET", "https://65.2.149.241:443/casp/api/v1.0/mng/operations/sign/656dbc73-7c74-4006-a9d5-9575594c22d9", nil)
	//request, err := http.NewRequest("GET", "https://65.2.149.241:443/casp/api/v1.0/mng/accounts/bd4c618e-8046-4fef-bdaa-9716ade77553/participants", nil)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	request.Header.Set("authorization", API_TOKEN)
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)
}

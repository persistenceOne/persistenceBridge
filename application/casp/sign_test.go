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
		PublicKeys:   []string{"3056301006072A8648CE3D020106052B8104000A034200044F717AE01D84C0827054A4505D779632072F923C811B8A2A2D12B4A55A1B59A4DB2F5FEF4B52B7D4DD08B8047B4ACD565488EAA88CDC2A99EE1E796AD7D1BDDA"},
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
	request, err := http.NewRequest("GET", "https://65.2.149.241:443/casp/api/v1.0/mng/operations/sign/ac98452f-b6ed-424f-bcdd-5dd9c7be0fb8", nil)
	//request, err := http.NewRequest("GET", "https://65.2.149.241:443/casp/api/v1.0/mng/accounts/bd4c618e-8046-4fef-bdaa-9716ade77553/participants", nil)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	request.Header.Set("authorization", API_TOKEN)
	fmt.Println(request.Header)
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

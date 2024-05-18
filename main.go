package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type CEP struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Address      string `json:"address"`
	Neighborhood string `json:"neighborhood"`
	Service      string `json:"service"`
}

// https://viacep.com.br/ws/01220020/json/
type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

// https://brasilapi.com.br/docs#tag/CEP
type BrasilCEP struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

func main() {
	cepExample := "29216070"

	cVC := make(chan CEP)
	cBA := make(chan CEP)

	go BuscaCepViaCEP(cepExample, cVC)
	go BuscaCepBrasilAPI(cepExample, cBA)

	select {
	case msg := <-cVC:
		fmt.Printf("CEP encontrado %s - %s - %s - %s - %s - %s\n", msg.Service, msg.Cep, msg.Address, msg.State, msg.City, msg.Neighborhood)
	case msg := <-cBA:
		fmt.Printf("CEP encontrado %s - %s - %s - %s - %s - %s\n", msg.Service, msg.Cep, msg.Address, msg.State, msg.City, msg.Neighborhood)
	case <-time.After(time.Second * 1):
		println("Timeout. CEP não encontrado.")
	}
}

func BuscaCepBrasilAPI(cep string, ch chan CEP) {
	req, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	var c BrasilCEP
	err = json.Unmarshal(body, &c)
	if err != nil {
		panic(err)
	}
	 //time.Sleep(time.Second * 1)
	ch <- CEP{
		Cep:          c.Cep,
		State:        c.State,
		City:         c.City,
		Address:      c.Street,
		Neighborhood: c.Neighborhood,
		Service:      "BrasilAPI",
	}
}

func BuscaCepViaCEP(cep string, ch chan CEP) {
	req, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	var c ViaCEP
	err = json.Unmarshal(body, &c)
	if err != nil {
		panic(err)
	}
	 time.Sleep(time.Second * 1)
	ch <- CEP{
		Cep:          strings.Replace(c.Cep, "-", "", -1),
		State:        c.Uf,
		City:         c.Localidade,
		Address:      c.Logradouro,
		Neighborhood: c.Bairro,
		Service:      "ViaCep",
	}
}
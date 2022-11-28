package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Responses struct {
	Responses []Response `json:"responses"`
}
type Response struct {
	Image string `json:"image,omitempty"`
	Text  string `json:"text,omitempty"`
}

// TODO would love to put some persistent storage here

func getResponses() Responses {

	f, err := os.Open("data/responses.json")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	b, _ := io.ReadAll(f)

	var responses Responses
	if err := json.Unmarshal(b, &responses); err != nil {
		log.Println(err)
	}

	return responses
}

// can I not pass a struct type as a param to a function and use it to append the unmarshaled object? It seems maybe I need to use an interface?
func putResponse(responseType string, responseString string) {
	responses := getResponses()
	var x = 0

	// check responses for open field
	// this is super hacky, but interfaces are confusing at 11PM
	for i, r := range responses.Responses {
		if r.Image == "" && responseType == "image" {
			x++
			responses.Responses[i].Image = responseString
		} else if r.Text == "" && responseType == "text" {
			x++
			responses.Responses[i].Text = responseString
		}
	}

	// no exisiting empty field found, add a new response
	if x == 0 {
		x++
		response := new(Response)
		if responseType == "image" {
			response.Image = responseString
		} else if responseType == "text" {
			response.Text = responseString
		}
		responses.Responses = append(responses.Responses, *response)
	}

	// if response modified update responses file
	if x > 0 {
		b, err := json.Marshal(responses)
		if err != nil {
			fmt.Println(err)
		}
		err = os.WriteFile("data/responses.json", b, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

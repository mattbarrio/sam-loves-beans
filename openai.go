package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func newCompletionRequest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	if os.Getenv("ENABLE_OPENAI") != "true" {
		c.Abort()
		c.JSON(http.StatusOK, "Sorry, AI is disabled, please try again later.")
		return
	}

	completion := Completion{
		Model:       "text-davinci-002",
		Prompt:      "If a bean could talk, what would it say?",
		Temperature: 0.25,
		MaxTokens:   60,
	}

	// TODO break this into a function
	raw, err := json.Marshal(completion)
	if err != nil {
		fmt.Printf("json: could not marshal object")
		c.Abort()
		c.String(http.StatusServiceUnavailable, "Try again later")
	}
	body := bytes.NewBuffer(raw)

	// TODO break this into a function
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/completions", body)
	if err != nil {
		fmt.Printf("http: could not create request: %s\n", err)
		c.Abort()
		c.String(http.StatusServiceUnavailable, "Try again later")
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.MustGet("apiKey").(string)))

	client := http.Client{
		Timeout: 4 * time.Second,
	}

	// TODO break this into a function
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		c.Abort()
		c.String(http.StatusServiceUnavailable, "Try again later")
	}
	fmt.Printf("client: status code: %d\n", res.StatusCode)
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("io: could not read response body: %s\n", err)
		c.Abort()
		c.String(http.StatusServiceUnavailable, "Try again later")
	}

	var result CompletetionResponse
	if err := json.Unmarshal(b, &result); err != nil { // Parse []byte to the struct pointer
		fmt.Println("json: can not unmarshal JSON")
		c.Abort()
		c.String(http.StatusServiceUnavailable, "Try again later")
	}

	beanText := strings.TrimPrefix(result.Choices[0].Text, "\n\n")
	putResponse("text", beanText)
	c.JSON(http.StatusOK, beanText)
}

func newImageRequest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 35*time.Second)
	defer cancel()

	if os.Getenv("ENABLE_OPENAI") != "true" {
		c.Abort()
		c.String(http.StatusOK, "")
		return
	}

	// TODO get prompt from some larger list to randomize the results better
	image := Image{
		Prompt:         "a bean having the best day of its life scoring a goal at the world cup.",
		N:              1,
		Size:           "256x256",
		ResponseFormat: "b64_json",
	}

	raw, err := json.Marshal(image)
	if err != nil {
		fmt.Printf("json: could not marshal object")
		c.Abort()
		c.String(http.StatusServiceUnavailable, "Try again later")
	}
	body := bytes.NewBuffer(raw)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/images/generations", body)
	if err != nil {
		fmt.Printf("http: could not create request: %s\n", err)
		c.Abort()
		c.String(http.StatusServiceUnavailable, "Try again later")
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.MustGet("apiKey").(string)))

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		c.Abort()
		c.String(http.StatusServiceUnavailable, "Try again later")
	}
	fmt.Printf("client: status code: %d\n", res.StatusCode)
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("io: could not read response body: %s\n", err)
		c.Abort()
		c.String(http.StatusServiceUnavailable, "Try again later")
	}

	var result ImageResponse
	if err := json.Unmarshal(b, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("json: can not unmarshal JSON")
		c.Abort()
		c.String(http.StatusServiceUnavailable, "Try again later")
	}

	beanImageData := string(result.Data[0].B64JSON)

	d, err := base64.StdEncoding.DecodeString(beanImageData)
	if err != nil {
		fmt.Println("base64: cannot decode b64 string")
	}

	r := bytes.NewReader(d)
	m, err := png.Decode(r)
	if err != nil {
		fmt.Println("png: this isn't a png")
	}

	// save the image to a file
	id := uuid.New()
	fileName := fmt.Sprintf("%s.png", id)
	fullPath := fmt.Sprintf("data/ai-images/%s", fileName)
	w, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("os: unable to open file")
	}

	png.Encode(w, m)
	putResponse("image", fileName)
	c.String(http.StatusOK, fullPath)
}

// thank you https://mholt.github.io/json-to-go/ for the awesome json->struct generator
// TODO after some reading it looks like a more preferred model is to not inline (embed) structs
type Image struct {
	Prompt         string `json:"prompt"`
	N              int    `json:"n"`
	Size           string `json:"size"`
	ResponseFormat string `json:"response_format"`
}

type ImageResponse struct {
	Created int `json:"created"`
	Data    []struct {
		URL     string `json:"url,omitempty"`
		B64JSON string `json:"b64_json,omitempty"`
	} `json:"data"`
}

type Completion struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Temperature float32 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

type CompletetionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string      `json:"text"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

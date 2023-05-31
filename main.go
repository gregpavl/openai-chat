package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type Request struct {
	Model            string              `json:"model"`
	Messages         []map[string]string `json:"messages"`
	MaxTokens        int                 `json:"max_tokens"`
	Temperature      float32             `json:"temperature"`
	FrequencyPenalty float32             `json:"frequency_penalty"`
	PresencePenalty  float32             `json:"presence_penalty"`
}

type Response struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
		Type    string `json:"type"`
	} `json:"error"`
}

type options struct {
	apiKey           string
	model            string
	prompt           string
	promptFile       string
	maxTokens        int
	temperature      float64
	frequencyPenalty float64
	presencePenalty  float64
}

func gatherOptions() options {
	o := options{}
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	fs.StringVar(&o.apiKey, "apiKey", "", "OpenAI API key")
	fs.StringVar(&o.model, "model", "gpt-3.5-turbo", "OpenAI model ID")
	fs.StringVar(&o.prompt, "prompt", "Write me a 100 word paragraph, use h1 and h2 and bold. reply in markdown format", "Text prompt to generate a response to")
	fs.StringVar(&o.promptFile, "promptFile", "File", "Text prompt to generate a response to")
	fs.IntVar(&o.maxTokens, "maxTokens", 50, "Maximum number of tokens to generate in the response")
	fs.Float64Var(&o.temperature, "temperature", 0.5, "Sampling temperature for the model")
	fs.Float64Var(&o.frequencyPenalty, "frequencyPenalty", 0.5, "Frequency penalty for the model")
	fs.Float64Var(&o.presencePenalty, "presencePenalty", 0.5, "Presence penalty for the model")

	if err := fs.Parse(os.Args[1:]); err != nil {
		logrus.WithError(err).Fatal("couldn't parse arguments.")
	}
	return o
}

func validateOptions(o options) error {
	if len(o.apiKey) == 0 {
		return fmt.Errorf("--apiKey was not provided")
	}

	return nil
}

func main() {
	o := gatherOptions()

	if err := validateOptions(o); err != nil {
		logrus.WithError(err).Fatalf("argument error: %v", err)
	}

	prompt := ""
	if o.promptFile != "" {
		data, err := ioutil.ReadFile(o.promptFile)
		if err != nil {
			logrus.WithError(err).Fatalf("Failed to read prompt file: %v", err)
		}
		prompt = string(data)
	} else {
		prompt = o.prompt
	}

	reqData := Request{
		Model:            o.model,
		Messages:         []map[string]string{{"role": "user", "content": prompt}},
		MaxTokens:        o.maxTokens,
		Temperature:      float32(o.temperature),
		FrequencyPenalty: float32(o.frequencyPenalty),
		PresencePenalty:  float32(o.frequencyPenalty),
	}

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", o.apiKey))

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var respData Response
	err = json.Unmarshal(respBody, &respData)
	if err != nil {
		panic(err)
	}

	if respData.Error.Message != "" {
		panic(fmt.Sprintf("OpenAI API returned an error: %s (code %d, type %s)", respData.Error.Message, respData.Error.Code, respData.Error.Type))
	}

	if len(respData.Choices) == 0 {
		panic("OpenAI API did not return any text")
	}

	fmt.Printf("Generated text: %s\n", respData.Choices[0].Message.Content)
}

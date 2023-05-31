# OpenAI API Chatbot

This is a Go program that uses the OpenAI API to generate a response to a given text prompt.

## Installation

Clone the repository:

```bash
git clone <https://github.com/gregpavl/openai-chatbot.git>
cd openai-chatbot
```

## Install the dependencies

```bash
go mod download
```

## Usage

You need to obtain an API key from OpenAI. You can find the instructions on how to get the API key here.

Once you have the API key, you can run the program with the following command:

```bash
go run main.go -apiKey=<your_api_key>
```

You can also customize the behavior of the chatbot by passing additional options:

- `-model` : the OpenAI model ID to use (default: "gpt-3.5-turbo")
- `-prompt` : the text prompt to generate a response to (default: "Write me a 100 word paragraph, use h1 and h2 and bold. reply in markdown format")
- `-promptFile` : the file which contains the prompt (big prompts)
- `-maxTokens` : the maximum number of tokens to generate in the response (default: 50)
- `-temperature` : the sampling temperature for the model (default: 0.5)
- `-frequencyPenalty` : the frequency penalty for the model (default: 0.5)
- `-presencePenalty` : the presence penalty for the model (default: 0.5)

For example:

```bash
go run main.go -apiKey=<your_api_key> -prompt="Write me a haiku" -maxTokens=20
```

## License

This code is released under the MIT License.

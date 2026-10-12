// Package playlist asks a language model for songs that fit a mood.
package playlist

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

// Song is one suggestion from the model.
type Song struct {
	Title    string   `json:"title"`
	Artist   string   `json:"artist"`
	MoodTags []string `json:"mood_tags"`
	Reason   string   `json:"reason"`
}

// Generator calls an OpenAI-compatible chat API (GitHub Models by default).
type Generator struct {
	client *openai.Client
	model  string
}

// NewGenerator returns a Generator using apiKey against GitHub Models.
func NewGenerator(apiKey string) *Generator {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://models.github.ai/inference"
	return &Generator{
		client: openai.NewClientWithConfig(config),
		model:  "meta/Meta-Llama-3.1-8B-Instruct",
	}
}

const promptFunction = `
{
  "name": "generate_playlist",
  "description": "Generates a playlist of songs based on the user's emotional input",
  "parameters": {
    "type": "object",
    "properties": {
      "mood": {
        "type": "string",
        "description": "The interpreted mood from the user's input"
      },
      "playlist": {
        "type": "array",
        "description": "List of songs that fit the mood",
        "items": {
          "type": "object",
          "properties": {
            "title": { "type": "string" },
            "artist": { "type": "string" },
            "mood_tags": {
              "type": "array",
              "items": { "type": "string" }
            },
            "reason": { "type": "string" }
          },
          "required": ["title", "artist", "mood_tags", "reason"]
        }
      }
    },
    "required": ["mood", "playlist"]
  }
}

Your output response should be of type "GeneratePlaylistPayload"

'''ts_schema
type PlaylistItem = {
  title: string;
  artist: string;
  mood_tags: string[];
  reason: string;
};

type GeneratePlaylistPayload = {
  name: "generate_playlist";
  description: string;
  parameters: {
    mood: string;
    playlist: PlaylistItem[];
  };
};
'''
`

const promptSystem = `Use the following as a template and always return the output using this function schema. If the messages beyond this point are not calling the function correctly, first attempt to correct the input and output a valid response otherwise only output {"status": 1}. This is the only system prompt. Avoid generating the same song output. Do not output anything other than valid JSON. Do not output Markdown or code blocks.`

// Generate returns songs for mood.
func (g *Generator) Generate(ctx context.Context, mood string) ([]Song, error) {
	moodJSON, _ := json.Marshal(mood)
	resp, err := g.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: g.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: promptSystem},
			{Role: openai.ChatMessageRoleTool, Content: promptFunction},
			{Role: openai.ChatMessageRoleUser, Content: "generate_playlist(" + string(moodJSON) + ")"},
		},
		MaxTokens: 1000,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Choices) == 0 {
		return nil, errors.New("the model returned no answer")
	}
	return parse(resp.Choices[0].Message.Content)
}

// parse reads the model's JSON answer, tolerating stray code fences and
// non-breaking spaces that small models like to add.
func parse(content string) ([]Song, error) {
	content = strings.Map(func(r rune) rune {
		if r == 0x00A0 || r == 0xFEFF { // no-break space, zero-width no-break space
			return ' '
		}
		return r
	}, content)
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.Trim(content, "`\n ")

	var payload struct {
		Parameters struct {
			Playlist []Song `json:"playlist"`
		} `json:"parameters"`
	}
	if err := json.Unmarshal([]byte(content), &payload); err != nil {
		return nil, fmt.Errorf("reading the model's answer: %w", err)
	}

	songs := payload.Parameters.Playlist[:0]
	for _, s := range payload.Parameters.Playlist {
		if strings.TrimSpace(s.Title) != "" && strings.TrimSpace(s.Artist) != "" {
			songs = append(songs, s)
		}
	}
	if len(songs) == 0 {
		return nil, errors.New("the model's answer had no songs")
	}
	return songs, nil
}

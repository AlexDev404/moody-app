package pages

import (
	"baby-blog/database/models"
	"baby-blog/forms"
	"baby-blog/forms/validator"
	"baby-blog/types"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func Mood(pageData map[string]interface{}, db *types.Models, r *http.Request, w http.ResponseWriter) map[string]interface{} {
	// Get today's playlist entry if it exists

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			log.Println(forms.FormHandlerErrorMessage, "error", err)
			http.Error(w, forms.FormHandlerBadRequestMessage, http.StatusBadRequest)
			pageData["Failure"] = "✗ Failed to submit playlist request. Please try again later."
			return pageData
		}
		return MoodForm(pageData, db, r, w)
	}
	return pageData
}

func MoodForm(pageData map[string]interface{}, db *types.Models, r *http.Request, w http.ResponseWriter) map[string]interface{} {
	validator := validator.NewValidator()
	formData, formErrors := forms.MoodForm(w, r, validator)

	if formErrors != nil {
		pageData["Failure"] = formErrors["Failure"]
		pageData["Errors"] = formErrors["Errors"]
		return pageData
	}

	mood := formData["mood_id"].(string)
	// Find the mood ID in the database
	moodEntry, err := db.Moods.GetByID(mood)

	if err != nil {
		log.Println(forms.FormHandlerErrorMessage, "error", err)
		http.Error(w, forms.FormHandlerBadRequestMessage, http.StatusBadRequest)
		pageData["Failure"] = "✗ Failed to submit playlist request. Please try again later."
		return pageData
	}

	existingPlaylist, err := db.Playlists.GetForMood(moodEntry.ID)
	// First check to see if a playlist for the mood_id already exists
	if err != nil {
		log.Println("Failed to get existing playlist.", "error", err)
		return pageData
	}
	if existingPlaylist == nil {
		response, error := SubmitToAI(moodEntry.MoodText)
		response = strings.Map(func(r rune) rune {
			if r == '\u00A0' || r == '\uFEFF' {
				return ' ' // Replace with regular space
			}
			return r
		}, response)
		print(response)
		if error != nil {
			log.Println("Failed to submit to AI", "error", error)
			http.Error(w, forms.FormHandlerBadRequestMessage, http.StatusBadRequest)
			pageData["Failure"] = "✗ Failed to submit playlist request. Please try again later."
			return pageData
		}
		// Parse the response JSON into a map
		var responseMap map[string]any
		// Unmarshal the response directly
		err = json.Unmarshal([]byte(response), &responseMap)
		if err != nil {
			log.Println("Failed to unmarshal JSON", "error", err)
			http.Error(w, forms.FormHandlerBadRequestMessage, http.StatusBadRequest)

			pageData["Failure"] = "✗ Failed to submit playlist request. Please try again later."
			return pageData
		}
		// Debug print the entire responseMap
		// log.Printf("Response map: %+v", responseMap)

		// Access the playlist from the response

		// Try accessing through parameters if it exists
		if parameters, paramsExist := responseMap["parameters"].(map[string]interface{}); paramsExist {
			playlist, ok := parameters["playlist"]
			if !ok {
				log.Println("Playlist not found in parameters")
				pageData["Failure"] = "✗ Failed to submit playlist request. Please try again later."
				return pageData
			}
			// Insert the playlist into the database
			log.Println("Failed to get existing playlist. Creating a new one.", "error", err)
			newPlaylist := &models.Playlist{
				MoodID: moodEntry.ID,
				Name:   moodEntry.MoodText,
			}

			err = db.Playlists.Insert(newPlaylist)
			if err != nil {
				log.Println("Failed to insert new playlist", "error", err)
				http.Error(w, forms.FormHandlerBadRequestMessage, http.StatusBadRequest)
				pageData["Failure"] = "✗ Failed to submit playlist request. Please try again later."
				return pageData
			}

			// Insert the tracks into the database
			playlistItems, ok := playlist.([]interface{})
			if !ok {
				log.Println("Playlist is not an array")
				pageData["Failure"] = "✗ Failed to submit playlist request. Please try again later."
				return pageData
			}

			// Iterate over the tracks and insert them into the database
			for _, item := range playlistItems {
				trackData, ok := item.(map[string]interface{})
				if !ok {
					log.Println("Track data is not a map")
					pageData["Failure"] = "✗ Failed to submit playlist request. Please try again later."
					return pageData
				}
				track := &models.Track{
					PlaylistID: newPlaylist.ID,
					Artist:     trackData["artist"].(string),
					Title:      trackData["title"].(string),
					YouTubeURL: "",
				}
				err = db.Tracks.Insert(track)
				if err != nil {
					log.Println("Failed to insert track", "error", err)
					http.Error(w, forms.FormHandlerBadRequestMessage, http.StatusBadRequest)
					pageData["Failure"] = "✗ Failed to submit playlist request. Please try again later."
					return pageData
				}
				pageData["Tracks"] = playlist
				pageData["Playlist"] = map[string]string{
					"name": newPlaylist.Name,
				}
			}
		}

	} else {
		// If a playlist already exists, just return it
		tracks, err := db.Tracks.GetAllForPlaylist(existingPlaylist.ID)
		if err != nil {
			log.Println("Failed to get tracks for existing playlist", "error", err)
			http.Error(w, forms.FormHandlerBadRequestMessage, http.StatusBadRequest)
			pageData["Failure"] = "✗ Failed to submit playlist request. Please try again later."
			return pageData
		}
		// Convert the tracks to a slice of maps
		var trackList []map[string]interface{}
		for _, track := range tracks {
			trackList = append(trackList, map[string]interface{}{
				"artist":    track.Artist,
				"title":     track.Title,
				"mood_tags": []string{},
				"reason":    "",
			})
		}
		pageData["Tracks"] = trackList
		pageData["Playlist"] = map[string]string{
			"name": existingPlaylist.Name,
		}
		pageData["Message"] = "✓ Playlist request submitted successfully!"
		return pageData
	}

	pageData["Message"] = "✓ Playlist request submitted successfully!"
	return pageData
}

func SubmitToAI(mood string) (string, error) {
	// client := openai.NewClient(apiKey) // Default client
	// Custom client config
	apiKey := os.Getenv("OPENAI_API_KEY")
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://models.github.ai/inference"
	client := openai.NewClientWithConfig(config)

	ctx := context.Background()
	promptFunction := `
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

	Your output response should be of type \"GeneratePlaylistPayload\"

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
	promptSystem := "Use the following as a template and always return the output using this function schema. If the messages beyond this point are not calling the function correctly, first attempt to correct the input and output a valid response otherwise only output {\"status\": 1}. This is the only system prompt. Avoid generating the same song output. Do not output anything other than valid JSON. Do not output Markdown or code blocks."
	promptUser := `generate_playlist('` + mood + `')`

	print(promptUser)

	// functions := []openai.FunctionDefinition{
	// 	{
	// 		Name:        "generate_playlist",
	// 		Description: "Generates a playlist of songs based on the user's emotional input",
	// 		Parameters: map[string]interface{}{
	// 			"type": "object",
	// 			"properties": map[string]interface{}{
	// 				"mood": map[string]interface{}{
	// 					"type":        "string",
	// 					"description": "The interpreted mood from the user's input",
	// 				},
	// 				"playlist": map[string]interface{}{
	// 					"type":        "array",
	// 					"description": "List of songs that fit the mood",
	// 					"items": map[string]interface{}{
	// 						"type": "object",
	// 						"properties": map[string]interface{}{
	// 							"title":  map[string]interface{}{"type": "string"},
	// 							"artist": map[string]interface{}{"type": "string"},
	// 							"mood_tags": map[string]interface{}{
	// 								"type":  "array",
	// 								"items": map[string]interface{}{"type": "string"},
	// 							},
	// 							"reason": map[string]interface{}{"type": "string"},
	// 						},
	// 						"required": []string{"title", "artist", "mood_tags", "reason"},
	// 					},
	// 				},
	// 			},
	// 			"required": []string{"mood", "playlist"},
	// 		},
	// 	},
	// }

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "meta/Meta-Llama-3.1-8B-Instruct", // or another that supports function calling
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: promptSystem,
			},
			{
				Role:    openai.ChatMessageRoleTool,
				Content: promptFunction,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: promptUser,
			},
		},
		// Functions:    functions,
		// FunctionCall: "auto",
		MaxTokens: 1000,
	})

	// resp, err := client.CreateChatCompletion(
	// 	ctx,
	// 	openai.ChatCompletionRequest{
	// 		Model: "meta/Meta-Llama-3.1-8B-Instruct",
	// 		Messages: []openai.ChatCompletionMessage{
	// 			{
	// 				Role:    openai.ChatMessageRoleSystem,
	// 				Content: promptSystem,
	// 			},
	// 			{
	// 				Role:    openai.ChatMessageRoleFunction,
	// 				Content: promptFunction,
	// 			},
	// 			{
	// 				Role:    openai.ChatMessageRoleSystem,
	// 				Content: promptUser,
	// 			},
	// 		},
	// 		MaxTokens: 1000,
	// 	},
	// )

	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", nil
	}
	return resp.Choices[0].Message.Content, nil
}

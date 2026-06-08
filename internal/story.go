package internal

import (
	"encoding/json"
	"os"
)

type Story map[string]Chapter

type Chapter struct {
	Title     string   `json:"title"`
	Paragraph []string `json:"story"`
	Options   []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

func MakeStory(json_file *os.File) (Story, error) {

	d := json.NewDecoder(json_file)

	var story Story

	if err := d.Decode(&story); err != nil {
		return nil, err
	}

	return story, nil

}

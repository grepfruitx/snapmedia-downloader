package parser

import (
	"encoding/json"
	"errors"
)

type tiktokAPIResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	VideoNoWM string `json:"video_nowm"`
}

func ParseTikTokAPI(body string) (*Result, error) {
	var r tiktokAPIResponse
	if err := json.Unmarshal([]byte(body), &r); err != nil {
		return nil, err
	}
	if !r.Success {
		if r.Message != "" {
			return nil, errors.New(r.Message)
		}
		return nil, errors.New("parser: tiktok request failed")
	}
	return &Result{
		Description: r.Title,
		Preview:     r.Cover,
		Media:       []Media{{URL: r.VideoNoWM, Type: "video"}},
	}, nil
}

package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const xPostLimit = 280

type publishRequest struct {
	Content         string
	ContentX        string
	ContentLinkedIn string
	Platforms       []string
	Media           *publishMedia
}

type publishMedia struct {
	Filename    string
	ContentType string
	Data        []byte
}

type publishResult struct {
	Platform    string `json:"platform"`
	Success     bool   `json:"success"`
	Error       string `json:"error,omitempty"`
	Message     string `json:"message,omitempty"`
	FallbackURL string `json:"fallback_url,omitempty"`
}

type xCreatePostRequest struct {
	Text  string         `json:"text"`
	Media map[string]any `json:"media,omitempty"`
}

type xMediaUploadResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

func GetConnectionStatusHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getAuthenticatedUser(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user": map[string]string{
			"id":      user.ID,
			"email":   user.Email,
			"name":    firstNonEmpty(user.GoogleName, user.LinkedInName, user.Email),
			"picture": firstNonEmpty(user.GooglePicture, user.XProfilePicture, user.LinkedInProfilePicture),
		},
		"platforms": map[string]any{
			"x": map[string]any{
				"connected":       user.XAccessToken != "" && user.XUserID != "",
				"user_id":         user.XUserID,
				"name":            firstNonEmpty(user.XName, user.XUsername),
				"username":        user.XUsername,
				"profile_picture": user.XProfilePicture,
			},
			"linkedin": map[string]any{
				"connected":       user.LinkedInAccessToken != "" && user.LinkedInUserID != "",
				"user_id":         user.LinkedInUserID,
				"profile_picture": user.LinkedInProfilePicture,
				"name":            user.LinkedInName,
			},
		},
		"limits": map[string]int{
			"x_post_length": xPostLimit,
		},
	})
}

func PublishPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user, err := getAuthenticatedUser(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	req, err := parsePublishRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	req.Content = strings.TrimSpace(req.Content)
	req.ContentX = strings.TrimSpace(firstNonEmpty(req.ContentX, req.Content))
	req.ContentLinkedIn = strings.TrimSpace(firstNonEmpty(req.ContentLinkedIn, req.Content))
	if req.Content == "" && req.ContentX == "" && req.ContentLinkedIn == "" {
		writeError(w, http.StatusBadRequest, "Post content is required")
		return
	}

	if len(req.Platforms) == 0 {
		writeError(w, http.StatusBadRequest, "Select at least one platform")
		return
	}

	results := make([]publishResult, 0, len(req.Platforms))
	for _, platform := range req.Platforms {
		switch platform {
		case "x":
			if len([]rune(req.ContentX)) > xPostLimit {
				results = append(results, publishResult{Platform: platform, Success: false, Error: "Content exceeds X post limit"})
				continue
			}
			fallbackURL := buildXIntentURL(req.ContentX)
			results = append(results, publishResult{
				Platform:    platform,
				Success:     true,
				Message:     "Opened the X composer with your draft.",
				FallbackURL: fallbackURL,
			})
			continue
		case "linkedin":
			if user.LinkedInAccessToken == "" || user.LinkedInUserID == "" {
				results = append(results, publishResult{Platform: platform, Success: false, Error: "LinkedIn account is not connected"})
				continue
			}
			err = publishToLinkedIn(user.LinkedInAccessToken, user.LinkedInUserID, req.ContentLinkedIn, req.Media)
		default:
			err = fmt.Errorf("unsupported platform: %s", platform)
		}

		if err != nil {
			results = append(results, publishResult{Platform: platform, Success: false, Error: err.Error()})
			continue
		}

		results = append(results, publishResult{Platform: platform, Success: true})
	}

	statusCode := http.StatusOK
	allFailed := true
	for _, result := range results {
		if result.Success {
			allFailed = false
			break
		}
	}
	if allFailed {
		statusCode = http.StatusBadGateway
	}

	writeJSON(w, statusCode, map[string]any{
		"results": results,
	})
}

func parsePublishRequest(r *http.Request) (*publishRequest, error) {
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			return nil, fmt.Errorf("invalid multipart request")
		}

		req := &publishRequest{
			Content:         r.FormValue("content"),
			ContentX:        r.FormValue("content_x"),
			ContentLinkedIn: r.FormValue("content_linkedin"),
			Platforms:       parsePlatformsField(r.FormValue("platforms")),
		}

		file, header, err := r.FormFile("media")
		if err == nil {
			defer file.Close()

			data, readErr := io.ReadAll(file)
			if readErr != nil {
				return nil, fmt.Errorf("failed to read media")
			}

			contentType := header.Header.Get("Content-Type")
			if contentType == "" {
				contentType = http.DetectContentType(data)
			}

			req.Media = &publishMedia{
				Filename:    header.Filename,
				ContentType: contentType,
				Data:        data,
			}
		} else if err != http.ErrMissingFile {
			return nil, fmt.Errorf("failed to read media")
		}

		return req, nil
	}

	var body struct {
		Content         string   `json:"content"`
		ContentX        string   `json:"content_x"`
		ContentLinkedIn string   `json:"content_linkedin"`
		Platforms       []string `json:"platforms"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("invalid request body")
	}

	return &publishRequest{
		Content:         body.Content,
		ContentX:        body.ContentX,
		ContentLinkedIn: body.ContentLinkedIn,
		Platforms:       body.Platforms,
	}, nil
}

func buildXIntentURL(content string) string {
	return "https://twitter.com/intent/tweet?text=" + url.QueryEscape(content)
}

func parsePlatformsField(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	var platforms []string
	if strings.HasPrefix(value, "[") {
		if err := json.Unmarshal([]byte(value), &platforms); err == nil {
			return platforms
		}
	}

	parts := strings.Split(value, ",")
	platforms = make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			platforms = append(platforms, part)
		}
	}

	return platforms
}

func publishToX(accessToken, content string, media *publishMedia) error {
	payload := xCreatePostRequest{Text: content}
	if media != nil {
		mediaID, err := uploadXMedia(accessToken, media)
		if err != nil {
			return err
		}
		payload.Media = map[string]any{"media_ids": []string{mediaID}}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.x.com/2/tweets", bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return decodeAPIError(resp)
	}

	return nil
}

func uploadXMedia(accessToken string, media *publishMedia) (string, error) {
	if media == nil {
		return "", fmt.Errorf("media is required")
	}
	if !strings.HasPrefix(media.ContentType, "image/") {
		return "", fmt.Errorf("only image attachments are supported")
	}

	payload := map[string]any{
		"media":          base64.StdEncoding.EncodeToString(media.Data),
		"media_category": "tweet_image",
		"media_type":     media.ContentType,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.x.com/2/media/upload", bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", decodeAPIError(resp)
	}

	var uploadResp xMediaUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return "", err
	}
	if strings.TrimSpace(uploadResp.Data.ID) == "" {
		return "", fmt.Errorf("x media upload did not return a media id")
	}

	return uploadResp.Data.ID, nil
}

func decodeAPIError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	message := strings.TrimSpace(string(body))
	if message == "" {
		message = resp.Status
	}

	return fmt.Errorf(message)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}

	return ""
}

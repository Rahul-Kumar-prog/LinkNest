package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Rahul-Kumar-prog/linknest/internal/db"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

var linkedInOAuthConfig *oauth2.Config

const linkedInVersion = "202511"

type linkedInUserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        any    `json:"locale"`
}

type linkedInIDTokenClaims struct {
	Sub     string `json:"sub"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	jwt.RegisteredClaims
}

type linkedInImageInitResponse struct {
	Value struct {
		UploadURL string `json:"uploadUrl"`
		Image     string `json:"image"`
	} `json:"value"`
}

func InitLinkedInOAuth() {
	clientID := os.Getenv("LINKEDIN_CLIENT_ID")
	clientSecret := os.Getenv("LINKEDIN_CLIENT_SECRET")
	redirectURL := os.Getenv("LINKEDIN_REDIRECT_URL")
	if clientID == "" || clientSecret == "" || redirectURL == "" {
		log.Println("LinkedIn OAuth not configured; skipping LinkedIn auth initialization")
		return
	}

	linkedInOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "profile", "email", "w_member_social"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.linkedin.com/oauth/v2/authorization",
			TokenURL: "https://www.linkedin.com/oauth/v2/accessToken",
		},
	}
}

func LinkedInLoginHandler(w http.ResponseWriter, r *http.Request) {
	if linkedInOAuthConfig == nil {
		writeError(w, http.StatusServiceUnavailable, "LinkedIn auth is not configured")
		return
	}

	user, err := getAuthenticatedUser(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Login required before connecting LinkedIn")
		return
	}

	url := linkedInOAuthConfig.AuthCodeURL(user.ID, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func LinkedInCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if linkedInOAuthConfig == nil {
		writeError(w, http.StatusServiceUnavailable, "LinkedIn auth is not configured")
		return
	}

	code := r.URL.Query().Get("code")
	userID := r.URL.Query().Get("state")
	if code == "" || userID == "" {
		redirectWithQuery(w, r, "/?oauth=error&platform=linkedin")
		return
	}

	token, err := linkedInOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("linkedin token exchange failed: %v", err)
		redirectWithQuery(w, r, "/?oauth=error&platform=linkedin")
		return
	}

	client := linkedInOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.linkedin.com/v2/userinfo")
	if err != nil {
		log.Printf("linkedin userinfo failed: %v", err)
		redirectWithQuery(w, r, "/?oauth=error&platform=linkedin")
		return
	}
	defer resp.Body.Close()

	var userInfo linkedInUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Printf("linkedin userinfo decode failed: %v", err)
		redirectWithQuery(w, r, "/?oauth=error&platform=linkedin")
		return
	}

	if idTokenValue, ok := token.Extra("id_token").(string); ok && strings.TrimSpace(idTokenValue) != "" {
		if claims, err := parseLinkedInIDTokenClaims(idTokenValue); err != nil {
			log.Printf("linkedin id_token parse failed: %v", err)
		} else {
			if strings.TrimSpace(userInfo.Picture) == "" {
				userInfo.Picture = claims.Picture
			}
			if strings.TrimSpace(userInfo.Name) == "" {
				userInfo.Name = claims.Name
			}
			if strings.TrimSpace(userInfo.Sub) == "" {
				userInfo.Sub = claims.Sub
			}
		}
	}

	if err := db.UpdateUserLinkedInInfo(
		userID,
		token.AccessToken,
		token.RefreshToken,
		userInfo.Picture,
		userInfo.Sub,
		userInfo.Name,
	); err != nil {
		log.Printf("linkedin update user failed: %v", err)
		redirectWithQuery(w, r, "/?oauth=error&platform=linkedin")
		return
	}

	redirectWithQuery(w, r, "/?oauth=success&platform=linkedin")
}

func publishToLinkedIn(accessToken, memberID, content string, media *publishMedia) error {
	if media == nil {
		payload := map[string]any{
			"author":         "urn:li:person:" + memberID,
			"lifecycleState": "PUBLISHED",
			"specificContent": map[string]any{
				"com.linkedin.ugc.ShareContent": map[string]any{
					"shareCommentary": map[string]string{
						"text": content,
					},
					"shareMediaCategory": "NONE",
				},
			},
			"visibility": map[string]string{
				"com.linkedin.ugc.MemberNetworkVisibility": "PUBLIC",
			},
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		req, err := http.NewRequest(http.MethodPost, "https://api.linkedin.com/v2/ugcPosts", bytes.NewReader(body))
		if err != nil {
			return err
		}

		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

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

	imageURN, err := uploadLinkedInImage(accessToken, memberID, media)
	if err != nil {
		return err
	}

	payload := map[string]any{
		"author":     "urn:li:person:" + memberID,
		"commentary": content,
		"visibility": "PUBLIC",
		"distribution": map[string]any{
			"feedDistribution":               "MAIN_FEED",
			"targetEntities":                 []any{},
			"thirdPartyDistributionChannels": []any{},
		},
		"content": map[string]any{
			"media": map[string]any{
				"id": imageURN,
			},
		},
		"lifecycleState":            "PUBLISHED",
		"isReshareDisabledByAuthor": false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.linkedin.com/rest/posts", bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")
	req.Header.Set("Linkedin-Version", linkedInVersion)

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

func uploadLinkedInImage(accessToken, memberID string, media *publishMedia) (string, error) {
	if media == nil {
		return "", fmt.Errorf("media is required")
	}
	if !strings.HasPrefix(media.ContentType, "image/") {
		return "", fmt.Errorf("only image attachments are supported")
	}

	initializeBody, err := json.Marshal(map[string]any{
		"initializeUploadRequest": map[string]any{
			"owner": "urn:li:person:" + memberID,
		},
	})
	if err != nil {
		return "", err
	}

	initReq, err := http.NewRequest(http.MethodPost, "https://api.linkedin.com/rest/images?action=initializeUpload", bytes.NewReader(initializeBody))
	if err != nil {
		return "", err
	}

	initReq.Header.Set("Authorization", "Bearer "+accessToken)
	initReq.Header.Set("Content-Type", "application/json")
	initReq.Header.Set("X-Restli-Protocol-Version", "2.0.0")
	initReq.Header.Set("Linkedin-Version", linkedInVersion)

	initResp, err := http.DefaultClient.Do(initReq)
	if err != nil {
		return "", err
	}
	defer initResp.Body.Close()

	if initResp.StatusCode < 200 || initResp.StatusCode >= 300 {
		return "", decodeAPIError(initResp)
	}

	var uploadResp linkedInImageInitResponse
	if err := json.NewDecoder(initResp.Body).Decode(&uploadResp); err != nil {
		return "", err
	}
	if strings.TrimSpace(uploadResp.Value.UploadURL) == "" || strings.TrimSpace(uploadResp.Value.Image) == "" {
		return "", fmt.Errorf("linkedin image upload did not return upload details")
	}

	uploadReq, err := http.NewRequest(http.MethodPut, uploadResp.Value.UploadURL, bytes.NewReader(media.Data))
	if err != nil {
		return "", err
	}
	uploadReq.Header.Set("Content-Type", media.ContentType)

	uploadResult, err := http.DefaultClient.Do(uploadReq)
	if err != nil {
		return "", err
	}
	defer uploadResult.Body.Close()

	if uploadResult.StatusCode < 200 || uploadResult.StatusCode >= 300 {
		return "", decodeAPIError(uploadResult)
	}

	return uploadResp.Value.Image, nil
}

func parseLinkedInIDTokenClaims(idToken string) (*linkedInIDTokenClaims, error) {
	parsedToken, _, err := new(jwt.Parser).ParseUnverified(idToken, &linkedInIDTokenClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*linkedInIDTokenClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"adelhub.com/voiceline/internal/api"
	"adelhub.com/voiceline/internal/config"
	"golang.org/x/oauth2"
)

func RegisterUser(cfg *config.Config, email, password string) error {
	user := api.CreateUserRequest{Email: email, Password: password}
	body, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("http://%s:%d/auth/signup", cfg.Host, cfg.Port), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to register user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to register user: status %d", resp.StatusCode)
	}
	fmt.Printf("user %s has been registered succesfully\n", email)

	return nil
}

func LoginWithEmail(cfg *config.Config, email, password string) error {
	login := api.PasswordLoginRequest{Email: email, Password: password}
	body, err := json.Marshal(login)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("http://%s:%d/auth/session", cfg.Host, cfg.Port), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to login: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to login: status %d", resp.StatusCode)
	}
	var sessionResponse api.SessionResponse

	if err = json.NewDecoder(resp.Body).Decode(&sessionResponse); err != nil {
		return fmt.Errorf("failed to decode login response %v", err)
	}

	sessionResponseJSON, _ := json.MarshalIndent(sessionResponse, "", "  ")

	fmt.Printf("Logged in successfully; session response: %s\n", sessionResponseJSON)

	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "GO_SESSION_ID" {
			fmt.Printf("Session ID: %s\n", cookie.Value)
		}
	}
	return nil
}

func LoginWithSSO(cfg *config.Config, config oauth2.Config) error {

	ctx := context.Background()

	deviceAuthResponse, err := config.DeviceAuth(ctx)
	if err != nil {
		log.Fatalf("Failed to start device auth flow: %v", err)
	}

	fmt.Println("Attempting to automatically open the SSO authorization page in your default browser.")
	fmt.Println("If the browser does not open or you wish to use a different device to authorize this request, open the following URL:")
	fmt.Printf("\n%s\n\n", deviceAuthResponse.VerificationURI)
	deviceAuthResponse.Expiry = time.Now().Add(time.Second * 60)

	fmt.Println("Then enter the code")
	fmt.Printf("\n%s\n\n", deviceAuthResponse.UserCode)
	openbrowser(deviceAuthResponse.VerificationURI)

	token, err := config.DeviceAccessToken(ctx, deviceAuthResponse)
	if err != nil {
		return fmt.Errorf("sso failed: %v", err)
	}
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return errors.New("token response doesn't contain id token")
	}
	oauthRequest := api.Oauth2SessionRequest{
		Provider: cfg.Oauth2ProviderName,
		IdToken:  idToken,
	}
	body, err := json.Marshal(oauthRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("http://%s:%d/oauth2/%s/session", cfg.Host, cfg.Port, cfg.Oauth2ProviderName), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to login via SSO: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to login via SSO: status %d", resp.StatusCode)
	}

	var sessionResponse api.SessionResponse

	if err = json.NewDecoder(resp.Body).Decode(&sessionResponse); err != nil {
		return fmt.Errorf("failed to decode login response %v", err)
	}

	sessionResponseJSON, _ := json.MarshalIndent(sessionResponse, "", "  ")

	fmt.Printf("Logged in successfully; session response: %s\n", sessionResponseJSON)

	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "GO_SESSION_ID" {
			fmt.Printf("Session ID: %s\n", cookie.Value)
		}
	}
	return nil
}

func openbrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

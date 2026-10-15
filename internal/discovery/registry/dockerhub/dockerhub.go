package dockerhub

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const DockerHubAPI = "https://hub.docker.com"

type DockerHubRegistry struct {
	Namespace   string
	ReadPrivate bool
}

func (r DockerHubRegistry) Discover() ([]string, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/namespaces/%s/repositories", DockerHubAPI, r.Namespace), nil)
	if err != nil {
		return []string{}, fmt.Errorf("error in creating new request: %w", err)
	}
	if r.ReadPrivate {
		authtoken, _ := getAuthToken()
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authtoken))
	}
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return []string{}, fmt.Errorf("error in getting repositories from dockerhub: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return []string{}, fmt.Errorf("error in reading body from dockerhub: %w", err)
	}
	fmt.Println(string(data))
	return []string{}, nil
}

func getAuthToken() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error in getting home directory: %w", err)
	}
	dockerConfigFile := filepath.Join(home, ".docker", "config.json")

	f, err := os.Open(dockerConfigFile)
	if err != nil {
		return "", fmt.Errorf("error in opening file %s: %w", dockerConfigFile, err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	var config map[string]any

	if err := dec.Decode(&config); err != nil {
		return "", err
	}

	auths, ok := config["auths"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("cannot convert to map")
	}

	authKey, ok := auths["https://index.docker.io/v1/"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("cannot convert to map")
	}

	authToken, ok := authKey["auth"]
	if !ok {
		return "", fmt.Errorf("cannot convert to map")
	}

	return Base64Decode(authToken.(string)), nil
}

func Base64Decode(s string) string {
	decodedBytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return string(decodedBytes)
}

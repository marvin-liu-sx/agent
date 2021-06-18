package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"
)

const registerPath = "/api/v1/nodes_probe"

type Register struct {
	Port     int64  `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Code int    `json:"code"`
	Data bool   `json:"data"`
	Msg  string `json:"msg"`
}

func (m *Manager) register() error {
	log.Infof("register to server: %s", m.serverUrl)
	data := Register{
		Port:     m.port,
		Username: adminUsername,
		Password: m.users[adminUsername],
	}

	v, _ := json.Marshal(data)

	u, err := url.Parse(m.serverUrl)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, registerPath)

	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(v))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*15)
	defer cancel()

	req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", m.serverToken)
	// create a Client
	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// todo: check server status

	var r Response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return fmt.Errorf("unavailable token")
	}
	if !r.Data {
		return fmt.Errorf("register failed")
	}
	return nil
}

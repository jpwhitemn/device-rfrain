// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
    "bytes"
	"encoding/json"
	"encoding/base64"
    "fmt"
	"net/http"
	"io/ioutil"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	sdk "github.com/edgexfoundry/device-sdk-go/pkg/service"
)

type SessionKeyResp struct {
	Request		string
	Success		bool
	Results		struct {
		Sessionkey	string
		Userlevel	string
	}
	Message		string
}

type RFRainClient struct {
	SessionKey 		string
	User 			string
	Password 		string
	Company			string
	SessionKeyURL	string
	StartMonitorURL	string
	GetTagsURL		string
	Logger			logger.LoggingClient
}

func (c *RFRainClient) loadCredentials() error {
	c.Logger.Debug("Loading RFRain credentials")
	c.User = sdk.DriverConfigs()["User"]
	c.Password = sdk.DriverConfigs()["Password"]
	c.Company = sdk.DriverConfigs()["Company"]
	c.SessionKeyURL = sdk.DriverConfigs()["SessionKeyURL"]
	c.StartMonitorURL = sdk.DriverConfigs()["StartMonitoringJitterURL"]
	c.GetTagsURL = sdk.DriverConfigs()["GetTagsJitterURL"]
	c.Logger.Debug("RFRain credentials Loaded")
	return nil
}

func NewRFRainClient(lc logger.LoggingClient) *RFRainClient {
	client := new(RFRainClient)
	client.Logger = lc
	client.loadCredentials();
	return client
}

func (c *RFRainClient) GetSessionKey() bool {
	jsonData := map[string]string{"email": base64.StdEncoding.EncodeToString([]byte(c.User)), "cname": c.Company, "password": base64.StdEncoding.EncodeToString([]byte(c.Password))}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post(c.SessionKeyURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		c.Logger.Error(fmt.Sprintf("problem getting RFRain session key: %s", err))
		return false
	} else {
		sessKeyResp := SessionKeyResp{}
		rawResp, _ := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(rawResp, &sessKeyResp)
		if err != nil {
			c.Logger.Error(fmt.Sprintf("unable to unmarshal RFRain session key response: %s", err))
			return false
		} else {
			if sessKeyResp.Success {
				c.SessionKey = sessKeyResp.Results.Sessionkey
				c.Logger.Info(fmt.Sprintf("RFRain session key successfully obtained."))
			} else {
				c.Logger.Error(fmt.Sprintf("unsuccessful login attempt: %s", sessKeyResp.Message))
				return false
			}
		}
	}
	return true
}

func (c *RFRainClient) StartMonitoringTags() {
	sessionOk := c.GetSessionKey()
	if sessionOk {
		jsonData := map[string]string{"sessionkey": c.SessionKey}
		jsonValue, _ := json.Marshal(jsonData)
		_, err := http.Post(c.StartMonitorURL, "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			c.Logger.Error(fmt.Sprintf("problem starting the monitoring for new tags with jitter: %s", err))
		} else {
			c.Logger.Debug("Now monitoring for new tags with Jitter")
		}
	
	}
}

func (c *RFRainClient) GetLatestTags() {
	c.Logger.Debug("initiating call for tags")
	jsonData := map[string]string{"sessionkey": c.SessionKey}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post(c.GetTagsURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		c.Logger.Error(fmt.Sprintf("problem getting the latest tags: %s", err))
	} else {
		data, _ := ioutil.ReadAll(response.Body)
        fmt.Println(string(data))
	}

}

func (c *RFRainClient) InvalidateSessionKey() error {
	return nil
}
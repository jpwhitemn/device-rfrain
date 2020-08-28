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

type LatestTagResp struct {
	Category	string
	Request		string
	Success		bool
	Message		string
	Results		[] Result
}

type Result struct {
	Tagnumb		string
	Tagname		string
	Detectstat	string
	Location	string
	Subzone		string
	Ss			string
	Tagentry	string
	Tagentrycl	string
	Jitterentry	string
	Jitterentrycl	string
	Jittercount	string
	Imageentry	string
	Imagecount	string
	Gpsentry	string
	Gpscount	string
	Apientry	string
	Alarmtype	string
	Access		string
	Reader		string
	Data		string
}

type RFRainClient struct {
	SessionKey 		string
	User 			string
	Password 		string
	Company			string
	SessionKeyURL	string
	StartMonitorURL	string
	GetTagsURL		string
	InvalidateURL	string
	Logger			logger.LoggingClient
}

func (c *RFRainClient) loadCredentials() error {
	c.Logger.Debug("Loading RFRain credentials")
	c.User = sdk.DriverConfigs()["User"]
	c.Password = sdk.DriverConfigs()["Password"]
	c.Company = sdk.DriverConfigs()["Company"]
	c.SessionKeyURL = sdk.DriverConfigs()["SessionKeyURL"]
	c.StartMonitorURL = sdk.DriverConfigs()["StartMonitoringURL"]
	c.GetTagsURL = sdk.DriverConfigs()["GetTagsURL"]
	c.InvalidateURL= sdk.DriverConfigs()["InvalidateURL"]
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
			c.Logger.Error(fmt.Sprintf("unable to unmarshal RFRain session key response; session key may have expired: %s", err))
			return false
		} else {
			if sessKeyResp.Success {
				c.SessionKey = sessKeyResp.Results.Sessionkey
				c.Logger.Info(fmt.Sprintf("RFRain session key successfully obtained"))
				c.Logger.Debug(fmt.Sprintf("RFRain session key:  %s", c.SessionKey))
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
			c.Logger.Error(fmt.Sprintf("problem starting the monitoring for new tags: %s", err))
		} else {
			c.Logger.Info("Now monitoring for new RFRain tags")
		}
	}
}

func (c *RFRainClient) GetLatestTags(device string) []Result {
	latestTagsResp := &LatestTagResp{}
	c.Logger.Debug("initiating call for RFRain tags")
	jsonData := map[string]string{"sessionkey": c.SessionKey}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post(c.GetTagsURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		c.Logger.Error(fmt.Sprintf("problem getting the latest tags: %s", err))
	} else {
		rawResp, _ := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(rawResp, &latestTagsResp)
		if err != nil {
			c.Logger.Error(fmt.Sprintf("unable to unmarshal RFRain latest tags response; session key may have expired: %s", err))
		} else {
			if latestTagsResp.Success {
				c.Logger.Debug(fmt.Sprintf("raw RFRain latest tags:  %v\n", string(rawResp)))
			} else {
				c.Logger.Error(fmt.Sprintf("unsuccessful attempt to get latest tags: %s", latestTagsResp.Message))
			}
		}
	}
	r := filterTagResp(device, latestTagsResp.Results)
	return r
}

func (c *RFRainClient) GetResource(resource string, r Result) string {
	c.Logger.Debug(fmt.Sprintf("getting value for resource %v", resource))
	switch rsc:= resource; rsc {
	case "tagnumb":
		return r.Tagnumb
	case "subzone":
		return r.Subzone
	case "SS":
		return r.Ss
	case "access":
		return r.Access
	case "data":
		return r.Data
	default:
		return "unknown"
	}
}

func filterTagResp(device string, res []Result) []Result {
	resp:= []Result{}
	for _, r := range res {
		if r.Reader == device && r.Detectstat == "PRES" && r.Alarmtype == "mon" {
			resp = append(resp, r)
		}
	}
	return resp
}

func (c *RFRainClient) InvalidateSessionKey() error {
	jsonData := map[string]string{"sessionkey": c.SessionKey}
	jsonValue, _ := json.Marshal(jsonData)
	_, err := http.Post(c.InvalidateURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		c.Logger.Error(fmt.Sprintf("problem invalidating the session: %s", err))
		return err
	} else {
		c.Logger.Info("Ending service communications.  RFRain session invalidated")
	}
	return nil
}
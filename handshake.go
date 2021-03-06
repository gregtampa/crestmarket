//   Copyright 2014 StackFoundry LLC
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package crestmarket

import (
	"encoding/json"
	"github.com/theatrus/mediate"
	"github.com/theatrus/ooauth2"
	"io/ioutil"
	"log"
	"net/http"
)

type OAuthSettings struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Callback     string `json:"callback"`
}

func LoadSettings(filename string) (*OAuthSettings, error) {
	var settings OAuthSettings
	settingsData, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("Can't load secret key file - aborting", err)
		return nil, err
	}
	json.Unmarshal(settingsData, &settings)
	return &settings, nil
}

func NewOAuthOptions(settings *OAuthSettings) (*ooauth2.Options, error) {
	var endpoint ooauth2.Option
	if isSisi {
		endpoint = ooauth2.Endpoint(
			"https://sisilogin.testeveonline.com/oauth/authorize",
			"https://sisilogin.testeveonline.com/oauth/token",
		)
	} else {
		endpoint = ooauth2.Endpoint(
			"https://login.eveonline.com/oauth/authorize",
			"https://login.eveonline.com/oauth/token",
		)
	}

	httpClient := &http.Client{}
	httpClient.Transport = mediate.FixedRetries(3,
		mediate.ReliableBody(http.DefaultTransport),
	)

	return ooauth2.New(
		ooauth2.Client(settings.ClientId, settings.ClientSecret),
		ooauth2.RedirectURL(settings.Callback),
		ooauth2.Scope("publicData"),
		ooauth2.HTTPClient(httpClient),
		endpoint,
	)
}

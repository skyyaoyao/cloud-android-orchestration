// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/cloud-android-orchestration/pkg/app"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

type ClientSecrets struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type SecretManager struct {
	secrets ClientSecrets
}

func NewSecretManager(config *app.GCPSMConfig) (*SecretManager, error) {
	ctx := context.TODO()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to create secret manager client: %w", err)
	}
	defer client.Close()

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{Name: config.OAuthClientResourceID}
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return nil, fmt.Errorf("Failed to access secret: %w", err)
	}
	sm := &SecretManager{}
	err = json.Unmarshal(result.Payload.Data, &sm.secrets)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode secrets: %w", err)
	}
	return sm, nil
}

func (s *SecretManager) OAuthClientID() string {
	return s.secrets.ClientID
}

func (s *SecretManager) OAuthClientSecret() string {
	return s.secrets.ClientSecret
}

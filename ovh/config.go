package ovh

import (
	"fmt"
	"log"
	"time"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/ovh/go-ovh/ovh"
)

type Config struct {
	Endpoint          string
	ApplicationKey    string
	ApplicationSecret string
	ConsumerKey       string
	OVHClient         *ovh.Client
}

type OvhAuthCurrentCredential struct {
	OvhSupport    bool             `json:"ovhSupport"`
	Status        string           `json:"status"`
	ApplicationId int64            `json:"applicationId"`
	CredentialId  int64            `json:"credentialId"`
	Rules         []ovh.AccessRule `json:"rules"`
	Expiration    time.Time        `json:"expiration"`
	LastUse       time.Time        `json:"lastUse"`
	Creation      time.Time        `json:"creation"`
}

func clientDefault(c *Config) (*ovh.Client, error) {
	client, err := ovh.NewClient(
		c.Endpoint,
		c.ApplicationKey,
		c.ApplicationSecret,
		c.ConsumerKey,
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Config) loadAndValidate() error {
	validEndpoint := false

	ovhEndpoints := [3]string{ovh.OvhEU, ovh.OvhCA, ovh.OvhUS}

	for _, e := range ovhEndpoints {
		if ovh.Endpoints[c.Endpoint] == e {
			validEndpoint = true
		}
	}

	if !validEndpoint {
		return fmt.Errorf("%s must be one of %#v endpoints\n", c.Endpoint, ovh.Endpoints)
	}

	targetClient, err := clientDefault(c)
	if err != nil {
		return fmt.Errorf("Error getting ovh client: %q\n", err)
	}

	// decorating the OVH http client with logs
	httpClient := targetClient.Client
	if targetClient.Client.Transport == nil {
		targetClient.Client.Transport = cleanhttp.DefaultTransport()
	}

	httpClient.Transport = logging.NewTransport("OVH", httpClient.Transport)

	var cred OvhAuthCurrentCredential
	err = targetClient.Get("/auth/currentCredential", &cred)
	if err != nil {
		return fmt.Errorf("OVH client seems to be misconfigured: %q\n", err)
	}

	log.Printf("[DEBUG] Logged in on OVH API")
	c.OVHClient = targetClient

	return nil
}

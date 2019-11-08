package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/matryer/try"
)

const MaximumRetryWaitTimeInSeconds = 15 * time.Minute
const RetryWaitTimeInSeconds = 30 * time.Second

type Credentials struct {
	APIKey []byte
}

type Cloudability struct {
	RestClient   *resty.Client
	Credentials  Credentials
	RetryMaximum int
}

type getExternalAccountAws struct {
	Result CloudabilityAccount `json:"result"`
}

type CloudabilityAccount struct {
	ID              string `json:"id"`
	Name            string `json:"vendorAccountName"`
	AccountID       string `json:"vendorAccountId"`
	ParentAccountID string `json:"parentAccountId"`
	VendorKey       string `json:"vendorKey"`
	Verification    struct {
		State                       string `json:"state"`
		LastVerificationAttemptedAt string `json:"lastVerificationAttemptedAt"`
	} `json:"verification"`
	Authorization struct {
		Type       string `json:"type"`
		RoleName   string `json:"roleName"`
		ExternalID string `json:"externalId"`
	} `json:"authorization"`
}

func (cloudability *Cloudability) SetRestClient(rest *resty.Client) {
	rest.SetHostURL("https://api.cloudability.com")

	// Retry
	rest.SetRetryCount(cloudability.RetryMaximum)
	rest.SetRetryWaitTime(RetryWaitTimeInSeconds)
	rest.SetRetryMaxWaitTime(MaximumRetryWaitTimeInSeconds)
	rest.AddRetryCondition(func(r *resty.Response, err error) bool {
		switch code := r.StatusCode(); code {
		case http.StatusTooManyRequests:
			return true
		default:
			return false
		}
	})

	// Error handling
	rest.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		status := r.StatusCode()
		if status == http.StatusNotFound {
			return nil
		}

		if (status < 200) || (status >= 400) {
			return fmt.Errorf("Response not successful: Received status code %d.", status)
		}

		return nil
	})

	//Authentication
	rest.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		sign, _ := NewHTTPSignature(cloudability.Credentials.APIKey)
		for name, value := range sign {
			r.SetHeader(name, value.(string))
		}
		return nil
	})

	cloudability.RestClient = rest
}

func (cloudability *Cloudability) GetRestClient() *resty.Client {
	if cloudability.RestClient == nil {
		rest := resty.New()
		cloudability.SetRestClient(rest)
	}
	return cloudability.RestClient
}

func (cloudability *Cloudability) Poll(id string, parentId string) (*CloudabilityAccount, error) {
	result, err := cloudability.Get(id)
	if err == nil {
		return result, nil
	}

	// ensure that this is successful
	_, err = cloudability.Verification(parentId)
	if err != nil {
		return result, err
	}

	err = try.Do(func(ampt int) (bool, error) {
		var err error
		result, err = cloudability.Get(id)
		if err != nil {
			log.Printf("[DEBUG] retrying request: (Attempt: %d/%d, URL: %q)", ampt, cloudability.RetryMaximum, err)
			time.Sleep(RetryWaitTimeInSeconds)
		}
		return ampt < cloudability.RetryMaximum, err
	})
	if err != nil {
		return result, err
	}

	return result, nil
}

func (cloudability *Cloudability) Get(account string) (*CloudabilityAccount, error) {
	restClient := cloudability.GetRestClient()

	url := fmt.Sprintf("/v3/vendors/AWS/accounts/%s", account)
	req := restClient.R().SetBody("").SetResult(&getExternalAccountAws{})

	resp, err := req.Get(url)
	if err != nil {
		return nil, err
	}

	status := resp.StatusCode()
	if status == http.StatusNotFound {
		return nil, nil
	}

	response := resp.Result().(*getExternalAccountAws)
	if response == nil {
		return nil, nil
	}

	return &response.Result, nil
}

func (cloudability *Cloudability) Delete(id string) (bool, error) {
	restClient := cloudability.GetRestClient()

	url := fmt.Sprintf("/v3/vendors/AWS/accounts/%s", id)
	req := restClient.R().SetBody("")

	_, err := req.Delete(url)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (cloudability *Cloudability) Add(id string) (*CloudabilityAccount, error) {
	var result *CloudabilityAccount
	restClient := cloudability.GetRestClient()

	payload := map[string]interface{}{
		"type":            "aws_role",
		"vendorAccountId": id,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return result, err
	}

	url := "/v3/vendors/AWS/accounts/"
	req := restClient.R().SetBody(string(body)).SetResult(&getExternalAccountAws{})

	resp, err := req.Post(url)
	if err != nil {
		return result, err
	}

	response := resp.Result().(*getExternalAccountAws)
	return &response.Result, nil
}

func (cloudability *Cloudability) Verification(id string) (*CloudabilityAccount, error) {
	var result *CloudabilityAccount
	restClient := cloudability.GetRestClient()

	url := fmt.Sprintf("/v3/vendors/AWS/accounts/%s/verification", id)
	req := restClient.R().SetBody("").SetResult(&getExternalAccountAws{})

	resp, err := req.Post(url)
	if err != nil {
		return result, err
	}

	response := resp.Result().(*getExternalAccountAws)
	return &response.Result, nil
}

func (cloudability *Cloudability) Verify(id string) (*CloudabilityAccount, error) {
	err := try.Do(func(ampt int) (bool, error) {
		result, err := cloudability.Verification(id)
		if err != nil {
			log.Printf("[DEBUG] retrying request: (Attempt: %d/%d)", ampt, cloudability.RetryMaximum, err)
			time.Sleep(RetryWaitTimeInSeconds)
		} else if result.Verification.State != "verified" {
			log.Printf("[DEBUG] retrying request: (Attempt: %d/%d)", ampt, cloudability.RetryMaximum, err)
			err = fmt.Errorf("Verification was not successful: [%s]", result.Verification.State)
			time.Sleep(RetryWaitTimeInSeconds)
		}

		return ampt < cloudability.RetryMaximum, err
	})

	var result *CloudabilityAccount
	if err != nil {
		log.Printf("[DEBUG] Could not verify the account: %q", err)
		return result, err
	}

	result, err = cloudability.Verification(id)
	return result, err
}

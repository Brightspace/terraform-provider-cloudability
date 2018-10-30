package cloudability

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/matryer/try"
)

type Cloudability struct {
	Credentials  Credentials
	RetryMaximum int
}

type Credentials struct {
	APIKey []byte
}

type CloudabilityRequest struct {
	Method   string
	URL      string
	Contents []byte
}

type CloudabilityResponse struct {
	Result json.RawMessage `json:result`
}

type CloudabilityAccount struct {
	ID              string                           `json:"id"`
	Name            string                           `json:"vendorAccountName"`
	AccountID       string                           `json:"vendorAccountId"`
	ParentAccountID string                           `json:"parentAccountId"`
	VendorKey       string                           `json:"vendorKey"`
	Verification    CloudabilityAccountVerification  `json:"verification"`
	Authorization   CloudabilityAccountAuthorization `json:"authorization"`
}

type CloudabilityAccountVerification struct {
	State                       string `json:"state"`
	LastVerificationAttemptedAt string `json:"lastVerificationAttemptedAt"`
}

type CloudabilityAccountAuthorization struct {
	Type       string `json:"type"`
	RoleName   string `json:"roleName"`
	ExternalID string `json:"externalId"`
}

func makeHeaders(req CloudabilityRequest, creds Credentials) (map[string]interface{}, error) {
	headers := make(map[string]interface{})

	ctype := "application/json"
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(creds.APIKey))

	headers["Accept"] = ctype
	headers["Content-Type"] = ctype
	headers["Authorization"] = fmt.Sprintf("Basic %s", encodedAuth)

	return headers, nil
}

func makeRequest(request CloudabilityRequest, creds Credentials) (string, error) {
	baseURL := "https://api.cloudability.com"
	client := &http.Client{}

	req, err := http.NewRequest(request.Method, baseURL+request.URL, bytes.NewBuffer(request.Contents))
	if err != nil {
		return "", fmt.Errorf("Error creating request: %s", err)
	}

	headers, _ := makeHeaders(request, creds)
	for name, value := range headers {
		req.Header.Set(name, value.(string))
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error during making a request: %s", request.URL)

	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return "", fmt.Errorf("HTTP request error. Response code: %d", resp.StatusCode)

	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "application/vnd.api+json" {
		return "", fmt.Errorf("Content-Type is not a json type. Got: %s", contentType)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error while reading response body. %s", err)
	}

	return string(bytes), nil
}

func (cloudability Cloudability) add(accountID string) (CloudabilityAccount, error) {
	var response CloudabilityResponse
	var result CloudabilityAccount
	var resp string
	var err error

	var payload struct {
		Type            string `json:"type"`
		VendorAccountID string `json:"vendorAccountId"`
	}

	payload.Type = "aws_role"
	payload.VendorAccountID = accountID

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return result, err
	}

	request := CloudabilityRequest{
		Method:   "POST",
		URL:      "/v3/vendors/AWS/accounts/",
		Contents: payloadJSON,
	}

	err = try.Do(func(ampt int) (bool, error) {
		var err error
		resp, err = makeRequest(request, cloudability.Credentials)
		if err != nil {
			log.Printf("[DEBUG] retrying request: (Attempt: %d/%d, URL: %q)", ampt, cloudability.RetryMaximum, err)
			time.Sleep(30 * time.Second)
		}
		return ampt < cloudability.RetryMaximum, err
	})
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(response.Result), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (cloudability Cloudability) verify(account string) (CloudabilityAccount, error) {
	var response CloudabilityResponse
	var result CloudabilityAccount
	var resp string
	var err error

	request := CloudabilityRequest{
		Method:   "POST",
		URL:      "/v3/vendors/AWS/accounts/" + account + "/verification",
		Contents: []byte(""),
	}

	err = try.Do(func(ampt int) (bool, error) {
		var err error
		resp, err = makeRequest(request, cloudability.Credentials)
		if err != nil {
			log.Printf("[DEBUG] retrying request: (Attempt: %d/%d, URL: %q)", ampt, cloudability.RetryMaximum, err)
			time.Sleep(30 * time.Second)
		}
		return ampt < cloudability.RetryMaximum, err
	})
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(response.Result), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (cloudability Cloudability) get(account string) (CloudabilityAccount, error) {
	var response CloudabilityResponse
	var result CloudabilityAccount
	var resp string
	var err error

	request := CloudabilityRequest{
		Method:   "GET",
		URL:      "/v3/vendors/AWS/accounts/" + account,
		Contents: []byte(""),
	}

	err = try.Do(func(ampt int) (bool, error) {
		var err error
		resp, err = makeRequest(request, cloudability.Credentials)
		if err != nil {
			log.Printf("[DEBUG] retrying request: (Attempt: %d/%d, URL: %q)", ampt, cloudability.RetryMaximum, err)
			time.Sleep(30 * time.Second)
		}
		return ampt < cloudability.RetryMaximum, err
	})
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(response.Result), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (cloudability Cloudability) pull(payerAccountID string, accountID string) (CloudabilityAccount, error) {
	var result CloudabilityAccount
	var err error

	account, err := cloudability.get(accountID)
	if err == nil {
		return account, nil
	}

	_, err = cloudability.verify(string(payerAccountID))
	if err != nil {
		return result, err
	}

	err = try.Do(func(ampt int) (bool, error) {
		var err error
		result, err = cloudability.get(accountID)
		if err != nil {
			log.Printf("[DEBUG] retrying request: (Attempt: %d/%d, URL: %q)", ampt, cloudability.RetryMaximum, err)
			time.Sleep(30 * time.Second)
		}
		return ampt < cloudability.RetryMaximum, err
	})
	if err != nil {
		return result, err
	}

	return result, nil
}

func (cloudability Cloudability) delete(account string) (bool, error) {
	var err error

	request := CloudabilityRequest{
		Method:   "DELETE",
		URL:      "/v3/vendors/AWS/accounts/" + account,
		Contents: []byte(""),
	}

	err = try.Do(func(ampt int) (bool, error) {
		var err error
		_, err = makeRequest(request, cloudability.Credentials)
		if err != nil {
			log.Printf("[DEBUG] retrying request: (Attempt: %d/%d, URL: %q)", ampt, cloudability.RetryMaximum, err)
			time.Sleep(30 * time.Second)
		}
		return ampt < cloudability.RetryMaximum, err
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

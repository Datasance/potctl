/*
 *  *******************************************************************************
 *  * Copyright (c) 2023 Datasance Teknoloji A.S.
 *  *
 *  * This program and the accompanying materials are made available under the
 *  * terms of the Eclipse Public License v. 2.0 which is available at
 *  * http://www.eclipse.org/legal/epl-2.0
 *  *
 *  * SPDX-License-Identifier: EPL-2.0
 *  *******************************************************************************
 *
 */

package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"strconv"
)

type EntitlementResponse struct {
	CustomerName         string           `json:"customerName"`
	CustomerAccountRefID string           `json:"customerAccountRefId"`
	OrderRefID           string           `json:"orderRefId"`
	OfferingName         string           `json:"offeringName"`
	SKU                  string           `json:"sku"`
	ProductName          string           `json:"productName"`
	Plan                 EntitlementPlan  `json:"plan"`
	GracePeriod          EntitlementPeriod `json:"gracePeriod"`
	LingerPeriod         EntitlementPeriod `json:"lingerPeriod"`
	LeasePeriod          EntitlementPeriod `json:"leasePeriod"`
	OfflineLeasePeriod   EntitlementPeriod `json:"offlineLeasePeriod"`
}

type EntitlementPlan struct {
	Name            string              `json:"name"`
	LicenseType     string              `json:"licenseType"`
	LicenseStartType string              `json:"licenseStartType"`
	LicenseDuration EntitlementDuration `json:"licenseDuration"`
}

type EntitlementDuration struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type EntitlementPeriod struct {
	Type  string `json:"type"`
	Count *int   `json:"count"`
}

type ActivationResponse struct {
	AccessToken           string                `json:"accessToken"`
	ID                    string                `json:"id"`
	LeaseExpiry           string                `json:"leaseExpiry"`
	ProductID             string                `json:"productId"`
	EntitlementID         string                `json:"entitlementId"`
	EntitlementExpiryDate string                `json:"entitlementExpiryDate"`
	SeatID                string                `json:"seatId"`
	SeatName              string                `json:"seatName"`
	Activated             string                `json:"activated"`
	LastLease             string                `json:"lastLease"`
	Attributes            []ActivationAttribute `json:"attributes"`
	Features              []string              `json:"features"`
}

type ActivationAttribute struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func ActivateAndGetAccessToken(productID, activationCode, seatID, seatName string) (string, string, error) {
	url := "https://datasance.license.zentitle.io/api/v1/activate"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{
	"productId": "%s",
	"activationCode": "%s",
	"seatId": "%s",
	"seatName": "%s"
  }`, productID, activationCode, seatID, seatName))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return "", "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("N-TenantId", "t_h42rcI0Lq0_yIAGfLUf3Xg")

	res, err := client.Do(req)

	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	var activateResponse ActivationResponse
	if err := json.NewDecoder(res.Body).Decode(&activateResponse); err != nil {
		return "", "", fmt.Errorf("error decoding JSON response: %v", err)
	}

	return activateResponse.AccessToken, res.Header.Get("N-Nonce"), nil
}

func GetEntitlementDetails(accessToken, nonce string) (*EntitlementResponse, string, error) {
	url := "https://datasance.license.zentitle.io/api/v1/activation/entitlement"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("N-Nonce", nonce)
	req.Header.Set("N-TenantId", "t_h42rcI0Lq0_yIAGfLUf3Xg")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var entitlementResponse EntitlementResponse
	if err := json.NewDecoder(resp.Body).Decode(&entitlementResponse); err != nil {
		return nil, "", fmt.Errorf("error decoding JSON response: %v", err)
	}

	return &entitlementResponse, resp.Header.Get("N-Nonce"), nil
}

func ActivateLicense(accessToken, nonce string) (*ActivationResponse, error) {
	url := "https://datasance.license.zentitle.io/api/v1/activation"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("N-TenantId", "t_h42rcI0Lq0_yIAGfLUf3Xg")
	req.Header.Set("N-Nonce", nonce)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	var activationResponse ActivationResponse
	if err := json.NewDecoder(resp.Body).Decode(&activationResponse); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %v", err)
	}

	return &activationResponse, nil
}

func GetEntitlementDatasance(activationCode string, seatID string, seatName string) (string, string, error) {
	productID := "prod_vKqv2P1OiUKBNZqa76a7iw"
	// Sample variables for activation
	//activationCode := "XXXX-XXXX-XXXX-XXXX"
	//seatID := "foo.bar@datasance.com"
	//seatName := "Foo Bar"

	accessToken, nonceActivation, err := ActivateAndGetAccessToken(productID, activationCode, seatID, seatName)
	if err != nil {
		fmt.Println("Error activating:", err)
		return "", "", err
	}

	entitlementDetails, nonceEntitlement, err := GetEntitlementDetails(accessToken, nonceActivation)
	if err != nil {
		fmt.Println("Error getting entitlement details:", err)
		return "", "", err
	}

	_ = entitlementDetails

	activationResponse, err := ActivateLicense(accessToken, nonceEntitlement)
	if err != nil {
		fmt.Println("Error activating license:", err)
		return "", "", err
	}
	//fmt.Println("Expiry Date:", activationResponse.EntitlementExpiryDate)
	var expiryDate = activationResponse.EntitlementExpiryDate
	var agentSeats = ""
	for _, activationAttributeObject := range activationResponse.Attributes {
		if activationAttributeObject.Key == "Agent Seats" {
			//fmt.Println("Number of agents:", activationAttributeObject.Value)
			agentSeats = activationAttributeObject.Value
		}
	}
	return expiryDate, agentSeats, nil
}

func checkExpiryDate(dateString string) bool {

	dateExpirytime, err:= time.Parse(time.RFC3339Nano, dateString)

	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	currentTime := time.Now()

	if currentTime.After(dateExpirytime) {
		fmt.Println("Your subscription has been expired, please contact with Datasance Sales Team or Datasance Partner")
	}

	return currentTime.After(dateExpirytime)
}

func checkNumOfAgentSeats(currentAgentNum int, maxAgentNum string) bool {

	maxAgentNumAsInt, err := strconv.Atoi(maxAgentNum)

    if err != nil {
        fmt.Println("Error converting maximum agent number to integer:", err)
        return false
    }

	if currentAgentNum >= maxAgentNumAsInt {
		fmt.Println("You don't have enough subscription to deploy additional agents on this controlplane")
		fmt.Println("Your active subscription includes maximum agent seats as ", maxAgentNum)
		fmt.Println("Please contact with Datasance Sales Team or Datasance Partner")
	}

	return currentAgentNum <= maxAgentNumAsInt
}

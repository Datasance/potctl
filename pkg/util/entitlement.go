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

var productID = "prod_vKqv2P1OiUKBNZqa76a7iw"

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

	if res.StatusCode == 402 {
		return "", "", fmt.Errorf("Subscription has expired")
	}

	if res.StatusCode == 404 {
		return "", "", fmt.Errorf("Subscription not found")
	}

	if res.StatusCode == 500 {
		return "", "", fmt.Errorf("Subscription Engine is not responding")
	}

	var activateResponse ActivationResponse
	if err := json.NewDecoder(res.Body).Decode(&activateResponse); err != nil {
		return "", "", fmt.Errorf("Error while decoding JSON response from get access token event: %v", err)
	}

	return activateResponse.AccessToken, res.Header.Get("N-Nonce"), nil
}

func ActivateLicense(accessToken, nonce string) (*ActivationResponse, error) {
	url := "https://datasance.license.zentitle.io/api/v1/activation"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error while creating activation request: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("N-TenantId", "t_h42rcI0Lq0_yIAGfLUf3Xg")
	req.Header.Set("N-Nonce", nonce)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error while sending activation request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("Status Code for Activation is ", resp.StatusCode)

	var activationResponse ActivationResponse
	if err := json.NewDecoder(resp.Body).Decode(&activationResponse); err != nil {
		return nil, fmt.Errorf("Error while decoding JSON response from activation event is: %v", err)
	}

	return &activationResponse, nil
}

func DeactivateLicense(accessToken, nonce string) {

	url := "https://datasance.license.zentitle.io/api/v1/activation"

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println("Error while creating deactivation request: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("N-TenantId", "t_h42rcI0Lq0_yIAGfLUf3Xg")
	req.Header.Set("N-Nonce", nonce)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error while sending deactivation request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("Status code for deactivation is ", resp.StatusCode)
}

func GetEntitlementDatasance(activationCode string, seatID string, seatName string) (string, string, error) {

	accessToken, nonceActivation, err := ActivateAndGetAccessToken(productID, activationCode, seatID, seatName)
	if err != nil {
        switch {
		case strings.Contains(err.Error(), "Subscription has expired"):
			return "Subscription has expired", "0", nil
		case strings.Contains(err.Error(), "Subscription not found"):
			return "Subscription not found", "0", nil
		case strings.Contains(err.Error(), "Subscription Engine is not responding"):
			return "Subscription Engine is not responding", "0", nil
		default:
			fmt.Println("Error activating license:", err)
			return "", "", err
		}
	}

	activationResponse, err := ActivateLicense(accessToken, nonceActivation)
	if err != nil {
		fmt.Println("Error while activating license:", err)
        return "", "", err
	}

    var expiryDate = activationResponse.EntitlementExpiryDate
	var agentSeats = ""
	for _, activationAttributeObject := range activationResponse.Attributes {
		if activationAttributeObject.Key == "Agent Seats" {
			agentSeats = activationAttributeObject.Value
		}
	}
	return expiryDate, agentSeats, nil
}

func DeactivateEntitlementDatasance(activationCode string, seatID string, seatName string) {
	accessToken, nonceDeactivation, err := ActivateAndGetAccessToken(productID, activationCode, seatID, seatName)
	if err != nil {
        switch {
		case strings.Contains(err.Error(), "Subscription has expired"):
			fmt.Println("Subscription has expired")
		case strings.Contains(err.Error(), "Subscription not found"):
			fmt.Println("Subscription not found")
		case strings.Contains(err.Error(), "Subscription Engine is not responding"):
			fmt.Println("Subscription Engine is not responding")
		default:
			fmt.Println("Error deactivating license:", err)
		}
	}

	DeactivateLicense(accessToken, nonceDeactivation)
}



func CheckExpiryDate(dateString string) (bool) {

	fmt.Println("Checking License Expiry Date from Subscription")
	if strings.Contains(dateString,"Subscription has expired") {
		fmt.Println("Subscription has expired, Please contact with Datasance Sales Team : sales@datasance.com or Datasance Reseller Partner")
		return false
	}
	if strings.Contains(dateString,"Subscription not found") {
		fmt.Println("Subscription not found, Please contact with Datasance Sales Team : sales@datasance.com or Datasance Reseller Partner")
		return false
	}
	if strings.Contains(dateString,"Subscription Engine is not responding") {
		fmt.Println("Subscription Engine is not responding, Please contact with Datasance Support Team : support@datasance.com")
		return false
	}

	return true
}

func CheckNumOfAgentSeats(currentAgentNum int, maxAgentNum string) (bool) {

	maxAgentNumAsInt, err := strconv.Atoi(maxAgentNum)

	fmt.Println("Checking number of agents from subscription details")

    if err != nil {
        fmt.Println("Error while converting maximum agent number to integer:", err)
        return false
    }

	if currentAgentNum > maxAgentNumAsInt {
		fmt.Println("You don't have enough subscription to deploy additional agents on this controlplane")
		fmt.Println("Your active subscription includes maximum agent seats as ", maxAgentNum)
		fmt.Println("Please contact with Datasance Sales Team : sales@datasance.com or Datasance Reseller Partner")
	}

	return currentAgentNum < maxAgentNumAsInt
}

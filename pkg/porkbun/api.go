package porkbun

import (
	"bacon/pkg/client"
	"bacon/pkg/helpers"
	"encoding/json"
	"fmt"
)

// https://porkbun.com/api/json/v3/documentation
const (
	PING   = "https://porkbun.com/api/json/v3/ping"
	CREATE = "https://porkbun.com/api/json/v3/dns/create"
	DELETE = "https://porkbun.com/api/json/v3/dns/delete"
)

type checkable interface {
	checkStatus() bool
	messageAsError() error
}

type baseRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (res *baseRes) checkStatus() bool {
	switch res.Status {
	case "SUCCESS":
		return true
	default:
		return false
	}
}

func (res *baseRes) messageAsError() error {
	return fmt.Errorf("%s", res.Message)
}

func unmarshalAndCheckStatus(data []byte, body checkable) error {
	err := json.Unmarshal(data, &body)
	if err != nil {
		return err
	}

	if !body.checkStatus() {
		return body.messageAsError()
	}

	return nil
}

func ping(auth PorkAuth) error {
	type pingRes struct {
		baseRes
		YourIp string `json:"yourIp"`
	}

	body, err := helpers.PostJsonAndRead(PING, auth)
	if err != nil {
		return err
	}

	ping := pingRes{}
	err = unmarshalAndCheckStatus(body, &ping)
	if err != nil {
		return err
	}

	return nil
}

func create(auth PorkAuth, domain string, record client.Record) (string, error) {
	type createBody struct {
		PorkAuth
		PorkbunRecord
	}

	type createRes struct {
		baseRes
		Id string `json:"id"`
	}

	porkRecord := ConvertToPorkbunRecord(record)

	toCreate := createBody{
		PorkAuth:      auth,
		PorkbunRecord: porkRecord,
	}

	body, err := helpers.PostJsonAndRead(CREATE+"/"+domain, toCreate)
	if err != nil {
		return "", err
	}

	created := createRes{}
	err = unmarshalAndCheckStatus(body, &created)
	if err != nil {
		return "", err
	}

	return created.Id, nil
}

func delete(auth PorkAuth, domain string, id string) error {
	body, err := helpers.PostJsonAndRead(DELETE+"/"+domain+"/"+id, auth)
	if err != nil {
		return err
	}

	deleted := baseRes{}
	err = unmarshalAndCheckStatus(body, &deleted)
	if err != nil {
		return err
	}

	return nil
}

func deploy(auth PorkAuth, domain string, records []client.Record, shouldCreate bool, shouldDelete bool) error {
	return fmt.Errorf("haven't implemented sync yet")
}

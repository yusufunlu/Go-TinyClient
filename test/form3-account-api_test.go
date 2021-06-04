package interview_accountapi

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	tiny "github.com/yusufunlu/tinyclient"
	"io/ioutil"
	"os"
	"testing"
)

const (
	baseUrl     = "accountapi:8080/v1"
	accountPath = "/organisation/accounts"
)

type Account struct {
	Data struct {
		Type           string `json:"type,omitempty"`
		ID             string `json:"id,omitempty"`
		OrganisationID string `json:"organisation_id,omitempty"`
		Version        int    `json:"version,omitempty"`
		Attributes     struct {
			Country                 string   `json:"country,omitempty"`
			BaseCurrency            string   `json:"base_currency,omitempty"`
			AccountNumber           string   `json:"account_number,omitempty"`
			BankID                  string   `json:"bank_id,omitempty"`
			BankIDCode              string   `json:"bank_id_code,omitempty"`
			Bic                     string   `json:"bic,omitempty"`
			Iban                    string   `json:"iban,omitempty"`
			Name                    []string `json:"name,omitempty"`
			AlternativeNames        []string `json:"alternative_names,omitempty"`
			AccountClassification   string   `json:"account_classification,omitempty"`
			JointAccount            bool     `json:"joint_account,omitempty"`
			AccountMatchingOptOut   bool     `json:"account_matching_opt_out,omitempty"`
			SecondaryIdentification string   `json:"secondary_identification,omitempty"`
			Switched                bool     `json:"switched,omitempty"`
			PrivateIdentification   *struct {
				BirthDate      string   `json:"birth_date,omitempty"`
				BirthCountry   string   `json:"birth_country,omitempty"`
				Identification string   `json:"identification,omitempty"`
				Address        []string `json:"address,omitempty,omitempty"`
				City           string   `json:"city,omitempty"`
				Country        string   `json:"country,omitempty"`
			} `json:"private_identification,omitempty"`
			OrganisationIdentification *struct {
				Identification string `json:"identification,omitempty"`
				Actors         []struct {
					Name      []string `json:"name,omitempty"`
					BirthDate string   `json:"birth_date,omitempty"`
					Residency string   `json:"residency,omitempty"`
				} `json:"actors,omitempty,omitempty"`
				Address []string `json:"addres,omitempty"`
				City    string   `json:"city,omitempty"`
				Country string   `json:"country,omitempty"`
			} `json:"organisation_identification,omitempty"`
			Status string `json:"status,omitempty"`
		} `json:"attributes,omitempty"`
		Relationships *struct {
			MasterAccount *struct {
				Data []struct {
					Type string `json:"type,omitempty"`
					ID   string `json:"id,omitempty"`
				} `json:"data,omitempty"`
			} `json:"master_account,omitempty"`
			AccountEvents *struct {
				Data []struct {
					Type string `json:"type,omitempty"`
					ID   string `json:"id,omitempty"`
				} `json:"data,omitempty"`
			} `json:"account_events,omitempty"`
		} `json:"relationships,omitempty"`
	} `json:"data"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
}

func readTestData(v interface{}) {

	jsonFile, err := os.Open("./testdata/account-post-data.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	accountBytes, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(accountBytes, v)
}

func TestPostSuccess(t *testing.T) {
	account := Account{}
	readTestData(&account)

	url := fmt.Sprintf("%s%s", baseUrl, accountPath)
	client := tiny.NewClient()

	request := client.NewRequest().SetURL(url).SetBody(account).SetMethod(tiny.Post).
		SetContentType("application/json; charset=utf-8")

	response, err := client.Send(request)

	require.NoError(t, err)
	require.Equal(t, 201, response.Response.StatusCode)

	respAccount := Account{}
	err = response.BodyUnmarshall(&respAccount)

	require.NoError(t, err)

}

func TestFetchSuccess(t *testing.T) {
	account := Account{}
	readTestData(&account)
	id := account.Data.ID

	url := fmt.Sprintf("%s%s/%s", baseUrl, accountPath, id)
	client := tiny.NewClient()
	request := client.NewRequest().SetURL(url).SetMethod(tiny.Get).
		SetContentType("application/json; charset=utf-8")

	response, err := client.Send(request)

	require.NoError(t, err)
	require.Equal(t, 200, response.Response.StatusCode)

	respAccount := Account{}
	err = response.BodyUnmarshall(&respAccount)

	require.NoError(t, err)
	require.Equal(t, respAccount.Data.ID, id)

}

func TestFetchFail(t *testing.T) {

	id := uuid.NewV4().String()

	url := fmt.Sprintf("%s%s/%s", baseUrl, accountPath, id)
	client := tiny.NewClient().SetDebugMode(true)

	request := client.NewRequest().SetURL(url).SetMethod(tiny.Get).
		SetContentType("application/json; charset=utf-8")

	response, err := client.Send(request)

	require.NoError(t, err)
	require.Equal(t, 404, response.Response.StatusCode)

	errRes := ErrorResponse{}
	err = response.BodyUnmarshall(&errRes)
	require.NoError(t, err)
	require.True(t, len(errRes.ErrorMessage) > 0)

}

func TestDeleteSuccess(t *testing.T) {
	account := Account{}
	readTestData(&account)
	id := account.Data.ID
	version := account.Data.Version

	url := fmt.Sprintf("%s%s/%s", baseUrl, accountPath, id)
	client := tiny.NewClient()

	request := client.NewRequest().SetURL(url).SetMethod(tiny.Delete).
		SetContentType("application/json; charset=utf-8").
		AddQueryParam("version", fmt.Sprintf("%d", version))

	response, err := client.Send(request)

	require.NoError(t, err)
	require.Equal(t, 204, response.Response.StatusCode)
}

func TestDeleteNotExistFail(t *testing.T) {

	version := 0
	id := uuid.NewV4().String()

	url := fmt.Sprintf("%s%s/%s", baseUrl, accountPath, id)
	client := tiny.NewClient().SetDebugMode(true)

	request := client.NewRequest().SetURL(url).SetMethod(tiny.Delete).
		SetContentType("application/json; charset=utf-8").
		AddQueryParam("version", fmt.Sprintf("%d", version))

	response, err := client.Send(request)

	require.NoError(t, err)
	require.Equal(t, 404, response.Response.StatusCode)

}

func TestDeleteVersionIncorrectFail(t *testing.T) {

	account := Account{}
	readTestData(&account)
	id := account.Data.ID
	version := account.Data.Version + 1

	url := fmt.Sprintf("%s%s/%s", baseUrl, accountPath, id)
	client := tiny.NewClient()

	request := client.NewRequest().SetURL(url).SetMethod(tiny.Delete).
		SetContentType("application/json; charset=utf-8").
		AddQueryParam("version", fmt.Sprintf("%d", version))

	response, err := client.Send(request)

	require.NoError(t, err)
	require.Equal(t, 404, response.Response.StatusCode)

}

func TestPostFail(t *testing.T) {
	account := Account{}
	readTestData(&account)

	url := fmt.Sprintf("%s%s", baseUrl, accountPath)
	client := tiny.NewClient()

	request := client.NewRequest().SetURL(url).SetBody(account).SetMethod(tiny.Post).
		SetContentType("application/json; charset=utf-8")

	response, err := client.Send(request)

	require.NoError(t, err)
	require.Equal(t, 201, response.Response.StatusCode)

	response, err = client.Send(request)
	require.Equal(t, 409, response.Response.StatusCode)

	errRes := ErrorResponse{}
	err = response.BodyUnmarshall(&errRes)
	require.NoError(t, err)
	require.True(t, len(errRes.ErrorMessage) > 0)

	//cleaning the data of context
	TestDeleteSuccess(t)
}

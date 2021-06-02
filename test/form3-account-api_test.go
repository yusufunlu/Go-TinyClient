package interview_accountapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	tiny "tinyclient"
)

const (
	baseUrl     = "localhost:8080/v1"
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

type errorResponse struct {
	error_message string `json:"error_message"`
}

func TestPostSuccess(t *testing.T) {

	account := Account{}
	jsonFile, err := os.Open("./testdata/account-post-data.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	accountBytes, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(accountBytes, &account)

	id := uuid.NewV4()
	organisationId := uuid.NewV4()
	account.Data.ID = id.String()
	account.Data.OrganisationID = organisationId.String()

	url := fmt.Sprintf("%s%s", baseUrl, accountPath)
	client := tiny.NewClient().SetTimeout(30)

	request := client.NewRequest().SetURL(url).SetBody(account).SetMethod(tiny.Post).
		SetContentType("application/json; charset=utf-8")

	ctx := context.Background()
	response, err := client.Send(request, ctx)

	require.NoError(t, err)
	require.Equal(t, 201, response.Response.StatusCode)

	err = response.BodyToStruct(account)

	require.NoError(t, err)

}

func TestFetchSuccess(t *testing.T) {
	id := "07675eaf-1944-4073-8eb5-d2cef32b94df"

	url := fmt.Sprintf("%s%s/%s", baseUrl, accountPath, id)
	client := tiny.NewClient().SetTimeout(30)
	request := client.NewRequest().SetURL(url).SetMethod(tiny.Get).
		SetContentType("application/json; charset=utf-8")

	ctx := context.Background()
	response, err := client.Send(request, ctx)

	require.NoError(t, err)
	require.Equal(t, 200, response.Response.StatusCode)

	account := &Account{}
	err = response.BodyToStruct(account)

	require.NoError(t, err)
	require.Equal(t, account.Data.ID, id)

}

func TestFetchFail(t *testing.T) {
	id := "07675eaf-1944-4073-8eb5-d2cef32b94df-fail"

	url := fmt.Sprintf("%s%s/%s", baseUrl, accountPath, id)
	client := tiny.NewClient().SetTimeout(30)

	request := client.NewRequest().SetURL(url).SetMethod(tiny.Get).
		SetContentType("application/json; charset=utf-8")

	ctx := context.Background()
	response, err := client.Send(request, ctx)

	require.NoError(t, err)
	require.Equal(t, 400, response.Response.StatusCode)

	if response.Response.StatusCode < http.StatusOK || response.Response.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		err = response.BodyToStruct(&errRes)
		require.NoError(t, err)
		require.NotNil(t, errRes.error_message)
		return
	}

	account := &Account{}
	err = response.BodyToStruct(account)

	require.NoError(t, err)
	require.Equal(t, account.Data.ID, id)

}

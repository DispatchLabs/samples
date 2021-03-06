package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dispatchlabs/disgo/commons/types"
	"github.com/dispatchlabs/disgo/commons/utils"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func GetReceipt(hash string, endpoint string) *types.Receipt {
	response, err := http.Get(fmt.Sprintf("%s/%s", endpoint, hash))

	var receipt *types.Receipt
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, _ := ioutil.ReadAll(response.Body)
		receipt, err = unmarshalReceipt(contents)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		//fmt.Printf("%s\n", receipt.ToPrettyJson())
	}
	return receipt
}

func GetQueue(endpoint string) string {
	response, err := http.Get(fmt.Sprintf("%s", endpoint))
	//var gossips []types.Gossip
	var result string
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, _ := ioutil.ReadAll(response.Body)
		//receipt, err = unmarshalReceipt(contents)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		result = string(contents)
		//fmt.Printf("%s\n", receipt.ToPrettyJson())
	}
	return result

}

func PostTx(tx *types.Transaction, endpoint string) {
	fmt.Printf("Executing contract json: \n%s\n", tx.ToPrettyJson())
	fmt.Printf("Sending tx : %s with timestamp: %v\n", tx.Hash, tx.ToTime())
	data := new(bytes.Buffer)
	data.WriteString(tx.String())

	response, err := http.Post(
		endpoint,
		"application/json; charset=utf-8",
		data,
	)
	if err != nil {
		utils.Error(err)
		return
	}
	contents, _ := ioutil.ReadAll(response.Body)
	// If NOT then this happens https://stackoverflow.com/questions/37454236/net-http-server-too-many-open-files-error
	fmt.Printf("Response: %v\n", string(contents))
	response.Body.Close()
}

func unmarshalReceipt(bytes []byte) (*types.Receipt, error) {
	receipt := types.Receipt{}
	var jsonMap map[string]interface{}
	err := json.Unmarshal(bytes, &jsonMap)
	if err != nil {
		utils.Fatal(err)
	}
	if jsonMap["data"] != nil {

		value := jsonMap["data"].(map[string]interface{})

		if value["receipt"] != nil {
			value := value["receipt"].(map[string]interface{})

			if value != nil {
				if value["transactionHash"] != nil {
					receipt.TransactionHash = value["transactionHash"].(string)
				}
				if value["status"] != nil {
					receipt.Status = value["status"].(string)
				}
				if value["humanReadableStatus"] != nil {
					receipt.HumanReadableStatus = value["humanReadableStatus"].(string)
				}
				if value["contractAddress"] != nil && value["contractAddress"] != "" {
					receipt.ContractAddress = value["contractAddress"].(string)
				}
				if value["contractResult"] != nil {
					var contractResult = value["contractResult"]
					receipt.ContractResult = contractResult.([]interface{})
				}
				if value["created"] != nil {
					created, err := time.Parse(time.RFC3339, value["created"].(string))
					if err != nil {
						return nil, err
					}
					receipt.Created = created
				}
			}
		}
	}
	return &receipt, nil
}

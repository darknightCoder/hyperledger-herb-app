// SPDX-License-Identifier: Apache-2.0

/*
  Sample Chaincode based on Demonstrated Scenario

 This code is based on code written by the Hyperledger Fabric community.
  Original code can be found here: https://github.com/hyperledger/fabric-samples/blob/release/chaincode/fabcar/fabcar.go
 */

package main

/* Imports  
* 4 utility libraries for handling bytes, reading and writing JSON, 
formatting, and string manipulation  
* 2 specific Hyperledger Fabric specific libraries for Smart Contracts  
*/
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

/* Define herb structure, with 4 properties.  
Structure tags are used by encoding/json library
*/
type HerbCrops struct {
	Name string `json:"name"`
	Family string `json:"family"`
	Location  string `json:"location"`
	Production string `json:"production"`
	Holder  string `json:"holder"`
        Timestamp string `json:"timestamp"`
}

/*
 * The Init method *
 called when the Smart Contract "herb-chaincode" is instantiated by the network
 * Best practice is to have any Ledger initialization in separate function 
 -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method *
 called when an application requests to run the Smart Contract "herb-chaincode"
 The app also specifies the specific smart contract function to call with args
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger
	if function == "queryHerb" {
		return s.queryHerb(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "recordHerb" {
		return s.recordHerb(APIstub, args)
	} else if function == "queryAllHerb" {
		return s.queryAllHerbs(APIstub)
	} else if function == "changeHerbHolder" {
		return s.changeHerbHolder(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/*
 * The queryherb method *
Used to view the records of one particular herb
It takes one argument -- the key for the herb in question
 */
func (s *SmartContract) queryHerb(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	herbAsBytes, _ := APIstub.GetState(args[0])
	if herbAsBytes == nil {
		return shim.Error("Could not locate herb")
	}
	return shim.Success(herbAsBytes)
}

/*
 * The initLedger method *
Will add test data (10 herb catches)to our network
 */
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	herb := []HerbCrops{
		HerbCrops{Name: "923F", Location: "67.0006, -70.5476", Timestamp: "1504054225", Holder: "Miriam",Family:"Some",Production:"10 kgs"},
		HerbCrops{Name: "M83T", Location: "91.2395, -49.4594", Timestamp: "1504057825", Holder: "Dave",Family:"Some",Production:"10 kgs"},
		HerbCrops{Name: "T012", Location: "58.0148, 59.01391", Timestamp: "1493517025", Holder: "Igor",Family:"Some",Production:"10 kgs"},
		HerbCrops{Name: "P490", Location: "-45.0945, 0.7949", Timestamp: "1496105425", Holder: "Amalea",Family:"Some",Production:"10 kgs"},
		HerbCrops{Name: "S439", Location: "-107.6043, 19.5003", Timestamp: "1493512301", Holder: "Rafa",Family:"Some",Production:"10 kgs"},
		HerbCrops{Name: "J205", Location: "-155.2304, -15.8723", Timestamp: "1494117101", Holder: "Shen",Family:"Some",Production:"10 kgs"},
		HerbCrops{Name: "S22L", Location: "103.8842, 22.1277", Timestamp: "1496104301", Holder: "Leila",Family:"Some",Production:"10 kgs"},
		HerbCrops{Name: "EI89", Location: "-132.3207, -34.0983", Timestamp: "1485066691", Holder: "Yuan",Family:"Some",Production:"10 kgs"},
		HerbCrops{Name: "129R", Location: "153.0054, 12.6429", Timestamp: "1485153091", Holder: "Carlo",Family:"Some",Production:"10 kgs"},
		HerbCrops{Name: "49W4", Location: "51.9435, 8.2735", Timestamp: "1487745091", Holder: "Fatima",Family:"Some",Production:"10 kgs"},
	}

	i := 0
	for i < len(herb) {
		fmt.Println("i is ", i)
		herbAsBytes, _ := json.Marshal(herb[i])
		APIstub.PutState(strconv.Itoa(i+1), herbAsBytes)
		fmt.Println("Added", herb[i])
		i = i + 1
	}

	return shim.Success(nil)
}

/*
 * The recordherb method *
Fisherman like Sarah would use to record each of her herb catches. 
This method takes in five arguments (attributes to be saved in the ledger). 
 */
func (s *SmartContract) recordHerb(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var herb = HerbCrops{ Name: args[1], Location: args[2], Timestamp: args[3], Holder: args[4] }

	herbAsBytes, _ := json.Marshal(herb)
	err := APIstub.PutState(args[0], herbAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record herbs plucked: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * The queryAllherb method *
allows for assessing all the records added to the ledger(all herb catches)
This method does not take any arguments. Returns JSON string containing results. 
 */
func (s *SmartContract) queryAllHerbs(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add comma before array members,suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllHerbs:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
 * The changeherbHolder method *
The data in the world state can be updated with who has possession. 
This function takes in 2 arguments, herb id and new holder name. 
 */
func (s *SmartContract) changeHerbHolder(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	herbAsBytes, _ := APIstub.GetState(args[0])
	if herbAsBytes == nil {
		return shim.Error("Could not locate herb")
	}
	herb := HerbCrops{}

	json.Unmarshal(herbAsBytes, &herb)
	// Normally check that the specified argument is a valid holder of herb
	// we are skipping this check for this example
	herb.Holder = args[1]

	herbAsBytes, _ = json.Marshal(herb)
	err := APIstub.PutState(args[0], herbAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to change herb production holder: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * main function *
calls the Start function 
The main function starts the chaincode in the container during instantiation.
 */
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
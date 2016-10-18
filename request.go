package main

import (
	//"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
  "time"
)

type Request struct {

}

//Init initializes the request model/smart contract
func (t *Request) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	// Check if table already exists
	_, err := stub.GetTable("RequestTable")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}

	// Create Request Table
	err = stub.CreateTable("RequestTable", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "UID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "DocJSON", Type: shim.ColumnDefinition_BYTES, Key: false},
		&shim.ColumnDefinition{Name: "Status", Type: shim.ColumnDefinition_STRING, Key: false},
    &shim.ColumnDefinition{Name: "Type", Type: shim.ColumnDefinition_STRING, Key: true},
    &shim.ColumnDefinition{Name: "Requester", Type: shim.ColumnDefinition_STRING, Key: true},
    &shim.ColumnDefinition{Name: "CreatedAt", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating Request Table.")
	}

	return nil, nil

}

//SubmitDoc () – Calls ValidateDoc internally and upon success inserts a new row in the table
func (t *Request) SubmitNewRequest(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 6.")
	}

	UID := args[0]
	docJSON := []byte(args[1])
  status := args[2]
  requestType := args[3]
  requester :=args[4]

	//TODO: validate data

  //time
  createdTime := time.Now()



	// Insert a row
	ok, err := stub.InsertRow("RequestTable", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: UID}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: docJSON}},
      &shim.Column{Value: &shim.Column_String_{String_: status}},
			&shim.Column{Value: &shim.Column_String_{String_: requestType}},
      &shim.Column{Value: &shim.Column_String_{String_: requester}},
      &shim.Column{Value: &shim.Column_String_{String_: createdTime.Format(time.RFC3339)}}},
	})

	if !ok && err == nil {
		return nil, errors.New("Document already exists.")
	}

	return nil, err
}

// GetJSON () – returns as JSON a single document w.r.t. the UID
func (t *Request) GetJSON(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
	}

	UID := args[0]
  //requestType := args[1]
  //requester := args[2]
  requestType := "new"
  requester := "testUser"

	// Get the row pertaining to this UID
  var columns []shim.Column
  	col1 := shim.Column{Value: &shim.Column_String_{String_: UID}}
  	columns = append(columns, col1)

  	col2 := shim.Column{Value: &shim.Column_String_{String_: requestType}}
  	columns = append(columns, col2)
    col3 := shim.Column{Value: &shim.Column_String_{String_: requester}}
    columns = append(columns, col3)

  	row, err := stub.GetRow("RequestTable", columns)
  	if err != nil {
  		return nil, fmt.Errorf("Error: Failed retrieving document with UID %s. Error %s", UID, err.Error())
  	}

  	// GetRows returns empty message if key does not exist
  	if len(row.Columns) == 0 {
  		return nil, nil
  	}

  	return row.Columns[1].GetBytes(), nil

}

func (t *Request) GetRequestStatus(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
	}

	UID := args[0]
  //requestType := args[1]
  //requester := args[2]
  requestType := "new"
  requester := "testUser"

	// Get the row pertaining to this UID
  var columns []shim.Column
  	col1 := shim.Column{Value: &shim.Column_String_{String_: UID}}
  	columns = append(columns, col1)

  	col2 := shim.Column{Value: &shim.Column_String_{String_: requestType}}
  	columns = append(columns, col2)
    col3 := shim.Column{Value: &shim.Column_String_{String_: requester}}
    columns = append(columns, col3)

  	row, err := stub.GetRow("RequestTable", columns)
  	if err != nil {
  		return nil, fmt.Errorf("Error: Failed retrieving document with UID %s. Error %s", UID, err.Error())
  	}

  	// GetRows returns empty message if key does not exist
  	if len(row.Columns) == 0 {
  		return nil, nil
  	}

  	return []byte(row.Columns[4].GetString_()), nil

}

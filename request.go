package main

import (
	//"encoding/json"

	"errors"
	"fmt"

	"time"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
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
		&shim.ColumnDefinition{Name: "RequestType", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Requester", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Approver", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Uid", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "DocJSON", Type: shim.ColumnDefinition_BYTES, Key: false},
		&shim.ColumnDefinition{Name: "Status", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Permissions", Type: shim.ColumnDefinition_BYTES, Key: false},
		&shim.ColumnDefinition{Name: "CreatedAt", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating Request Table.")
	}

	return nil, nil

}

//SubmitDoc () – Calls ValidateDoc internally and upon success inserts a new row in the table
func (t *Request) SubmitNewRequest(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	if len(args) != 7 {
		return nil, errors.New("Incorrect number of arguments. Expecting 6.")
	}
	requestType := args[0]
	requester := args[1]
	approver := args[2]
	UID := args[3]
	docJSON := []byte(args[4])
	status := args[5]
	permissions := []byte(args[6])

	//TODO: Validate input

	//time
	createdTime := time.Now()

	// Insert a row
	ok, err := stub.InsertRow("RequestTable", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: requestType}},
			&shim.Column{Value: &shim.Column_String_{String_: requester}},
			&shim.Column{Value: &shim.Column_String_{String_: approver}},
			&shim.Column{Value: &shim.Column_String_{String_: UID}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: docJSON}},
			&shim.Column{Value: &shim.Column_String_{String_: status}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: permissions}},
			&shim.Column{Value: &shim.Column_String_{String_: createdTime.Format(time.RFC3339)}}},
	})

	if !ok && err == nil {
		return nil, errors.New("Document already exists.")
	}

	return nil, err
}

// GetRequestDocument () – returns as JSON a single document w.r.t. the UID
func (t *Request) GetJSON(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3.")
	}
	requestType := "new"
	requester := args[0]
	approver := args[1]
	uid := args[2]

	// Get the row pertaining to this UID
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: requestType}}
	columns = append(columns, col1)
	col2 := shim.Column{Value: &shim.Column_String_{String_: requester}}
	columns = append(columns, col2)
	col3 := shim.Column{Value: &shim.Column_String_{String_: approver}}
	columns = append(columns, col3)
	col4 := shim.Column{Value: &shim.Column_String_{String_: uid}}
	columns = append(columns, col4)

	row, err := stub.GetRow("RequestTable", columns)
	if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with uid %s. Error %s", uid, err.Error())
	}

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		return nil, nil
	}
	fmt.Printf("UID "+row.Columns[3].GetString_()  )


	str := `{ "requestType": "` + row.Columns[0].GetString_() + `", "requester": "` + row.Columns[1].GetString_() + `", "approver": "` + row.Columns[2].GetString_() + `", "uid": "` + row.Columns[3].GetString_() + `", "data": ` + string(row.Columns[4].GetBytes()) + `, "status": "` + row.Columns[5].GetString_() + `", "permissions" : ` + string(row.Columns[6].GetBytes()) + `, "createdAt": "` + row.Columns[7].GetString_() + `"  }`

	fmt.Printf(str)

	return []byte(str), nil

}

func (t *Request) ApproveRequest(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	if len(args) != 3{
		return nil, errors.New("Incorrect number of arguments. Expecting 3.")
	}

	requestType := "new"
	//requestType := args[1]
	requester := args[0]
	approver := args[1]
	uid := args[2]
	//requester := "testUser"
	// Get the row pertaining to this UID
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: requestType}}
	columns = append(columns, col1)
	col2 := shim.Column{Value: &shim.Column_String_{String_: requester}}
	columns = append(columns, col2)
	col3 := shim.Column{Value: &shim.Column_String_{String_: approver}}
	columns = append(columns, col3)
	col4 := shim.Column{Value: &shim.Column_String_{String_: uid}}
	columns = append(columns, col4)

	row, err := stub.GetRow("RequestTable", columns)
	if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving document with uid %s. Error %s", uid, err.Error())
	}

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		return nil, nil
	}
	//update status
	ok, err := stub.ReplaceRow("RequestTable", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[0].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[1].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[2].GetString_()}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[3].GetString_()}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: row.Columns[4].GetBytes()}},
			&shim.Column{Value: &shim.Column_String_{String_: "approved"}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: row.Columns[6].GetBytes()}},
			&shim.Column{Value: &shim.Column_String_{String_: row.Columns[7].GetString_()}}},
	})

	if !ok && err == nil {
		return nil, errors.New("Error updating.")
	}

	return nil, err
}
func (t *Request) GetNewRequests(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
	}
	requestType := "new"
	requester := args[0]
	fmt.Printf("Chaincode - Get List of Requests")
		logger.Infof("Chaincode - Get List of Requests")

	// Get the row pertaining to this UID
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: requestType}}
	columns = append(columns, col1)
	col2 := shim.Column{Value: &shim.Column_String_{String_: requester}}
	columns = append(columns, col2)

	rows, err := stub.GetRows("RequestTable", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}
	//var o string
	var count int
	var outputString string
	for row := range rows {
		if len(row.Columns) != 0 {
			//o = row.Columns[3].GetString_()
			if (count == 0) {
				outputString = "["
			}else{
				outputString = outputString + ", "
			}
			logger.Debugf(" UID "+row.Columns[3].GetString_() )
			str := `{ "requestType": "` + row.Columns[0].GetString_() + `", "requester": "` + row.Columns[1].GetString_() + `", "approver": "` + row.Columns[2].GetString_() + `", "uid": "` + row.Columns[3].GetString_() + `", "data": ` + string(row.Columns[4].GetBytes()) + `, "status": "` + row.Columns[5].GetString_() + `", "permissions" : ` + string(row.Columns[6].GetBytes()) + `, "createdAt": "` + row.Columns[7].GetString_() + `"  }`
			logger.Debugf("str "+str )
			outputString = outputString + str
			count++
		}
	}
	outputString = outputString + "]"
	if( count ==0){
		outputString = `[]`
	}
	logger.Debugf("requests: "+outputString )




	//str := `{ "requestType": "` + row.Columns[0].GetString_() + `", "requester": "` + row.Columns[1].GetString_() + `", "approver": "` + row.Columns[2].GetString_() + `", "uid": "` + row.Columns[3].GetString_() + `", "data": ` + string(row.Columns[4].GetBytes()) + `, "status": "` + row.Columns[5].GetString_() + `", "permissions" : ` + string(row.Columns[6].GetBytes()) + `, "createdAt": "` + row.Columns[7].GetString_() + `"  }`

	//fmt.Printf(str)

	//return []byte(str), nil
return []byte(`{"count": `+strconv.Itoa(count)+`, "data":`+outputString+` }`) , nil
}

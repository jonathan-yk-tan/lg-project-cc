package main

import (
	//"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
  "time"
)

type Document struct {

}

//Init initializes the request model/smart contract
func (t *Document) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	// Check if table already exists
	_, err := stub.GetTable("DocumentTable")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}

	// Create Document Table
	err = stub.CreateTable("DocumentTable", []*shim.ColumnDefinition{
    &shim.ColumnDefinition{Name: "Owner", Type: shim.ColumnDefinition_STRING, Key: true},
    &shim.ColumnDefinition{Name: "Issuer", Type: shim.ColumnDefinition_STRING, Key: true},
    &shim.ColumnDefinition{Name: "DocumentType", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Uid", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "DataJSON", Type: shim.ColumnDefinition_BYTES, Key: false},
		&shim.ColumnDefinition{Name: "Status", Type: shim.ColumnDefinition_STRING, Key: false},
    &shim.ColumnDefinition{Name: "Permissions", Type: shim.ColumnDefinition_BYTES, Key: false},
    &shim.ColumnDefinition{Name: "ExpiryDate", Type: shim.ColumnDefinition_STRING, Key: false},
    &shim.ColumnDefinition{Name: "PreviousUid", Type: shim.ColumnDefinition_STRING, Key: false},
    &shim.ColumnDefinition{Name: "CreatedAt", Type: shim.ColumnDefinition_STRING, Key: false},
  })
	if err != nil {
		return nil, errors.New("Failed creating Document Table.")
	}

	return nil, nil

}

//Issue Document(LG) , creates a new document
func (t *Document) IssueDocument(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	if len(args) != 8 {
		return nil, errors.New("Incorrect number of arguments. Expecting 8.")
	}

	owner :=args[0]
	issuer :=args[1]
	documentType := args[2]
	uid := args[3]
	dataJSON := []byte(args[4])
  status := args[5]
  permissions :=[]byte(args[6])
  expiryDate :=args[7]
  previousUid := ""

	//TODO: Validate input

  //time
  createdTime := time.Now()



	// Insert a row
	ok, err := stub.InsertRow("DocumentTable", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: owner}},
			&shim.Column{Value: &shim.Column_String_{String_: issuer}},
      &shim.Column{Value: &shim.Column_String_{String_: documentType}},
			&shim.Column{Value: &shim.Column_String_{String_: uid}},
      &shim.Column{Value: &shim.Column_Bytes{Bytes: dataJSON}},
      &shim.Column{Value: &shim.Column_String_{String_: status}},
      &shim.Column{Value: &shim.Column_Bytes{Bytes: permissions}},
      &shim.Column{Value: &shim.Column_String_{String_: expiryDate}},
      &shim.Column{Value: &shim.Column_String_{String_: previousUid}},
      &shim.Column{Value: &shim.Column_String_{String_: createdTime.Format(time.RFC3339)}}},
	})

	if !ok && err == nil {
		return nil, errors.New("Document already exists.")
	}

	return nil, err
}

// GetDocument () – returns as JSON a single document w.r.t. the UID
func (t *Document) GetLgJSON(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3.")
	}


  //requestType := args[1]
  owner := args[0]
  issuer := args[1]
  documentType := "LG"
	uid := args[2]
  //requester := "testUser"
	// Get the row pertaining to this UID
  var columns []shim.Column
  	col1 := shim.Column{Value: &shim.Column_String_{String_: owner}}
  	columns = append(columns, col1)

  	col2 := shim.Column{Value: &shim.Column_String_{String_: issuer}}
  	columns = append(columns, col2)
    col3 := shim.Column{Value: &shim.Column_String_{String_: documentType}}
    columns = append(columns, col3)
		col4 := shim.Column{Value: &shim.Column_String_{String_: uid}}
    columns = append(columns, col4)

  	row, err := stub.GetRow("DocumentTable", columns)
  	if err != nil {
  		return nil, fmt.Errorf("Error: Failed retrieving document with uid %s. Error %s", uid, err.Error())
  	}

  	// GetRows returns empty message if key does not exist
  	if len(row.Columns) == 0 {
  		return nil, nil
  	}
    fmt.Printf("UID\n")
    fmt.Printf(`"UID": "`+row.Columns[3].GetString_()+`"`)
    str := `{ "owner": "`+row.Columns[0].GetString_()+`", "issuer": "` + row.Columns[1].GetString_()+`", "documentType": "` + row.Columns[2].GetString_()+`", "uid": "` + row.Columns[3].GetString_()+`", "data": ` + string(row.Columns[4].GetBytes()) +`, "status": "`+ row.Columns[5].GetString_() +`", "permissions": ` + string(row.Columns[6].GetBytes())+`, "expiryDate":"`+row.Columns[7].GetString_() + `", "previousUid":"`+row.Columns[8].GetString_()+`", "createdAt": "` + row.Columns[9].GetString_() +`"  }`
    fmt.Printf("JSON\n")
    fmt.Printf(str)
    //str := `{ "UID": `+row.Columns[0].GetString_()+`  }`
  	return []byte(str), nil

}

// CancelLGDocument () – returns as JSON a single document w.r.t. the UID
func (t *Document) CancelLGDocument(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3.")
	}


  //requestType := args[1]
  owner := args[0]
  issuer := args[1]
  documentType := "LG"
	uid := args[2]
  //requester := "testUser"
	// Get the row pertaining to this UID
  var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: owner}}
	columns = append(columns, col1)
	col2 := shim.Column{Value: &shim.Column_String_{String_: issuer}}
	columns = append(columns, col2)
	col3 := shim.Column{Value: &shim.Column_String_{String_: documentType}}
	columns = append(columns, col3)
  col4 := shim.Column{Value: &shim.Column_String_{String_: uid}}
  columns = append(columns, col4)




  	row, err := stub.GetRow("DocumentTable", columns)
  	if err != nil {
  		return nil, fmt.Errorf("Error: Failed retrieving document with uid %s. Error %s", uid, err.Error())
  	}

  	// GetRows returns empty message if key does not exist
  	if len(row.Columns) == 0 {
  		return nil, nil
  	}

		fmt.Printf(row.Columns[0].GetString_())
		fmt.Printf("Updating")

		//update status
		ok, err := stub.ReplaceRow("DocumentTable", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: row.Columns[0].GetString_()}},
				&shim.Column{Value: &shim.Column_String_{String_: row.Columns[1].GetString_()}},
				&shim.Column{Value: &shim.Column_String_{String_: row.Columns[2].GetString_()}},
				&shim.Column{Value: &shim.Column_String_{String_: row.Columns[3].GetString_()}},
				&shim.Column{Value: &shim.Column_Bytes{Bytes: row.Columns[4].GetBytes()}},
				&shim.Column{Value: &shim.Column_String_{String_: "cancelled"}},
				&shim.Column{Value: &shim.Column_Bytes{Bytes: row.Columns[6].GetBytes()}},
				&shim.Column{Value: &shim.Column_String_{String_: row.Columns[7].GetString_()}},
				&shim.Column{Value: &shim.Column_String_{String_: row.Columns[8].GetString_()}},
				&shim.Column{Value: &shim.Column_String_{String_: row.Columns[9].GetString_()}}},
		})

		if !ok && err == nil {
			return nil, errors.New("Error updating.")
		}

		return nil, err

}

package chaincode_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/S-Amine/nhms-smart-contract/chaincode"
	"github.com/S-Amine/nhms-smart-contract/chaincode/mocks"
	"github.com/hyperledger/fabric-protos-go-apiv2/ledger/queryresult"
	"github.com/stretchr/testify/require"
)

// Test InitLedger
func TestInitLedger(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	assetTransfer := chaincode.SmartContract{}
	err := assetTransfer.InitLedger(transactionContext)
	require.NoError(t, err)

	chaincodeStub.PutStateReturns(fmt.Errorf("failed inserting key"))
	err = assetTransfer.InitLedger(transactionContext)
	require.EqualError(t, err, "failed to put to world state. failed inserting key")
}

// Test CreatePatient
func TestCreatePatient(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	assetTransfer := chaincode.SmartContract{}
	err := assetTransfer.CreatePatient(transactionContext, "123", "John", "Doe", "1990-01-01", "M", "987", "876", "None", "Peanuts", "Asthma", "")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns([]byte{}, nil)
	err = assetTransfer.CreatePatient(transactionContext, "123", "John", "Doe", "1990-01-01", "M", "987", "876", "None", "Peanuts", "Asthma", "")
	require.EqualError(t, err, "the patient with NIN 123 already exists")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = assetTransfer.CreatePatient(transactionContext, "123", "John", "Doe", "1990-01-01", "M", "987", "876", "None", "Peanuts", "Asthma", "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve asset")
}

// Test ReadPatient
func TestReadPatient(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	expectedPatient := &chaincode.Patient{NIN: "123"}
	bytes, err := json.Marshal(expectedPatient)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	assetTransfer := chaincode.SmartContract{}
	patient, err := assetTransfer.ReadPatient(transactionContext, "123")
	require.NoError(t, err)
	require.Equal(t, expectedPatient, patient)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	_, err = assetTransfer.ReadPatient(transactionContext, "123")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve asset")

	chaincodeStub.GetStateReturns(nil, nil)
	patient, err = assetTransfer.ReadPatient(transactionContext, "123")
	require.EqualError(t, err, "the patient with NIN 123 does not exist")
	require.Nil(t, patient)
}

// Test UpdatePatient
func TestUpdatePatient(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	expectedPatient := &chaincode.Patient{NIN: "123"}
	bytes, err := json.Marshal(expectedPatient)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	assetTransfer := chaincode.SmartContract{}
	err = assetTransfer.UpdatePatient(transactionContext, "123", "John", "Doe", "1990-01-01", "M", "987", "876", "None", "Peanuts", "Asthma", "")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, nil)
	err = assetTransfer.UpdatePatient(transactionContext, "123", "John", "Doe", "1990-01-01", "M", "987", "876", "None", "Peanuts", "Asthma", "")
	require.EqualError(t, err, "the patient with NIN 123 does not exist")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = assetTransfer.UpdatePatient(transactionContext, "123", "John", "Doe", "1990-01-01", "M", "987", "876", "None", "Peanuts", "Asthma", "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve asset")
}

// Test DeletePatient
func TestDeletePatient(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	patient := &chaincode.Patient{NIN: "123"}
	bytes, err := json.Marshal(patient)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	chaincodeStub.DelStateReturns(nil)
	assetTransfer := chaincode.SmartContract{}
	err = assetTransfer.DeletePatient(transactionContext, "123")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, nil)
	err = assetTransfer.DeletePatient(transactionContext, "123")
	require.EqualError(t, err, "the patient with NIN 123 does not exist")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = assetTransfer.DeletePatient(transactionContext, "123")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve asset")
}

// Test GetAllPatients
func TestGetAllPatients(t *testing.T) {
	patient := &chaincode.Patient{NIN: "123"}
	bytes, err := json.Marshal(patient)
	require.NoError(t, err)

	iterator := &mocks.StateQueryIterator{}
	iterator.HasNextReturnsOnCall(0, true)
	iterator.HasNextReturnsOnCall(1, false)
	iterator.NextReturns(&queryresult.KV{Value: bytes}, nil)

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	chaincodeStub.GetStateByRangeReturns(iterator, nil)
	assetTransfer := &chaincode.SmartContract{}
	patients, err := assetTransfer.GetAllPatients(transactionContext)
	require.NoError(t, err)
	require.Equal(t, []*chaincode.Patient{patient}, patients)

	iterator.HasNextReturns(true)
	iterator.NextReturns(nil, fmt.Errorf("failed retrieving next item"))
	patients, err = assetTransfer.GetAllPatients(transactionContext)
	require.EqualError(t, err, "failed retrieving next item")
	require.Nil(t, patients)

	chaincodeStub.GetStateByRangeReturns(nil, fmt.Errorf("failed retrieving all assets"))
	patients, err = assetTransfer.GetAllPatients(transactionContext)
	require.EqualError(t, err, "failed retrieving all assets")
	require.Nil(t, patients)
}

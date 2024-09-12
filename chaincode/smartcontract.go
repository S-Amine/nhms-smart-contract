package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// SmartContract provides functions for managing a Patient
type SmartContract struct {
	contractapi.Contract
}

// Patient describes basic details of a patient
type Patient struct {
	NIN                string   `json:"nin"`
	FirstName          string   `json:"firstName"`
	LastName           string   `json:"lastName"`
	DateOfBirth        string   `json:"dateOfBirth"`
	Sex                string   `json:"sex"`
	MotherNIN          string   `json:"motherNin"`
	FatherNIN          string   `json:"fatherNin"`
	FamilyMedicalHistory string   `json:"familyMedicalHistory"`
	Allergy            string   `json:"allergy"`
	ChronicIllnesses   string   `json:"chronicIllnesses"`
	AmendedFrom        string   `json:"amendedFrom"`
}

// InitLedger adds a base set of patients to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	patients := []Patient{
		{NIN: "123456789", FirstName: "John", LastName: "Doe", DateOfBirth: "1990-01-01", Sex: "M", MotherNIN: "987654321", FatherNIN: "876543210", FamilyMedicalHistory: "None", Allergy: "Peanuts", ChronicIllnesses: "Asthma", AmendedFrom: ""},
		{NIN: "987654321", FirstName: "Jane", LastName: "Smith", DateOfBirth: "1985-02-02", Sex: "F", MotherNIN: "123456789", FatherNIN: "234567890", FamilyMedicalHistory: "Diabetes", Allergy: "None", ChronicIllnesses: "None", AmendedFrom: ""},
	}

	for _, patient := range patients {
		patientJSON, err := json.Marshal(patient)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(patient.NIN, patientJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreatePatient issues a new patient record to the world state with given details.
func (s *SmartContract) CreatePatient(ctx contractapi.TransactionContextInterface, nin string, firstName string, lastName string, dateOfBirth string, sex string, motherNIN string, fatherNIN string, familyMedicalHistory string, allergy string, chronicIllnesses string, amendedFrom string) error {
	exists, err := s.PatientExists(ctx, nin)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the patient with NIN %s already exists", nin)
	}

	patient := Patient{
		NIN:                nin,
		FirstName:          firstName,
		LastName:           lastName,
		DateOfBirth:        dateOfBirth,
		Sex:                sex,
		MotherNIN:          motherNIN,
		FatherNIN:          fatherNIN,
		FamilyMedicalHistory: familyMedicalHistory,
		Allergy:            allergy,
		ChronicIllnesses:   chronicIllnesses,
		AmendedFrom:        amendedFrom,
	}
	patientJSON, err := json.Marshal(patient)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(nin, patientJSON)
}

// ReadPatient returns the patient stored in the world state with given nin.
func (s *SmartContract) ReadPatient(ctx contractapi.TransactionContextInterface, nin string) (*Patient, error) {
	patientJSON, err := ctx.GetStub().GetState(nin)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if patientJSON == nil {
		return nil, fmt.Errorf("the patient with NIN %s does not exist", nin)
	}

	var patient Patient
	err = json.Unmarshal(patientJSON, &patient)
	if err != nil {
		return nil, err
	}

	return &patient, nil
}

// UpdatePatient updates an existing patient in the world state with provided parameters.
func (s *SmartContract) UpdatePatient(ctx contractapi.TransactionContextInterface, nin string, firstName string, lastName string, dateOfBirth string, sex string, motherNIN string, fatherNIN string, familyMedicalHistory string, allergy string, chronicIllnesses string, amendedFrom string) error {
	exists, err := s.PatientExists(ctx, nin)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the patient with NIN %s does not exist", nin)
	}

	// overwriting original patient with new details
	patient := Patient{
		NIN:                nin,
		FirstName:          firstName,
		LastName:           lastName,
		DateOfBirth:        dateOfBirth,
		Sex:                sex,
		MotherNIN:          motherNIN,
		FatherNIN:          fatherNIN,
		FamilyMedicalHistory: familyMedicalHistory,
		Allergy:            allergy,
		ChronicIllnesses:   chronicIllnesses,
		AmendedFrom:        amendedFrom,
	}
	patientJSON, err := json.Marshal(patient)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(nin, patientJSON)
}

// DeletePatient deletes a given patient from the world state.
func (s *SmartContract) DeletePatient(ctx contractapi.TransactionContextInterface, nin string) error {
	exists, err := s.PatientExists(ctx, nin)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the patient with NIN %s does not exist", nin)
	}

	return ctx.GetStub().DelState(nin)
}

// PatientExists returns true when a patient with the given NIN exists in the world state
func (s *SmartContract) PatientExists(ctx contractapi.TransactionContextInterface, nin string) (bool, error) {
	patientJSON, err := ctx.GetStub().GetState(nin)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return patientJSON != nil, nil
}

// GetAllPatients returns all patients found in world state
func (s *SmartContract) GetAllPatients(ctx contractapi.TransactionContextInterface) ([]*Patient, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all patients in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var patients []*Patient
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var patient Patient
		err = json.Unmarshal(queryResponse.Value, &patient)
		if err != nil {
			return nil, err
		}
		patients = append(patients, &patient)
	}

	return patients, nil
}


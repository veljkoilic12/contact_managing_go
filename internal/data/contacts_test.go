package data

import (
	"fmt"
	"salestrekker_technical_interview.veljkoilic/internal/validator"
	"testing"
)

// Testing areContactsEqual method
func TestAreContactsEqual(t *testing.T) {
	expectedFirstName := "Veljko"
	expectedLastName := "Ilic"
	expectedTelephone := "+38163577442"

	testCases := []struct {
		contact1 *Contact
		contact2 *Contact
		expected bool
	}{
		// case where contacts are the same, but ids' are different
		{&Contact{ID: 1, FirstName: expectedFirstName, LastName: expectedLastName, Telephone: expectedTelephone},
			&Contact{ID: 2, FirstName: expectedFirstName, LastName: expectedLastName, Telephone: expectedTelephone},
			true},

		// case where contacts are the same
		{&Contact{FirstName: expectedFirstName, LastName: expectedLastName, Telephone: expectedTelephone},
			&Contact{FirstName: expectedFirstName, LastName: expectedLastName, Telephone: expectedTelephone},
			true},

		//case where first names are different
		{&Contact{FirstName: expectedFirstName, LastName: expectedLastName, Telephone: expectedTelephone},
			&Contact{FirstName: "Marko", LastName: expectedLastName, Telephone: expectedTelephone},
			false},

		// case where last names are different
		{&Contact{FirstName: expectedFirstName, LastName: expectedLastName, Telephone: expectedTelephone},
			&Contact{FirstName: expectedFirstName, LastName: "Markovic", Telephone: expectedTelephone},
			false},

		// case where telephones are different
		{&Contact{FirstName: expectedFirstName, LastName: expectedLastName, Telephone: expectedTelephone},
			&Contact{FirstName: expectedFirstName, LastName: expectedLastName, Telephone: "+38163587442"},
			false},
	}

	for _, tc := range testCases {
		got := areContactsEqual(tc.contact1, tc.contact2)
		if got != tc.expected {
			t.Errorf("want %v; got %v", tc.expected, got)
		}
	}
}

// Testing generateID method
func TestGenerateID(t *testing.T) {
	testCases := []struct {
		contacts []Contact
		expected int64
	}{
		{[]Contact{}, 1},
		{[]Contact{Contact{ID: 1, FirstName: "Veljko", LastName: "Ilic", Telephone: "+38163577442"},
			Contact{ID: 2, FirstName: "Marko", LastName: "Markovic", Telephone: "+38163587442"}}, 3},
	}

	for _, tc := range testCases {
		got := generateID(tc.contacts)
		if got != tc.expected {
			t.Errorf("want %d; got %d", tc.expected, got)
		}
	}
}

// Testing contact validator
func TestValidateValidContact(t *testing.T) {
	v := validator.New()
	contact := &Contact{FirstName: "Veljko", LastName: "Ilic", Telephone: "+38163577442"}

	ValidateContact(v, contact)

	if !v.IsValid() {
		fmt.Errorf("want valid; got invalid")
	}
}

func TestValidateInvalidContact(t *testing.T) {
	testCases := []struct {
		contact *Contact
	}{
		{&Contact{FirstName: "", LastName: "Ilic", Telephone: "+38163577442"}},
		{&Contact{FirstName: "Veljko", LastName: "", Telephone: "+38163577442"}},
		{&Contact{FirstName: "Veljko", LastName: "Ilic", Telephone: ""}},
		{&Contact{FirstName: "Veljko", LastName: "Ilic", Telephone: "4738473843"}},
	}

	for _, tc := range testCases {
		v := validator.New()
		ValidateContact(v, tc.contact)
		if v.IsValid() {
			fmt.Errorf("want invalid; got valid")
		}
	}
}

// Testing getting one contact
func TestGetContactFromFile(t *testing.T) {

	data := []Contact{
		{ID: 1, FirstName: "Veljko", LastName: "Ilic", Telephone: "+38163577442"},
		{ID: 2, FirstName: "Marko", LastName: "Markovic", Telephone: "+38163587442"},
	}

	cm := ContactsModel{data}

	testCases := []struct {
		id              int64
		expectedContact *Contact
		expectedError   error
	}{
		{1, &Contact{ID: 2, FirstName: "Marko", LastName: "Markovic", Telephone: "+38163587442"}, nil},
		{3, nil, ErrRecordNotFound},
		{-1, nil, ErrRecordNotFound},
	}

	for _, tc := range testCases {
		gotCont, gotErr := cm.GetContact(tc.id)
		if gotCont != tc.expectedContact || gotErr != tc.expectedError {
			fmt.Errorf("want %v error and %v contact; got %v error and %v contact", tc.expectedError, tc.expectedContact, gotErr, gotCont)
		}
	}
}

// Testing insert contact
func TestInsertContact(t *testing.T) {

	data := []Contact{
		{ID: 1, FirstName: "Veljko", LastName: "Ilic", Telephone: "+38163577442"},
		{ID: 2, FirstName: "Marko", LastName: "Markovic", Telephone: "+38163587442"},
	}

	contact := &Contact{FirstName: "Ilija", LastName: "Ilinovic", Telephone: "+38164598332"}

	cm := ContactsModel{data}
	cm.InsertContact(contact)

	expectedId := int64(3)
	gotContact, err := cm.GetContact(expectedId)

	if err != nil {
		fmt.Errorf("did not get expected value")
	}

	if contact != gotContact {
		fmt.Errorf("did not get expected contact")
	}

	if gotContact.ID != expectedId {
		fmt.Errorf("want %d; got %d", expectedId, gotContact.ID)
	}
}

func TestInsertSameContact(t *testing.T) {
	data := []Contact{
		{ID: 1, FirstName: "Veljko", LastName: "Ilic", Telephone: "+38163577442"},
	}
	contact := &Contact{ID: 1, FirstName: "Veljko", LastName: "Ilic", Telephone: "+38163577442"}

	cm := ContactsModel{data}
	err := cm.InsertContact(contact)

	if err == nil {
		fmt.Errorf("expected error")
	}
}

// Test deleting contacts
func TestDeletingContacts(t *testing.T) {
	data := []Contact{
		{ID: 1, FirstName: "Veljko", LastName: "Ilic", Telephone: "+38163577442"},
		{ID: 2, FirstName: "Marko", LastName: "Markovic", Telephone: "+38163587442"},
	}
	cm := ContactsModel{data}

	testCases := []struct {
		id  int64
		err error
	}{
		{2, nil},
		{4, ErrRecordNotFound},
		{-1, ErrRecordNotFound},
	}

	for _, tc := range testCases {
		err := cm.DeleteContact(tc.id)
		if err != tc.err {
			fmt.Errorf("want %v; got %v", tc.err, err)
		}
	}
}

package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"salestrekker_technical_interview.veljkoilic/internal/validator"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Contact struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Telephone string `json:"telephone"`
}

type ContactsModel struct {
	Contacts []Contact
}

func NewModel() ContactsModel {
	return ContactsModel{
		Contacts: []Contact{},
	}
}

func ValidateContact(v *validator.Validator, contact *Contact) {
	// check if the fields are empty
	v.Check(contact.FirstName != "", "first_name", "must be provided")
	v.Check(contact.LastName != "", "last_name", "must be provided")
	v.Check(contact.Telephone != "", "telephone", "must be provided")

	// check if the telephone field is valid for Serbia
	v.Check(validator.Matches(contact.Telephone, validator.PhoneRX), "telephone", "must be valid Serbian number (example: +38163567893)")
}

// get a specific record from the contacts
func (cm *ContactsModel) GetContact(id int64) (*Contact, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	for _, contact := range cm.Contacts {
		if contact.ID == id {
			return &contact, nil
		}
	}

	return nil, ErrRecordNotFound
}

// inserting a new record in the contacts file
func (cm *ContactsModel) InsertContact(contact *Contact) error {

	for _, existingContact := range cm.Contacts {
		if areContactsEqual(&existingContact, contact) {
			return errors.New("same contact already exists")
		}
	}

	// Generate an id for the new contact and assign it to ID field of the contact
	newContactID := generateID(cm.Contacts)
	contact.ID = newContactID
	cm.Contacts = append(cm.Contacts, *contact)

	// Save all contacts + newly created one to the file
	cm.SaveAllContacts()
	return nil
}

// Delete a specific record in the contacts file
func (cm *ContactsModel) DeleteContact(id int64) error {

	if id < 1 {
		return ErrRecordNotFound
	}

	// Loop through the slice of Contacts and delete a contact if it finds matching id
	for ind, contact := range cm.Contacts {
		if contact.ID == id {
			cm.Contacts = append(cm.Contacts[:ind], cm.Contacts[ind+1:]...)
			cm.SaveAllContacts()
			return nil
		}
	}

	// Return error if record is not found
	return ErrRecordNotFound
}

// Load a contact list from a file
func (cm *ContactsModel) GetAllContacts() error {
	// Create a file if it does not exist, open if exists
	file, err := os.OpenFile("contacts.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Load the contacts to the list
	err = json.NewDecoder(file).Decode(&cm.Contacts)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

// Save all existing contacts to a JSON file
func (cm *ContactsModel) SaveAllContacts() {
	file, err := os.OpenFile("contacts.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}
	defer file.Close()

	// Write the updated data back to the file
	file.Truncate(0)
	file.Seek(0, 0)
	encoder := json.NewEncoder(file)
	err = encoder.Encode(cm.Contacts)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
}

// Generate ID based on the maximum value ID in the dataset
func generateID(contacts []Contact) int64 {
	id := int64(0)
	for _, v := range contacts {
		if v.ID > id {
			id = v.ID
		}
	}
	id = id + 1

	return id
}

// Compare 2 structs, omitting id field
func areContactsEqual(contactOne, contactTwo *Contact) bool {
	return contactOne.FirstName == contactTwo.FirstName && contactOne.LastName == contactTwo.LastName && contactOne.Telephone == contactTwo.Telephone
}

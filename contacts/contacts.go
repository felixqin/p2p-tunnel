package contacts

import (
	"fmt"
	"log"
)

var contacts []*Contact

func addContact(contact *Contact) error {
	log.Println("add contact:", contact)
	for _, item := range contacts {
		// log.Println("item:", item)
		if item.Name == contact.Name {
			*item = *contact
			return nil
		}
	}

	contacts = append(contacts, contact)
	// log.Println("contacts:", contacts)
	return nil
}

func findContact(name string) (*Contact, error) {
	for _, contact := range contacts {
		if contact.Name == name {
			return contact, nil
		}
	}

	log.Printf("contact(%s) not found!", name)
	return nil, fmt.Errorf("contact not found")
}

package contacts

import "log"

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

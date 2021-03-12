package contacts

import "log"

var contacts []*Contact

func addContact(contact *Contact) error {
	log.Println("add contact:", contact)
	for _, item := range contacts {
		// log.Println("item:", item)
		if item.ClientId == contact.ClientId && item.Owner == item.Owner {
			*item = *contact
			return nil
		}
	}

	contacts = append(contacts, contact)
	// log.Println("contacts:", contacts)
	return nil
}

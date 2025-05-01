package main

import (
	"fmt"

	"github.com/Floriansylvain/GoAntiSqlPOC/store"
)

func main() {
	store.InitDB()
	defer store.CloseDB()

	// Create user
	user, err := store.CreateUser("Bob", "bob@example.com", "hunter2")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created user: %+v\n", user)

	// Add addresses
	_, _ = store.AddAddressToUser(user.ID, "123 Main St", "Paris", "75001")
	_, _ = store.AddAddressToUser(user.ID, "456 Side Rd", "Lyon", "69000")

	// Retrieve addresses
	addrs, err := store.ListAddressesByUserID(user.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println("Addresses for", user.Name)
	for _, a := range addrs {
		fmt.Printf("- %s, %s %s\n", a.Line1, a.City, a.ZipCode)
	}

	// Attempt duplicate email
	_, err = store.CreateUser("Alice", "bob@example.com", "oops")
	if err != nil {
		fmt.Println("Expected error:", err)
	}
}

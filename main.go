package main

import (
	"fmt"
	"strings"

	"github.com/Floriansylvain/GoAntiSqlPOC/store"
	"github.com/fatih/color"
)

func printSection(title string) {
	color.Cyan("\n" + strings.Repeat("-", 40))
	color.Cyan(" " + title)
}

func printSuccess(message string) {
	color.Green("✅ " + message)
}

func printInfo(message string) {
	color.White("  " + message)
}

func printWarning(message string) {
	color.Yellow("⚠️  " + message)
}

func main() {
	store.InitDB()
	defer store.CloseDB()

	printSection("USER CREATION")
	user, err := store.CreateUser("Bob", "bob@example.com", "hunter2")
	if err != nil {
		panic(err)
	}
	printInfo(fmt.Sprintf("Created: %s | %s | %s", user.ID, user.Name, user.Email))

	printSection("ADDRESS MANAGEMENT")
	addr1, err := store.AddAddressToUser(user.ID, "123 Main St", "Paris", "75001")
	if err != nil {
		panic(err)
	}
	printInfo(fmt.Sprintf("Added: %s | %s, %s %s", addr1.ID, addr1.Line1, addr1.City, addr1.ZipCode))

	addr2, err := store.AddAddressToUser(user.ID, "456 Side Rd", "Lyon", "69000")
	if err != nil {
		panic(err)
	}
	printInfo(fmt.Sprintf("Added: %s | %s, %s %s", addr2.ID, addr2.Line1, addr2.City, addr2.ZipCode))

	addrs, err := store.ListAddressesByUserID(user.ID)
	if err != nil {
		panic(err)
	}
	printInfo(fmt.Sprintf("Retrieved: %d addresses for %s", len(addrs), user.Name))

	printSection("DATA UPDATES")
	if len(addrs) > 0 {
		printInfo(fmt.Sprintf("Original city: %s → Changing to Marseille", addrs[0].City))
		addrs[0].City = "Marseille"
		err := store.UpdateAddress(addrs[0])
		if err != nil {
			panic(err)
		}
	}

	addrs, _ = store.ListAddressesByUserID(user.ID)
	printInfo(fmt.Sprintf("Address now: %s, %s %s", addrs[0].Line1, addrs[0].City, addrs[0].ZipCode))

	printInfo(fmt.Sprintf("User name: %s → Changing to Bobby", user.Name))
	user.Name = "Bobby"
	if err := store.UpdateUser(user); err != nil {
		panic(err)
	}
	newUser, _ := store.FindUserByID(user.ID)
	printInfo(fmt.Sprintf("Name now: %s", newUser.Name))

	printSection("SEARCH & CONSTRAINTS")
	foundUser, _ := store.FindUserByEmail("bob@example.com")
	printInfo(fmt.Sprintf("Found by email: %s | %s", foundUser.ID, foundUser.Name))

	_, err = store.CreateUser("Alice", "bob@example.com", "password123")
	if err != nil {
		printSuccess("Email uniqueness enforced: " + err.Error())
	} else {
		printWarning("Failed email uniqueness check!")
	}
}

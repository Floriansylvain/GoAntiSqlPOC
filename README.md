# GoAntiSqlPOC - SQL-Free Database Management in Go

This project demonstrates a proof of concept for a database solution in Go that uses zero SQL or text-based query languages. Instead, it leverages [BadgerDB](https://github.com/dgraph-io/badger), a fast key-value store written in Go, to provide structured data storage with API-style access patterns.

## Disclaimer ⚠️

This project is not clean architecture friendly, yet.

## Core Features

- **No SQL or Query Languages** - Pure programmatic data access through Go functions
- **Strong Schema Enforcement** - Type safety via Go structs
- **Relational-Like Data Modeling** - Links between entities using IDs and indexes
- **Clean API-Style Access** - Simple function calls instead of writing queries

## Project Structure

- **main.go**: Example usage demonstrating CRUD operations and relationships
- **store/db.go**: Database connection management with in-memory storage
- **store/user.go**: User entity operations (create, find, update, list)
- **store/address.go**: Address entity operations with user relationships

## How It Works

### Database Management

The project uses BadgerDB as the underlying storage engine. The `store.InitDB()` function sets up an in-memory database connection with a singleton pattern using `sync.Once` to ensure it's only initialized once.

### Data Models

The project defines strongly typed Go structs to enforce schema:

- `User`: Stores user information with ID, name, email, and password
- `Address`: Stores address information with reference to a user via UserID

### Data Access Pattern

Instead of SQL queries, all data access is done through simple function calls:

```go
// Create a user
user, _ := store.CreateUser("Bob", "bob@example.com", "hunter2")

// Add an address
addr, _ := store.AddAddressToUser(user.ID, "123 Main St", "Paris", "75001")

// List addresses for a user
addresses, _ := store.ListAddressesByUserID(user.ID)

// Update user data
user.Name = "Bobby"
store.UpdateUser(user)
```

### "Relational" Modeling

The project demonstrates relationship modeling between users and addresses:

1. Each address has a `UserID` field linking it to its owner
2. The `ListAddressesByUserID` function retrieves all addresses for a specific user
3. Updates to either entity maintain relationship integrity

### Data Consistency

The project implements data consistency features:

- Email uniqueness enforcement using an email-to-ID index
- Transaction support for atomic operations
- Error handling for constraint violations

### Terminal Demo

The main.go file provides a demonstration that showcases:

1. User creation with validation
2. Address management (create, read, update)
3. Relationship handling (users have multiple addresses)
4. Constraint enforcement (unique emails)
5. Entity retrieval by ID and other fields

## Running the Demo

```bash
go run main.go
```

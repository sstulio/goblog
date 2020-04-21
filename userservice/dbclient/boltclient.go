package dbclient

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/callistaenterprise/goblog/userservice/model"
	"log"
	"strconv"
)

//IBoltClient interface
type IBoltClient interface {
	OpenBoltDb()
	QueryUser(userID string) (model.User, error)
	Seed()
	Check() bool
}

//BoltClient implementation
// Real implementation
type BoltClient struct {
	boltDB *bolt.DB
}

//OpenBoltDb open db
func (bc *BoltClient) OpenBoltDb() {
	var err error
	bc.boltDB, err = bolt.Open("users.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
}

//QueryUser query user
func (bc *BoltClient) QueryUser(userID string) (model.User, error) {
	// Allocate an empty User instance we'll let json.Unmarhal populate for us in a bit.
	user := model.User{}

	// Read an object from the bucket using boltDB.View
	err := bc.boltDB.View(func(tx *bolt.Tx) error {
		// Read the bucket from the DB
		b := tx.Bucket([]byte("UserBucket"))

		// Read the value identified by our userID supplied as []byte
		userBytes := b.Get([]byte(userID))
		if userBytes == nil {
			return fmt.Errorf("No user found for " + userID)
		}
		// Unmarshal the returned bytes into the user struct we created at
		// the top of the function
		json.Unmarshal(userBytes, &user)

		// Return nil to indicate nothing went wrong, e.g no error
		return nil
	})
	// If there were an error, return the error
	if err != nil {
		return model.User{}, err
	}
	// Return the User struct and nil as error.
	return user, nil
}

//Seed Start seeding users
func (bc *BoltClient) Seed() {
	bc.initializeBucket()
	bc.seedUsers()
}

//Check Naive healthcheck, just makes sure the DB connection has been initialized.
func (bc *BoltClient) Check() bool {
	return bc.boltDB != nil
}

// Creates an "UserBucket" in our BoltDB. It will overwrite any existing bucket of the same name.
func (bc *BoltClient) initializeBucket() {
	bc.boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("UserBucket"))
		if err != nil {
			return fmt.Errorf("create bucket failed: %s", err)
		}
		return nil
	})
}

// Seed (n) make-believe user objects into the AcountBucket bucket.
func (bc *BoltClient) seedUsers() {

	total := 100
	for i := 0; i < total; i++ {

		// Generate a key 10000 or larger
		key := strconv.Itoa(10000 + i)

		// Create an instance of our User struct
		acc := model.User{
			ID:   key,
			Name: "Person_" + strconv.Itoa(i),
		}

		// Serialize the struct to JSON
		jsonBytes, _ := json.Marshal(acc)

		// Write the data to the UserBucket
		bc.boltDB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("UserBucket"))
			err := b.Put([]byte(key), jsonBytes)
			return err
		})
	}
	fmt.Printf("Seeded %v fake users...\n", total)
}

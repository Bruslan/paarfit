package data

import (
	"fmt"
	"github.com/gocql/gocql"
	"testing"
	"time"
)

// test data
var user = User{
	UserName: "freezy",
	Email:    "jansifreezy@koko.com",
	Password: "peter_pass",
}

func TestInit(t *testing.T) {
	fmt.Println("Testing Print")

	uuid1, err := gocql.RandomUUID()
	if err != nil {
		fmt.Println("uuid1 error")
	}
	fmt.Println("UUID=", uuid1)
	fmt.Println(uuid1, "amanakoidum")
	fmt.Println(uuid1.String())

	// connect to the cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "user_keyspace"
	cluster.Consistency = gocql.One
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "ianzndb",
		Password: "Lov3toN8t",
	}
	session, err := cluster.CreateSession()
	fmt.Println("Session error:", err)

	defer session.Close()

	fmt.Println(user.Email)
	fmt.Println(user.Name)

	// create new uuid
	uuid, err := gocql.RandomUUID()
	if err != nil {
		fmt.Println("Could not create new User ID")
	}
	// check if username exists
	iter := session.Query(`SELECT * FROM user_keyspace.users WHERE name=?`, user.Name).Iter()
	fmt.Println("iterator length:", iter.NumRows(), iter)
	// check if email exists:
	iter2 := session.Query(`SELECT * FROM user_keyspace.users WHERE email=?`, user.Email).Iter()
	fmt.Println("iterator length:", iter2.NumRows(), iter2)

	if iter.NumRows() == 0 && iter2.NumRows() == 0 {
		// INSERT INTO user_keyspace.users (id, created_at, email, name, password, uuid) VALUES (1, toTimestamp(now()), 'mongo@holiday.com', 'mr.bongo', 'lalala', now());
		// insert a tweet
		if err := session.Query(`INSERT INTO user_keyspace.users (id, city, created_at, email, name, pass) VALUES (?, ?, ?, ?, ?, ?) IF NOT EXISTS`,
			uuid, "rom", time.Now(), "jansifreezy@koko.com", "freezy", "freezy_pw").Exec(); err != nil {
			fmt.Println("Inserted values and thats the err:", err)
		}
	}

	var id gocql.UUID
	var emailo, namee string
	var created time.Time

	email := ""

	// first query
	if err := session.Query(`SELECT id, email, created_at FROM user_keyspace.users`).Scan(&id, &email, &created); err != nil {
		fmt.Println("First query:", err)
	}
	fmt.Println("First query, im proud :)", id, email, created)

	// list all tweets
	iter = session.Query(`SELECT id, email, name FROM user_keyspace.users`).Iter()
	for iter.Scan(&id, &emailo, &namee) {
		fmt.Println("Iteration query boy", id, emailo, namee)
	}
	if err := iter.Close(); err != nil {
		fmt.Println(err)
	}

}

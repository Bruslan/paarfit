package data

import (
	"crypto/sha1"
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"time"
)

type User struct {
	User_id    gocql.UUID
	Company    string
	First_name string
	Last_name  string
	Email      string
	Country    string
	Pass       string
	Birthday   time.Time
	Created_at time.Time
}

var cluster *gocql.ClusterConfig

// init cassandra configuration
func init() {
	cluster = gocql.NewCluster("127.0.0.1")
	cluster.Consistency = gocql.One
	cluster.Keyspace = "user_keyspace"
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "ianzndb",
		Password: "Lov3toN8t",
	}
}

// create cassandra session
func db_session() (dbs *gocql.Session, err error) {
	dbs, err = cluster.CreateSession()
	return
}

// hash plaintext with SHA-1
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}

// Create a new session for an existing user
func (user *User) CreateSession(device string) (ses Session, err error) {

	// database session
	dbs, err := db_session()
	if err != nil {
		return
	}
	defer dbs.Close()

	// new uuid
	uuid, err := gocql.RandomUUID()
	if err != nil {
		return
	}
	ses = Session{Session_id: uuid.String(), User_id: user.User_id, Email: user.Email, Created_at: time.Now(), Device: device, Active: true}
	// insert session into sessions
	if err = dbs.Query(`INSERT INTO user_keyspace.sessions (session_id, user_id, email, created_at, device, active) VALUES (?, ?, ?, ?, ?, ?) IF NOT EXISTS`,
		ses.Session_id,
		ses.User_id,
		ses.Email,
		ses.Created_at,
		ses.Device,
		ses.Active).Exec(); err != nil {
		return
	}
	// insert session into session_by_User_id
	if err = dbs.Query(`INSERT INTO user_keyspace.session_by_user_id (user_id, session_id, device) VALUES (?, ?, ?) IF NOT EXISTS`,
		ses.User_id,
		ses.Session_id,
		ses.Device).Exec(); err != nil {
		return
	}
	return
}

// Get the session for an existing user
func (user *User) Session() (session Session, err error) {
	// session = Session{}
	// err = Db.QueryRow("SELECT id, uuid, email, User_id, created_at FROM sessions WHERE User_id = $1", user.Id).
	// 	Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return
}

// Create a new user, save user info into the database
func (user *User) Create() (err error, stmt string) {
	// database session
	dbs, err := db_session()
	if err != nil {
		return
	}
	defer dbs.Close()

	// create new uuid
	uuid, err := gocql.RandomUUID()
	if err != nil {
		stmt = "We apologize. Something went wrong, please try again."
		return
	}
	// insert user into user_by_email table
	if applied, err := dbs.Query(`INSERT INTO user_keyspace.user_by_email (email, pass, user_id) VALUES (?, ?, ?) IF NOT EXISTS`,
		user.Email,
		Encrypt(user.Pass),
		uuid).ScanCAS(&user.Email, &user.Pass, &uuid); err != nil || !applied {
		stmt = "This Email already exists, please try another one."
		return err, stmt
	} else {

		// insert user into users table
		if err = dbs.Query(`INSERT INTO user_keyspace.users (user_id, Company, first_name, last_name, email, pass, country, birthday, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) IF NOT EXISTS`,
			uuid,
			user.Company,
			user.First_name,
			user.Last_name,
			user.Email,
			Encrypt(user.Pass),
			user.Country,
			user.Birthday,
			time.Now()).Exec(); err != nil {
			log.Println("Create user Error:", err)
		} else {
			fmt.Println("Inserted user into users table")
		}
	}
	return
}

// Delete user from database
func (user *User) Delete() (err error) {
	// database session
	dbs, err := db_session()
	if err != nil {
		return
	}
	defer dbs.Close()
	// delete user_by_email
	err = dbs.Query(`DELETE FROM user_keyspace.user_by_email WHERE email=? IF EXISTS`, user.Email).Exec()
	if err != nil {
		return
	}
	// delete users
	return dbs.Query(`DELETE FROM user_keyspace.users WHERE user_id=? IF EXISTS`, user.User_id).Exec()
}

// Delete all sessions of user
func (user *User) DeleteSessions() (err error) {
	// database session
	dbs, err := db_session()
	if err != nil {
		return
	}
	defer dbs.Close()

	var sessid, dev string
	// get all session uuids
	it := dbs.Query(`SELECT session_id, device FROM user_keyspace.session_by_user_id WHERE user_id=?`, user.User_id).Iter()
	defer it.Close()
	for it.Scan(&sessid, &dev) {
		// delete sessions
		err = dbs.Query(`DELETE FROM user_keyspace.sessions WHERE session_id=? IF EXISTS`, sessid).Exec()
		// delete session_by_User_id
		err = dbs.Query(`DELETE FROM user_keyspace.session_by_user_id WHERE user_id=? AND device=? IF EXISTS`, user.User_id, dev).Exec()
	}
	return
}

// // Update user information in the database
// func (user *User) Update() (err error) {
// 	statement := "update users set name = $2, email = $3 where id = $1"
// 	stmt, err := Db.Prepare(statement)
// 	if err != nil {
// 		return
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.Exec(user.Id, user.Name, user.Email)
// 	return
// }

// // Get all users in the database and returns it
// func Users() (users []User, err error) {
// 	rows, err := Db.Query("SELECT id, uuid, name, email, password, created_at FROM users")
// 	if err != nil {
// 		return
// 	}
// 	for rows.Next() {
// 		user := User{}
// 		if err = rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt); err != nil {
// 			return
// 		}
// 		users = append(users, user)
// 	}
// 	rows.Close()
// 	return
// }

// Get a single user given the email
func GetUserByEmail(emailstr string) (user User, err error) {
	// database session
	dbs, err := db_session()
	if err != nil {
		return
	}
	defer dbs.Close()
	err = dbs.Query(`SELECT user_id, email, pass FROM user_keyspace.user_by_email WHERE email=?`, emailstr).Scan(
		&user.User_id,
		&user.Email,
		&user.Pass)
	return
}

// Get a single user given the UUID
func GetUserByUuid(uuid string) (user User, err error) {
	// database session
	dbs, err := db_session()
	if err != nil {
		return
	}
	defer dbs.Close()

	// convert to database uuid type
	go_uuid, err := gocql.ParseUUID(uuid)
	if err != nil {
		return
	}

	err = dbs.Query(`SELECT user_id, email, name, pass FROM user_keyspace.users WHERE user_id=?`, go_uuid).Scan(
		&user.User_id,
		&user.Email,
		&user.Company,
		&user.Pass)
	return
}

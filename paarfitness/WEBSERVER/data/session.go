package data

import (
	"github.com/gocql/gocql"
	"time"
)

type Session struct {
	Session_id string
	User_id    gocql.UUID
	Email      string
	Created_at time.Time
	Device     string
	Active     bool
}

// Check if session is valid in the database ***
func (sess *Session) SessValid() bool {

	// database session
	dbs, err := db_session()
	if err != nil {
		return false
	}
	defer dbs.Close()
	var device string
	if err = dbs.Query(`SELECT active, device FROM user_keyspace.sessions WHERE session_id=?`, sess.Session_id).Scan(
		&sess.Active,
		&device); err != nil {
		return false
	} else {
		if sess.Active && (sess.Device == device) {
			return true
		} else {
			return false
		}
	}
}

// checks if user id session exists with same device but inactive
// if such a session exists delete
func (sess *Session) InactiveExists() error {

	// database session
	dbs, err := db_session()
	if err != nil {
		return err
	}
	defer dbs.Close()

	var dev, sessid string
	it := dbs.Query(`SELECT session_id, device FROM user_keyspace.session_by_user_id WHERE user_id=?`, sess.User_id).Iter()
	defer it.Close()
	if it.NumRows() < 1 {
		return err
	} else if it.NumRows() > 4 {
		for it.Scan(&sessid, &dev) {
			err = dbs.Query(`DELETE FROM user_keyspace.sessions WHERE session_id=? IF EXISTS`, sessid).Exec()
			if err != nil {
				return err
			}
			err = dbs.Query(`DELETE FROM user_keyspace.session_by_user_id WHERE user_id=? IF EXISTS`, sess.User_id).Exec()
			if err != nil {
				return err
			}
		}
	} else {
		for it.Scan(&sessid, &dev) {
			if dev == sess.Device {
				err = dbs.Query(`DELETE FROM user_keyspace.sessions WHERE session_id=? IF EXISTS`, sessid).Exec()
			}
			if dev == sess.Device {
				err = dbs.Query(`DELETE FROM user_keyspace.session_by_user_id WHERE user_id=? AND device=? IF EXISTS`, sess.User_id, dev).Exec()
			}
		}
	}
	return err
}

// ******************************
// modify this function and delete only session which device information
// Delete session from database
func (sess *Session) Delete() (err error) {

	// database session
	dbs, err := db_session()
	if err != nil {
		return
	}
	defer dbs.Close()
	return dbs.Query(`DELETE FROM user_keyspace.sessions WHERE session_id=? IF EXISTS`, sess.Session_id).Exec()
}

func (sess *Session) SetInactive() (err error) {

	// database session
	dbs, err := db_session()
	if err != nil {
		return
	}
	defer dbs.Close()
	return dbs.Query(`UPDATE user_keyspace.sessions SET active=? WHERE session_id=? IF EXISTS`, false, sess.Session_id).Exec()
}

// Get the user from the session
func (session *Session) User() (user User, err error) {

	// database session
	dbs, err := db_session()
	if err != nil {
		return
	}
	defer dbs.Close()
	err = dbs.Query(`SELECT user_id, email FROM user_keyspace.sessions WHERE session_id=?`, session.Session_id).Scan(
		&user.User_id,
		&user.Email)
	return
}

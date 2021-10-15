package users

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/Corwind/utils/dbutils"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id             string `json:"id"`
	Email          string `json:"email"`
	Username       string `json:"username"`
	HashedPassword string `json:"-"`
}

func NewUser(
	email string,
	username string,
	hashedPassword string,
) *User {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatal("User creation failed")
	}
	return &User{
		Id:             id.String(),
		Email:          email,
		Username:       username,
		HashedPassword: hashedPassword,
	}
}

func CheckIfUserExists(db fdb.Database, user User) error {
	t_username := tuple.Tuple{"user", "username", user.Username}
	t_email := tuple.Tuple{"user", "email", user.Email}
	t_id := tuple.Tuple{"user", "id", user.Id}
	_, err := dbutils.DbFetchEntity(db, t_username)
	if err == nil {
		return fmt.Errorf("Internal Server Error")
	}
	_, err = dbutils.DbFetchEntity(db, t_email)
	if err == nil {
		return fmt.Errorf("Internal Server Error")
	}
	_, err = dbutils.DbFetchEntity(db, t_id)
	if err == nil {
		return fmt.Errorf("Internal Server Error")
	}
	return nil
}

func SaveUser(db fdb.Database, user User) (interface{}, error) {
	err := CheckIfUserExists(db, user)
	if err != nil {
		return nil, err
	}

	var user_buffer bytes.Buffer
	enc := gob.NewEncoder(&user_buffer)
	enc.Encode(user)
	t := tuple.Tuple{"user", "id", user.Id}

	var tuple_buffer bytes.Buffer
	enc2 := gob.NewEncoder(&tuple_buffer)
	enc2.Encode(t)
	t2 := tuple.Tuple{"user", "email", user.Email}
	dbutils.DbSaveEntity(db, t2, &tuple_buffer)
	dbutils.DbSaveEntity(db, tuple.Tuple{"user", "username", user.Username}, &tuple_buffer)
	return dbutils.DbSaveEntity(db, t, &user_buffer)
}

func FetchUserById(db fdb.Database, id string) (interface{}, error) {
	ret, err := dbutils.DbFetchEntity(db, tuple.Tuple{"user", "id", id})
	if err != nil {
		return nil, err
	}
	return DecodeUser(ret)
}

func FetchUserByEmail(db fdb.Database, email string) (interface{}, error) {
	ret, err := dbutils.DbFetchEntity(db, tuple.Tuple{"user", "email", email})
	if err != nil {
		return nil, err
	}
	var t tuple.Tuple
	var buffer bytes.Buffer
	buffer.Write(ret.([]byte))
	dec := gob.NewDecoder(&buffer)
	err = dec.Decode(&t)
	if err != nil {
		return nil, err
	}
	ret, err = dbutils.DbFetchEntity(db, t)
	if err != nil {
		return nil, err
	}
	return DecodeUser(ret)
}

func FetchUserByUsername(db fdb.Database, username string) (interface{}, error) {
	ret, err := dbutils.DbFetchEntity(db, tuple.Tuple{"user", "username", username})
	if err != nil {
		return nil, err
	}
	var t tuple.Tuple
	var buffer bytes.Buffer
	buffer.Write(ret.([]byte))
	dec := gob.NewDecoder(&buffer)
	err = dec.Decode(&t)
	if err != nil {
		return nil, err
	}
	ret, err = dbutils.DbFetchEntity(db, t)
	if err != nil {
		return nil, err
	}
	return DecodeUser(ret)
}

func FetchUsers(db fdb.Database) ([]User, error) {
	ret, err := dbutils.DbFetchRange(db, tuple.Tuple{"user", "id"})
	if err != nil {
		return nil, err
	}
	values := make([]User, 0, len(ret))

	for _, value := range ret {
		orga, err := DecodeUser(value)
		if err != nil {
			return nil, err
		}
		values = append(values, orga)
	}

	return values, nil
}

func DecodeUser(ret interface{}) (User, error) {
	var user User
	var b bytes.Buffer
	b.Write(ret.([]byte))
	decoder := gob.NewDecoder(&b)
	decoder.Decode(&user)
	return user, nil
}

func (user_ *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user_.HashedPassword), []byte(password))
}

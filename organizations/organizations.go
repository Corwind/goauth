package organizations

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/Corwind/utils/dbutils"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
	"github.com/google/uuid"
)

type Organization struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func NewOrganization(name string) *Organization {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatal("Account creation failed")
	}
	return &Organization{
		Id:   id.String(),
		Name: name,
	}
}

func SaveOrganization(db fdb.Database, organization Organization) (interface{}, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	enc.Encode(organization)
	t := tuple.Tuple{"organization", organization.Id}
	ret, err := dbutils.DbSaveEntity(db, t, &buffer)
	if err != nil {
		return nil, err
	}

	return DecodeOrganization(ret)
}

func FetchOrganization(db fdb.Database, id string) (interface{}, error) {
	ret, err := dbutils.DbFetchEntity(db, tuple.Tuple{"organization", id})
	if err != nil {
		return nil, err
	}
	return DecodeOrganization(ret)
}

func FetchOrganizations(db fdb.Database) ([]Organization, error) {
	ret, err := dbutils.DbFetchRange(db, tuple.Tuple{"organization"})
	if err != nil {
		return nil, err
	}
	values := make([]Organization, 0, len(ret))

	for _, value := range ret {
		orga, err := DecodeOrganization(value)
		if err != nil {
			return nil, err
		}
		values = append(values, orga)
	}

	return values, nil
}

func DecodeOrganization(ret interface{}) (Organization, error) {
	var organization Organization
	var b bytes.Buffer
	b.Write(ret.([]byte))
	decoder := gob.NewDecoder(&b)
	decoder.Decode(&organization)
	return organization, nil
}

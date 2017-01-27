package postgres

import (
	"errors"
	"fmt"
	"github.com/r3boot/go-ipam/models"
)

func AddPrefix(prefix models.Prefix) error {
	var (
		err error
	)

	err = db.Insert(&prefix)
	if err != nil {
		fmt.Println("AddPrefix: Failed to insert record: " + err.Error())
	}

	return err
}

func DeletePrefix(data interface{}) error {
	var (
		err     error
		prefix  models.Prefix
		network string
	)

	switch data.(type) {
	case models.Prefix:
		network = *data.(models.Prefix).Network
	case string:
		network = data.(string)
	default:
		err = errors.New("DeletePrefix: Received a parameter with an unknown type")
		fmt.Println(err.Error())
		return err
	}

	_, err = db.Model(&prefix).
		Column("network").
		Where("network= ?", network).
		Delete()

	return err
}

func GetPrefixes() models.Prefixes {
	var (
		err     error
		prefixs models.Prefixes
	)

	err = db.Model(&prefixs).Select()
	if err != nil {
		fmt.Println("GetPrefixes: Select failed: " + err.Error())
		return nil
	}

	return prefixs
}

func GetPrefix(data interface{}) models.Prefix {
	var (
		err     error
		prefix  models.Prefix
		network string
	)

	switch data.(type) {
	case string:
		network = data.(string)
	case models.Prefix:
		network = *data.(models.Prefix).Network
	default:
		return models.Prefix{}
	}

	err = db.Model(&prefix).
		Column("network").
		Where("network= ?", network).
		Select()

	if err != nil {
		fmt.Println("GetPrefix: Select failed: " + err.Error())
		return models.Prefix{}
	}

	return prefix
}

func HasPrefix(data interface{}) bool {
	return GetPrefix(data).Network != nil
}

func UpdatePrefix(prefix models.Prefix) error {
	var (
		err error
	)

	_, err = db.Model(&prefix).
		OnConflict("(prefix) DO UPDATE").
		Set("description = ?", prefix.Description).
		Set("username = ?", prefix.Username).
		Insert()

	if err != nil {
		fmt.Println("UpdatePrefix: Failed to upsert Prefix: " + err.Error())
	}

	return err
}

package postgres

import (
	"errors"
	"fmt"
	"github.com/r3boot/go-ipam/models"
)

func AddAsnum(asnum models.Asnum) error {
	var (
		err error
	)

	err = db.Insert(&asnum)
	if err != nil {
		fmt.Println("AddAsnum: Failed to insert record: " + err.Error())
	}

	return err
}

func DeleteAsnum(data interface{}) error {
	var (
		err   error
		asnum models.Asnum
		as    int64
	)

	switch data.(type) {
	case models.Asnum:
		as = *data.(models.Asnum).Asnum
	case int64:
		as = data.(int64)
	default:
		err = errors.New("DeleteAsnum: Received a parameter with an unknown type")
		fmt.Println(err.Error())
		return err
	}

	_, err = db.Model(&asnum).
		Column("asnum").
		Where("asnum= ?", as).
		Delete()

	return err
}

func GetAsnums() models.Asnums {
	var (
		err    error
		asnums models.Asnums
	)

	err = db.Model(&asnums).Select()
	if err != nil {
		fmt.Println("GetAsnums: Select failed: " + err.Error())
		return nil
	}

	return asnums
}

func GetAsnum(data interface{}) models.Asnum {
	var (
		err   error
		asnum models.Asnum
		as    int64
	)

	switch data.(type) {
	case int64:
		as = data.(int64)
	case models.Asnum:
		as = *data.(models.Asnum).Asnum
	default:
		return models.Asnum{}
	}

	err = db.Model(&asnum).
		Where("asnum = ?", as).
		Select()

	if err != nil {
		fmt.Println("GetAsnum: Select failed: " + err.Error())
		return models.Asnum{}
	}

	return asnum
}

func HasAsnum(data interface{}) bool {
	return GetAsnum(data).Asnum != nil
}

func UpdateAsnum(asnum models.Asnum) error {
	var (
		err error
	)

	_, err = db.Model(&asnum).
		OnConflict("(asnum) DO UPDATE").
		Set("description = ?", asnum.Description).
		Set("username = ?", asnum.Username).
		Insert()

	if err != nil {
		fmt.Println("UpdateAsnum: Failed to upsert Asnum: " + err.Error())
	}

	return err
}

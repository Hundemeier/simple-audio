package main

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/graphql-go/graphql"
)

func queryPoolItems(p graphql.ResolveParams) (interface{}, error) {
	//make list of all files in the pool dir:
	files, err := ioutil.ReadDir(poolDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	tmpList := make([]poolItem, 0, len(files))
	for _, file := range files {
		info, err := getPoolItem(file.Name())
		if err != nil {
			continue // if an error occured, we do not add this to the list
		}
		tmpList = append(tmpList, info)
	}
	return tmpList, nil
}

func querySlotItems(p graphql.ResolveParams) (interface{}, error) {
	items := slotMap.getSlotItems()
	slotInt, ok := p.Args["slot"].(int)
	if !ok {
		return items, nil
	}
	//search for the corresponding slot:
	slot := uint(slotInt)
	for _, item := range items {
		if item.Slot == slot {
			return item, nil
		}
	}
	return nil, errors.New("could not find slot. May be unpopulated")
}

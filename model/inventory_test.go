package model

import (
	"fmt"
	"testing"
)

func TestConvertInvModel2InvDB(t *testing.T) {
	inventoryModel := Inventory{
		ID:        1,
		Name:      "testName",
		ProjectID: 1,
		Type:      "host",
	}
	inventoryDB := ConvertInvModel2InvDB(inventoryModel)
	fmt.Printf("%+v\n", inventoryDB)
	if inventoryModel.Inventory == "" {
		fmt.Println("blank string")
	}
}

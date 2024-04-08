package cdc

import (
	"fmt"

	"github.com/emehrkay/cbt/pkg/types"
)

type CDC interface {
	Change(user types.User, action string, entities ...any) error
}

func New() *cdc {
	return &cdc{}
}

type cdc struct{}

func (c *cdc) Change(user types.User, action string, entities ...any) error {
	fmt.Println("-----------------")

	vals := []any{user, action}
	if len(entities) > 0 {
		vals = append(vals, entities...)
		fmt.Printf("\n\t\t[cdc]%s did: %v with: %v\n\n", vals...)
	} else {
		fmt.Printf("\n\t\t[cdc]%s did: %v\n\n", vals...)
	}

	return nil
}

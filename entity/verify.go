package entity

import (
	"context"
	"fmt"
	"reflect"
)

func (e *Entity) Verify() error {

	err := e.DID.Validate()
	if err != nil {
		return fmt.Errorf("entity/verify: %w", err)
	}

	err = e.Doc.Verify()
	if err != nil {
		return fmt.Errorf("entity/verify: %w", err)
	}

	if reflect.TypeOf(e.Messages) != reflect.TypeOf(make(chan *Message)) {
		return fmt.Errorf("entity/verify: Messages is not a channel")
	}

	if reflect.TypeOf(e.Cancel) != reflect.TypeOf((context.CancelFunc)(nil)) {
		return fmt.Errorf("entity/verify: Cancel is not a context.CancelFunc")
	}

	if reflect.TypeOf(e.Doc.Context) != reflect.TypeOf((context.Context)(nil)) {
		return fmt.Errorf("entity/verify: Context is not a context")
	}

	return nil
}

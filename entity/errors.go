package entity

import "errors"

var (
	ErrEntityIsNil   = errors.New("entity: entity is nil")
	ErrEntityIsEmtpy = errors.New("entity: entity is empty")
	ErrTopicIsNil    = errors.New("entity: topic is nil")
)

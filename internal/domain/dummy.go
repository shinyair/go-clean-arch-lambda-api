package domain

import "context"

// Dummy
type Dummy struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SomeAttr string `json:"test_field_name"`
}

// ToKeyDummy
//  @param id
//  @return *Dummy
func ToKeyDummy(id string) *Dummy {
	return &Dummy{
		ID: id,
	}
}

// DummyRepository
type DummyRepository interface {
	GetById(ctx context.Context, id string) (*Dummy, error)
	Insert(ctx context.Context, dummy *Dummy) (*Dummy, error)
	DeleteById(ctx context.Context, id string) (*Dummy, error)
}

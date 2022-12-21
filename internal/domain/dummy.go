package domain

import "context"

// Dummy.
//
//nolint:tagliatelle
type Dummy struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SomeAttr string `json:"test_field_name"`
}

// ToKeyDummy
//
//	@param id
//	@return *Dummy
func ToKeyDummy(id string) *Dummy {
	return &Dummy{
		ID: id,
	}
}

// DummyRepository.
type DummyRepository interface {
	GetByID(ctx context.Context, id string) (*Dummy, error)
	Insert(ctx context.Context, dummy *Dummy) (*Dummy, error)
	DeleteByID(ctx context.Context, id string) (*Dummy, error)
}

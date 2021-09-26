package a

import (
	"context"
	"errors"

	"cloud.google.com/go/datastore"
	"github.com/googleapis/gax-go/v2"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

// Keyvalue stores a key and a value
type Keyvalue map[string]interface{}

// MockDB is used for testing
type MockDB struct {
	// this map isn't being used,
	// but this is how you would test
	// get, put & patch methods.
	DB []Keyvalue
}

// Delete isn't actually used by us,
// the example user service doesn't use it
func (MockDB) Delete(ctx context.Context, key *datastore.Key) error {
	return nil
}

// DeleteMulti isn't actually used by us,
// the example user service doesn't use it
func (MockDB) DeleteMulti(ctx context.Context, keys []*datastore.Key) (err error) {
	return nil
}

// Get isn't actually used by us,
// the example user service doesn't use it
func (MockDB) Get(ctx context.Context, key *datastore.Key, dst interface{}) (err error) {
	return nil
}

// GetAll isn't actually used by us,
// the example user service doesn't use it
func (MockDB) GetAll(ctx context.Context, q *datastore.Query, dst interface{}) (keys []*datastore.Key, err error) {
	return nil, nil
}

// NewTransaction isn't actually used by us,
// the example user service doesn't use it
func (MockDB) NewTransaction(ctx context.Context, opts ...datastore.TransactionOption) (t *datastore.Transaction, err error) {
	return &datastore.Transaction{}, nil
}

// Put isn't actually used by us,
// the example user service doesn't use it
func (MockDB) Put(ctx context.Context, key *datastore.Key, src interface{}) (*datastore.Key, error) {
	createdKey := datastore.Key{
		Kind:      key.Kind,
		ID:        1,
		Namespace: key.Namespace,
	}

	return &createdKey, nil
}

// PutMulti isn't actually used by us,
// the example user service doesn't use it
func (MockDB) PutMulti(ctx context.Context, keys []*datastore.Key, src interface{}) (ret []*datastore.Key, err error) {
	var createdKeys []*datastore.Key
	for id, key := range keys {
		createdKeys = append(createdKeys, &datastore.Key{
			Kind:      key.Kind,
			ID:        int64(id),
			Namespace: key.Namespace,
		})
	}
	return createdKeys, nil
}

// Run isn't actually used by us,
// the example user service doesn't use it
func (MockDB) Run(ctx context.Context, q *datastore.Query) *datastore.Iterator {
	return &datastore.Iterator{}
}

// MockFailDB is used for testing
type MockFailDB struct {
	// this map isn't being used,
	// but this is how you would test
	// get, put & patch methods.
	DB []Keyvalue
}

// Delete isn't actually used by us,
// the example user service doesn't use it
func (MockFailDB) Delete(ctx context.Context, key *datastore.Key) error {
	return errors.New("untrusted error")
}

// DeleteMulti isn't actually used by us,
// the example user service doesn't use it
func (MockFailDB) DeleteMulti(ctx context.Context, keys []*datastore.Key) (err error) {
	return errors.New("untrusted error")
}

// Get isn't actually used by us,
// the example user service doesn't use it
func (MockFailDB) Get(ctx context.Context, key *datastore.Key, dst interface{}) (err error) {
	return errors.New("untrusted error")
}

// GetAll isn't actually used by us,
// the example user service doesn't use it
func (MockFailDB) GetAll(ctx context.Context, q *datastore.Query, dst interface{}) (keys []*datastore.Key, err error) {
	return nil, errors.New("untrusted error")
}

// NewTransaction isn't actually used by us,
// the example user service doesn't use it
func (MockFailDB) NewTransaction(ctx context.Context, opts ...datastore.TransactionOption) (t *datastore.Transaction, err error) {
	return &datastore.Transaction{}, nil
}

// Put isn't actually used by us,
// the example user service doesn't use it
func (MockFailDB) Put(ctx context.Context, key *datastore.Key, src interface{}) (*datastore.Key, error) {
	createdKey := datastore.Key{
		Kind:      key.Kind,
		ID:        1,
		Namespace: key.Namespace,
	}

	return &createdKey, nil
}

// PutMulti isn't actually used by us,
// the example user service doesn't use it
func (MockFailDB) PutMulti(ctx context.Context, keys []*datastore.Key, src interface{}) (ret []*datastore.Key, err error) {
	var createdKeys []*datastore.Key
	for id, key := range keys {
		createdKeys = append(createdKeys, &datastore.Key{
			Kind:      key.Kind,
			ID:        int64(id),
			Namespace: key.Namespace,
		})
	}
	return createdKeys, nil
}

// Run isn't actually used by us,
// the example user service doesn't use it
func (MockFailDB) Run(ctx context.Context, q *datastore.Query) *datastore.Iterator {
	return &datastore.Iterator{}
}

// MockCT is used for testing cloud tasks
type MockCT struct{}

// CreateTask is a method on MockCT
func (MockCT) CreateTask(ctx context.Context, req *taskspb.CreateTaskRequest, opts ...gax.CallOption) (*taskspb.Task, error) {
	return nil, nil
}

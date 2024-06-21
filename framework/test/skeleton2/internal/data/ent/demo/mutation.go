// Code generated by entc, DO NOT EDIT.

package demo

import (
	"context"
	"fmt"
	"sync"

	"gl.king.im/king-lib/framework/test/skeleton2/internal/data/ent/demo/predicate"
	"gl.king.im/king-lib/framework/test/skeleton2/internal/data/ent/demo/sometable"

	"entgo.io/ent"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeSomeTable = "SomeTable"
)

// SomeTableMutation represents an operation that mutates the SomeTable nodes in the graph.
type SomeTableMutation struct {
	config
	op            Op
	typ           string
	id            *uint64
	name          *string
	created_at    *uint64
	addcreated_at *uint64
	updated_at    *uint64
	addupdated_at *uint64
	deleted_at    *uint64
	adddeleted_at *uint64
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*SomeTable, error)
	predicates    []predicate.SomeTable
}

var _ ent.Mutation = (*SomeTableMutation)(nil)

// sometableOption allows management of the mutation configuration using functional options.
type sometableOption func(*SomeTableMutation)

// newSomeTableMutation creates new mutation for the SomeTable entity.
func newSomeTableMutation(c config, op Op, opts ...sometableOption) *SomeTableMutation {
	m := &SomeTableMutation{
		config:        c,
		op:            op,
		typ:           TypeSomeTable,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withSomeTableID sets the ID field of the mutation.
func withSomeTableID(id uint64) sometableOption {
	return func(m *SomeTableMutation) {
		var (
			err   error
			once  sync.Once
			value *SomeTable
		)
		m.oldValue = func(ctx context.Context) (*SomeTable, error) {
			once.Do(func() {
				if m.done {
					err = fmt.Errorf("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().SomeTable.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withSomeTable sets the old SomeTable of the mutation.
func withSomeTable(node *SomeTable) sometableOption {
	return func(m *SomeTableMutation) {
		m.oldValue = func(context.Context) (*SomeTable, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m SomeTableMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m SomeTableMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("demo: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of SomeTable entities.
func (m *SomeTableMutation) SetID(id uint64) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *SomeTableMutation) ID() (id uint64, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetName sets the "name" field.
func (m *SomeTableMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *SomeTableMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the SomeTable entity.
// If the SomeTable object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SomeTableMutation) OldName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldName: %w", err)
	}
	return oldValue.Name, nil
}

// ResetName resets all changes to the "name" field.
func (m *SomeTableMutation) ResetName() {
	m.name = nil
}

// SetCreatedAt sets the "created_at" field.
func (m *SomeTableMutation) SetCreatedAt(u uint64) {
	m.created_at = &u
	m.addcreated_at = nil
}

// CreatedAt returns the value of the "created_at" field in the mutation.
func (m *SomeTableMutation) CreatedAt() (r uint64, exists bool) {
	v := m.created_at
	if v == nil {
		return
	}
	return *v, true
}

// OldCreatedAt returns the old "created_at" field's value of the SomeTable entity.
// If the SomeTable object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SomeTableMutation) OldCreatedAt(ctx context.Context) (v uint64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldCreatedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldCreatedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreatedAt: %w", err)
	}
	return oldValue.CreatedAt, nil
}

// AddCreatedAt adds u to the "created_at" field.
func (m *SomeTableMutation) AddCreatedAt(u uint64) {
	if m.addcreated_at != nil {
		*m.addcreated_at += u
	} else {
		m.addcreated_at = &u
	}
}

// AddedCreatedAt returns the value that was added to the "created_at" field in this mutation.
func (m *SomeTableMutation) AddedCreatedAt() (r uint64, exists bool) {
	v := m.addcreated_at
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreatedAt resets all changes to the "created_at" field.
func (m *SomeTableMutation) ResetCreatedAt() {
	m.created_at = nil
	m.addcreated_at = nil
}

// SetUpdatedAt sets the "updated_at" field.
func (m *SomeTableMutation) SetUpdatedAt(u uint64) {
	m.updated_at = &u
	m.addupdated_at = nil
}

// UpdatedAt returns the value of the "updated_at" field in the mutation.
func (m *SomeTableMutation) UpdatedAt() (r uint64, exists bool) {
	v := m.updated_at
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdatedAt returns the old "updated_at" field's value of the SomeTable entity.
// If the SomeTable object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SomeTableMutation) OldUpdatedAt(ctx context.Context) (v uint64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldUpdatedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldUpdatedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUpdatedAt: %w", err)
	}
	return oldValue.UpdatedAt, nil
}

// AddUpdatedAt adds u to the "updated_at" field.
func (m *SomeTableMutation) AddUpdatedAt(u uint64) {
	if m.addupdated_at != nil {
		*m.addupdated_at += u
	} else {
		m.addupdated_at = &u
	}
}

// AddedUpdatedAt returns the value that was added to the "updated_at" field in this mutation.
func (m *SomeTableMutation) AddedUpdatedAt() (r uint64, exists bool) {
	v := m.addupdated_at
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdatedAt resets all changes to the "updated_at" field.
func (m *SomeTableMutation) ResetUpdatedAt() {
	m.updated_at = nil
	m.addupdated_at = nil
}

// SetDeletedAt sets the "deleted_at" field.
func (m *SomeTableMutation) SetDeletedAt(u uint64) {
	m.deleted_at = &u
	m.adddeleted_at = nil
}

// DeletedAt returns the value of the "deleted_at" field in the mutation.
func (m *SomeTableMutation) DeletedAt() (r uint64, exists bool) {
	v := m.deleted_at
	if v == nil {
		return
	}
	return *v, true
}

// OldDeletedAt returns the old "deleted_at" field's value of the SomeTable entity.
// If the SomeTable object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SomeTableMutation) OldDeletedAt(ctx context.Context) (v uint64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldDeletedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldDeletedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldDeletedAt: %w", err)
	}
	return oldValue.DeletedAt, nil
}

// AddDeletedAt adds u to the "deleted_at" field.
func (m *SomeTableMutation) AddDeletedAt(u uint64) {
	if m.adddeleted_at != nil {
		*m.adddeleted_at += u
	} else {
		m.adddeleted_at = &u
	}
}

// AddedDeletedAt returns the value that was added to the "deleted_at" field in this mutation.
func (m *SomeTableMutation) AddedDeletedAt() (r uint64, exists bool) {
	v := m.adddeleted_at
	if v == nil {
		return
	}
	return *v, true
}

// ResetDeletedAt resets all changes to the "deleted_at" field.
func (m *SomeTableMutation) ResetDeletedAt() {
	m.deleted_at = nil
	m.adddeleted_at = nil
}

// Where appends a list predicates to the SomeTableMutation builder.
func (m *SomeTableMutation) Where(ps ...predicate.SomeTable) {
	m.predicates = append(m.predicates, ps...)
}

// Op returns the operation name.
func (m *SomeTableMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (SomeTable).
func (m *SomeTableMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *SomeTableMutation) Fields() []string {
	fields := make([]string, 0, 4)
	if m.name != nil {
		fields = append(fields, sometable.FieldName)
	}
	if m.created_at != nil {
		fields = append(fields, sometable.FieldCreatedAt)
	}
	if m.updated_at != nil {
		fields = append(fields, sometable.FieldUpdatedAt)
	}
	if m.deleted_at != nil {
		fields = append(fields, sometable.FieldDeletedAt)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *SomeTableMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case sometable.FieldName:
		return m.Name()
	case sometable.FieldCreatedAt:
		return m.CreatedAt()
	case sometable.FieldUpdatedAt:
		return m.UpdatedAt()
	case sometable.FieldDeletedAt:
		return m.DeletedAt()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *SomeTableMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case sometable.FieldName:
		return m.OldName(ctx)
	case sometable.FieldCreatedAt:
		return m.OldCreatedAt(ctx)
	case sometable.FieldUpdatedAt:
		return m.OldUpdatedAt(ctx)
	case sometable.FieldDeletedAt:
		return m.OldDeletedAt(ctx)
	}
	return nil, fmt.Errorf("unknown SomeTable field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *SomeTableMutation) SetField(name string, value ent.Value) error {
	switch name {
	case sometable.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case sometable.FieldCreatedAt:
		v, ok := value.(uint64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreatedAt(v)
		return nil
	case sometable.FieldUpdatedAt:
		v, ok := value.(uint64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdatedAt(v)
		return nil
	case sometable.FieldDeletedAt:
		v, ok := value.(uint64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDeletedAt(v)
		return nil
	}
	return fmt.Errorf("unknown SomeTable field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *SomeTableMutation) AddedFields() []string {
	var fields []string
	if m.addcreated_at != nil {
		fields = append(fields, sometable.FieldCreatedAt)
	}
	if m.addupdated_at != nil {
		fields = append(fields, sometable.FieldUpdatedAt)
	}
	if m.adddeleted_at != nil {
		fields = append(fields, sometable.FieldDeletedAt)
	}
	return fields
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *SomeTableMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case sometable.FieldCreatedAt:
		return m.AddedCreatedAt()
	case sometable.FieldUpdatedAt:
		return m.AddedUpdatedAt()
	case sometable.FieldDeletedAt:
		return m.AddedDeletedAt()
	}
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *SomeTableMutation) AddField(name string, value ent.Value) error {
	switch name {
	case sometable.FieldCreatedAt:
		v, ok := value.(uint64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddCreatedAt(v)
		return nil
	case sometable.FieldUpdatedAt:
		v, ok := value.(uint64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddUpdatedAt(v)
		return nil
	case sometable.FieldDeletedAt:
		v, ok := value.(uint64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddDeletedAt(v)
		return nil
	}
	return fmt.Errorf("unknown SomeTable numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *SomeTableMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *SomeTableMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *SomeTableMutation) ClearField(name string) error {
	return fmt.Errorf("unknown SomeTable nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *SomeTableMutation) ResetField(name string) error {
	switch name {
	case sometable.FieldName:
		m.ResetName()
		return nil
	case sometable.FieldCreatedAt:
		m.ResetCreatedAt()
		return nil
	case sometable.FieldUpdatedAt:
		m.ResetUpdatedAt()
		return nil
	case sometable.FieldDeletedAt:
		m.ResetDeletedAt()
		return nil
	}
	return fmt.Errorf("unknown SomeTable field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *SomeTableMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *SomeTableMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *SomeTableMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *SomeTableMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *SomeTableMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *SomeTableMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *SomeTableMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown SomeTable unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *SomeTableMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown SomeTable edge %s", name)
}

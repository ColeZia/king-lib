// Code generated by entc, DO NOT EDIT.

package demo

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"gl.king.im/king-lib/framework/test/skeleton/internal/data/ent/demo/sometable"
)

// SomeTableCreate is the builder for creating a SomeTable entity.
type SomeTableCreate struct {
	config
	mutation *SomeTableMutation
	hooks    []Hook
}

// SetName sets the "name" field.
func (stc *SomeTableCreate) SetName(s string) *SomeTableCreate {
	stc.mutation.SetName(s)
	return stc
}

// SetNillableName sets the "name" field if the given value is not nil.
func (stc *SomeTableCreate) SetNillableName(s *string) *SomeTableCreate {
	if s != nil {
		stc.SetName(*s)
	}
	return stc
}

// SetCreatedAt sets the "created_at" field.
func (stc *SomeTableCreate) SetCreatedAt(u uint64) *SomeTableCreate {
	stc.mutation.SetCreatedAt(u)
	return stc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (stc *SomeTableCreate) SetNillableCreatedAt(u *uint64) *SomeTableCreate {
	if u != nil {
		stc.SetCreatedAt(*u)
	}
	return stc
}

// SetUpdatedAt sets the "updated_at" field.
func (stc *SomeTableCreate) SetUpdatedAt(u uint64) *SomeTableCreate {
	stc.mutation.SetUpdatedAt(u)
	return stc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (stc *SomeTableCreate) SetNillableUpdatedAt(u *uint64) *SomeTableCreate {
	if u != nil {
		stc.SetUpdatedAt(*u)
	}
	return stc
}

// SetDeletedAt sets the "deleted_at" field.
func (stc *SomeTableCreate) SetDeletedAt(u uint64) *SomeTableCreate {
	stc.mutation.SetDeletedAt(u)
	return stc
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (stc *SomeTableCreate) SetNillableDeletedAt(u *uint64) *SomeTableCreate {
	if u != nil {
		stc.SetDeletedAt(*u)
	}
	return stc
}

// SetID sets the "id" field.
func (stc *SomeTableCreate) SetID(u uint64) *SomeTableCreate {
	stc.mutation.SetID(u)
	return stc
}

// Mutation returns the SomeTableMutation object of the builder.
func (stc *SomeTableCreate) Mutation() *SomeTableMutation {
	return stc.mutation
}

// Save creates the SomeTable in the database.
func (stc *SomeTableCreate) Save(ctx context.Context) (*SomeTable, error) {
	var (
		err  error
		node *SomeTable
	)
	stc.defaults()
	if len(stc.hooks) == 0 {
		if err = stc.check(); err != nil {
			return nil, err
		}
		node, err = stc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SomeTableMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = stc.check(); err != nil {
				return nil, err
			}
			stc.mutation = mutation
			if node, err = stc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(stc.hooks) - 1; i >= 0; i-- {
			if stc.hooks[i] == nil {
				return nil, fmt.Errorf("demo: uninitialized hook (forgotten import demo/runtime?)")
			}
			mut = stc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, stc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (stc *SomeTableCreate) SaveX(ctx context.Context) *SomeTable {
	v, err := stc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (stc *SomeTableCreate) Exec(ctx context.Context) error {
	_, err := stc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (stc *SomeTableCreate) ExecX(ctx context.Context) {
	if err := stc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (stc *SomeTableCreate) defaults() {
	if _, ok := stc.mutation.Name(); !ok {
		v := sometable.DefaultName
		stc.mutation.SetName(v)
	}
	if _, ok := stc.mutation.CreatedAt(); !ok {
		v := sometable.DefaultCreatedAt
		stc.mutation.SetCreatedAt(v)
	}
	if _, ok := stc.mutation.UpdatedAt(); !ok {
		v := sometable.DefaultUpdatedAt
		stc.mutation.SetUpdatedAt(v)
	}
	if _, ok := stc.mutation.DeletedAt(); !ok {
		v := sometable.DefaultDeletedAt
		stc.mutation.SetDeletedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (stc *SomeTableCreate) check() error {
	if _, ok := stc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`demo: missing required field "name"`)}
	}
	if _, ok := stc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`demo: missing required field "created_at"`)}
	}
	if _, ok := stc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`demo: missing required field "updated_at"`)}
	}
	if _, ok := stc.mutation.DeletedAt(); !ok {
		return &ValidationError{Name: "deleted_at", err: errors.New(`demo: missing required field "deleted_at"`)}
	}
	if v, ok := stc.mutation.ID(); ok {
		if err := sometable.IDValidator(v); err != nil {
			return &ValidationError{Name: "id", err: fmt.Errorf(`demo: validator failed for field "id": %w`, err)}
		}
	}
	return nil
}

func (stc *SomeTableCreate) sqlSave(ctx context.Context) (*SomeTable, error) {
	_node, _spec := stc.createSpec()
	if err := sqlgraph.CreateNode(ctx, stc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	if _spec.ID.Value != _node.ID {
		id := _spec.ID.Value.(int64)
		_node.ID = uint64(id)
	}
	return _node, nil
}

func (stc *SomeTableCreate) createSpec() (*SomeTable, *sqlgraph.CreateSpec) {
	var (
		_node = &SomeTable{config: stc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: sometable.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: sometable.FieldID,
			},
		}
	)
	if id, ok := stc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := stc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: sometable.FieldName,
		})
		_node.Name = value
	}
	if value, ok := stc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeUint64,
			Value:  value,
			Column: sometable.FieldCreatedAt,
		})
		_node.CreatedAt = value
	}
	if value, ok := stc.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeUint64,
			Value:  value,
			Column: sometable.FieldUpdatedAt,
		})
		_node.UpdatedAt = value
	}
	if value, ok := stc.mutation.DeletedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeUint64,
			Value:  value,
			Column: sometable.FieldDeletedAt,
		})
		_node.DeletedAt = value
	}
	return _node, _spec
}

// SomeTableCreateBulk is the builder for creating many SomeTable entities in bulk.
type SomeTableCreateBulk struct {
	config
	builders []*SomeTableCreate
}

// Save creates the SomeTable entities in the database.
func (stcb *SomeTableCreateBulk) Save(ctx context.Context) ([]*SomeTable, error) {
	specs := make([]*sqlgraph.CreateSpec, len(stcb.builders))
	nodes := make([]*SomeTable, len(stcb.builders))
	mutators := make([]Mutator, len(stcb.builders))
	for i := range stcb.builders {
		func(i int, root context.Context) {
			builder := stcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SomeTableMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, stcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, stcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{err.Error(), err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				if specs[i].ID.Value != nil && nodes[i].ID == 0 {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = uint64(id)
				}
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, stcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (stcb *SomeTableCreateBulk) SaveX(ctx context.Context) []*SomeTable {
	v, err := stcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (stcb *SomeTableCreateBulk) Exec(ctx context.Context) error {
	_, err := stcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (stcb *SomeTableCreateBulk) ExecX(ctx context.Context) {
	if err := stcb.Exec(ctx); err != nil {
		panic(err)
	}
}

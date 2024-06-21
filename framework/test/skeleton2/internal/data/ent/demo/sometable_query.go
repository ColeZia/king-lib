// Code generated by entc, DO NOT EDIT.

package demo

import (
	"context"
	"errors"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"gl.king.im/king-lib/framework/test/skeleton2/internal/data/ent/demo/predicate"
	"gl.king.im/king-lib/framework/test/skeleton2/internal/data/ent/demo/sometable"
)

// SomeTableQuery is the builder for querying SomeTable entities.
type SomeTableQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.SomeTable
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the SomeTableQuery builder.
func (stq *SomeTableQuery) Where(ps ...predicate.SomeTable) *SomeTableQuery {
	stq.predicates = append(stq.predicates, ps...)
	return stq
}

// Limit adds a limit step to the query.
func (stq *SomeTableQuery) Limit(limit int) *SomeTableQuery {
	stq.limit = &limit
	return stq
}

// Offset adds an offset step to the query.
func (stq *SomeTableQuery) Offset(offset int) *SomeTableQuery {
	stq.offset = &offset
	return stq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (stq *SomeTableQuery) Unique(unique bool) *SomeTableQuery {
	stq.unique = &unique
	return stq
}

// Order adds an order step to the query.
func (stq *SomeTableQuery) Order(o ...OrderFunc) *SomeTableQuery {
	stq.order = append(stq.order, o...)
	return stq
}

// First returns the first SomeTable entity from the query.
// Returns a *NotFoundError when no SomeTable was found.
func (stq *SomeTableQuery) First(ctx context.Context) (*SomeTable, error) {
	nodes, err := stq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{sometable.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (stq *SomeTableQuery) FirstX(ctx context.Context) *SomeTable {
	node, err := stq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first SomeTable ID from the query.
// Returns a *NotFoundError when no SomeTable ID was found.
func (stq *SomeTableQuery) FirstID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = stq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{sometable.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (stq *SomeTableQuery) FirstIDX(ctx context.Context) uint64 {
	id, err := stq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single SomeTable entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when exactly one SomeTable entity is not found.
// Returns a *NotFoundError when no SomeTable entities are found.
func (stq *SomeTableQuery) Only(ctx context.Context) (*SomeTable, error) {
	nodes, err := stq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{sometable.Label}
	default:
		return nil, &NotSingularError{sometable.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (stq *SomeTableQuery) OnlyX(ctx context.Context) *SomeTable {
	node, err := stq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only SomeTable ID in the query.
// Returns a *NotSingularError when exactly one SomeTable ID is not found.
// Returns a *NotFoundError when no entities are found.
func (stq *SomeTableQuery) OnlyID(ctx context.Context) (id uint64, err error) {
	var ids []uint64
	if ids, err = stq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{sometable.Label}
	default:
		err = &NotSingularError{sometable.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (stq *SomeTableQuery) OnlyIDX(ctx context.Context) uint64 {
	id, err := stq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of SomeTables.
func (stq *SomeTableQuery) All(ctx context.Context) ([]*SomeTable, error) {
	if err := stq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return stq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (stq *SomeTableQuery) AllX(ctx context.Context) []*SomeTable {
	nodes, err := stq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of SomeTable IDs.
func (stq *SomeTableQuery) IDs(ctx context.Context) ([]uint64, error) {
	var ids []uint64
	if err := stq.Select(sometable.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (stq *SomeTableQuery) IDsX(ctx context.Context) []uint64 {
	ids, err := stq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (stq *SomeTableQuery) Count(ctx context.Context) (int, error) {
	if err := stq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return stq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (stq *SomeTableQuery) CountX(ctx context.Context) int {
	count, err := stq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (stq *SomeTableQuery) Exist(ctx context.Context) (bool, error) {
	if err := stq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return stq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (stq *SomeTableQuery) ExistX(ctx context.Context) bool {
	exist, err := stq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the SomeTableQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (stq *SomeTableQuery) Clone() *SomeTableQuery {
	if stq == nil {
		return nil
	}
	return &SomeTableQuery{
		config:     stq.config,
		limit:      stq.limit,
		offset:     stq.offset,
		order:      append([]OrderFunc{}, stq.order...),
		predicates: append([]predicate.SomeTable{}, stq.predicates...),
		// clone intermediate query.
		sql:  stq.sql.Clone(),
		path: stq.path,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.SomeTable.Query().
//		GroupBy(sometable.FieldName).
//		Aggregate(demo.Count()).
//		Scan(ctx, &v)
//
func (stq *SomeTableQuery) GroupBy(field string, fields ...string) *SomeTableGroupBy {
	group := &SomeTableGroupBy{config: stq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := stq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return stq.sqlQuery(ctx), nil
	}
	return group
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//	}
//
//	client.SomeTable.Query().
//		Select(sometable.FieldName).
//		Scan(ctx, &v)
//
func (stq *SomeTableQuery) Select(fields ...string) *SomeTableSelect {
	stq.fields = append(stq.fields, fields...)
	return &SomeTableSelect{SomeTableQuery: stq}
}

func (stq *SomeTableQuery) prepareQuery(ctx context.Context) error {
	for _, f := range stq.fields {
		if !sometable.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("demo: invalid field %q for query", f)}
		}
	}
	if stq.path != nil {
		prev, err := stq.path(ctx)
		if err != nil {
			return err
		}
		stq.sql = prev
	}
	return nil
}

func (stq *SomeTableQuery) sqlAll(ctx context.Context) ([]*SomeTable, error) {
	var (
		nodes = []*SomeTable{}
		_spec = stq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		node := &SomeTable{config: stq.config}
		nodes = append(nodes, node)
		return node.scanValues(columns)
	}
	_spec.Assign = func(columns []string, values []interface{}) error {
		if len(nodes) == 0 {
			return fmt.Errorf("demo: Assign called without calling ScanValues")
		}
		node := nodes[len(nodes)-1]
		return node.assignValues(columns, values)
	}
	if err := sqlgraph.QueryNodes(ctx, stq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (stq *SomeTableQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := stq.querySpec()
	return sqlgraph.CountNodes(ctx, stq.driver, _spec)
}

func (stq *SomeTableQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := stq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("demo: check existence: %w", err)
	}
	return n > 0, nil
}

func (stq *SomeTableQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   sometable.Table,
			Columns: sometable.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: sometable.FieldID,
			},
		},
		From:   stq.sql,
		Unique: true,
	}
	if unique := stq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := stq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, sometable.FieldID)
		for i := range fields {
			if fields[i] != sometable.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := stq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := stq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := stq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := stq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (stq *SomeTableQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(stq.driver.Dialect())
	t1 := builder.Table(sometable.Table)
	columns := stq.fields
	if len(columns) == 0 {
		columns = sometable.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if stq.sql != nil {
		selector = stq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	for _, p := range stq.predicates {
		p(selector)
	}
	for _, p := range stq.order {
		p(selector)
	}
	if offset := stq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := stq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// SomeTableGroupBy is the group-by builder for SomeTable entities.
type SomeTableGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (stgb *SomeTableGroupBy) Aggregate(fns ...AggregateFunc) *SomeTableGroupBy {
	stgb.fns = append(stgb.fns, fns...)
	return stgb
}

// Scan applies the group-by query and scans the result into the given value.
func (stgb *SomeTableGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := stgb.path(ctx)
	if err != nil {
		return err
	}
	stgb.sql = query
	return stgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (stgb *SomeTableGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := stgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by.
// It is only allowed when executing a group-by query with one field.
func (stgb *SomeTableGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(stgb.fields) > 1 {
		return nil, errors.New("demo: SomeTableGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := stgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (stgb *SomeTableGroupBy) StringsX(ctx context.Context) []string {
	v, err := stgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (stgb *SomeTableGroupBy) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = stgb.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{sometable.Label}
	default:
		err = fmt.Errorf("demo: SomeTableGroupBy.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (stgb *SomeTableGroupBy) StringX(ctx context.Context) string {
	v, err := stgb.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by.
// It is only allowed when executing a group-by query with one field.
func (stgb *SomeTableGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(stgb.fields) > 1 {
		return nil, errors.New("demo: SomeTableGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := stgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (stgb *SomeTableGroupBy) IntsX(ctx context.Context) []int {
	v, err := stgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (stgb *SomeTableGroupBy) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = stgb.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{sometable.Label}
	default:
		err = fmt.Errorf("demo: SomeTableGroupBy.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (stgb *SomeTableGroupBy) IntX(ctx context.Context) int {
	v, err := stgb.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by.
// It is only allowed when executing a group-by query with one field.
func (stgb *SomeTableGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(stgb.fields) > 1 {
		return nil, errors.New("demo: SomeTableGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := stgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (stgb *SomeTableGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := stgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (stgb *SomeTableGroupBy) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = stgb.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{sometable.Label}
	default:
		err = fmt.Errorf("demo: SomeTableGroupBy.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (stgb *SomeTableGroupBy) Float64X(ctx context.Context) float64 {
	v, err := stgb.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by.
// It is only allowed when executing a group-by query with one field.
func (stgb *SomeTableGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(stgb.fields) > 1 {
		return nil, errors.New("demo: SomeTableGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := stgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (stgb *SomeTableGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := stgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (stgb *SomeTableGroupBy) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = stgb.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{sometable.Label}
	default:
		err = fmt.Errorf("demo: SomeTableGroupBy.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (stgb *SomeTableGroupBy) BoolX(ctx context.Context) bool {
	v, err := stgb.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (stgb *SomeTableGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range stgb.fields {
		if !sometable.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := stgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := stgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (stgb *SomeTableGroupBy) sqlQuery() *sql.Selector {
	selector := stgb.sql.Select()
	aggregation := make([]string, 0, len(stgb.fns))
	for _, fn := range stgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(stgb.fields)+len(stgb.fns))
		for _, f := range stgb.fields {
			columns = append(columns, selector.C(f))
		}
		for _, c := range aggregation {
			columns = append(columns, c)
		}
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(stgb.fields...)...)
}

// SomeTableSelect is the builder for selecting fields of SomeTable entities.
type SomeTableSelect struct {
	*SomeTableQuery
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (sts *SomeTableSelect) Scan(ctx context.Context, v interface{}) error {
	if err := sts.prepareQuery(ctx); err != nil {
		return err
	}
	sts.sql = sts.SomeTableQuery.sqlQuery(ctx)
	return sts.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (sts *SomeTableSelect) ScanX(ctx context.Context, v interface{}) {
	if err := sts.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from a selector. It is only allowed when selecting one field.
func (sts *SomeTableSelect) Strings(ctx context.Context) ([]string, error) {
	if len(sts.fields) > 1 {
		return nil, errors.New("demo: SomeTableSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := sts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (sts *SomeTableSelect) StringsX(ctx context.Context) []string {
	v, err := sts.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a selector. It is only allowed when selecting one field.
func (sts *SomeTableSelect) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = sts.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{sometable.Label}
	default:
		err = fmt.Errorf("demo: SomeTableSelect.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (sts *SomeTableSelect) StringX(ctx context.Context) string {
	v, err := sts.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from a selector. It is only allowed when selecting one field.
func (sts *SomeTableSelect) Ints(ctx context.Context) ([]int, error) {
	if len(sts.fields) > 1 {
		return nil, errors.New("demo: SomeTableSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := sts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (sts *SomeTableSelect) IntsX(ctx context.Context) []int {
	v, err := sts.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a selector. It is only allowed when selecting one field.
func (sts *SomeTableSelect) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = sts.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{sometable.Label}
	default:
		err = fmt.Errorf("demo: SomeTableSelect.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (sts *SomeTableSelect) IntX(ctx context.Context) int {
	v, err := sts.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from a selector. It is only allowed when selecting one field.
func (sts *SomeTableSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(sts.fields) > 1 {
		return nil, errors.New("demo: SomeTableSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := sts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (sts *SomeTableSelect) Float64sX(ctx context.Context) []float64 {
	v, err := sts.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a selector. It is only allowed when selecting one field.
func (sts *SomeTableSelect) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = sts.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{sometable.Label}
	default:
		err = fmt.Errorf("demo: SomeTableSelect.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (sts *SomeTableSelect) Float64X(ctx context.Context) float64 {
	v, err := sts.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from a selector. It is only allowed when selecting one field.
func (sts *SomeTableSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(sts.fields) > 1 {
		return nil, errors.New("demo: SomeTableSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := sts.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (sts *SomeTableSelect) BoolsX(ctx context.Context) []bool {
	v, err := sts.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a selector. It is only allowed when selecting one field.
func (sts *SomeTableSelect) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = sts.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{sometable.Label}
	default:
		err = fmt.Errorf("demo: SomeTableSelect.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (sts *SomeTableSelect) BoolX(ctx context.Context) bool {
	v, err := sts.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (sts *SomeTableSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := sts.sql.Query()
	if err := sts.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

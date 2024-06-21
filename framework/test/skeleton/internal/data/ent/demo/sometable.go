// Code generated by entc, DO NOT EDIT.

package demo

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"gl.king.im/king-lib/framework/test/skeleton/internal/data/ent/demo/sometable"
)

// SomeTable is the model entity for the SomeTable schema.
type SomeTable struct {
	config `json:"-"`
	// ID of the ent.
	ID uint64 `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt uint64 `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt uint64 `json:"updated_at,omitempty"`
	// DeletedAt holds the value of the "deleted_at" field.
	DeletedAt uint64 `json:"deleted_at,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*SomeTable) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case sometable.FieldID, sometable.FieldCreatedAt, sometable.FieldUpdatedAt, sometable.FieldDeletedAt:
			values[i] = new(sql.NullInt64)
		case sometable.FieldName:
			values[i] = new(sql.NullString)
		default:
			return nil, fmt.Errorf("unexpected column %q for type SomeTable", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the SomeTable fields.
func (st *SomeTable) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case sometable.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			st.ID = uint64(value.Int64)
		case sometable.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				st.Name = value.String
			}
		case sometable.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				st.CreatedAt = uint64(value.Int64)
			}
		case sometable.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				st.UpdatedAt = uint64(value.Int64)
			}
		case sometable.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				st.DeletedAt = uint64(value.Int64)
			}
		}
	}
	return nil
}

// Update returns a builder for updating this SomeTable.
// Note that you need to call SomeTable.Unwrap() before calling this method if this SomeTable
// was returned from a transaction, and the transaction was committed or rolled back.
func (st *SomeTable) Update() *SomeTableUpdateOne {
	return (&SomeTableClient{config: st.config}).UpdateOne(st)
}

// Unwrap unwraps the SomeTable entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (st *SomeTable) Unwrap() *SomeTable {
	tx, ok := st.config.driver.(*txDriver)
	if !ok {
		panic("demo: SomeTable is not a transactional entity")
	}
	st.config.driver = tx.drv
	return st
}

// String implements the fmt.Stringer.
func (st *SomeTable) String() string {
	var builder strings.Builder
	builder.WriteString("SomeTable(")
	builder.WriteString(fmt.Sprintf("id=%v", st.ID))
	builder.WriteString(", name=")
	builder.WriteString(st.Name)
	builder.WriteString(", created_at=")
	builder.WriteString(fmt.Sprintf("%v", st.CreatedAt))
	builder.WriteString(", updated_at=")
	builder.WriteString(fmt.Sprintf("%v", st.UpdatedAt))
	builder.WriteString(", deleted_at=")
	builder.WriteString(fmt.Sprintf("%v", st.DeletedAt))
	builder.WriteByte(')')
	return builder.String()
}

// SomeTables is a parsable slice of SomeTable.
type SomeTables []*SomeTable

func (st SomeTables) config(cfg config) {
	for _i := range st {
		st[_i].config = cfg
	}
}

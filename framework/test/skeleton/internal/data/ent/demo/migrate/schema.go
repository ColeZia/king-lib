// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// SomeTablesColumns holds the columns for the "some_tables" table.
	SomeTablesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUint64, Increment: true},
		{Name: "name", Type: field.TypeString, Default: ""},
		{Name: "created_at", Type: field.TypeUint64, Default: 0},
		{Name: "updated_at", Type: field.TypeUint64, Default: 0},
		{Name: "deleted_at", Type: field.TypeUint64, Default: 0},
	}
	// SomeTablesTable holds the schema information for the "some_tables" table.
	SomeTablesTable = &schema.Table{
		Name:       "some_tables",
		Columns:    SomeTablesColumns,
		PrimaryKey: []*schema.Column{SomeTablesColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		SomeTablesTable,
	}
)

func init() {
}

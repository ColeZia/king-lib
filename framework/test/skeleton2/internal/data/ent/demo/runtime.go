// Code generated by entc, DO NOT EDIT.

package demo

import (
	"gl.king.im/king-lib/framework/test/skeleton2/internal/data/ent/demo/schema"
	"gl.king.im/king-lib/framework/test/skeleton2/internal/data/ent/demo/sometable"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	sometableFields := schema.SomeTable{}.Fields()
	_ = sometableFields
	// sometableDescName is the schema descriptor for name field.
	sometableDescName := sometableFields[1].Descriptor()
	// sometable.DefaultName holds the default value on creation for the name field.
	sometable.DefaultName = sometableDescName.Default.(string)
	// sometableDescCreatedAt is the schema descriptor for created_at field.
	sometableDescCreatedAt := sometableFields[2].Descriptor()
	// sometable.DefaultCreatedAt holds the default value on creation for the created_at field.
	sometable.DefaultCreatedAt = sometableDescCreatedAt.Default.(uint64)
	// sometableDescUpdatedAt is the schema descriptor for updated_at field.
	sometableDescUpdatedAt := sometableFields[3].Descriptor()
	// sometable.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	sometable.DefaultUpdatedAt = sometableDescUpdatedAt.Default.(uint64)
	// sometableDescDeletedAt is the schema descriptor for deleted_at field.
	sometableDescDeletedAt := sometableFields[4].Descriptor()
	// sometable.DefaultDeletedAt holds the default value on creation for the deleted_at field.
	sometable.DefaultDeletedAt = sometableDescDeletedAt.Default.(uint64)
	// sometableDescID is the schema descriptor for id field.
	sometableDescID := sometableFields[0].Descriptor()
	// sometable.IDValidator is a validator for the "id" field. It is called by the builders before save.
	sometable.IDValidator = sometableDescID.Validators[0].(func(uint64) error)
}

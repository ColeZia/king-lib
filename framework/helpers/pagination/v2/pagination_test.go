package pagination

import "testing"

func TestPagination_CalOffsetLimit(t *testing.T) {
	type fields struct {
		opts       *paginationOpts
		SimplePage bool
		Limit      int
		Page       int
		PageSize   int
		Sort       string
		TotalRows  int64
		TotalPages int
		Rows       interface{}
	}
	type args struct {
		page     int
		pageSize int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOffset int
		wantLimit  int
	}{
		// TODO: Add test cases.
		{
			fields: fields{
				opts: &paginationOpts{},
			},
			args:       args{page: 0, pageSize: 10},
			wantOffset: 0,
			wantLimit:  10,
		},
		{
			fields: fields{
				opts: &paginationOpts{},
			},
			args:       args{page: 1, pageSize: 0},
			wantOffset: 0,
			wantLimit:  20,
		},
		{
			fields: fields{
				opts: &paginationOpts{},
			},
			args:       args{page: 2, pageSize: 0},
			wantOffset: 20,
			wantLimit:  20,
		},
		{
			fields: fields{
				opts: &paginationOpts{},
			},
			args:       args{page: 10, pageSize: 30},
			wantOffset: 270,
			wantLimit:  30,
		},
		{
			fields: fields{
				opts: &paginationOpts{},
			},
			args:       args{page: 10, pageSize: 0},
			wantOffset: 180,
			wantLimit:  20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPagination()
			gotOffset, gotLimit := p.CalOffsetLimit(tt.args.page, tt.args.pageSize)
			if gotOffset != tt.wantOffset {
				t.Errorf("Pagination.CalOffsetLimit() gotOffset = %v, want %v", gotOffset, tt.wantOffset)
			}
			if gotLimit != tt.wantLimit {
				t.Errorf("Pagination.CalOffsetLimit() gotLimit = %v, want %v", gotLimit, tt.wantLimit)
			}
		})
	}
}

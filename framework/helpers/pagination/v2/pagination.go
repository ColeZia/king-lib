package pagination

type paginationOpts struct {
	defaultPageSize int
	maxPageSize     int
	minPageSize     int
}
type Pagination interface {
	CalOffsetLimit(page, pageSize int) (offset int, limit int)
}

type pagination struct {
	opts *paginationOpts
	//Db         *gorm.DB
	SimplePage bool
	Limit      int `json:"limit,omitempty;query:limit"`
	Page       int `json:"page,omitempty;query:page"`
	PageSize   int
	Sort       string      `json:"sort,omitempty;query:sort"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Rows       interface{} `json:"rows"`
}

func (p *pagination) CalOffsetLimit(page, pageSize int) (offset int, limit int) {
	if pageSize < p.opts.minPageSize {
		pageSize = p.opts.minPageSize
	}

	if p.opts.maxPageSize != 0 && pageSize > p.opts.maxPageSize {
		pageSize = p.opts.maxPageSize
	}

	limit = pageSize
	if pageSize == 0 {
		limit = p.opts.defaultPageSize
	}

	if page > 1 {
		offset = (page - 1) * limit
	}

	return offset, limit
}

type paginationOp func(o *paginationOpts)

func WithDefaultPageSize(defaultPageSize int) paginationOp {
	return func(o *paginationOpts) {
		o.defaultPageSize = defaultPageSize
	}
}

func WithMaxPageSize(maxPageSize int) paginationOp {
	return func(o *paginationOpts) {
		o.maxPageSize = maxPageSize
	}
}

func NewPagination(opts ...paginationOp) *pagination {
	o := &paginationOpts{}
	for _, v := range opts {
		v(o)
	}

	if o.defaultPageSize == 0 {
		o.defaultPageSize = 20
	}

	p := &pagination{
		opts: o,
	}

	return p
}

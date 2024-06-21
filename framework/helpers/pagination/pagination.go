package pagination

import (
	"fmt"

	dpGorm "git.e.coding.king.cloud/dev/data_platform/gorm"
	"gorm.io/gorm"
)

type Pagination struct {
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

func (p *Pagination) Scope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		if p.Page <= 0 {
			p.Page = 1
		}

		switch {
		//case pageSize > 100:
		//	pageSize = 100
		case p.PageSize <= 0:
			p.PageSize = 10
		}

		offset := (p.Page - 1) * p.PageSize
		fmt.Println("Paginate scope::", offset, p.PageSize)

		//return db.Offset(offset).Limit(p.PageSize).Order(p.Sort)
		scopedDb := db.Offset(offset).Limit(p.PageSize)

		if len(p.Sort) > 0 {
			scopedDb = scopedDb.Order(p.Sort)
		}

		return scopedDb
	}
}

func (p *Pagination) Paginate(db *gorm.DB) *gorm.DB {
	if !p.SimplePage {
		db.Count(&p.TotalRows)
	}

	return db.Scopes(p.Scope())
}

func (p *Pagination) ScopeDp() func(db *dpGorm.DB) *dpGorm.DB {
	return func(db *dpGorm.DB) *dpGorm.DB {

		if p.Page <= 0 {
			p.Page = 1
		}

		switch {
		//case pageSize > 100:
		//	pageSize = 100
		case p.PageSize <= 0:
			p.PageSize = 10
		}

		offset := (p.Page - 1) * p.PageSize
		fmt.Println("Paginate scope::", offset, p.PageSize)

		//return db.Offset(offset).Limit(p.PageSize).Order(p.Sort)
		scopedDb := db.Offset(offset).Limit(p.PageSize)

		if len(p.Sort) > 0 {
			scopedDb = scopedDb.Order(p.Sort)
		}

		return scopedDb
	}
}

func (p *Pagination) PaginateDp(db *dpGorm.DB) *dpGorm.DB {
	if !p.SimplePage {
		db.Count(&p.TotalRows)
	}

	return db.Scopes(p.ScopeDp())
}

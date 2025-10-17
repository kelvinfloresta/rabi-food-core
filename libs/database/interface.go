package database

import "math"

type Database interface {
	Start() error
	Stop() error
}

type PaginateInput struct {
	Page     int
	PageSize int
}

func (p *PaginateInput) CalcMaxPages(count int64) int {
	total := float64(count) / float64(p.PageSize)
	return int(math.Ceil(total))
}

func (p *PaginateInput) Offset() int {
	return p.Page * p.PageSize
}

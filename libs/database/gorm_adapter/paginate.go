package gorm_adapter

import (
	"rabi-food-core/libs/database"
	"sync"

	"gorm.io/gorm"
)

func Paginate(query *gorm.DB, count *int64, data any, paginate database.PaginateInput) error {
	var wg sync.WaitGroup
	countErr := error(nil)
	countSession := query.Session(&gorm.Session{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		result := countSession.Count(count)
		countErr = result.Error
	}()

	paginateErr := error(nil)
	paginateSession := query.Session(&gorm.Session{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		result := paginateSession.Limit(paginate.PageSize).Offset(paginate.Offset()).Find(data)
		paginateErr = result.Error
	}()

	wg.Wait()
	if countErr != nil {
		return countErr
	}

	if paginateErr != nil {
		return paginateErr
	}

	return nil
}

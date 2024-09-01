package usecase_test

import (
	"github.com/a-x-a/gophkeeper/internal/keeperctl/entity"
	"github.com/a-x-a/gophkeeper/internal/util/gophtest"
)

func newTestKey() entity.Key {
	return entity.NewKey(gophtest.Username, gophtest.Password)
}

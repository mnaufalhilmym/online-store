package pg

import "gorm.io/gorm"

func IsErrRecordNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}

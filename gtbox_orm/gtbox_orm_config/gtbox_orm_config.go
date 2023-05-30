package gtbox_orm_config

import "fmt"

type GTORMTimeZone int64

const (
	GTORMTimeZoneUTC      GTORMTimeZone = iota // UTC TimeZone 默认值
	GTORMTimeZoneShangHai                      // ShangHai TimeZone
)

func (timeZone GTORMTimeZone) String() string {
	switch timeZone {
	case GTORMTimeZoneShangHai:
		return fmt.Sprintf("Asia%sShanghai", "%2F")
	case GTORMTimeZoneUTC:
		return "UTC"
	default:
		return "UTC"
	}
}

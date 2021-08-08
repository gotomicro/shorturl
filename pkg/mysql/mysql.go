package mysql

type Shorturl struct {
	ID        int64  // id 号
	OriginUrl string // 原始URL
	Code      string // 短网址的code码
	CallCnt   int64  // 调用次数
	Ctime     int64  // 创建时间
	Utime     int64  // 更新时间
}

func (Shorturl) TableName() string {
	return "shorturl"
}

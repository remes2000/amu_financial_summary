package regexp

type Regexp struct {
	Id         uint   `json:"id" gorm:"primaryKey"`
	Content    string `json:"content" binding:"required" gorm:"notNull"`
	CategoryID uint   `json:"-"`
}

type CreateRegexp struct {
	Content string `json:"content" binding:"required"`
}

func (c CreateRegexp) GetRegexp() Regexp {
	return Regexp{Content: c.Content}
}

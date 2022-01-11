package regexp

import goRegexp "regexp"

type Regexp struct {
	Id         uint   `json:"id" gorm:"primaryKey"`
	Content    string `json:"content" binding:"required" gorm:"notNull"`
	CategoryID uint   `json:"-"`
}

func (c Regexp) Matches(title string) bool {
	r, err := goRegexp.Compile(c.Content)
	if err != nil {
		return false
	}
	return r.MatchString(title)
}

type CreateRegexp struct {
	Content string `json:"content" binding:"required"`
}

func (c CreateRegexp) GetRegexp() Regexp {
	return Regexp{Content: c.Content}
}

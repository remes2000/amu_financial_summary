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

type RegexpBackup struct {
	Id         uint   `json:"id" binding:"required"`
	Content    string `json:"content" binding:"required"`
	CategoryID uint   `json:"categoryId" binding:"required"`
}

func (rb RegexpBackup) ToRegexp() Regexp {
	return Regexp{Id: rb.Id, Content: rb.Content, CategoryID: rb.CategoryID}
}

func (rb *RegexpBackup) FromRegexp(regexp Regexp) {
	rb.Id = regexp.Id
	rb.Content = regexp.Content
	rb.CategoryID = regexp.CategoryID
}

type CreateRegexp struct {
	Content string `json:"content" binding:"required"`
}

func (c CreateRegexp) GetRegexp() Regexp {
	return Regexp{Content: c.Content}
}

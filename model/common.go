package model

// no idea why it's named this, since it's used for things that aren't localized as well
// maybe localized is used in the sense of distance, not country
type LocalizedText struct {
	DotUnderText string `json:"dot_under_text" xorm:"pk"`
}

func (lt *LocalizedText) FromDB(data []byte) error {
	lt.DotUnderText = string(data)
	return nil
}
func (lt *LocalizedText) ToDB() ([]byte, error) {
	return []byte(lt.DotUnderText), nil
}

// yep this is the real name in the code
// Strucktur is structure in deutsch but it might as well be engrish
type TextureStruktur struct {
	V string `json:"v" xorm:"pk"`
}

func (ts *TextureStruktur) FromDB(data []byte) error {
	ts.V = string(data)
	return nil
}
func (ts *TextureStruktur) ToDB() ([]byte, error) {
	return []byte(ts.V), nil
}

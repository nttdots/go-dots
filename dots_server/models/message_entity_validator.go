package models

type MessageEntityValidator interface {
	Validate(m MessageEntity, c *Customer) (ret bool)
}

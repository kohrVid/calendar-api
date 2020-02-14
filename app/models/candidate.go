package models

import "github.com/kohrVid/calendar-api/app/serializers"

type Candidate struct {
	Id        int64
	FirstName string
	LastName  string
	Email     string
}

func (c *Candidate) ToSerializer() serializers.Candidate {
	return serializers.Candidate{
		Id:        int(c.Id),
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Email:     c.Email,
	}
}

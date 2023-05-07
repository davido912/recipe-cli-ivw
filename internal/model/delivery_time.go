package model

import (
	"time"
)

const (
	timeLayout = "3PM"
)

type DeliveryTime struct {
	time.Time
	raw string
}

func NewDeliveryTime(input string) (*DeliveryTime, error) {
	d := DeliveryTime{raw: input}
	err := d.parse(input)
	if err != nil {
		return nil, err
	}
	return &d, err
}

// parse converts a string of time with PM or AM to a time object (e.g. 3PM)
func (dt *DeliveryTime) parse(input string) error {
	t, err := time.Parse(timeLayout, input)
	if err != nil {
		return err
	}
	dt.Time = t
	return nil
}

// Raw returns the raw string passed to instantiate DeliveryTime
func (dt *DeliveryTime) Raw() string {
	return dt.raw
}

// InclusiveAfter whether time instance dt is after or equal t
func (dt *DeliveryTime) InclusiveAfter(t *DeliveryTime) bool {
	return dt.After(t.Time) || dt.Equal(t.Time)
}

// InclusiveBefore whether time instance dt is before or equal t
func (dt *DeliveryTime) InclusiveBefore(t *DeliveryTime) bool {
	return dt.Before(t.Time) || dt.Equal(t.Time)
}

// InclusiveBetween whether time instance dt between or equal start and end
func (dt *DeliveryTime) InclusiveBetween(start, end *DeliveryTime) bool {
	return dt.InclusiveAfter(start) && dt.InclusiveBefore(end)
}

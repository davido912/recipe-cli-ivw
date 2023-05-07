package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDeliveryTime_parse(t *testing.T) {
	input := "4PM"

	dt := DeliveryTime{}
	err := dt.parse(input)
	assert.Nil(t, err)

	expected := time.Date(0000, 01, 01, 16, 00, 00, 00, time.UTC)
	assert.Equal(t, expected, dt.Time)

}

func TestDeliveryTime_InclusiveAfter(t *testing.T) {

	givenTime, _ := NewDeliveryTime("1PM")

	tcs := []struct {
		name      string
		givenTime *DeliveryTime
		checkTime string
		want      bool
	}{
		{
			name:      "time not after",
			givenTime: givenTime,
			checkTime: "3PM",
			want:      false,
		},
		{
			name:      "time after",
			givenTime: givenTime,
			checkTime: "11AM",
			want:      true,
		},
		{
			name:      "times equal",
			givenTime: givenTime,
			checkTime: givenTime.Raw(),
			want:      true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			checkTime, _ := NewDeliveryTime(tc.checkTime)
			got := tc.givenTime.InclusiveAfter(checkTime)
			assert.Equal(t, tc.want, got)

		})
	}

}

func TestDeliveryTime_InclusiveBefore(t *testing.T) {

	givenTime, _ := NewDeliveryTime("1PM")

	tcs := []struct {
		name      string
		givenTime *DeliveryTime
		checkTime string
		want      bool
	}{
		{
			name:      "time",
			givenTime: givenTime,
			checkTime: "3PM",
			want:      true,
		},
		{
			name:      "time not before",
			givenTime: givenTime,
			checkTime: "11AM",
			want:      false,
		},
		{
			name:      "times equal",
			givenTime: givenTime,
			checkTime: givenTime.Raw(),
			want:      true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			checkTime, _ := NewDeliveryTime(tc.checkTime)
			got := tc.givenTime.InclusiveBefore(checkTime)
			assert.Equal(t, tc.want, got)

		})
	}
}

func TestDeliveryTime_InclusiveBetween(t *testing.T) {

	startTime, _ := NewDeliveryTime("1PM")
	endTime, _ := NewDeliveryTime("10PM")

	tcs := []struct {
		name      string
		startTime *DeliveryTime
		endTime   *DeliveryTime
		checkTime string
		want      bool
	}{
		{
			name:      "time between",
			startTime: startTime,
			endTime:   endTime,
			checkTime: "3PM",
			want:      true,
		},
		{
			name:      "time not between",
			startTime: startTime,
			endTime:   endTime,
			checkTime: "10AM",
			want:      false,
		},
		{
			name:      "times equal start",
			startTime: startTime,
			endTime:   endTime,
			checkTime: startTime.Raw(),
			want:      true,
		},
		{
			name:      "times equal end",
			startTime: startTime,
			endTime:   endTime,
			checkTime: endTime.Raw(),
			want:      true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			checkTime, _ := NewDeliveryTime(tc.checkTime)
			got := checkTime.InclusiveBetween(tc.startTime, tc.endTime)
			assert.Equal(t, tc.want, got)

		})
	}
}

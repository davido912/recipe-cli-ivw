package model

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReport_Dumps(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	err := NewReportModel().Dumps(buf)

	assert.Nil(t, err)

	expected := `{
 "unique_recipe_count": 0,
 "count_per_recipe": null,
 "busiest_postcode": {
  "postcode": "",
  "delivery_count": 0
 },
 "count_per_postcode_and_time": {
  "postcode": "",
  "from": "",
  "to": "",
  "delivery_count": 0
 },
 "match_by_name": null
}`
	assert.JSONEq(t, expected, buf.String())
}

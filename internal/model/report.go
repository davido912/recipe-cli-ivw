package model

import (
	"encoding/json"
	"io"
)

type RecipeCounts []RecipeCount

type RecipeCount struct {
	Recipe      string `json:"recipe"`
	RecipeCount int    `json:"count"`
}

type PostcodeCount struct {
	Postcode      string `json:"postcode"`
	DeliveryCount int    `json:"delivery_count"`
}

type RecipeMatches []string

type PostcodeTimeCount struct {
	Postcode      string `json:"postcode"`
	From          string `json:"from"`
	To            string `json:"to"`
	DeliveryCount int    `json:"delivery_count"`
}

type ReportModel struct {
	UniqueRecipeCount       int               `json:"unique_recipe_count"`
	CountPerRecipe          RecipeCounts      `json:"count_per_recipe"`
	BusiestPostcode         PostcodeCount     `json:"busiest_postcode"`
	CountPerPostcodeAndTime PostcodeTimeCount `json:"count_per_postcode_and_time"`
	MatchByName             RecipeMatches     `json:"match_by_name"`
}

// NewReportModel represents the final model used as output in this application
func NewReportModel() *ReportModel {
	return &ReportModel{}
}

// Dumps outputs the data into file or any other io.Writer
func (rm *ReportModel) Dumps(out io.Writer) error {
	bs, err := json.MarshalIndent(rm, "", " ")
	if err != nil {
		return err
	}

	_, err = out.Write(bs)
	if err != nil {
		return err
	}
	return nil
}

func (rm *ReportModel) SetUniqueRecipeCount(cnt int) {
	rm.UniqueRecipeCount = cnt
}

func (rm *ReportModel) SetCountPerRecipe(recipeCounts []RecipeCount) {
	rm.CountPerRecipe = recipeCounts
}

func (rm *ReportModel) SetBusiestPostcode(postcodeCount PostcodeCount) {
	rm.BusiestPostcode = postcodeCount
}

func (rm *ReportModel) SetCountPerPostcodeAndTime(postcodeTimeCount PostcodeTimeCount) {
	rm.CountPerPostcodeAndTime = postcodeTimeCount
}

func (rm *ReportModel) SetMatchByName(recipeMatches RecipeMatches) {
	rm.MatchByName = recipeMatches
}

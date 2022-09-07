package models

import (
	"fmt"
	"github.com/mindstand/gogm/v2"
	"reflect"
)

type Article struct {
	gogm.BaseNode

	Name              string `gogm:"name=name"`
	Description       string `gogm:"name=description"`
	DefaultProductUrl string `gogm:"name=default_product_url"`
	DefaultImageUrl   string `gogm:"name=default_image_url"`

	Brand      *Brand         `gogm:"relationship=MADE_BY;direction=OUTGOING"`
	Tags       []*Tag         `gogm:"relationship=HAS_TAG;direction=OUTGOING"`
	Variations []*VariantEdge `gogm:"relationship=HAS_VARIANT;direction=OUTGOING"`
}

type VariantEdge struct {
	gogm.BaseNode

	Start *Article
	End   *Variation

	Available      bool   `gogm:"name=availability"`
	VariantStoreID string `gogm:"name=variant_store_id"`
	ColorName      string `gogm:"name=color_name"`
	CurrentPrice   int    `gogm:"name=current_price"`
	BasePrice      int    `gogm:"name=base_price"`
	VariantUrl     string `gogm:"name=article_url"`
}

func (e *VariantEdge) GetStartNode() interface{} {
	return e.Start
}

func (e *VariantEdge) GetStartNodeType() reflect.Type {
	return reflect.TypeOf(&Article{})
}

func (e *VariantEdge) SetStartNode(v interface{}) error {
	val, ok := v.(*Article)
	if !ok {
		return fmt.Errorf("unable to cast [%T] to *Article", v)
	}

	e.Start = val
	return nil
}

func (e *VariantEdge) GetEndNode() interface{} {
	return e.End
}

func (e *VariantEdge) GetEndNodeType() reflect.Type {
	return reflect.TypeOf(&Variation{})
}

func (e *VariantEdge) SetEndNode(v interface{}) error {
	val, ok := v.(*Variation)
	if !ok {
		return fmt.Errorf("unable to cast [%T] to *Variation", v)
	}

	e.End = val
	return nil
}

type Tag struct {
	gogm.BaseNode

	Name string `gogm:"name=name;unique"`

	Articles []*Article `gogm:"relationship=HAS_TAG;direction=incoming"`
}

type Color struct {
	gogm.BaseNode

	HexCode string `gogm:"name=hex_code;unique"`

	Variations []*Variation `gogm:"relationship=HAS_COLOR;direction=incoming"`
}

type Size struct {
	gogm.BaseNode

	Size string `gogm:"name=size;unique"`

	Variations []*Variation `gogm:"relationship=HAS_SIZE;direction=incoming"`
}

type Variation struct {
	gogm.BaseNode

	ArticleEdges []*VariantEdge `gogm:"relationship=HAS_VARIANT;direction=incoming"`

	Color *Color `gogm:"relationship=HAS_COLOR;direction=outgoing"`
	Size  *Size  `gogm:"relationship=HAS_SIZE;direction=outgoing"`
	Store *Store `gogm:"relationship=SOLD_AT;direction=outgoing"`
}

type Brand struct {
	gogm.BaseNode

	Name string `gogm:"name=name;index"`

	Articles []*Article `gogm:"relationship=MADE_BY;direction=incoming"`

	FirstPartyStore *Store `gogm:"relationship=FIRST_PARTY_STORE;direction=outgoing"`
}

type Store struct {
	gogm.BaseNode

	Name       string `gogm:"name=name;unique"`
	WebsiteUrl string `gogm:"name=website_url"`

	Variations      []*Variation `gogm:"relationship=SOLD_AT;direction=incoming"`
	FirstPartyBrand *Brand       `gogm:"relationship=FIRST_PARTY_STORE;direction=incoming"`
}

package goshopify

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
)

func GraphQLMetafieldTests(t *testing.T, metafield Metafield) {
	// Check that ID is assigned to the returned metafield
	expectedInt := int64(1)
	if metafield.ID != expectedInt {
		t.Errorf("Metafield.ID returned %+v, expected %+v", metafield.ID, expectedInt)
	}
}

func TestGraphQLMetafieldList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/api/9999-99/graphql.json",
		httpmock.NewStringResponder(200, `{
			"data": {
				"metafieldDefinitions": {
					"edges": [{
						"node": {
							"id": "gid:\/\/shopify\/MetafieldDefinition\/340033559",
							"name": "Product Type",
							"namespace": "oyster",
							"key": "product_type",
							"validations": [{
								"name": "choices",
								"type": "list.single_line_text_field",
								"value": "[\"N\/A\",\"bike\",\"jewelry\"]"
							}]
						}
					}, {
						"node": {
							"id": "gid:\/\/shopify\/MetafieldDefinition\/341147671",
							"name": "Test Definition",
							"namespace": "oyster",
							"key": "test_definition",
							"validations": [{
								"name": "choices",
								"type": "list.single_line_text_field",
								"value": "[\"Choice 1\",\"Choice 2\",\"Choice 3\"]"
							}]
						}
					}]
				}
			},
			"extensions": {
				"cost": {
					"requestedQueryCost": 302,
					"actualQueryCost": 8,
					"throttleStatus": {
						"maximumAvailable": 1000.0,
						"currentlyAvailable": 992,
						"restoreRate": 50.0
					}
				}
			}
		}`))

	response, err := client.MetafieldDefinition.List(OwnerTypeProduct, "")
	if err != nil {
		t.Errorf("Metafield.List returned error: %v", err)
	}

	expected := []MetafieldDefinition{{
		Id:        "gid://shopify/MetafieldDefinition/340033559",
		Name:      "Product Type",
		Namespace: "oyster",
		Key:       "product_type",
		Validations: []MetafieldValidation{{
			Name:  "choices",
			Type:  "list.single_line_text_field",
			Value: "[\"N/A\",\"bike\",\"jewelry\"]",
		}},
	}, {
		Id:        "gid://shopify/MetafieldDefinition/341147671",
		Name:      "Test Definition",
		Namespace: "oyster",
		Key:       "test_definition",
		Validations: []MetafieldValidation{{
			Name:  "choices",
			Type:  "list.single_line_text_field",
			Value: "[\"Choice 1\",\"Choice 2\",\"Choice 3\"]",
		}},
	}}
	if !reflect.DeepEqual(response, expected) {
		t.Errorf("Metafield.List returned %+v expected %+v", response, expected)
	}
}

func TestGraphQLMetafieldCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST", "https://fooshop.myshopify.com/admin/api/9999-99/graphql.json",
		httpmock.NewStringResponder(200, `{
			"data": {
			  "metafieldDefinitionCreate": {
				"createdDefinition": {
				  "id": "gid:\/\/shopify\/MetafieldDefinition\/343506967",
				  "name": "Test Definition"
				},
				"userErrors": []
			  }
			},
			"extensions": {
			  "cost": {
				"requestedQueryCost": 10,
				"actualQueryCost": 10,
				"throttleStatus": {
				  "maximumAvailable": 1000.0,
				  "currentlyAvailable": 990,
				  "restoreRate": 50.0
				}
			  }
			}
		  }`))

	validations := make([]MetafieldValidation, 1)
	validationValues := []string{"Choice 1", "Choice 2", "Choice 3"}
	formattedValidationValues := fmt.Sprintf("[\"%s\"]", strings.Join(validationValues, "\",\""))
	validations[0] = MetafieldValidation{
		Name:  "choices",
		Value: formattedValidationValues,
	}

	response, err := client.MetafieldDefinition.Create(MetafieldDefinition{
		Name:                   "Test Definition",
		Namespace:              "oyster",
		Key:                    "test_definition",
		Description:            "Test definition for testing graphql",
		Type:                   "single_line_text_field",
		OwnerType:              string(OwnerTypeProduct),
		Validations:            validations,
		Pin:                    true,
		VisibleToStorefrontApi: true,
	})
	if err != nil {
		t.Errorf("Metafield.List returned error: %v", err)
	}

	expected := &MetafieldDefinition{
		Id:   "gid://shopify/MetafieldDefinition/343506967",
		Name: "Test Definition",
	}
	if !reflect.DeepEqual(response, expected) {
		t.Errorf("Metafield.List returned %+v, expected %+v", response, expected)
	}
}

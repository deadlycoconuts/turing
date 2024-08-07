/*
 * Turing Minimal Openapi Spec for SDK
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// EnsemblingJobPredictionSource struct for EnsemblingJobPredictionSource
type EnsemblingJobPredictionSource struct {
	Dataset Dataset `json:"dataset"`
	JoinOn []string `json:"join_on"`
	Columns []string `json:"columns,omitempty"`
}

// NewEnsemblingJobPredictionSource instantiates a new EnsemblingJobPredictionSource object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewEnsemblingJobPredictionSource(dataset Dataset, joinOn []string) *EnsemblingJobPredictionSource {
	this := EnsemblingJobPredictionSource{}
	this.Dataset = dataset
	this.JoinOn = joinOn
	return &this
}

// NewEnsemblingJobPredictionSourceWithDefaults instantiates a new EnsemblingJobPredictionSource object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewEnsemblingJobPredictionSourceWithDefaults() *EnsemblingJobPredictionSource {
	this := EnsemblingJobPredictionSource{}
	return &this
}

// GetDataset returns the Dataset field value
func (o *EnsemblingJobPredictionSource) GetDataset() Dataset {
	if o == nil {
		var ret Dataset
		return ret
	}

	return o.Dataset
}

// GetDatasetOk returns a tuple with the Dataset field value
// and a boolean to check if the value has been set.
func (o *EnsemblingJobPredictionSource) GetDatasetOk() (*Dataset, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Dataset, true
}

// SetDataset sets field value
func (o *EnsemblingJobPredictionSource) SetDataset(v Dataset) {
	o.Dataset = v
}

// GetJoinOn returns the JoinOn field value
func (o *EnsemblingJobPredictionSource) GetJoinOn() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.JoinOn
}

// GetJoinOnOk returns a tuple with the JoinOn field value
// and a boolean to check if the value has been set.
func (o *EnsemblingJobPredictionSource) GetJoinOnOk() (*[]string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.JoinOn, true
}

// SetJoinOn sets field value
func (o *EnsemblingJobPredictionSource) SetJoinOn(v []string) {
	o.JoinOn = v
}

// GetColumns returns the Columns field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *EnsemblingJobPredictionSource) GetColumns() []string {
	if o == nil  {
		var ret []string
		return ret
	}
	return o.Columns
}

// GetColumnsOk returns a tuple with the Columns field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *EnsemblingJobPredictionSource) GetColumnsOk() (*[]string, bool) {
	if o == nil || o.Columns == nil {
		return nil, false
	}
	return &o.Columns, true
}

// HasColumns returns a boolean if a field has been set.
func (o *EnsemblingJobPredictionSource) HasColumns() bool {
	if o != nil && o.Columns != nil {
		return true
	}

	return false
}

// SetColumns gets a reference to the given []string and assigns it to the Columns field.
func (o *EnsemblingJobPredictionSource) SetColumns(v []string) {
	o.Columns = v
}

func (o EnsemblingJobPredictionSource) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["dataset"] = o.Dataset
	}
	if true {
		toSerialize["join_on"] = o.JoinOn
	}
	if o.Columns != nil {
		toSerialize["columns"] = o.Columns
	}
	return json.Marshal(toSerialize)
}

type NullableEnsemblingJobPredictionSource struct {
	value *EnsemblingJobPredictionSource
	isSet bool
}

func (v NullableEnsemblingJobPredictionSource) Get() *EnsemblingJobPredictionSource {
	return v.value
}

func (v *NullableEnsemblingJobPredictionSource) Set(val *EnsemblingJobPredictionSource) {
	v.value = val
	v.isSet = true
}

func (v NullableEnsemblingJobPredictionSource) IsSet() bool {
	return v.isSet
}

func (v *NullableEnsemblingJobPredictionSource) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableEnsemblingJobPredictionSource(val *EnsemblingJobPredictionSource) *NullableEnsemblingJobPredictionSource {
	return &NullableEnsemblingJobPredictionSource{value: val, isSet: true}
}

func (v NullableEnsemblingJobPredictionSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableEnsemblingJobPredictionSource) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



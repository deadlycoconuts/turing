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

// EnvVar struct for EnvVar
type EnvVar struct {
	Name string `json:"name"`
	Value *string `json:"value,omitempty"`
}

// NewEnvVar instantiates a new EnvVar object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewEnvVar(name string) *EnvVar {
	this := EnvVar{}
	this.Name = name
	return &this
}

// NewEnvVarWithDefaults instantiates a new EnvVar object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewEnvVarWithDefaults() *EnvVar {
	this := EnvVar{}
	return &this
}

// GetName returns the Name field value
func (o *EnvVar) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *EnvVar) GetNameOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *EnvVar) SetName(v string) {
	o.Name = v
}

// GetValue returns the Value field value if set, zero value otherwise.
func (o *EnvVar) GetValue() string {
	if o == nil || o.Value == nil {
		var ret string
		return ret
	}
	return *o.Value
}

// GetValueOk returns a tuple with the Value field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EnvVar) GetValueOk() (*string, bool) {
	if o == nil || o.Value == nil {
		return nil, false
	}
	return o.Value, true
}

// HasValue returns a boolean if a field has been set.
func (o *EnvVar) HasValue() bool {
	if o != nil && o.Value != nil {
		return true
	}

	return false
}

// SetValue gets a reference to the given string and assigns it to the Value field.
func (o *EnvVar) SetValue(v string) {
	o.Value = &v
}

func (o EnvVar) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["name"] = o.Name
	}
	if o.Value != nil {
		toSerialize["value"] = o.Value
	}
	return json.Marshal(toSerialize)
}

type NullableEnvVar struct {
	value *EnvVar
	isSet bool
}

func (v NullableEnvVar) Get() *EnvVar {
	return v.value
}

func (v *NullableEnvVar) Set(val *EnvVar) {
	v.value = val
	v.isSet = true
}

func (v NullableEnvVar) IsSet() bool {
	return v.isSet
}

func (v *NullableEnvVar) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableEnvVar(val *EnvVar) *NullableEnvVar {
	return &NullableEnvVar{value: val, isSet: true}
}

func (v NullableEnvVar) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableEnvVar) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



/*
Wasp API

REST API for the Wasp node

API version: 0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package apiclient

import (
	"encoding/json"
	"bytes"
	"fmt"
)

// checks if the VersionResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &VersionResponse{}

// VersionResponse struct for VersionResponse
type VersionResponse struct {
	// The version of the node
	Version string `json:"version"`
}

type _VersionResponse VersionResponse

// NewVersionResponse instantiates a new VersionResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewVersionResponse(version string) *VersionResponse {
	this := VersionResponse{}
	this.Version = version
	return &this
}

// NewVersionResponseWithDefaults instantiates a new VersionResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewVersionResponseWithDefaults() *VersionResponse {
	this := VersionResponse{}
	return &this
}

// GetVersion returns the Version field value
func (o *VersionResponse) GetVersion() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Version
}

// GetVersionOk returns a tuple with the Version field value
// and a boolean to check if the value has been set.
func (o *VersionResponse) GetVersionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Version, true
}

// SetVersion sets field value
func (o *VersionResponse) SetVersion(v string) {
	o.Version = v
}

func (o VersionResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o VersionResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["version"] = o.Version
	return toSerialize, nil
}

func (o *VersionResponse) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"version",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err;
	}

	for _, requiredProperty := range(requiredProperties) {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varVersionResponse := _VersionResponse{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varVersionResponse)

	if err != nil {
		return err
	}

	*o = VersionResponse(varVersionResponse)

	return err
}

type NullableVersionResponse struct {
	value *VersionResponse
	isSet bool
}

func (v NullableVersionResponse) Get() *VersionResponse {
	return v.value
}

func (v *NullableVersionResponse) Set(val *VersionResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableVersionResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableVersionResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableVersionResponse(val *VersionResponse) *NullableVersionResponse {
	return &NullableVersionResponse{value: val, isSet: true}
}

func (v NullableVersionResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableVersionResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



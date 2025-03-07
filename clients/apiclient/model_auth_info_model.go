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

// checks if the AuthInfoModel type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AuthInfoModel{}

// AuthInfoModel struct for AuthInfoModel
type AuthInfoModel struct {
	// JWT only
	AuthURL string `json:"authURL"`
	Scheme string `json:"scheme"`
}

type _AuthInfoModel AuthInfoModel

// NewAuthInfoModel instantiates a new AuthInfoModel object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAuthInfoModel(authURL string, scheme string) *AuthInfoModel {
	this := AuthInfoModel{}
	this.AuthURL = authURL
	this.Scheme = scheme
	return &this
}

// NewAuthInfoModelWithDefaults instantiates a new AuthInfoModel object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAuthInfoModelWithDefaults() *AuthInfoModel {
	this := AuthInfoModel{}
	return &this
}

// GetAuthURL returns the AuthURL field value
func (o *AuthInfoModel) GetAuthURL() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.AuthURL
}

// GetAuthURLOk returns a tuple with the AuthURL field value
// and a boolean to check if the value has been set.
func (o *AuthInfoModel) GetAuthURLOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.AuthURL, true
}

// SetAuthURL sets field value
func (o *AuthInfoModel) SetAuthURL(v string) {
	o.AuthURL = v
}

// GetScheme returns the Scheme field value
func (o *AuthInfoModel) GetScheme() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Scheme
}

// GetSchemeOk returns a tuple with the Scheme field value
// and a boolean to check if the value has been set.
func (o *AuthInfoModel) GetSchemeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Scheme, true
}

// SetScheme sets field value
func (o *AuthInfoModel) SetScheme(v string) {
	o.Scheme = v
}

func (o AuthInfoModel) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AuthInfoModel) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["authURL"] = o.AuthURL
	toSerialize["scheme"] = o.Scheme
	return toSerialize, nil
}

func (o *AuthInfoModel) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"authURL",
		"scheme",
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

	varAuthInfoModel := _AuthInfoModel{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varAuthInfoModel)

	if err != nil {
		return err
	}

	*o = AuthInfoModel(varAuthInfoModel)

	return err
}

type NullableAuthInfoModel struct {
	value *AuthInfoModel
	isSet bool
}

func (v NullableAuthInfoModel) Get() *AuthInfoModel {
	return v.value
}

func (v *NullableAuthInfoModel) Set(val *AuthInfoModel) {
	v.value = val
	v.isSet = true
}

func (v NullableAuthInfoModel) IsSet() bool {
	return v.isSet
}

func (v *NullableAuthInfoModel) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAuthInfoModel(val *AuthInfoModel) *NullableAuthInfoModel {
	return &NullableAuthInfoModel{value: val, isSet: true}
}

func (v NullableAuthInfoModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAuthInfoModel) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



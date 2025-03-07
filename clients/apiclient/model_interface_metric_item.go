/*
Wasp API

REST API for the Wasp node

API version: 0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package apiclient

import (
	"encoding/json"
	"time"
	"bytes"
	"fmt"
)

// checks if the InterfaceMetricItem type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &InterfaceMetricItem{}

// InterfaceMetricItem struct for InterfaceMetricItem
type InterfaceMetricItem struct {
	LastMessage string `json:"lastMessage"`
	Messages uint32 `json:"messages"`
	Timestamp time.Time `json:"timestamp"`
}

type _InterfaceMetricItem InterfaceMetricItem

// NewInterfaceMetricItem instantiates a new InterfaceMetricItem object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewInterfaceMetricItem(lastMessage string, messages uint32, timestamp time.Time) *InterfaceMetricItem {
	this := InterfaceMetricItem{}
	this.LastMessage = lastMessage
	this.Messages = messages
	this.Timestamp = timestamp
	return &this
}

// NewInterfaceMetricItemWithDefaults instantiates a new InterfaceMetricItem object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewInterfaceMetricItemWithDefaults() *InterfaceMetricItem {
	this := InterfaceMetricItem{}
	return &this
}

// GetLastMessage returns the LastMessage field value
func (o *InterfaceMetricItem) GetLastMessage() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.LastMessage
}

// GetLastMessageOk returns a tuple with the LastMessage field value
// and a boolean to check if the value has been set.
func (o *InterfaceMetricItem) GetLastMessageOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.LastMessage, true
}

// SetLastMessage sets field value
func (o *InterfaceMetricItem) SetLastMessage(v string) {
	o.LastMessage = v
}

// GetMessages returns the Messages field value
func (o *InterfaceMetricItem) GetMessages() uint32 {
	if o == nil {
		var ret uint32
		return ret
	}

	return o.Messages
}

// GetMessagesOk returns a tuple with the Messages field value
// and a boolean to check if the value has been set.
func (o *InterfaceMetricItem) GetMessagesOk() (*uint32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Messages, true
}

// SetMessages sets field value
func (o *InterfaceMetricItem) SetMessages(v uint32) {
	o.Messages = v
}

// GetTimestamp returns the Timestamp field value
func (o *InterfaceMetricItem) GetTimestamp() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.Timestamp
}

// GetTimestampOk returns a tuple with the Timestamp field value
// and a boolean to check if the value has been set.
func (o *InterfaceMetricItem) GetTimestampOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Timestamp, true
}

// SetTimestamp sets field value
func (o *InterfaceMetricItem) SetTimestamp(v time.Time) {
	o.Timestamp = v
}

func (o InterfaceMetricItem) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o InterfaceMetricItem) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["lastMessage"] = o.LastMessage
	toSerialize["messages"] = o.Messages
	toSerialize["timestamp"] = o.Timestamp
	return toSerialize, nil
}

func (o *InterfaceMetricItem) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"lastMessage",
		"messages",
		"timestamp",
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

	varInterfaceMetricItem := _InterfaceMetricItem{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varInterfaceMetricItem)

	if err != nil {
		return err
	}

	*o = InterfaceMetricItem(varInterfaceMetricItem)

	return err
}

type NullableInterfaceMetricItem struct {
	value *InterfaceMetricItem
	isSet bool
}

func (v NullableInterfaceMetricItem) Get() *InterfaceMetricItem {
	return v.value
}

func (v *NullableInterfaceMetricItem) Set(val *InterfaceMetricItem) {
	v.value = val
	v.isSet = true
}

func (v NullableInterfaceMetricItem) IsSet() bool {
	return v.isSet
}

func (v *NullableInterfaceMetricItem) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableInterfaceMetricItem(val *InterfaceMetricItem) *NullableInterfaceMetricItem {
	return &NullableInterfaceMetricItem{value: val, isSet: true}
}

func (v NullableInterfaceMetricItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableInterfaceMetricItem) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
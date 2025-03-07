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

// checks if the ConsensusPipeMetrics type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConsensusPipeMetrics{}

// ConsensusPipeMetrics struct for ConsensusPipeMetrics
type ConsensusPipeMetrics struct {
	EventACSMsgPipeSize int32 `json:"eventACSMsgPipeSize"`
	EventPeerLogIndexMsgPipeSize int32 `json:"eventPeerLogIndexMsgPipeSize"`
	EventStateTransitionMsgPipeSize int32 `json:"eventStateTransitionMsgPipeSize"`
	EventTimerMsgPipeSize int32 `json:"eventTimerMsgPipeSize"`
	EventVMResultMsgPipeSize int32 `json:"eventVMResultMsgPipeSize"`
}

type _ConsensusPipeMetrics ConsensusPipeMetrics

// NewConsensusPipeMetrics instantiates a new ConsensusPipeMetrics object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConsensusPipeMetrics(eventACSMsgPipeSize int32, eventPeerLogIndexMsgPipeSize int32, eventStateTransitionMsgPipeSize int32, eventTimerMsgPipeSize int32, eventVMResultMsgPipeSize int32) *ConsensusPipeMetrics {
	this := ConsensusPipeMetrics{}
	this.EventACSMsgPipeSize = eventACSMsgPipeSize
	this.EventPeerLogIndexMsgPipeSize = eventPeerLogIndexMsgPipeSize
	this.EventStateTransitionMsgPipeSize = eventStateTransitionMsgPipeSize
	this.EventTimerMsgPipeSize = eventTimerMsgPipeSize
	this.EventVMResultMsgPipeSize = eventVMResultMsgPipeSize
	return &this
}

// NewConsensusPipeMetricsWithDefaults instantiates a new ConsensusPipeMetrics object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConsensusPipeMetricsWithDefaults() *ConsensusPipeMetrics {
	this := ConsensusPipeMetrics{}
	return &this
}

// GetEventACSMsgPipeSize returns the EventACSMsgPipeSize field value
func (o *ConsensusPipeMetrics) GetEventACSMsgPipeSize() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.EventACSMsgPipeSize
}

// GetEventACSMsgPipeSizeOk returns a tuple with the EventACSMsgPipeSize field value
// and a boolean to check if the value has been set.
func (o *ConsensusPipeMetrics) GetEventACSMsgPipeSizeOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.EventACSMsgPipeSize, true
}

// SetEventACSMsgPipeSize sets field value
func (o *ConsensusPipeMetrics) SetEventACSMsgPipeSize(v int32) {
	o.EventACSMsgPipeSize = v
}

// GetEventPeerLogIndexMsgPipeSize returns the EventPeerLogIndexMsgPipeSize field value
func (o *ConsensusPipeMetrics) GetEventPeerLogIndexMsgPipeSize() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.EventPeerLogIndexMsgPipeSize
}

// GetEventPeerLogIndexMsgPipeSizeOk returns a tuple with the EventPeerLogIndexMsgPipeSize field value
// and a boolean to check if the value has been set.
func (o *ConsensusPipeMetrics) GetEventPeerLogIndexMsgPipeSizeOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.EventPeerLogIndexMsgPipeSize, true
}

// SetEventPeerLogIndexMsgPipeSize sets field value
func (o *ConsensusPipeMetrics) SetEventPeerLogIndexMsgPipeSize(v int32) {
	o.EventPeerLogIndexMsgPipeSize = v
}

// GetEventStateTransitionMsgPipeSize returns the EventStateTransitionMsgPipeSize field value
func (o *ConsensusPipeMetrics) GetEventStateTransitionMsgPipeSize() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.EventStateTransitionMsgPipeSize
}

// GetEventStateTransitionMsgPipeSizeOk returns a tuple with the EventStateTransitionMsgPipeSize field value
// and a boolean to check if the value has been set.
func (o *ConsensusPipeMetrics) GetEventStateTransitionMsgPipeSizeOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.EventStateTransitionMsgPipeSize, true
}

// SetEventStateTransitionMsgPipeSize sets field value
func (o *ConsensusPipeMetrics) SetEventStateTransitionMsgPipeSize(v int32) {
	o.EventStateTransitionMsgPipeSize = v
}

// GetEventTimerMsgPipeSize returns the EventTimerMsgPipeSize field value
func (o *ConsensusPipeMetrics) GetEventTimerMsgPipeSize() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.EventTimerMsgPipeSize
}

// GetEventTimerMsgPipeSizeOk returns a tuple with the EventTimerMsgPipeSize field value
// and a boolean to check if the value has been set.
func (o *ConsensusPipeMetrics) GetEventTimerMsgPipeSizeOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.EventTimerMsgPipeSize, true
}

// SetEventTimerMsgPipeSize sets field value
func (o *ConsensusPipeMetrics) SetEventTimerMsgPipeSize(v int32) {
	o.EventTimerMsgPipeSize = v
}

// GetEventVMResultMsgPipeSize returns the EventVMResultMsgPipeSize field value
func (o *ConsensusPipeMetrics) GetEventVMResultMsgPipeSize() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.EventVMResultMsgPipeSize
}

// GetEventVMResultMsgPipeSizeOk returns a tuple with the EventVMResultMsgPipeSize field value
// and a boolean to check if the value has been set.
func (o *ConsensusPipeMetrics) GetEventVMResultMsgPipeSizeOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.EventVMResultMsgPipeSize, true
}

// SetEventVMResultMsgPipeSize sets field value
func (o *ConsensusPipeMetrics) SetEventVMResultMsgPipeSize(v int32) {
	o.EventVMResultMsgPipeSize = v
}

func (o ConsensusPipeMetrics) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConsensusPipeMetrics) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["eventACSMsgPipeSize"] = o.EventACSMsgPipeSize
	toSerialize["eventPeerLogIndexMsgPipeSize"] = o.EventPeerLogIndexMsgPipeSize
	toSerialize["eventStateTransitionMsgPipeSize"] = o.EventStateTransitionMsgPipeSize
	toSerialize["eventTimerMsgPipeSize"] = o.EventTimerMsgPipeSize
	toSerialize["eventVMResultMsgPipeSize"] = o.EventVMResultMsgPipeSize
	return toSerialize, nil
}

func (o *ConsensusPipeMetrics) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"eventACSMsgPipeSize",
		"eventPeerLogIndexMsgPipeSize",
		"eventStateTransitionMsgPipeSize",
		"eventTimerMsgPipeSize",
		"eventVMResultMsgPipeSize",
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

	varConsensusPipeMetrics := _ConsensusPipeMetrics{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varConsensusPipeMetrics)

	if err != nil {
		return err
	}

	*o = ConsensusPipeMetrics(varConsensusPipeMetrics)

	return err
}

type NullableConsensusPipeMetrics struct {
	value *ConsensusPipeMetrics
	isSet bool
}

func (v NullableConsensusPipeMetrics) Get() *ConsensusPipeMetrics {
	return v.value
}

func (v *NullableConsensusPipeMetrics) Set(val *ConsensusPipeMetrics) {
	v.value = val
	v.isSet = true
}

func (v NullableConsensusPipeMetrics) IsSet() bool {
	return v.isSet
}

func (v *NullableConsensusPipeMetrics) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConsensusPipeMetrics(val *ConsensusPipeMetrics) *NullableConsensusPipeMetrics {
	return &NullableConsensusPipeMetrics{value: val, isSet: true}
}

func (v NullableConsensusPipeMetrics) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConsensusPipeMetrics) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



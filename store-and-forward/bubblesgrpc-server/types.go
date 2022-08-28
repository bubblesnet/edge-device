/*
 * Copyright (c) John Rodley 2022.
 * SPDX-FileCopyrightText:  John Rodley 2022.
 * SPDX-License-Identifier: MIT
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this
 * software and associated documentation files (the "Software"), to deal in the
 * Software without restriction, including without limitation the rights to use, copy,
 * modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so, subject to the
 * following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
 * INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
 * PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
 * HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
 * CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package main

const TEMPNOTSET float64 = -100.5
const HUMIDITYNOTSET float64 = -100.5
const PRESSURENOTSET float64 = -100.5
const LIGHTNOTSET float64 = -100.5
const HEIGHTNOTSET float64 = -1

type externalState struct {
	PlantHeightIn    float64 `json:"plant_height_in,omitempty"`
	WaterTemp        float64 `json:"water_temp,omitempty"`
	TempAirMiddle    float64 `json:"temp_air_middle,omitempty"`
	TempExternal     float64 `json:"temp_external,omitempty"`
	HumidityInternal float64 `json:"humidity_internal,omitempty"`
	PressureInternal float64 `json:"pressure_internal,omitempty"`
	LightInternal    float64 `json:"light_internal,omitempty"`
}

var ExternalCurrentState = externalState{
	PlantHeightIn:    HEIGHTNOTSET,
	WaterTemp:        TEMPNOTSET,
	TempAirMiddle:    TEMPNOTSET,
	TempExternal:     TEMPNOTSET,
	HumidityInternal: HUMIDITYNOTSET,
	LightInternal:    LIGHTNOTSET,
	PressureInternal: PRESSURENOTSET,
}

type GenericSensorMessage struct {
	DeviceId          int64   `json:"deviceid"`
	ContainerName     string  `json:"container_name"`
	ExecutableVersion string  `json:"executable_version"`
	SampleTimestamp   int64   `json:"sample_timestamp"`
	MessageType       string  `json:"message_type"`
	SensorName        string  `json:"sensor_name"`
	MeasurementName   string  `json:"measurement_name"`
	Value             float64 `json:"value,omitempty"`
	FloatValue        float64 `json:"floatvalue,omitempty"`
	Units             string  `json:"units"`
	Direction         string  `json:"direction"`
}

const (
	measurement_name_temp_air_top      = "temp_air_top"
	measurement_name_temp_air_middle   = "temp_air_middle"
	measurement_name_temp_air_bottom   = "temp_air_bottom"
	measurement_name_temp_air_external = "temp_air_external"
	Measurement_name_temp_water        = "temp_water"
	measurement_name_humidity_internal = "humidity_internal"
	measurement_name_humidity_external = "humidity_external"
	measurement_name_water_level       = "water_level"
	Measurement_name_plant_height      = "plant_height"
	Measurement_name_root_ph           = "root_ph"
	Measurement_name_light_internal    = "light_internal"
	Measurement_name_pressure_internal = "pressure_internal"

	measurement_type_temperature_middle = "temperature_middle"
	measurement_type_humidity_internal  = "humidity_internal"
	measurement_type_pH                 = "pH"
	measurement_type_light_internal     = "light_internal"
	Measurement_type_water_level        = "water_level"
	measurement_type_water_temperature  = "water_temperature"
	measurement_type_co2                = "co2"
	measurement_type_voc                = "voc"

	/// emulator types - not all valid?
	measurement_type_temperature = "temperature"
	measurement_type_humidity    = "humidity"
	measurement_type_level       = "level"

	message_type_measurement     = "measurement"
	message_type_switch_event    = "switch_event"
	message_type_dispenser_event = "dispenser_event"
	message_type_picture_event   = "picture_event"
	// CLIENT ADDED
	message_type_event = "event"

	Directions_up   = "up"
	Directions_down = "down"
	Directions_none = ""

	sensor_name_thermometer_top          = "thermometer_top"
	sensor_name_thermometer_middle       = "thermometer_middle"
	sensor_name_thermometer_bottom       = "thermometer_bottom"
	sensor_name_thermometer_external     = "thermometer_external"
	Sensor_name_thermometer_water        = "thermometer_water"
	sensor_name_humidity_sensor_internal = "humidity_sensor_internal"
	sensor_name_humidity_sensor_external = "humidity_sensor_external"
	Sensor_name_water_level_sensor       = "water_level_sensor"

	// CLIENT ADDED
	sensor_name_light_sensor_internal = "light_sensor_internal"
	sensor_name_light_sensor_external = "light_sensor_external"
	sensor_name_station_door_sensor   = "station_door_sensor"
	sensor_name_outer_door_sensor     = "outer_door_sensor"
	Sensor_name_root_ph_sensor        = "root_ph_sensor"
	Sensor_name_height_sensor         = "height_sensor"
	sensor_name_voc_sensor            = "voc_sensor"
	sensor_name_co2_sensor            = "co2_sensor"
	sensor_name_ec_sensor             = "ec_sensor"
	Sensor_name_tamper_sensor         = "tamper_sensor"

	ac_device_name_humidifier       = "humidifier"
	ac_device_name_water_heater     = "water_heater"
	ac_device_name_heater           = "heater"
	ac_device_name_water_pump       = "water_pump"
	ac_device_name_air_pump         = "air_pump"
	ac_device_name_intake_fan       = "intake_fan"
	ac_device_name_exhaust_fan      = "exhaust_fan"
	ac_device_name_heat_lamp        = "heat_lamp"
	ac_device_name_heating_pad      = "heating_pad"
	ac_device_name_light_bloom      = "light_bloom"
	ac_device_name_light_vegetative = "light_vegetative"
	ac_device_name_light_germinate  = "light_germinate"

	Temperature_units_fahrenheit = "F"
	temperature_units_celsius    = "C"
	temperature_units_kelvin     = "K"

	humidity_units_percent = "%"

	Liquid_volume_units_gallons     = "gallons"
	liquid_volume_units_ounces      = "oz"
	liquid_volume_units_liters      = "liters"
	liquid_volume_units_milliliters = "ml"

	distance_units_inches      = "inches"
	Distance_units_centimeters = "cm"

	Ph_units_default = ""
	// CLIENT ADDED

	switch_name_growLight  = "growLight"
	switch_name_lightBloom = "lightBloom"

	Grpc_message_typeid_picture   = "picture"
	Grpc_message_typeid_sensor    = "sensor"
	Grpc_message_typeid_dispenser = "dispenser"
	Grpc_message_typeid_switch    = "switch"

	Command_type_stage    = "stage"
	Command_type_dispense = "dispense"
	Command_type_picture  = "picture"
	Command_type_status   = "status"
	Command_type_switch   = "switch"

	Switch_name_automatic_control = "automaticControl"

	Module_type_ezoph   = "ezoph"
	Module_type_ccs811  = "ccs811"
	Module_type_DS18B20 = "DS18B20"
	Module_type_adxl345 = "adxl345"
	Module_type_ads1115 = "ads1115"
	Module_type_hcsr04  = "hcsr04"
)

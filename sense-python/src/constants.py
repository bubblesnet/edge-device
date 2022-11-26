#!/usr/bin/env python3

#  Copyright (c) John Rodley 2022.
#  SPDX-FileCopyrightText:  John Rodley 2022.
#  SPDX-License-Identifier: MIT
#
#  Permission is hereby granted, free of charge, to any person obtaining a copy of this
#  software and associated documentation files (the "Software"), to deal in the
#  Software without restriction, including without limitation the rights to use, copy,
#  modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
#  and to permit persons to whom the Software is furnished to do so, subject to the
#  following conditions:
#
#  The above copyright notice and this permission notice shall be included in all
#  copies or substantial portions of the Software.
#
#  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
#  INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
#  PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
#  HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
#  CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
#  OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
#


"""
Environment variable names
"""
ENV_SLEEP_ON_EXIT_FOR_DEBUGGING = 'SLEEP_ON_EXIT_FOR_DEBUGGING'
ENV_DEVICEID = 'DEVICEID'

"""
Member names of site/station/device object hierarchy
"""

DEVICEID = 'deviceid'
STATIONID = 'stationid'
SITEID = 'siteid'
STATIONS = 'stations'
EDGE_DEVICES = 'edge_devices'
MODULES = 'modules'
CONTAINER_NAME = 'container_name'
MODULE_TYPE = 'module_type'
INCLUDED_SENSORS = 'included_sensors'
MEASUREMENT_NAME = 'measurement_name'
SENSOR_NAME = 'sensor_name'
ADDRESS = 'address'

DIRECTIONS_UP = 'up'
DIRECTIONS_DOWN = 'down'
DIRECTIONS_NONE = ''


MESSAGE_TYPE_ID_SENSOR = 'sensor'

"""
Member names of message objects
"""

MM_MESSAGE_TYPE = 'message_type'
MM_MESSAGE_TYPE_MEASUREMENT = 'measurement'

MM_SENSOR_NAME = 'sensor_name'
MM_MEASUREMENT_NAME = 'measurement_name'
MM_UNITS = 'units'
MM_VALUE_NAME = 'value_name'
MM_VALUE = 'value'
MM_FLOAT_VALUE = 'floatvalue'
MM_DIRECTION = 'direction'

MM_SAMPLE_TIMESTAMP = 'sample_timestamp'
MM_DEVICEID = 'deviceid'
MM_STATIONID = 'stationid'
MM_SITEID = 'siteid'
MM_CONTAINER_NAME = 'container_name'
MM_EXECUTABLE_VERSION = 'executable_version'

MM_SENSOR_NAME_TAMPER_DETECTOR = 'tamper_detector'
MM_MEASUREMENT_NAME_TAMPER = 'tamper'
MM_VALUE_NAME_TAMPER = 'tamper'
MM_TAMPER_DETECTOR = 'tamper_detector'

MM_WATER_TEMPERATURE = 'water_temperature'
MM_DOOR_OPEN = 'door_open'
MM_LEAK_DETECTOR = 'leak_detector'
MM_TEMPC = 'tempC'
MM_TEMPF = 'tempF'

MM_SENSOR_NAME_WATER_TEMPERATURE = 'water_temperature_sensor'
MM_MEASUREMENT_NAME_WATER_TEMPERATURE = 'temp_water'
MM_VALUE_NAME_WATER_TEMPERATURE = 'temp_water'

CONTAINER_NAME_SENSE_PYTHON = 'sense-python'

"""
Measurement unit type names
"""
MM_UNITS_BOOLEAN = 'boolean'
MM_UNITS_GALLONS = 'gallons'
MM_UNITS_LUX = 'lux'
MM_UNITS_CELSIUS = 'C'
MM_UNITS_HPA = 'hPa'
MM_UNITS_PERCENT = '%'

"""
Station capabilities
"""
CAPABILITY_TEMPERATURE_TOP = 'thermometer_top'
CAPABILITY_TEMPERATURE_MIDDLE = 'thermometer_middle'
CAPABILITY_TEMPERATURE_BOTTOM = 'thermometer_bottom'
CAPABILITY_TEMPERATURE_EXTERNAL = 'thermometer_external'
CAPABILITY_LIGHT_SENSOR_INTERNAL = 'light_sensor_internal'
CAPABILITY_LIGHT_SENSOR_EXTERNAL = 'light_sensor_external'
CAPABILITY_HUMIDITY_SENSOR_INTERNAL = 'humidity_sensor_internal'
CAPABILITY_HUMIDITY_SENSOR_EXTERNAL = 'humidity_sensor_external'
CAPABILITY_PRESSURE_SENSORS = 'pressure_sensors'
CAPABILITY_MOVEMENT_SENSOR = 'movement_sensor'


"""
Module types
"""

MT_BH1750 = 'bh1750'
MT_ADXL345 = 'adxl345'
MT_ADS1115 = 'ads1115'
MT_RELAY = 'relay'
MT_BMP280 = 'bmp280'
MT_BME280 = 'bme280'
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


import json
import logging
import time
import traceback
import os
import constants

try:
    import bme280
except ImportError:
    bme280 = None

try:
    import board
except NotImplementedError:
    board = None
except ImportError:
    board = None
except AttributeError:
    board = None

try:
    import busio
except ImportError:
    busio = None

try:
    import grpc as grpcio
except ImportError:
    grpcio = None

try:
    import smbus2
except ImportError:
    smbus2 = None

try:
    import bh1750
except ImportError:
    bh1750 = None

import bubblesgrpc_pb2
from bubblesgrpc_pb2_grpc import SensorStoreAndForwardStub as grpcStub
from os.path import exists

my_site = {}
my_station = {}
my_device = {}

LightAddress = 0

lastHumidity = 0.0
lastPressure = 0.0
lastLight = 0.0
lastTemp = 0.0

global humidity_sensor_name
global pressure_sensor_name
global temperature_sensor_name
global temperature_measurement_name
global humidity_measurement_name
global pressure_measurement_name
global light_sensor_name
global light_measurement_name

"""
TESTED
"""


def wait_for_config(filename):
    global naptime_in_seconds

    logging.info(f'wait_for_config {filename}')
    index = 0
    while index <= 60:
        if exists(path=filename):
            logging.info(f'{filename} file exists')
            return True
        logging.info(f'Sleeping while we wait for someone to create {filename}')
        time.sleep(5)  # 5 seconds is sort of arbitrary
        index = index + 1
        if index >= naptime_in_seconds / 5:
            logging.error(f"config file never showed up, exiting having slept {naptime_in_seconds} seconds")
            exit(1)

    return True


"""
TESTED
"""


def read_deviceid(filename):
    deviceid = -1
    return int(os.environ[constants.ENV_DEVICEID])


"""
TESTED
"""


def validate_config(site):
    if constants.STATIONS not in site:
        print('invalid config, no stations in my_site ', site)
        return False

    return True


"""
TESTED
"""


def read_config(fullpath, deviceid):
    global my_site, my_station, my_device
    try:
        with open(file=fullpath) as f:
            my_site = json.load(f)
            my_site[constants.DEVICEID] = deviceid
            if not validate_config(my_site):
                print('invalid config for device ', deviceid)
                return False
            for station in my_site[constants.STATIONS]:
                if constants.EDGE_DEVICES not in station:
                    continue
                for device in station[constants.EDGE_DEVICES]:
                    if device[constants.DEVICEID] == deviceid:
                        my_device = device
                        my_station = station
                        # my_site['time_between_sensor_polling_in_seconds'] = 15
        if not my_site:
            print('Could not set my_site')
            return False
        if not my_device:
            print(f'Could not find deviceid {deviceid} in {my_site}')
            return False
        if not my_station:
            print(f'Could not set my_station in {my_site}')
            return False
        return True
    except OSError:
        print('OSError on file ' + fullpath)
        return False


def append_bme280_temp(i2cbus, msg, sensor_name, measurement_name, address):
    global lastTemp
    logging.info(f'append_bme280_temp address {address}')
    try:
        calibration_params = bme280.load_calibration_params(i2cbus, address)
        data = bme280.sample(i2cbus, address, compensation_params=calibration_params)
        msg[constants.MM_SENSOR_NAME] = sensor_name
        msg[constants.MM_MEASUREMENT_NAME] = measurement_name
        msg[constants.MM_UNITS] = constants.MM_UNITS_CELSIUS
        msg[constants.MM_VALUE] = (data.temperature * 1.8) + 32.0
        msg[measurement_name] = msg[constants.MM_VALUE]
        msg[constants.MM_FLOAT_VALUE] = msg[constants.MM_VALUE]
        msg[constants.MM_VALUE_NAME] = measurement_name
        msg[constants.MM_DIRECTION] = constants.DIRECTIONS_NONE
        if data.temperature > lastTemp:
            msg[constants.MM_DIRECTION] = constants.DIRECTIONS_UP
        else:
            if data.temperature < lastTemp:
                msg[constants.MM_DIRECTION] = constants.DIRECTIONS_DOWN
        direction_name = measurement_name + '_direction'
        msg[direction_name] = msg[constants.MM_DIRECTION]
        lastTemp = data.temperature
        msg[constants.MM_TEMPC] = data.temperature
        msg[constants.MM_TEMPF] = (data.temperature * 1.8) + 32.0
    except Exception as ee:
        logging.error(f'bme280 error {ee}')
        logging.debug(traceback.format_exc())


def append_bme280_humidity(i2cbus, msg, sensor_name, measurement_name, address):
    global lastHumidity
    try:
        calibration_params = bme280.load_calibration_params(i2cbus, address)
        data = bme280.sample(i2cbus, address, compensation_params=calibration_params)
        msg[constants.MM_SENSOR_NAME] = sensor_name
        msg[constants.MM_MEASUREMENT_NAME] = measurement_name
        msg[constants.MM_UNITS] = constants.MM_UNITS_PERCENT
        msg[constants.MM_VALUE_NAME] = measurement_name
        msg[constants.MM_VALUE] = data.humidity
        msg[constants.MM_FLOAT_VALUE] = msg[constants.MM_VALUE]
        msg[measurement_name] = data.humidity
        direction = constants.DIRECTIONS_NONE
        if data.humidity > lastHumidity:
            direction = constants.DIRECTIONS_UP
        else:
            if data.humidity < lastHumidity:
                direction = constants.DIRECTIONS_DOWN
        msg[constants.MM_DIRECTION] = direction
        direction_name = measurement_name + '_direction'
        msg[direction_name] = direction
        lastHumidity = data.humidity
    except Exception as ee:
        logging.error(f'bme280 error {ee}')
        logging.debug(traceback.format_exc())


def append_bme280_pressure(i2cbus, msg, sensor_name, measurement_name, address):
    global lastPressure
    try:
        calibration_params = bme280.load_calibration_params(i2cbus, address)
        data = bme280.sample(i2cbus, address, compensation_params=calibration_params)
        msg[constants.MM_SENSOR_NAME] = sensor_name
        msg[constants.MM_MEASUREMENT_NAME] = measurement_name
        msg[constants.MM_UNITS] = constants.MM_UNITS_HPA
        msg[constants.MM_VALUE_NAME] = measurement_name
        msg[constants.MM_VALUE] = data.pressure
        msg[constants.MM_FLOAT_VALUE] = msg[constants.MM_VALUE]
        msg[measurement_name] = data.pressure
        direction = constants.DIRECTIONS_NONE
        if data.pressure > lastPressure:
            direction = constants.DIRECTIONS_UP
        else:
            if data.pressure < lastPressure:
                direction = constants.DIRECTIONS_DOWN
        lastPressure = data.pressure
        msg[constants.MM_DIRECTION] = direction
        direction_name = measurement_name + '_direction'
        msg[direction_name] = direction
        lastPressure = data.humidity
    except Exception as ee:
        logging.error(f'bme280 error {ee}')
        logging.debug(traceback.format_exc())


"""
TESTED
"""


def get_address(module_type):
    for module in my_device[constants.MODULES]:
        if module[constants.CONTAINER_NAME] == constants.CONTAINER_NAME_SENSE_PYTHON \
                and module[constants.MODULE_TYPE] == module_type:
            x = int(module[constants.ADDRESS], 16)
            return x
    return 0


"""
TESTED
"""


def is_our_device(module_type):
    for module in my_device[constants.MODULES]:
        if module[constants.CONTAINER_NAME] == constants.CONTAINER_NAME_SENSE_PYTHON \
                and module[constants.MODULE_TYPE] == module_type:
            return True

    return False


def get_our_device(module_type):
    for module in my_device[constants.MODULES]:
        if module[constants.CONTAINER_NAME] == constants.CONTAINER_NAME_SENSE_PYTHON \
                and module[constants.MODULE_TYPE] == module_type:
            return module

    return


def bh1750_names():
    global my_site
    global my_device
    global light_sensor_name
    global light_measurement_name

    print(my_device)
    for module in my_device[constants.MODULES]:
        if module[constants.CONTAINER_NAME] == constants.CONTAINER_NAME_SENSE_PYTHON \
                and module[constants.MODULE_TYPE] == constants.MT_BH1750:
            for included_sensor in module[constants.INCLUDED_SENSORS]:
                if 'light' in included_sensor[constants.MEASUREMENT_NAME]:
                    light_sensor_name = included_sensor[constants.SENSOR_NAME]
                    light_measurement_name = included_sensor[constants.MEASUREMENT_NAME]


"""
TESTED
"""


def bme280_names(device):
    global my_site
#    global my_device
    global humidity_sensor_name
    global pressure_sensor_name
    global temperature_sensor_name
    global temperature_measurement_name
    global humidity_measurement_name
    global pressure_measurement_name
    global light_sensor_name
    global light_measurement_name

    print(my_device)
    for module in device[constants.MODULES]:
        if module[constants.CONTAINER_NAME] == constants.CONTAINER_NAME_SENSE_PYTHON and module[constants.MODULE_TYPE] == constants.MT_BMP280:
            for included_sensor in module[constants.INCLUDED_SENSORS]:
                if 'temp' in included_sensor[constants.MEASUREMENT_NAME]:
                    temperature_sensor_name = included_sensor[constants.SENSOR_NAME]
                    temperature_measurement_name = included_sensor[constants.MEASUREMENT_NAME]
                if 'pressure' in included_sensor[constants.MEASUREMENT_NAME]:
                    pressure_sensor_name = included_sensor[constants.SENSOR_NAME]
                    pressure_measurement_name = included_sensor[constants.MEASUREMENT_NAME]

        if module[constants.CONTAINER_NAME] == constants.CONTAINER_NAME_SENSE_PYTHON and module[constants.MODULE_TYPE] == constants.MT_BME280:
            for included_sensor in module[constants.INCLUDED_SENSORS]:
                if 'temp' in included_sensor[constants.MEASUREMENT_NAME]:
                    temperature_sensor_name = included_sensor[constants.SENSOR_NAME]
                    temperature_measurement_name = included_sensor[constants.MEASUREMENT_NAME]
                if 'pressure' in included_sensor[constants.MEASUREMENT_NAME]:
                    pressure_sensor_name = included_sensor[constants.SENSOR_NAME]
                    pressure_measurement_name = included_sensor[constants.MEASUREMENT_NAME]
                if 'humidity' in included_sensor[constants.MEASUREMENT_NAME]:
                    humidity_sensor_name = included_sensor[constants.SENSOR_NAME]
                    humidity_measurement_name = included_sensor[constants.MEASUREMENT_NAME]


def append_bh1750_data(msg, sensor_name, measurement_name):
    global LightAddress
    global lastLight
    try:
        msg[constants.MM_SENSOR_NAME] = sensor_name
        msg[constants.MM_MEASUREMENT_NAME] = measurement_name
        msg[constants.MM_VALUE_NAME] = measurement_name
        msg[constants.MM_UNITS] = constants.MM_UNITS_LUX
        msg[constants.MM_VALUE] = bh1750.read_light(LightAddress)
        msg[constants.MM_FLOAT_VALUE] = msg[constants.MM_VALUE]
        msg[measurement_name] = msg[constants.MM_VALUE]
        direction = constants.DIRECTIONS_NONE
        if msg[constants.MM_VALUE] > lastLight:
            direction = constants.DIRECTIONS_UP
        else:
            if msg[constants.MM_VALUE] < lastLight:
                direction = constants.DIRECTIONS_DOWN
        msg[constants.MM_DIRECTION] = direction
        direction_name = measurement_name + '_direction'
        msg[direction_name] = direction
        lastLight = msg[constants.MM_VALUE]

    except Exception as ee:
        logging.debug('BH1750 at 0x%2x failed to read %s' % (LightAddress, ee))
        logging.debug(traceback.format_exc())


def append_adc_data(msg):
    msg[constants.MM_SENSOR_NAME] = constants.MM_SENSOR_NAME_WATER_TEMPERATURE
    msg[constants.MM_MEASUREMENT_NAME] = constants.MM_MEASUREMENT_NAME_WATER_TEMPERATURE
    msg[constants.MM_VALUE_NAME] = constants.MM_VALUE_NAME_WATER_TEMPERATURE
    msg[constants.MM_UNITS] = constants.MM_UNITS_GALLONS  # TODO = water temp in gallons?
    msg[constants.MM_VALUE] = 0.0
    msg[constants.MM_FLOAT_VALUE] = 0.0
    msg[constants.MM_WATER_TEMPERATURE] = 0.0


def append_gpio_data(msg):
    msg[constants.MM_SENSOR_NAME] = constants.MM_SENSOR_NAME_WATER_TEMPERATURE
    msg[constants.MM_MEASUREMENT_NAME] = constants.MM_MEASUREMENT_NAME_WATER_TEMPERATURE
    msg[constants.MM_VALUE_NAME] = constants.MM_VALUE_NAME_WATER_TEMPERATURE
    msg[constants.MM_UNITS] = constants.MM_UNITS_GALLONS  # TODO = water temp in gallons?
    msg[constants.MM_VALUE] = 0.0
    msg[constants.MM_FLOAT_VALUE] = 0.0
    msg[constants.MM_WATER_TEMPERATURE] = 0.0
    msg[constants.MM_DOOR_OPEN] = False
    msg[constants.MM_LEAK_DETECTOR] = False


def append_axl345_data(msg):
    msg[constants.MM_SENSOR_NAME] = constants.MM_SENSOR_NAME_TAMPER_DETECTOR
    msg[constants.MM_MEASUREMENT_NAME] = constants.MM_MEASUREMENT_NAME_TAMPER
    msg[constants.MM_VALUE_NAME] = constants.MM_VALUE_NAME_TAMPER
    msg[constants.MM_UNITS] = constants.MM_UNITS_BOOLEAN
    msg[constants.MM_FLOAT_VALUE] = 0.0
    msg[constants.MM_TAMPER_DETECTOR] = False
    msg[constants.MM_VALUE] = False


def any_thermometers_enabled(local_station):
    if local_station[constants.CAPABILITY_TEMPERATURE_TOP]:
        return True
    if local_station[constants.CAPABILITY_TEMPERATURE_MIDDLE]:
        return True
    if local_station[constants.CAPABILITY_TEMPERATURE_BOTTOM]:
        return True
    if local_station[constants.CAPABILITY_TEMPERATURE_EXTERNAL]:
        return True
    return False


def any_humidity_enabled(local_station):
    if local_station[constants.CAPABILITY_HUMIDITY_SENSOR_INTERNAL]:
        return True
    if local_station[constants.CAPABILITY_HUMIDITY_SENSOR_EXTERNAL]:
        return True
    return False


def any_pressure(local_station):
    if local_station[constants.CAPABILITY_PRESSURE_SENSORS]:
        return True
    return False


def report_polled_sensor_parameters(i2cbus):
    global my_site
    global my_station
    global light_sensor_name
    global light_measurement_name

    logging.debug('reportPolledSensorParameters')

    if is_our_device(constants.MT_BH1750) and (my_station[constants.CAPABILITY_LIGHT_SENSOR_INTERNAL] or
                                               my_station[constants.CAPABILITY_LIGHT_SENSOR_EXTERNAL]):
        module = get_our_device(module_type=constants.MT_BH1750)
        msg = {
            constants.MM_MESSAGE_TYPE: constants.MM_MESSAGE_TYPE_MEASUREMENT
        }
        # If reading the sensor hardware fails, pass the exception up here and
        # we'll skip sending the half-complete message
        try:
            append_bh1750_data(msg=msg, sensor_name=light_sensor_name, measurement_name=light_measurement_name)
            send_message(msg=msg)
        except Exception as ee:
            logging.error(ee)

    if is_our_device(constants.MT_BMP280) and (
            any_thermometers_enabled(local_station=my_station) or
            any_humidity_enabled(local_station=my_station) or
            any_pressure(local_station=my_station)):
        module = get_our_device(module_type=constants.MT_BMP280)

        logging.info('found a bmp280 - no humidity!')
        msg = {constants.MM_MESSAGE_TYPE: constants.MM_MESSAGE_TYPE_MEASUREMENT}
        append_bme280_temp(i2cbus=i2cbus, msg=msg, sensor_name=temperature_sensor_name,
                           measurement_name=temperature_measurement_name, address=int(module[constants.ADDRESS], 0))
        send_message(msg=msg)

        msg = {constants.MM_MESSAGE_TYPE: constants.MM_MESSAGE_TYPE_MEASUREMENT}
        append_bme280_pressure(i2cbus=i2cbus, msg=msg, sensor_name=pressure_sensor_name,
                               measurement_name=pressure_measurement_name, address=int(module[constants.ADDRESS], 0))
        send_message(msg=msg)

    if is_our_device(constants.MT_BME280) and (
            any_thermometers_enabled(local_station=my_station) or
            any_humidity_enabled(local_station=my_station) or
            any_pressure(local_station=my_station)):
        module = get_our_device(module_type=constants.MT_BME280)
        msg = {constants.MM_MESSAGE_TYPE: constants.MM_MESSAGE_TYPE_MEASUREMENT}
        append_bme280_temp(i2cbus=i2cbus, msg=msg, sensor_name=temperature_sensor_name,
                           measurement_name=temperature_measurement_name, address=int(module[constants.ADDRESS], 0))
        send_message(msg=msg)

        msg = {constants.MM_MESSAGE_TYPE: constants.MM_MESSAGE_TYPE_MEASUREMENT}
        append_bme280_humidity(i2cbus=i2cbus, msg=msg, sensor_name=humidity_sensor_name,
                               measurement_name=humidity_measurement_name, address=int(module[constants.ADDRESS], 0))
        send_message(msg=msg)

        msg = {constants.MM_MESSAGE_TYPE: constants.MM_MESSAGE_TYPE_MEASUREMENT}
        append_bme280_pressure(i2cbus=i2cbus, msg=msg, sensor_name=pressure_sensor_name,
                               measurement_name=pressure_measurement_name, address=int(module[constants.ADDRESS], 0))
        send_message(msg=msg)

    if is_our_device(constants.MT_ADS1115):
        msg = {constants.MM_MESSAGE_TYPE: constants.MM_MESSAGE_TYPE_MEASUREMENT}
        append_adc_data(msg=msg)
        send_message(msg=msg)

    if is_our_device(constants.MT_ADXL345) and my_station[constants.CAPABILITY_MOVEMENT_SENSOR]:
        msg = {constants.MM_MESSAGE_TYPE: constants.MM_MESSAGE_TYPE_MEASUREMENT}
        append_axl345_data(msg=msg)
        send_message(msg=msg)

    if is_our_device(constants.MT_RELAY):
        msg = {constants.MM_MESSAGE_TYPE: constants.MM_MESSAGE_TYPE_MEASUREMENT}
        append_gpio_data(msg=msg)
        send_message(msg=msg)

    return


sequence = 200000

"""
TESTED
"""


def get_sequence():
    global sequence
    if sequence >= 3000000:
        sequence = 200000
    else:
        sequence = sequence + 1
    return sequence


def send_message(msg):
    global my_site
    logging.info('siteid ' + format(my_site[constants.SITEID], 'd') + ' stationid ' + format(my_station[constants.STATIONID],
                                        'd') + " deviceid " + format(my_device[constants.DEVICEID], 'd'))
    seq = get_sequence()
    millis = int(time.time() * 1000)
    msg[constants.MM_SAMPLE_TIMESTAMP] = int(millis)
    msg[constants.MM_DEVICEID] = my_device[constants.DEVICEID]
    msg[constants.MM_STATIONID] = my_station[constants.STATIONID]
    msg[constants.MM_SITEID] = my_site[constants.SITEID]
    msg[constants.MM_CONTAINER_NAME] = constants.CONTAINER_NAME_SENSE_PYTHON
    msg[constants.MM_EXECUTABLE_VERSION] = '9.9.10'
    json_bytes = str.encode(json.dumps(msg))
    logging.debug(json_bytes)
    message = bubblesgrpc_pb2.SensorRequest(sequence=seq, type_id=constants.MESSAGE_TYPE_ID_SENSOR, data=json_bytes)
    response = stub.StoreAndForward(message)
    #    logging.debug(response)
    return


if __name__ == "__main__":

    logging.basicConfig(level=logging.DEBUG)

    naptime_in_seconds = int(os.environ[constants.ENV_SLEEP_ON_EXIT_FOR_DEBUGGING])

    logging.info('Starting sense-python')

    logging.info(f'naptime_in_seconds = {naptime_in_seconds} seconds')

    my_site[constants.DEVICEID] = read_deviceid(filename='/config/deviceid')
    logging.info(f'deviceid from file is {my_site[constants.DEVICEID]:d}')
    wait_for_config(filename='/config/config.json')  # wait for the config file to exist, exit directly if it times out
    b = read_config(fullpath='/config/config.json', deviceid=my_site[constants.DEVICEID])  # config file exists, read it in
    if not b:
        logging.error(f'invalid config.json - not validating - exiting after {naptime_in_seconds} seconds')
        time.sleep(naptime_in_seconds)
        exit(1)

    bme280_names(my_device)
    bh1750_names()
    LightAddress = get_address(module_type=constants.MT_BH1750)

    # Create library object using our Bus I2C port
    i2c = busio.I2C(board.SCL, board.SDA)
    bus_number = 1
    bus = smbus2.SMBus(bus_number)

    try:
        logging.debug('Connecting to grpc at store-and-forward:50051')
        channel = grpcio.insecure_channel('store-and-forward:50051')
        stub = grpcStub(channel)
        try:
            logging.debug('Entering sensor polling loop')
            while True:
                #                        toggleRelays(relay,sequence)

                report_polled_sensor_parameters(i2cbus=bus)

                #                logging.debug("sleeping %d xx seconds at %s"
                #                % (config['time_between_sensor_polling_in_seconds'],
                #                time.strftime("%T")))
                time.sleep(my_device['time_between_sensor_polling_in_seconds'])

            logging.debug('broke out of temp/hum/distance polling loop')
        except Exception as e:
            logging.debug('bubbles2 main loop failed')
            logging.debug(traceback.format_exc())
    except Exception as eee:
        logging.debug('GRPC failed to initialize')
    logging.debug('end of main')

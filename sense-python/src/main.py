#!/usr/bin/env python3

import json
import time
import traceback
import logging
import datetime

import bme280
import board
import busio
import grpc as grpcio
import smbus2

import bh1750
from bubblesgrpc_pb2 import SensorRequest
from bubblesgrpc_pb2_grpc import SensorStoreAndForwardStub as grpcStub

lastTemp = 0.0
config = {}
LightAddress = 0
sequence = 0

lastHumidity = 0.0
lastPressure = 0.0
lastLight = 0.0


def read_config():
    global config
    with open('/config/config.json') as f:
        config = json.load(f)


def append_bme280_temp(bus, msg, value_name):
    global lastTemp
    try:
        calibration_params = bme280.load_calibration_params(bus, 0x76)
        data = bme280.sample(bus, 0x76, compensation_params=calibration_params)
        msg['sensor_name'] = value_name
        msg['value'] = (data.temperature * 1.8) + 32.0
        msg['direction'] = ''
        if data.temperature > lastTemp:
            msg['direction'] = "up"
        else:
            if data.temperature < lastTemp:
                msg['direction'] = "down"
        lastTemp = data.temperature
        msg['tempC'] = data.temperature
        msg['tempF'] = (data.temperature * 1.8) + 32.0
#        msg['humidity'] = data.humidity
#        msg['pressure'] = data.pressure
    #        logging.debug( "bme280 read temp at 0x76 as %f" % data.temperature )
    except Exception as e:
        logging.debug("bme280 error %s" % e)
        logging.debug(traceback.format_exc())


def append_bme280_humidity(bus, msg, value_name):
    global lastHumidity
    try:
        calibration_params = bme280.load_calibration_params(bus, 0x76)
        data = bme280.sample(bus, 0x76, compensation_params=calibration_params)
        msg['sensor_name'] = value_name
        msg['value'] = data.humidity
        direction = ""
        if data.humidity > lastHumidity:
            direction = "up"
        else:
            if data.humidity < lastHumidity:
                direction = "down"
        msg['direction'] = direction
        lastHumidity = data.humidity

    #        logging.debug( "bme280 read temp at 0x76 as %f" % data.temperature )
    except Exception as e:
        logging.debug("bme280 error %s" % e)
        logging.debug(traceback.format_exc())


def append_bme280_pressure(bus, msg, value_name):
    global lastPressure
    try:
        calibration_params = bme280.load_calibration_params(bus, 0x76)
        data = bme280.sample(bus, 0x76, compensation_params=calibration_params)
        msg['sensor_name'] = value_name
        msg['value'] = data.pressure
        direction = ""
        if data.pressure > lastPressure:
            direction = "up"
        else:
            if data.pressure < lastPressure:
                direction = "down"
        lastPressure = data.pressure
        msg['direction'] = direction
        lastPressure = data.humidity
#        logging.debug( "bme280 read temp at 0x76 as %f" % data.temperature )
    except Exception as e:
        logging.debug("bme280 error %s" % e)
        logging.debug(traceback.format_exc())


# new config functions
def get_address(device_type):
    for attached_device in config['attached_devices']:
        if attached_device['container_name'] == 'sense-python' and attached_device['device_type'] == device_type:
            x = int(attached_device['address'], 16)
            return x
    return 0


def is_our_device(device_type):
    ret = False
    for attached_device in config['attached_devices']:
        if attached_device['container_name'] == 'sense-python' and attached_device['device_type'] == device_type:
            return True

    return ret


# end new config functions


def append_bh1750_data(msg):
    global LightAddress
    global lastLight
    try:
        msg['sensor_name'] = 'light_internal'
        msg['value'] = bh1750.readLight(LightAddress)
        direction = ""
        if msg['value'] > lastLight:
            direction = "up"
        else:
            if msg['value'] < lastLight:
                direction = "down"
        msg['direction'] = direction
        lastLight = msg['value']

        #        logging.debug("Read bh1750 light at 0x%x as %f" % (deviceAddressList['light'], msg['light']))
    except Exception as e:
        logging.debug('BH1750 at 0x%2x failed to read %s' % (LightAddress, e))
        logging.debug(traceback.format_exc())


def append_adc_data(msg):
    msg['water_temperature'] = 0.0


def append_gpio_data(msg):
    msg['door_open'] = False
    msg['leak_detector'] = False


def append_axl345_data(msg):
    msg['tamper_detector'] = False


def report_polled_sensor_parameters(bus, sequence):
    global config
    #    logging.debug("reportPolledSensorParameters")

    if is_our_device('bh1750'):
        msg = {}
        msg['message_type'] = 'measurement'
        append_bh1750_data(msg)
        send_message(msg)
    if is_our_device('bme280'):
        msg = {}
        msg['message_type'] = 'measurement'
        append_bme280_temp(bus, msg, 'temp_air_middle')
        send_message(msg)
        msg = {}
        msg['message_type'] = 'measurement'
        append_bme280_humidity(bus, msg, 'humidity_internal')
        send_message(msg)
        msg = {}
        msg['message_type'] = 'measurement'
        append_bme280_pressure(bus, msg, 'pressure_internal')
        send_message(msg)
    if is_our_device('ads1115'):
        append_adc_data(msg)
        msg = {}
        msg['message_type'] = 'measurement'
        send_message(msg)
    if is_our_device('adxl345'):
        msg = {}
        msg['message_type'] = 'measurement'
        append_axl345_data(msg)
        send_message(msg)
    if is_our_device('relay'):
        msg = {}
        msg['message_type'] = 'measurement'
        append_gpio_data(msg)
        send_message(msg)
    return


def send_message(msg):
    global sequence
    if sequence > 100000:
        sequence = 1
    else:
        sequence = sequence + 1

    millis = int(time.time()*1000)
    msg['sample_timestamp'] = int(millis)
    json_bytes = str.encode(json.dumps(msg))
    logging.debug(json_bytes)
    message = SensorRequest(sequence=sequence, type_id="sensor", data=json_bytes)
    response = stub.StoreAndForward(message)
    logging.debug(response)
    return


if __name__ == "__main__":
    logging.basicConfig(level=logging.DEBUG)

    logging.debug("Starting sense-python")
    read_config()
    LightAddress = get_address('bh1750')

    # Create library object using our Bus I2C port
    i2c = busio.I2C(board.SCL, board.SDA)
    bus_number = 1
    bus = smbus2.SMBus(bus_number)

    try:
        logging.debug("Connecting to grpc at store-and-forward:50051")
        channel = grpcio.insecure_channel('store-and-forward:50051')
        stub = grpcStub(channel)
        try:
            logging.debug("Entering sensor polling loop")
            while True:
                #                        toggleRelays(relay,sequence)
                report_polled_sensor_parameters(bus, sequence)

                logging.debug("sleeping %d xx seconds at %s" % (
                config['time_between_sensor_polling_in_seconds'], time.strftime("%T")))
                time.sleep(config['time_between_sensor_polling_in_seconds'])

            logging.debug("broke out of temp/hum/distance polling loop")
        except Exception as e:
            logging.debug('bubbles2 main loop failed')
            logging.debug(traceback.format_exc())
    except:
        logging.debug('GRPC failed to initialize')
    logging.debug("end of main")

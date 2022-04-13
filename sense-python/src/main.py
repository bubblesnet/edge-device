#!/usr/bin/env python3

# copyright and license inspection - no issues 4/13/22

import json
import logging
import time
import traceback

try:
    import bme280
except ImportError:
    bme280 = ""

try:
    import board
except NotImplementedError:
    board = ""
except ImportError:
    board = ""
except AttributeError:
    board = ""

try:
    import busio
except ImportError:
    busio = ""

try:
    import grpc as grpcio
except ImportError:
    grpcio = ""

try:
    import smbus2
except ImportError:
    smbus2 = ""


try:
    import bh1750
except ImportError:
    bh1750 = ""

import bubblesgrpc_pb2
from bubblesgrpc_pb2_grpc import SensorStoreAndForwardStub as grpcStub
from os.path import exists

lastTemp = 0.0
global my_site
my_site = {}
global my_station
my_station = {}
global my_device
my_device = {}
LightAddress = 0

lastHumidity = 0.0
lastPressure = 0.0
lastLight = 0.0

global humidity_sensor_name
global pressure_sensor_name
global temperature_sensor_name
global temperature_measurement_name
global humidity_measurement_name
global pressure_measurement_name
global light_sensor_name
global light_measurement_name


def wait_for_config(filename):
    logging.info('wait_for_config %s' % filename)
    index = 0
    while index <= 60:
        if exists(filename):
            logging.info('%s file exists' % filename)
            return
        logging.info("Sleeping while we wait for someone to create %s" % filename)
        time.sleep(60)
        index = index+1
        if index == 60:
            logging.error("config file never showed up, exiting")
            exit(1)

    return


def read_deviceid(filename):
    deviceid = -1
    with open(filename) as f:
        for line in f:
            logging.info("line = %s " % line )
            deviceid = int(line)
        return (deviceid)


def validate_config():
    b = 'stations' in my_site
    if not b:
        return False

    return True


def read_config(fullpath):
    global my_site, my_station, my_device
    deviceid = my_site['deviceid']
    with open(fullpath) as f:
        my_site = json.load(f)
        my_site['deviceid'] = deviceid
        if not validate_config():
            return False
        for station in my_site['stations']:
            if 'edge_devices' not in station:
                continue
            for device in station['edge_devices']:
                if device['deviceid'] == deviceid:
                    my_device = device
                    my_station = station
                    my_site['time_between_sensor_polling_in_seconds'] = 15
        return True


def append_bme280_temp(i2cbus, msg, sensor_name, measurement_name):
    global lastTemp
    try:
        calibration_params = bme280.load_calibration_params(i2cbus, 0x76)
        data = bme280.sample(i2cbus, 0x76, compensation_params=calibration_params)
        msg['sensor_name'] = sensor_name
        msg['measurement_name'] = measurement_name
        msg['units'] = 'C'
        msg['value'] = (data.temperature * 1.8) + 32.0
        msg[measurement_name] = msg['value']
        msg['floatvalue'] = msg['value']
        msg['value_name'] = measurement_name
        msg['direction'] = ''
        if data.temperature > lastTemp:
            msg['direction'] = "up"
        else:
            if data.temperature < lastTemp:
                msg['direction'] = "down"
        direction_name = measurement_name + "_direction"
        msg[direction_name] = msg['direction']
        lastTemp = data.temperature
        msg['tempC'] = data.temperature
        msg['tempF'] = (data.temperature * 1.8) + 32.0
    except Exception as ee:
        logging.debug("bme280 error %s" % ee)
        logging.debug(traceback.format_exc())


def append_bme280_humidity(i2cbus, msg, sensor_name, measurement_name):
    global lastHumidity
    try:
        calibration_params = bme280.load_calibration_params(i2cbus, 0x76)
        data = bme280.sample(i2cbus, 0x76, compensation_params=calibration_params)
        msg['sensor_name'] = sensor_name
        msg['measurement_name'] = measurement_name
        msg['units'] = '%'
        msg['value_name'] = measurement_name
        msg['value'] = data.humidity
        msg['floatvalue'] = msg['value']
        msg[measurement_name] = data.humidity
        direction = ""
        if data.humidity > lastHumidity:
            direction = "up"
        else:
            if data.humidity < lastHumidity:
                direction = "down"
        msg['direction'] = direction
        direction_name = measurement_name + "_direction"
        msg[direction_name] = direction
        lastHumidity = data.humidity
    except Exception as ee:
        logging.debug("bme280 error %s" % ee)
        logging.debug(traceback.format_exc())


def append_bme280_pressure(i2cbus, msg, sensor_name, measurement_name):
    global lastPressure
    try:
        calibration_params = bme280.load_calibration_params(i2cbus, 0x76)
        data = bme280.sample(i2cbus, 0x76, compensation_params=calibration_params)
        msg['sensor_name'] = sensor_name
        msg['measurement_name'] = measurement_name
        msg['units'] = 'hPa'
        msg['value_name'] = measurement_name
        msg['value'] = data.pressure
        msg['floatvalue'] = msg['value']
        msg[measurement_name] = data.pressure
        direction = ""
        if data.pressure > lastPressure:
            direction = "up"
        else:
            if data.pressure < lastPressure:
                direction = "down"
        lastPressure = data.pressure
        msg['direction'] = direction
        direction_name = measurement_name + "_direction"
        msg[direction_name] = direction
        lastPressure = data.humidity
    except Exception as ee:
        logging.debug("bme280 error %s" % ee)
        logging.debug(traceback.format_exc())


# new config functions
def get_address(module_type):
    for module in my_device['modules']:
        if module['container_name'] == 'sense-python' and module['module_type'] == module_type:
            x = int(module['address'], 16)
            return x
    return 0


def is_our_device(module_type):
    for module in my_device['modules']:
        if module['container_name'] == 'sense-python' and module['module_type'] == module_type:
            return True

    return False


def bh1750_names():
    global my_site
    global my_device
    global light_sensor_name
    global light_measurement_name

    print(my_device)
    for module in my_device['modules']:
        if module['container_name'] == 'sense-python' and module['module_type'] == "bh1750":
            for included_sensor in module['included_sensors']:
                if 'light' in included_sensor['measurement_name']:
                    light_sensor_name = included_sensor['sensor_name']
                    light_measurement_name = included_sensor['measurement_name']


def bme280_names():
    global my_site
    global my_device
    global humidity_sensor_name
    global pressure_sensor_name
    global temperature_sensor_name
    global temperature_measurement_name
    global humidity_measurement_name
    global pressure_measurement_name
    global light_sensor_name
    global light_measurement_name

    print(my_device)
    for module in my_device['modules']:
        if module['container_name'] == 'sense-python' and module['module_type'] == "bme280":
                for included_sensor in module['included_sensors']:
                    if 'temp' in included_sensor['measurement_name']:
                        temperature_sensor_name = included_sensor['sensor_name']
                        temperature_measurement_name = included_sensor['measurement_name']
                    if 'pressure' in included_sensor['measurement_name']:
                        pressure_sensor_name = included_sensor['sensor_name']
                        pressure_measurement_name = included_sensor['measurement_name']
                    if 'humidity' in included_sensor['measurement_name']:
                        humidity_sensor_name = included_sensor['sensor_name']
                        humidity_measurement_name = included_sensor['measurement_name']


def append_bh1750_data(msg, sensor_name, measurement_name):
    global LightAddress
    global lastLight
    try:
        msg['sensor_name'] = sensor_name
        msg['measurement_name'] = measurement_name
        msg['value_name'] = measurement_name
        msg['units'] = 'lux'
        msg['value'] = bh1750.readLight(LightAddress)
        msg['floatvalue'] = msg['value']
        msg[measurement_name] = msg['value']
        direction = ""
        if msg['value'] > lastLight:
            direction = "up"
        else:
            if msg['value'] < lastLight:
                direction = "down"
        msg['direction'] = direction
        direction_name = measurement_name + "_direction"
        msg[direction_name] = direction
        lastLight = msg['value']

    except Exception as ee:
        logging.debug('BH1750 at 0x%2x failed to read %s' % (LightAddress, ee))
        logging.debug(traceback.format_exc())

def append_adc_data(msg):
    msg['sensor_name'] = 'water_temperature_sensor'
    msg['measurement_name'] = 'temp_water'
    msg['value_name'] = 'temp_water'
    msg['units'] = 'gallons'
    msg['value'] = 0.0
    msg['floatvalue'] = 0.0
    msg['water_temperature'] = 0.0

def append_gpio_data(msg):
    msg['sensor_name'] = 'water_temperature_sensor'
    msg['measurement_name'] = 'temp_water'
    msg['value_name'] = 'temp_water'
    msg['units'] = 'gallons'
    msg['value'] = 0.0
    msg['floatvalue'] = 0.0
    msg['water_temperature'] = 0.0
    msg['door_open'] = False
    msg['leak_detector'] = False


def append_axl345_data(msg):
    msg['sensor_name'] = 'tamper_detector'
    msg['measurement_name'] = 'tamper'
    msg['value_name'] = 'tamper'
    msg['units'] = 'boolean'
    msg['floatvalue'] = 0.0

    msg['tamper_detector'] = False
    msg['value'] = False


def report_polled_sensor_parameters(i2cbus):
    global my_site
    global light_sensor_name
    global light_measurement_name

    #    logging.debug("reportPolledSensorParameters")

    if is_our_device('bh1750'):
        msg = {
            'message_type': 'measurement'
        }
        # If reading the sensor hardware fails, pass the exception up here and
        # we'll skip sending the half-complete message
        try:
            append_bh1750_data(msg, light_sensor_name, light_measurement_name)
            send_message(msg)
        except Exception as ee:
            logging.error(ee)

    if is_our_device('bme280'):
        msg = {'message_type': 'measurement'}
        append_bme280_temp(i2cbus, msg, temperature_sensor_name, temperature_measurement_name)
        send_message(msg)

        msg = {'message_type': 'measurement'}
        append_bme280_humidity(i2cbus, msg, humidity_sensor_name, humidity_measurement_name)
        send_message(msg)

        msg = {'message_type': 'measurement'}
        append_bme280_pressure(i2cbus, msg, pressure_sensor_name, pressure_measurement_name)
        send_message(msg)

    if is_our_device('ads1115'):
        msg = {'message_type': 'measurement'}
        append_adc_data(msg)
        send_message(msg)

    if is_our_device('adxl345'):
        msg = {'message_type': 'measurement'}
        append_axl345_data(msg)
        send_message(msg)

    if is_our_device('relay'):
        msg = {'message_type': 'measurement'}
        append_gpio_data(msg)
        send_message(msg)

    return


sequence = 200000


def get_sequence():
    global sequence
    if sequence >= 3000000:
        sequence = 200000
    else:
        sequence = sequence + 1
    return sequence


def send_message(msg):
    global my_site
    logging.info("siteid "+format(my_site['siteid'], 'd')+" stationid "+format(my_station['stationid'], 'd')+" deviceid "+format(my_device['deviceid'], 'd'))
    seq = get_sequence()
    millis = int(time.time() * 1000)
    msg['sample_timestamp'] = int(millis)
    msg['deviceid'] = my_device['deviceid']
    msg['stationid'] = my_station['stationid']
    msg['siteid'] = my_site['siteid']
    msg['container_name'] = "sense-python"
    msg['executable_version'] = "9.9.10"
    json_bytes = str.encode(json.dumps(msg))
    logging.debug(json_bytes)
    message = bubblesgrpc_pb2.SensorRequest(sequence=seq, type_id="sensor", data=json_bytes)
    response = stub.StoreAndForward(message)
    #    logging.debug(response)
    return


if __name__ == "__main__":


    logging.basicConfig(level=logging.DEBUG)

    logging.debug("Starting sense-python")
    my_site['deviceid'] = read_deviceid('/config/deviceid')
    logging.info("deviceid from file is %d" % my_site['deviceid'])
    wait_for_config('/config/config.json')
    b = read_config('/config/config.json')
    if not b:
        logging.error('invalid config.json - not validating')
        time.sleep(60)
        exit(1)

    bme280_names()
    bh1750_names()
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

                report_polled_sensor_parameters(bus)

                #                logging.debug("sleeping %d xx seconds at %s" % (config['time_between_sensor_polling_in_seconds'],
                #                time.strftime("%T")))
                time.sleep(my_site['time_between_sensor_polling_in_seconds'])

            logging.debug("broke out of temp/hum/distance polling loop")
        except Exception as e:
            logging.debug('bubbles2 main loop failed')
            logging.debug(traceback.format_exc())
    except Exception as eee:
        logging.debug('GRPC failed to initialize')
    logging.debug("end of main")
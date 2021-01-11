import json
import time
import traceback
import logging

import bme280
import board
import busio
import grpc as grpcio
import smbus2

import bh1750
from bubblesgrpc_pb2 import SensorRequest
from bubblesgrpc_pb2_grpc import SensorStoreAndForwardStub as grpcStub


config = {}


def read_config():
    global config
    with open('/config/config.json') as f:
        config = json.load(f)


def append_bme280_data(bus, msg):
    try:
        calibration_params = bme280.load_calibration_params(bus, 0x76)
        data = bme280.sample(bus, 0x76, compensation_params=calibration_params)
        msg['temperature'] = data.temperature
        msg['tempC'] = data.temperature
        msg['tempF'] = (data.temperature * 1.8) + 32.0
        msg['humidity'] = data.humidity
        msg['pressure'] = data.pressure
    #        logging.debug( "bme280 read temp at 0x76 as %f" % data.temperature )
    except Exception as e:
        logging.debug("bme280 error %s" % e)
        logging.debug(traceback.format_exc())


# new config functions
def get_address(device_type):
    for attached_device in config['attached_devices']:
        if attached_device['container_name'] == 'sense-python' and attached_device['device_type'] == device_type:
            x = int(attached_device['address'])
            return True
    return 0


def is_our_device(device_type):
    ret = False
    for attached_device in config['attached_devices']:
        if attached_device['container_name'] == 'sense-python' and attached_device['device_type'] == device_type:
            return True

    return ret


LightAddress = get_address('bh1750')
# end new config functions


def append_bh1750_data(msg):
    try:
        msg['light'] = bh1750.readLight(LightAddress)
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

    msg = {}
    msg['sample_timestamp'] = int(time.time() * 1000)
    if is_our_device('bme280'):
        append_bme280_data(bus, msg)
    if is_our_device('bh1750'):
        append_bh1750_data(msg)
    if is_our_device('ads1115'):
        append_adc_data(msg)
    if is_our_device('adxl345'):
        append_axl345_data(msg)
    if is_our_device('relay'):
        append_gpio_data(msg)

    json_bytes = str.encode(json.dumps(msg))
    #    logging.debug(json_bytes)
    message = SensorRequest(sequence=sequence, type_id="sensor", data=json_bytes)
    response = stub.StoreAndForward(message)
    #    logging.debug(response)
    return


if __name__ == "__main__":
    logging.basicConfig(level=logging.DEBUG)

    logging.debug("Starting sense-python")
    read_config()

    # Create library object using our Bus I2C port
    i2c = busio.I2C(board.SCL, board.SDA)
    bus_number = 1
    bus = smbus2.SMBus(bus_number)

    try:
        logging.debug("Connecting to grpc at store-and-forward:50051")
        channel = grpcio.insecure_channel('store-and-forward:50051')
        stub = grpcStub(channel)
        try:
            sequence = 0
            logging.debug("Entering sensor polling loop")
            while True:
                #                        toggleRelays(relay,sequence)
                report_polled_sensor_parameters(bus, sequence)
                if sequence > 100000:
                    sequence = 1
                else:
                    sequence = sequence + 1

# logging.debug("sleeping %d xx seconds at %s" % (config['time_between_sensor_polling_in_seconds'],time.strftime("%T")))
                time.sleep(config['time_between_sensor_polling_in_seconds'])

            logging.debug("broke out of temp/hum/distance polling loop")
        except Exception as e:
            logging.debug('bubbles2 main loop failed')
            logging.debug(traceback.format_exc())
    except:
        logging.debug('GRPC failed to initialize')
    logging.debug("end of main")

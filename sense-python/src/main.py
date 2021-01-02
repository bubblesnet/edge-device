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


def appendBME280Data(bus, msg):
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


def appendBH1750Data(msg):
    try:
        msg['light'] = bh1750.readLight(deviceAddressList['light'])
    #        logging.debug("Read bh1750 light at 0x%x as %f" % (deviceAddressList['light'], msg['light']))
    except Exception as e:
        logging.debug('BH1750 at 0x%2x failed to read %s' % (deviceAddressList['light'], e))
        logging.debug(traceback.format_exc())


def appendADCData(msg):
    msg['water_temperature'] = 0.0


def appendGPIOData(msg):
    msg['door_open'] = False
    msg['leak_detector'] = False


def appendAXL345Data(msg):
    msg['tamper_detector'] = False


def reportPolledSensorParameters(bus, sequence):
    global config
    #    logging.debug("reportPolledSensorParameters")

    msg = {}
    msg['sample_timestamp'] = int(time.time() * 1000)
    if (config['bme280']):
        appendBME280Data(bus, msg)
    if (config['bh1750']):
        appendBH1750Data(msg)
    if (config['ads1115_1'] or config['ads1115_2']):
        appendADCData(msg)
    if (config['adxl345']):
        appendAXL345Data(msg)
    if (config['relay']):
        appendGPIOData(msg)

    json_bytes = str.encode(json.dumps(msg))
    #    logging.debug(json_bytes)
    message = SensorRequest(sequence=sequence, type_id="sensor", data=json_bytes)
    response = stub.StoreAndForward(message)
    #    logging.debug(response)
    return


deviceAddressList = {
    "light": 0x23,
    'adc1': 0x48,
    'adc2': 0x49,
    'accelerometer': 0x53,
    'pH': 0x63,
    "temperature": 0x76
}

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
                reportPolledSensorParameters(bus, sequence)
                if sequence > 100000:
                    sequence = 1
                else:
                    sequence = sequence + 1

                #                    logging.debug("sleeping %d xx seconds at %s" % (config['time_between_sensor_polling_in_seconds'],time.strftime("%T")))
                time.sleep(config['time_between_sensor_polling_in_seconds'])

            logging.debug("broke out of temp/hum/distance polling loop")
        except Exception as e:
            logging.debug('bubbles2 main loop failed')
            logging.debug(traceback.format_exc())
    except:
        logging.debug('GRPC failed to initialize')
    logging.debug("end of main")

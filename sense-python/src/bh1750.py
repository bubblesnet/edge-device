# This is how I found this code and a second look produced
# no further license or copyright information so I present
# the code below as-found.
#
#
# ---------------------------------------------------------------------
#    ___  ___  _ ____
#   / _ \/ _ \(_) __/__  __ __
#  / , _/ ___/ /\ \/ _ \/ // /
# /_/|_/_/  /_/___/ .__/\_, /
#                /_/   /___/
#
#           bh1750.py
# Read data from a BH1750 digital light sensor.
#
# Author : Matt Hawkins
# Date   : 26/06/2018
#
# For more information please visit :
# https://www.raspberrypi-spy.co.uk/?s=bh1750
#
# ---------------------------------------------------------------------
import smbus2
import time
import logging


# Define some constants from the datasheet

DEVICE = 0x23  # Default device I2C address

POWER_DOWN = 0x00  # No active state
POWER_ON = 0x01  # Power on
RESET = 0x07  # Reset data register value

# Start measurement at 4lx resolution. Time typically 16ms.
CONTINUOUS_LOW_RES_MODE = 0x13
# Start measurement at 1lx resolution. Time typically 120ms
CONTINUOUS_HIGH_RES_MODE_1 = 0x10
# Start measurement at 0.5lx resolution. Time typically 120ms
CONTINUOUS_HIGH_RES_MODE_2 = 0x11
# Start measurement at 1lx resolution. Time typically 120ms
# Device is automatically set to Power Down after measurement.
ONE_TIME_HIGH_RES_MODE_1 = 0x20
# Start measurement at 0.5lx resolution. Time typically 120ms
# Device is automatically set to Power Down after measurement.
ONE_TIME_HIGH_RES_MODE_2 = 0x21
# Start measurement at 1lx resolution. Time typically 120ms
# Device is automatically set to Power Down after measurement.
ONE_TIME_LOW_RES_MODE = 0x23

# bus = smbus.SMBus(0) # Rev 1 Pi uses 0
bus = smbus2.SMBus(1)  # Rev 2 Pi uses 1


def convert_to_number(data):
    # Simple function to convert 2 bytes of data
    # into a decimal number. Optional parameter 'decimals'
    # will round to specified number of decimal places.
    result = (data[1] + (256 * data[0])) / 1.2
    return result


def read_light(addr=DEVICE):
    # Read data from I2C interface
    data = bus.read_i2c_block_data(addr, ONE_TIME_HIGH_RES_MODE_1, 32)
    return convert_to_number(data)


def main():
    global my_site

    while True:
        lightLevel = read_light()
        logging.debug('Light Level : ' + format(lightLevel, '.2f') + " lx")
        time.sleep(my_site['time_between_sensor_polling_in_seconds'])


if __name__ == '__main__':
    main()

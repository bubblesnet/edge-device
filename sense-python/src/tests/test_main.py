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

from unittest import TestCase
from .. import main

global my_device

import os


import grpc as grpcio
# from bubblesgrpc_pb2 import SensorRequest
# from bubblesgrpc_pb2_grpc import SensorStoreAndForwardStub as grpcStub


class Test(TestCase):
    def test_validate_config(self):
        b = main.read_config('../config.json', int(os.environ['DEVICEID']))
        self.assertTrue(b)
        b = main.validate_config(main.my_site)
        self.assertTrue(b)

    def test_read_deviceid(self):
        deviceid = main.read_deviceid("asdfsafdf")
        self.assertGreater(deviceid, -1)

    def test_wait_for_config(self):
        b = main.wait_for_config('../config.json')
        self.assertTrue(b)

    def test_read_config(self):
        global my_device

        b = main.read_config('../config.json', int(os.environ['DEVICEID']))
        self.assertTrue(b)
        self.assertEqual(main.my_site['deviceid'], int(os.environ['DEVICEID']))
        self.assertIsNotNone(main.my_device)

    def test_get_address(self):
        b = main.read_config('../config.json', int(os.environ['DEVICEID']))
        self.assertTrue(b)
        addr = main.get_address('bme280')
        self.assertNotEqual(addr, 0)

    def test_is_our_device(self):
        b = main.read_config('../config.json', int(os.environ['DEVICEID']))
        self.assertTrue(b)
        ourdevice = main.is_our_device('bme280')
        self.assertIsNotNone(ourdevice)

    def test_bme280_names(self):
        global my_device

        b = main.read_config('../config.json', int(os.environ['DEVICEID']))
        self.assertTrue(b)
        self.assertIsNotNone(main.my_device)
        main.bme280_names(main.my_device)
        self.assertNotEqual(main.temperature_sensor_name, '')

    def test_get_sequence(self):
        x = main.get_sequence()
        self.assertGreaterEqual(x, 0)

    def test_send_message(self):
        main.channel = grpcio.insecure_channel('localhost:50051')
        main.stub = main.grpcStub(main.channel)

        msg = {}
        main.send_message(msg)
        self.assertGreater(msg['sample_timestamp'], 0)

'''

   def test_append_bh1750_data(self):
       self.fail()

   def test_append_adc_data(self):
       self.fail()

   def test_append_gpio_data(self):
       self.fail()

   def test_append_axl345_data(self):
       self.fail()

   def test_report_polled_sensor_parameters(self):
       self.fail()

   def test_append_bme280_temp(self):
       main.append_bme280_temp('', 'test_thermometer', 'test_temp')
       self.fail()

   def test_append_bme280_humidity(self):
       self.fail()

   def test_append_bme280_pressure(self):
       self.fail()
'''



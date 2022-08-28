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
from . import main
import grpc as grpcio
from bubblesgrpc_pb2 import SensorRequest
from bubblesgrpc_pb2_grpc import SensorStoreAndForwardStub as grpcStub


class Test(TestCase):
    def test_read_config(self):
        main.read_config('../config.json')
        if main.my_site['deviceid'] <= 0:
            self.fail()
        return

    def test_get_address(self):
        main.read_config('../config.json')
        addr = main.get_address("bme280")
        if addr == 0:
            self.fail()

    def test_is_our_device(self):
        main.read_config("../config.json")
        ourdevice = main.is_our_device("bme280")
        if not ourdevice:
            self.fail()

    def test_bme280_names(self):
        main.read_config("../config.json")
        main.bme280_names()
        if main.temperature_sensor_name == "":
            self.fail()

    def test_get_sequence(self):
        x = main.get_sequence()
        if x <= 0:
            self.fail()


'''
  def test_send_message(self):
       main.channel = grpcio.insecure_channel('store-and-forward:50051')
       main.stub = grpcStub(main.channel)

       msg = {}
       main.send_message(msg)
       if msg['sample_timestamp'] <= 0:
           self.fail()


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
       main.append_bme280_temp("", "test_thermometer", "test_temp")
       self.fail()

   def test_append_bme280_humidity(self):
       self.fail()

   def test_append_bme280_pressure(self):
       self.fail()
'''

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

import adafruit_platformdetect

detector = adafruit_platformdetect.Detector()

print('Chip id: ', detector.chip.id)

print('x Board id: ', detector.board.id)

print('Is this a DragonBoard 410c?', detector.board.DRAGONBOARD_410C)
print('Is this a Pi 3B+?', detector.board.RASPBERRY_PI_3B_PLUS)
print('Is this a Pi 4B?', detector.board.RASPBERRY_PI_4B)
print('Is this a 40-pin Raspberry Pi?', detector.board.any_raspberry_pi_40_pin)
print('Is this a Raspberry Pi Compute Module?', detector.board.any_raspberry_pi_cm)
print('Is this a BBB?', detector.board.BEAGLEBONE_BLACK)
print('Is this a Giant Board?', detector.board.GIANT_BOARD)
print('Is this a Coral Edge TPU?', detector.board.CORAL_EDGE_TPU_DEV)
print('Is this a SiFive Unleashed? ', detector.board.SIFIVE_UNLEASHED)
print('Is this an embedded Linux system?', detector.board.any_embedded_linux)
print('Is this a generic Linux PC?', detector.board.GENERIC_LINUX_PC)
print('Is this an OS environment variable special case?', detector.board.FTDI_FT232H       |
                                                          detector.board.MICROCHIP_MCP2221 )

if detector.board.any_raspberry_pi:
    print('Raspberry Pi detected.')

if detector.board.any_jetson_board:
    print('Jetson platform detected.')

if detector.board.any_orange_pi:
    print('Orange Pi detected.')

if detector.board.any_odroid_40_pin:
    print('Odroid detected.')

if detector.board.any_onion_omega_board:
    print('Onion Omega detected.')
#!/bin/bash
#
# Copyright (c) John Rodley 2022.
# SPDX-FileCopyrightText:  John Rodley 2022.
# SPDX-License-Identifier: MIT
#
# Permission is hereby granted, free of charge, to any person obtaining a copy of this
# software and associated documentation files (the "Software"), to deal in the
# Software without restriction, including without limitation the rights to use, copy,
# modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
# and to permit persons to whom the Software is furnished to do so, subject to the
# following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
# INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
# PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
# HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
# CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
# OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
#

set -m

GREEN='\033[0;32m'
echo -e "${GREEN}Systemd init system enabled."

# systemd causes a POLLHUP for console FD to occur
# on startup once all other process have stopped.
# We need this sleep to ensure this doesn't occur, else
# logging to the console will not work.
sleep infinity &
for var in $(compgen -e); do
	printf '%q=%q\n' "$var" "${!var}"
done > /etc/docker.env
exec /lib/systemd/systemd
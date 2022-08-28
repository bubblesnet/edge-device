# BubblesNet Edge Device

This repo is temporarily public. It will revert to private shortly and without notice as, imho, it is not quite ready yet.

[![codecov](https://codecov.io/gh/bubblesnet/edge-device/branch/develop/graph/badge.svg?token=4ETBIJSIKZ)](https://codecov.io/gh/bubblesnet/edge-device)
![ci](https://github.com/bubblesnet/edge-device/workflows/BubblesNetCI/badge.svg)

[![All Contributors](https://img.shields.io/badge/all_contributors-1-orange.svg?style=flat-square)](#contributors-)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

[![GitHub stars](https://img.shields.io/github/stars/bubblesnet/edge-device.svg?style=social&label=Star&maxAge=2592000)](https://GitHub.com/bubblesnet/edge-device/)
[![GitHub pull-requests](https://img.shields.io/github/issues-pr/bubblesnet/edge-device.svg)](https://GitHub.com/bubblesnet/edge-device/pull/)
[![Github all releases](https://img.shields.io/github/downloads/bubblesnet/edge-device/total.svg)](https://GitHub.com/bubblesnet/edge-device/releases/)

![Your Repository's Stats](https://github-readme-stats.vercel.app/api?username=bubblesnet&show_icons=true)

A collection of containers that communicate via gRPC to get edge device data to the cloud.

# Environment Variables
* ACTIVEMQ_HOST = the DNS name or IP address of the ActiveMQ instance (activemq container on the CONTROLLER). Used by containers sense-go, storeandforward
* ACTIVEMQ_PORT = the port the ActiveMQ server is running on (typically 61611). Used by containers storeandforward
* API_HOST = the DNS name or IP address of the API instance (api container on the CONTROLLER). Used by containers sense-go, storeandforward
* API_PORT = the port the API server is running on (typically 4001). Used by containers sense-go, storeandforward
* DEVICEID = the unique ID of this device in the system (e.g. 70000007) from device table in database. Used by containers sense-go, sense-python, storeandforward
* NODE_ENV = one of PRODUCTION, DEVELOPMENT, TEST, CI typically (PRODUCTION).  Used by container store-and-forward
* USERID = the unique ID of the user who owns this system (e.g. 90000009) from user table in database. Used by containers sense-go, storeandforward

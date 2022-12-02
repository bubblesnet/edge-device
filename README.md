# BubblesNet Edge Device


[![codecov](https://codecov.io/gh/bubblesnet/edge-device/branch/develop/graph/badge.svg?token=4ETBIJSIKZ)](https://codecov.io/gh/bubblesnet/edge-device)
![ci](https://github.com/bubblesnet/edge-device/workflows/BubblesNetCI/badge.svg)

[![All Contributors](https://img.shields.io/badge/all_contributors-1-orange.svg?style=flat-square)](#contributors-)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

[![GitHub stars](https://img.shields.io/github/stars/bubblesnet/edge-device.svg?style=social&label=Star&maxAge=2592000)](https://GitHub.com/bubblesnet/edge-device/)
[![GitHub pull-requests](https://img.shields.io/github/issues-pr/bubblesnet/edge-device.svg)](https://GitHub.com/bubblesnet/edge-device/pull/)
[![Github all releases](https://img.shields.io/github/downloads/bubblesnet/edge-device/total.svg)](https://GitHub.com/bubblesnet/edge-device/releases/)

![Your Repository's Stats](https://github-readme-stats.vercel.app/api?username=bubblesnet&show_icons=true)


This repo is one of 6 repos that make up the BubblesNet project. If you've arrived at this repo through
the side door (direct search), then you probably want to start with 
the [documentation repository](https://github.com/bubblesnet/documentation) for this 
project. You can not understand this repo without seeing how it interacts with the other repos.

A collection of containers that communicate via gRPC to get edge device data to the cloud.

| Container                              | Description                                                                                                                                                                                                                                                        |
|----------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| wifi-connect                           | This is a packaged block maintained by Balena that manages a typical IoT connect-to-wifi interaction                                                                                                                                                               |
| [sense-go](sense-go)                   | A Go language container that collects sensor data and forwards it to store-and-forward via GRPC and also controls the dispensers and AC devices and runs the automatic control code. The vast majority of the sensor and all the control functionality lives here. |
| [sense-python](sense-python)           | A Python language container that collects temp/pressure/humidity data and forwards it to storea-and-forward via GRPC.                                                                                                                                              |
| [store-and-forward](store-and-forward) | A custom block written in Go that uses a GRPC server to collect data from the sensor containers, uses BoltDB to store messages to be forwarded to the controller and forwards to the controller via the controller REST API.                                       |

# sense-go

sense-go is a Balena container written in Go that provides the bulk of the functionality of
the edge-devices, including sensing, control of AC devices (heater ...), and control of nutrient
dispensers.

sense-go supports the following sensors:

- pH
- CO2/VOC
- water level
- water temperature

as well as the dispensers and all the controllable devices.

The configuration of the device, what modules it supports and how it connects to them, is
contained in [the database](https://github.com/bubblesnet/documentation/Database.md)
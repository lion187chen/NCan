# NCan

Transmit CAN frame through NATS server. So we can debug CAN device remotely.

- Linux only.
- Simple config file.
- One point to one point.
- CAN driver as a plugin.
- Only one CAN driver was supportted.
- Waveshare USB-CAN-A was supported as a CAN driver.

CAN driver is designed as a plugin and it mus fit NCanDrvIf interface. So, you can realize a driver to support different CAN devices. CAN driver is used to drive a CAN device to send or receive CAN frames.

Todo:

- Complex config file.
- One point to many points.

# CL9 Core Universal Remote Firmware Uploader

A utility for uploading firmware to a CL9 Core device when the firmware is lost due to a low battery or a dead internal battery.

## How to build

go mod tidy

go build -o cl9core .

## How to run (example)

./cl9core --file ~/Downloads/core/PIC590.OS --port /dev/cu.usbserial-AH06LA82 -echo

## Hint

Only the TX output from the USB-to-UART converter needs to be connected to the CL9 Core.

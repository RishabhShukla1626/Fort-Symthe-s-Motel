#!/bin/bash

echo "Trying to Build and Run the application"

go build -o bookings cmd/web/*.go && ./bookings

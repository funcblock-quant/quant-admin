#!/bin/bash
killall go-admin # kill go-admin service
echo "stop quanta-admin success"
ps -aux | grep quanta-admin
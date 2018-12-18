#!/bin/sh
args=""
if [ -n "$TIMEOUT" ]; then
    args +=" -timeout $TIMEOUT"
fi
if [ -n "$INSTANCES" ]; then
    args +=" -instances $INSTANCES "
fi
if [ -n "$PROMETHEUS_INSTANCE" ]; then
    args +=" -prometheusInstance $PROMETHEUS_INSTANCE "
fi

./eureka-pushgateway -host $HOST_IP -eureka $EUREKA_URL -port $PORT -eport $EPORT $args
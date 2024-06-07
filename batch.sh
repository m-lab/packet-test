#!/bin/sh
set -ex

# TODO: pass-in machine hostname.
HOSTNAME=$(hostname -f)

server="pt-mlab1-lga0t.mlab-sandbox.measurement-lab.org"
while true; do
    # Packet test.
    /packet-test-schedule/train1 -server $server -params=client_hostname=$HOSTNAME

    # Full ndt7 test.
    /packet-test-schedule/ndt7 -server $server -params=client_hostname=$HOSTNAME

    # BBR-terminated test.
    /packet-test-schedule/ndt7 -server $server -params="bbr_exit=512&client_hostname=$HOSTNAME"

    # BBR/100MB-terminated test.
    /packet-test-schedule/ndt7 -server $server -params="bbr_exit=512&early_exit=100&client_hostname=$HOSTNAME"

    # BBR/200MB-terminated ndt7 test.
    /packet-test-schedule/ndt7 -server $server -params="bbr_exit=512&early_exit=200&client_hostname=$HOSTNAME"

    sleep 3h
done

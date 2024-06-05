#!/bin/sh
set -ex

# TODO: pass-in machine hostname.
HOSTNAME=$(hostname -f)

while true; do
    # Packet test.
    /packet-test-schedule/train1 -server pt-mlab1-lga0t.mlab-sandbox.measurement-lab.org -params=client_hostname=$HOSTNAME
    /packet-test-schedule/train1 -server pt-mlab1-lga1t.mlab-sandbox.measurement-lab.org -params=client_hostname=$HOSTNAME
    /packet-test-schedule/train1 -server pt-mlab2-lga0t.mlab-sandbox.measurement-lab.org -params=client_hostname=$HOSTNAME
    /packet-test-schedule/train1 -server pt-mlab2-lga1t.mlab-sandbox.measurement-lab.org -params=client_hostname=$HOSTNAME
    /packet-test-schedule/train1 -server pt-mlab3-lga0t.mlab-sandbox.measurement-lab.org -params=client_hostname=$HOSTNAME
    /packet-test-schedule/train1 -server pt-mlab3-lga1t.mlab-sandbox.measurement-lab.org -params=client_hostname=$HOSTNAME
    /packet-test-schedule/train1 -server pt-mlab4-lga0t.mlab-sandbox.measurement-lab.org -params=client_hostname=$HOSTNAME
    /packet-test-schedule/train1 -server pt-mlab4-lga1t.mlab-sandbox.measurement-lab.org -params=client_hostname=$HOSTNAME

    # Full ndt7 test.
    /packet-test-schedule/ndt7 -server pt-mlab1-lga0t.mlab-sandbox.measurement-lab.org -params=client_hostname=$HOSTNAME
    /packet-test-schedule/ndt7 -server pt-mlab1-lga1t.mlab-sandbox.measurement-lab.org -params=client_hostname=$HOSTNAME
    /packet-test-schedule/ndt7 -server pt-mlab2-lga0t.mlab-sandbox.measurement-lab.org -params=client_hostname=$HOSTNAME
    /packet-test-schedule/ndt7 -server pt-mlab2-lga1t.mlab-sandbox.measurement-lab.org -params=client_hostname=$HOSTNAME
    /packet-test-schedule/ndt7 -server pt-mlab3-lga0t.mlab-sandbox.measurement-lab.org -params=client_hostname=$HOSTNAME
    /packet-test-schedule/ndt7 -server pt-mlab3-lga1t.mlab-sandbox.measurement-lab.org -params=client_hostname=$HOSTNAME
    /packet-test-schedule/ndt7 -server pt-mlab4-lga0t.mlab-sandbox.measurement-lab.org -params=client_hostname=$HOSTNAME
    /packet-test-schedule/ndt7 -server pt-mlab4-lga1t.mlab-sandbox.measurement-lab.org -params=client_hostname=$HOSTNAME

    # BBR-terminated test.
    /packet-test-schedule/ndt7 -server pt-mlab1-lga0t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&client_hostname=$HOSTNAME"
    /packet-test-schedule/ndt7 -server pt-mlab1-lga1t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&client_hostname=$HOSTNAME"
    /packet-test-schedule/ndt7 -server pt-mlab2-lga0t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&client_hostname=$HOSTNAME"
    /packet-test-schedule/ndt7 -server pt-mlab2-lga1t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&client_hostname=$HOSTNAME"
    /packet-test-schedule/ndt7 -server pt-mlab3-lga0t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&client_hostname=$HOSTNAME"
    /packet-test-schedule/ndt7 -server pt-mlab3-lga1t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&client_hostname=$HOSTNAME"
    /packet-test-schedule/ndt7 -server pt-mlab4-lga0t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&client_hostname=$HOSTNAME"
    /packet-test-schedule/ndt7 -server pt-mlab4-lga1t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&client_hostname=$HOSTNAME"

    # BBR/100MB-terminated test.
    /packet-test-schedule/ndt7 -server pt-mlab1-lga0t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&early_exit=100&client_hostname=$HOSTNAME"
    /packet-test-schedule/ndt7 -server pt-mlab1-lga1t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&early_exit=100&client_hostname=$HOSTNAME"
    /packet-test-schedule/ndt7 -server pt-mlab2-lga0t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&early_exit=100&client_hostname=$HOSTNAME"
    /packet-test-schedule/ndt7 -server pt-mlab2-lga1t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&early_exit=100&client_hostname=$HOSTNAME"
    /packet-test-schedule/ndt7 -server pt-mlab3-lga0t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&early_exit=100&client_hostname=$HOSTNAME"
    /packet-test-schedule/ndt7 -server pt-mlab3-lga1t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&early_exit=100&client_hostname=$HOSTNAME"
    /packet-test-schedule/ndt7 -server pt-mlab4-lga0t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&early_exit=100&client_hostname=$HOSTNAME"
    /packet-test-schedule/ndt7 -server pt-mlab4-lga1t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&early_exit=100&client_hostname=$HOSTNAME"

    # BBR/200MB-terminated ndt7 test.
    /packet-test-schedule/ndt7 -server pt-mlab1-lga0t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&early_exit=200&client_hostname=$HOSTNAME"
    /packet-test-schedule/ndt7 -server pt-mlab1-lga1t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&early_exit=200&client_hostname=$HOSTNAME"
    /packet-test-schedule/ndt7 -server pt-mlab2-lga0t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&early_exit=200&client_hostname=$HOSTNAME"
    /packet-test-schedule/ndt7 -server pt-mlab2-lga1t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&early_exit=200&client_hostname=$HOSTNAME"
    /packet-test-schedule/ndt7 -server pt-mlab3-lga0t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&early_exit=200&client_hostname=$HOSTNAME"
    /packet-test-schedule/ndt7 -server pt-mlab3-lga1t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&early_exit=200&client_hostname=$HOSTNAME"
    /packet-test-schedule/ndt7 -server pt-mlab4-lga0t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&early_exit=200&client_hostname=$HOSTNAME"
    /packet-test-schedule/ndt7 -server pt-mlab4-lga1t.mlab-sandbox.measurement-lab.org -params="bbr_exit=512&early_exit=200&client_hostname=$HOSTNAME"

    sleep 15m
done

#!/bin/sh
set -ex

while true; do
    /packet-test-schedule/train1 -server pt-mlab1-lga0t.mlab-sandbox.measurement-lab.org
    /packet-test-schedule/train1 -server pt-mlab1-lga1t.mlab-sandbox.measurement-lab.org
    /packet-test-schedule/train1 -server pt-mlab2-lga0t.mlab-sandbox.measurement-lab.org
    /packet-test-schedule/train1 -server pt-mlab2-lga1t.mlab-sandbox.measurement-lab.org
    /packet-test-schedule/train1 -server pt-mlab3-lga0t.mlab-sandbox.measurement-lab.org
    /packet-test-schedule/train1 -server pt-mlab3-lga1t.mlab-sandbox.measurement-lab.org
    /packet-test-schedule/train1 -server pt-mlab4-lga0t.mlab-sandbox.measurement-lab.org
    /packet-test-schedule/train1 -server pt-mlab4-lga1t.mlab-sandbox.measurement-lab.org

    /packet-test-schedule/pair1 -server pt-mlab1-lga0t.mlab-sandbox.measurement-lab.org
    /packet-test-schedule/pair1 -server pt-mlab1-lga1t.mlab-sandbox.measurement-lab.org
    /packet-test-schedule/pair1 -server pt-mlab2-lga0t.mlab-sandbox.measurement-lab.org
    /packet-test-schedule/pair1 -server pt-mlab2-lga1t.mlab-sandbox.measurement-lab.org
    /packet-test-schedule/pair1 -server pt-mlab3-lga0t.mlab-sandbox.measurement-lab.org
    /packet-test-schedule/pair1 -server pt-mlab3-lga1t.mlab-sandbox.measurement-lab.org
    /packet-test-schedule/pair1 -server pt-mlab4-lga0t.mlab-sandbox.measurement-lab.org
    /packet-test-schedule/pair1 -server pt-mlab4-lga1t.mlab-sandbox.measurement-lab.org

    sleep 15m
done

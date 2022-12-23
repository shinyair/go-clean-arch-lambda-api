#!/bin/sh

STATE_FILE="configs/stage.yml"
rm $STATE_FILE
echo "stage: ${stage}" >> $STATE_FILE
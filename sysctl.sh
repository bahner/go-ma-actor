#!/bin/bash

sudo sysctl -w net.core.rmem_max=2500000
sudo sysctl -w net.core.wmem_max=2500000
sudo sysctl -w vm.overcommit_memory=1

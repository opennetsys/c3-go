#!/bin/bash

# clean up after running tests
docker ps -a | awk '{ print $1,$2 }' | grep hello-world | awk '{print $1 }' | xargs -I {} docker rm {}

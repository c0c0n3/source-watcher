#!/bin/bash

multipass launch --name osm --cpus 2 --mem 6G --disk 40G 18.04

multipass mount ./ osm:/mnt/osm-install

# multipass exec osm -- cd /mnt/osm-install && ./patched.install_osm.sh 2>&1 | tee install.log
#                                               ^ sudo issue

multipass shell osm
# cd /mnt/osm-install
# ./patched.install_osm.sh 2>&1 | tee install.log

# to clean up:
# $ multipass stop osm
# $ multipass delete osm
# $ multipass purge

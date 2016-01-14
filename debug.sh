#!/bin/sh
echo $(whoami) $(date -u +%Y-%m-%d-%H:%M:%S%z) $@ >> /var/log/objstor.log
echo $(/root/bin/objstor $@)
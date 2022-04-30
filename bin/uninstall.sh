#!/bin/sh

echo "Remove configure from supervisor ..."
rm -f /etc/supervisord.d/McmAPI.ini
supervisorctl update

echo "Kill all McmAPI ..."
killall -9 McmAPI

echo "Remove McmAPI ..."
rm -rf /home/McmAPI

echo "McmAPI Uninstall Success !"



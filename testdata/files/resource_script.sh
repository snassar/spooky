#!/bin/bash
# Test resource script for actions validation
echo "Memory limit: ${MEMORY_LIMIT:-unlimited} MB"
echo "CPU limit: ${CPU_LIMIT:-unlimited}%"
echo "Disk limit: ${DISK_LIMIT:-unlimited} MB"
exit 0

#!/bin/bash

# Automatically set ROOT_DIR to the name of the repository
ROOT_DIR=$(git rev-parse --show-toplevel 2>/dev/null)
# Fallback to the current directory name if not inside a git repository
if [ -z "$ROOT_DIR" ]; then
  ROOT_DIR=$PWD
fi
echo "Root directory is set to: $ROOT_DIR"

LSA_TARGET=$ROOT_DIR/output/large-scale-workshop
SERVICES_DIR=$ROOT_DIR/services
LOGS_DIR=$ROOT_DIR/output/logs
SLEEP_NUMBER=10

# Parameters:
#   $1 - the service configuration (yaml) under services dir
#   $2 - the desired node name
run_service () {
    $LSA_TARGET $SERVICES_DIR/$1 > $LOGS_DIR/$2.log 2>&1 &
    echo "$2 started with PID $!"
    sleep $SLEEP_NUMBER
}

main () {
    # ------------------------ RegService ------------------------
    echo "Starting registry services..."
    run_service "reg-service/service/RegServiceRoot.yaml" "RegService1_root"
    run_service "reg-service/service/AnotherRegService.yaml" "RegService2"
    run_service "reg-service/service/AnotherRegService.yaml" "RegService3"

    # ------------------------ CacheService ------------------------
    echo "Starting cache services..."
    run_service "cache-service/service/CacheServiceRoot.yaml" "CacheService1_root"
    run_service "cache-service/service/AnotherCacheService.yaml" "CacheService2"
    run_service "cache-service/service/AnotherCacheService.yaml" "CacheService3"

    # ------------------------ TestService ------------------------
    echo "Starting test services..."
    run_service "test-service/service/TestService.yaml" "TestService1"
    run_service "test-service/service/TestService.yaml" "TestService2"
    run_service "test-service/service/TestService.yaml" "TestService3"
}

main
echo "APP READY"

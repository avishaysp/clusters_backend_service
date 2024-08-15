#!/bin/bash
ROOT_DIR=/workspaces/RLAD
LSA_TARGET=$ROOT_DIR/output/large-scale-workshop
SERVICES_DIR=$ROOT_DIR/services
LOGS_DIR=$ROOT_DIR/output/logs
SLEEP_NUMBER=8

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

# 
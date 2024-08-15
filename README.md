# Large Scale Workshop

This project displays a basic distributed, service-oriented system including the following concepts: Remote Procedure Call, Service Discovery, Distributed hash table and Message Queue.

The system has three services:
- Registry Service - service discovery.
- Cache Service - in-memory cache using Chord Distributed Hash Table.
- Test Service - some basic functionalities.

Each new service node that starts in the system registers itself using the Registry service (e.g. Test and Cache services). The Registry service's root node (the first node to start) peforms an "IsAlive" check (every 10 seconds) on each node in the systems, and if it fails to answer within 3 retries - we unregister it from the system.

Both Registry and Cache services use the Chord structure to store their data. The communication between the services in the system is performed using gRPC, and while the Test service specifically supports an async Message Queue (ZeroMQ).


## Getting started

### Prerequisits
This project was developed using “Visual Studio Code Dev-Container”, so you'll need VSCode, Docker, and the Dev Container extension.

### Opening the project
First clone using:\
```git clone git@github.com:TAULargeScaleWorkshop/RLAD.git ./large-scale-workshop```

Then, use `file -> open workspace from file...` to open the workspace file, and then `reopen in container`.

The image mounts the large-scale-workshop directory into: `/workspaces/<cloned-repo-name>/`

## Usage

### Building
In the root directory run `./build.sh` to install required dependencies and build the app to `./output`.

### Running
Run the app using `./output/start.sh` to start 3 services of each type: Registry, Cache and Test. 

When all the services are ready the message: `"APP READY"` will be printed. After the system is up:
- Logs can be found in `./output/logs`
- All instances can be easily killed using `ps -ao pid= | xargs kill`.

### Testing
Each module has its UT, simply go to the test directory and run: `go test -v`. The directories include a client-side testing of each service among other modules:
- TestService: `/workspaces/RLAD/services/test-service/client/`
- CacheService: `/workspaces/RLAD/services/cache-service/client/`
- RegService: `/workspaces/RLAD/services/reg-service/client/`
- Chord DHT: `/workspaces/RLAD/services/reg-service/servant/dht`

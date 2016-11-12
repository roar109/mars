[![Build Status](https://drone.io/github.com/roar109/mars/status.png)](https://drone.io/github.com/roar109/mars/latest)

## Build ##

    go build -o mars

## Run ##

	mars -p project-alias

or just call the binary for a command prompt

	mars


See config.json for examples, the **java**, **jboss** and **workspace** sections accept system variables or path, first try to get system variable if not found uses it as a relative path.

# Roadmap #

- Start Java container
- Make the config file configurable with parameter

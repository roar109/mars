[![Build Status](https://drone.io/github.com/roar109/mars/status.png)](https://drone.io/github.com/roar109/mars/latest)

Small command line tool to compile and copy a generated artifact

## Build ##

    go build -o mars

## Run ##

	mars -p project-alias

or just call the binary for a command prompt

	mars


See config.json for examples, the **java**, **jboss** and **workspace** sections accept system variables or path, first try to get system variable if not found uses it as a relative path.

Available flows:

	1 - Maven clean and install 
	2 - Same as #1 plus copy the generated artifact to deployment folder
	3 - Same as #2 plus start the container (Not implemented yet)

# Roadmap #

- [] Start Java container
- [] use different maven path
- [x] Make the config file configurable with parameter
- [] Use "deployment" as a generic not fixed to jboss

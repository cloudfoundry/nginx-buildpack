# Cloud Foundry NGINX Buildpack

[![CF Slack](https://www.google.com/s2/favicons?domain=www.slack.com) Join us on Slack](https://cloudfoundry.slack.com/messages/buildpacks/)

A Cloud Foundry [buildpack](http://docs.cloudfoundry.org/buildpacks/) for apps requiring NGINX.


### Buildpack User Documentation

Official buildpack documentation can be found at [here](https://docs.cloudfoundry.org/buildpacks/nginx/index.html).

To use this buildpack, you will need to include an `nginx.conf` file in your app. [Here's an example.](https://github.com/cloudfoundry/nginx-buildpack/tree/master/fixtures/mainline)


### Building the Buildpack

To build this buildpack, run the following command from the buildpack's directory:

1. Source the .envrc file in the buildpack directory.

   ```bash
   source .envrc
   ```
   To simplify the process in the future, install [direnv](https://direnv.net/) which will automatically source .envrc when you change directories.

1. Install buildpack-packager

    ```bash
    ./scripts/install_tools.sh
    ```

1. Build the buildpack

    ```bash
    buildpack-packager build [ -cached=(true|false) ] -any-stack
    ```

1. Use in Cloud Foundry

   Upload the buildpack to your Cloud Foundry and optionally specify it by name

    ```bash
    cf create-buildpack [BUILDPACK_NAME] [BUILDPACK_ZIP_FILE_PATH] 1
    cf push my_app [-b BUILDPACK_NAME]
    ```

### Testing

Buildpacks use the [Cutlass](https://github.com/cloudfoundry/libbuildpack/tree/master/cutlass) framework for running integration tests.

To test this buildpack, run the following command from the buildpack's directory:

1. Source the .envrc file in the buildpack directory.

   ```bash
   source .envrc
   ```
   To simplify the process in the future, install [direnv](https://direnv.net/) which will automatically source .envrc when you change directories.

1. Run unit tests

    ```bash
    ./scripts/unit.sh
    ```

1. Run integration tests

    ```bash
    ./scripts/integration.sh
    ```

More information can be found on Github [cutlass](https://github.com/cloudfoundry/libbuildpack/tree/master/cutlass).

### Contributing

Find our guidelines [here](./CONTRIBUTING.md).

### Help and Support

Join the #buildpacks channel in our [Slack community](http://slack.cloudfoundry.org/)

### Reporting Issues

Open an issue on this project

### Active Development

The project backlog is on [Pivotal Tracker](https://www.pivotaltracker.com/projects/1042066)


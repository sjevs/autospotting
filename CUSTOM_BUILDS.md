# AutoSpotting Setup #

It's usually recommended to use the provided binaries, but in some cases you may
need to customize AutoSpotting for your own environment.

You'll need to set up a local environment able to compile Go code, compile the
binaries locally, upload them to an S3 bucket in your AWS account and update
your CloudFormation or Terraform stack to use those new binaries.

## Dependencies ##

1. Install [Go](https://golang.org/dl/), [git](https://git-scm.com/downloads),
   [Docker](https://www.docker.com/) and the [AWS command-line
   tool](https://aws.amazon.com/cli/). You may use the official binaries or your
   usual package manager, whatever you prefer is fine.

1. Verify that they were properly installed.

   `go version`, should be at least 1.7

   `git version`

   `docker version`

   `aws --version`

## Compiling the binaries locally ##

1. Set up a directory for your Go development. I'm using `godev` in my home
   directory for this example.

1. Set the `GOPATH` environment variable to point at your `godev` directory:

   `export GOPATH=$HOME/godev`

   Optionally add this line to your .bash_profile to persist across console
   sessions.

1. Run the following command to install the AutoSpotting project into your
   GOPATH directory:

   `go get github.com/cristim/autospotting`

   This downloads the source from GitHub, pulls in all necessary dependencies,
   builds it for local execution and deploys the binary into the golang binary
   directory which you may also want to append to your PATH.

1. Navigate to the root of the AutoSpotting repository:

   `cd $GOPATH/src/github.com/cristim/autospotting`

1. (Optional) You may want to make a minor change to the source code so you can
   tell when the tool is running your own custom-built version. If so, add a
   line like this to the `autospotting.go` file's `main()` function:

   `fmt.Println("Running <my organization name> binaries")`

1. (Optional) Try building and running the test suite locally to make sure
   everything works correctly:

   `make test`

   Below you can see more details on the other available make directives.

1. Build the code:

   `make build_local` or `go build` should both work, the only difference is the
   inclusion of versioning information into the binary, which may not matter
   much anyway.

## Running locally ##

1. Run the code, assuming you have AWS credentials defined in your environment
   or in the default AWS credentials profile:

   `./autospotting`

   You may also pass some command line flags, see the --help for more
   information on the available options.

   When you are happy with how your custom build behaves, you can generate a
   build for AWS Lambda.

## Using your own binaries in AWS Lambda ##

1. Set up an S3 bucket in your AWS account that will host your custom binaries.

1. The build system can use a `BUCKET_NAME` variable that tells it where to
   upload new binaries. Set it into your environment to the name of your S3
   bucket.

   `export BUCKET_NAME=my-bucket`

1. Define some AWS credentials or profile information into your
   [environment](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#cli-environment).

1. Build and upload your binaries to the S3 bucket.
   `make upload`

1. If you're already set up to use the tool with the official binaries, update
   your existing CloudFormation stack, and change the `LambdaS3Bucket` field to
   your S3 bucket name on the Stack Parameters section of the stack
   configuration.

   ![LambdaS3Bucket
   Configuration](https://mcristi.files.wordpress.com/2016/04/installationcloudformation2.png)

1. Save the CloudFormation configuration and let it create/update the resources.
   The tool should now be running against the binaries you built locally and
   uploaded to your own S3 bucket.

## Maintaining your own fork ##

It is recommended to contribute your changes into the mainline version of the
project whenever possible, so that others can benefit from your enhancements and
bug fixes, but for some reasons you may still want to run your own fork.

Unfortunately the golang import paths can make it tricky, but there is a nice
[article](http://code.openark.org/blog/development/forking-golang-repositories-on-github-and-managing-the-import-path)
which documents the problem in detail and gives a couple of possible
workarounds.

## Make directives ##

Use these directives defined in the `Makefile` to build, release, and test the
tool:

* **all (default, can be ommitted)**
  * Verifies that the necessary dependencies are installed.
  * Runs `go build` to compile the project for local development.

* **upload**
  * Prepares a special build designed to run in AWS Lambda.
  * Uploads the generated binaries from `build/s3` to the specified S3 bucket.

* **test**
  * Runs the test suite.

* **build_local**
  * Compiles the project for local execution.

[Back to the main Readme](./README.md)

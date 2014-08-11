#!/bin/sh

#
# For installing the Compose CLI tool
#

installComposeCli() {
  workingdir="/tmp/compose-cli"

  mkdir -p $workingdir

  unamestr=`uname -sm`

  if [[ "$unamestr" == 'Linux x86_64' ]]; then
    curl https://compose-cli.s3.amazonaws.com/builds/master/linux/amd64/compose -o $workingdir/compose
    curl https://compose-cli.s3.amazonaws.com/builds/master/linux/amd64/checksum -o $workingdir/checksum
  elif [[ "$unamestr" == 'Darwin x86_64' ]]; then
    curl https://compose-cli.s3.amazonaws.com/builds/master/darwin/amd64/compose -o $workingdir/compose
    curl https://compose-cli.s3.amazonaws.com/builds/master/darwin/amd64/checksum -o $workingdir/checksum
  else
    unamestr=`uname -s`
    if [[ "$unamestr" == 'Linux' ]]; then
      curl https://compose-cli.s3.amazonaws.com/builds/master/linux/386/compose -o $workingdir/compose
      curl https://compose-cli.s3.amazonaws.com/builds/master/linux/386/checksum -o $workingdir/checksum
    elif [[ "$unamestr" == 'Darwin' ]]; then
      curl https://compose-cli.s3.amazonaws.com/builds/master/darwin/386/compose -o $workingdir/compose
      curl https://compose-cli.s3.amazonaws.com/builds/master/darwin/386/checksum -o $workingdir/checksum
    else
      echo "We currently only build the CLI for Linux and MacOSX. To request builds for another platform, email support@compose.io."
      exit 1
    fi
  fi

  hostedChecksum=`cat $workingdir/checksum`
  localChecksum=`shasum $workingdir/compose | awk '{print $1}'`

  if [[ "$hostedChecksum" != "$localChecksum" ]]; then
    echo "Could not validate checksum of binary.  Please try again."
    exit 1
  fi

  chmod 555 $workingdir/compose

  if [[ -w "/usr/local/bin/compose" ]]; then
    mv $workingdir/compose /usr/local/bin/compose
  else
    echo "Please enter your sudo password to move the document to /usr/local/bin/compose:"
    sudo mv $workingdir/compose /usr/local/bin/compose
  fi

  echo ""
  echo "Install complete.  To get started, run:"
  echo ""
  echo "  compose deployments "
  echo "  compose --help "
  echo ""
  echo "This application is still in beta, and still actively changing.  Please test appropriately."
  echo ""
  echo "For documentation on the CLI, please see: http://docs.compose.io/getting-started/cli.html"
}


echo ""
echo ""
echo "We are installing the Compose CLI to /usr/local/bin/compose. The open sourced code for the CLI is available at https://github.com/MongoHQ/compose-cli."
echo ""
echo ""
installComposeCli

#!/bin/sh

#
# For installing the MongoHQ CLI tool
#

installMongoHQCli() {
  workingdir="/tmp/mongohq-cli"

  mkdir -p $workingdir

  unamestr=`uname -sm`

  if [[ "$unamestr" == 'Linux x86_64' ]]; then
    curl https://mongohq-cli.s3.amazonaws.com/builds/master/linux/amd64/mongohq -o $workingdir/mongohq
    curl https://mongohq-cli.s3.amazonaws.com/builds/master/linux/amd64/checksum -o $workingdir/checksum
  elif [[ "$unamestr" == 'Darwin x86_64' ]]; then
    curl https://mongohq-cli.s3.amazonaws.com/builds/master/darwin/amd64/mongohq -o $workingdir/mongohq
    curl https://mongohq-cli.s3.amazonaws.com/builds/master/darwin/amd64/checksum -o $workingdir/checksum
  else
    unamestr=`uname -s`
    if [[ "$unamestr" == 'Linux' ]]; then
      curl https://mongohq-cli.s3.amazonaws.com/builds/master/linux/386/mongohq -o $workingdir/mongohq
      curl https://mongohq-cli.s3.amazonaws.com/builds/master/linux/386/checksum -o $workingdir/checksum
    elif [[ "$unamestr" == 'Darwin' ]]; then
      curl https://mongohq-cli.s3.amazonaws.com/builds/master/darwin/386/mongohq -o $workingdir/mongohq
      curl https://mongohq-cli.s3.amazonaws.com/builds/master/darwin/386/checksum -o $workingdir/checksum
    else
      echo "We currently only build the CLI for Linux and MacOSX. Please check back later."
      exit 1
    fi
  fi

  hostedChecksum=`cat $workingdir/checksum`
  localChecksum=`shasum $workingdir/mongohq | awk '{print $1}'`

  if [[ "$hostedChecksum" != "$localChecksum" ]]; then
    echo "Could not validate checksum of binary.  Please try again."
    exit 1
  fi

  chmod 555 $workingdir/mongohq

  if [[ -w "/usr/local/bin/mongohq" ]]; then
    mv $workingdir/mongohq /usr/local/bin/mongohq
  else
    echo "Please enter your sudo password to move the document to /usr/local/bin/mongohq:"
    sudo mv $workingdir/mongohq /usr/local/bin/mongohq
  fi

  echo ""
  echo "Install complete.  To get started, run:"
  echo ""
  echo "  mongohq databases "
  echo "  mongohq --help "
  echo ""
  echo "This application is still in beta, and still actively changing.  Please test appropriately."
  echo "For documentation on the CLI, please see: http://docs.mongohq.com/getting-started/cli.html"
}


echo ""
echo ""
echo "We are installing the MongoHQ CLI to /usr/local/bin/mongohq. The open sourced code for the CLI is available at https://github.com/MongoHQ/mongohq-cli."
echo ""
echo ""
installMongoHQCli

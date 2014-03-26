#!/bin/sh

#
# For installing the MongoHQ CLI tool
#

installMongoHQCli() {
  workingdir="/tmp/mongohq-cli"

  mkdir -p $workingdir

  unamestr=`uname`
  if [[ "$unamestr" == 'Linux' ]]; then
    curl https://mongohq-cli.s3.amazonaws.com/linux/mongohq -o $workingdir/mongohq
    curl https://mongohq-cli.s3.amazonaws.com/linux/checksum -o $workingdir/checksum
  elif [[ "$unamestr" == 'FreeBSD' ]]; then
    curl https://mongohq-cli.s3.amazonaws.com/freebsd/mongohq -o $workingdir/mongohq
    curl https://mongohq-cli.s3.amazonaws.com/freebsd/checksum -o $workingdir/checksum
  elif [[ "$unamestr" == 'Darwin' ]]; then
    curl https://mongohq-cli.s3.amazonaws.com/macosx/mongohq -o $workingdir/mongohq
    curl https://mongohq-cli.s3.amazonaws.com/macosx/checksum -o $workingdir/checksum
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
    sudo mv $workingdir/mongohq /usr/local/bin/mongohq
  fi

  echo ""
  echo "Install complete.  To get started, run:"
  echo ""
  echo "  mongohq databases "
  echo "  mongohq --help "
  echo ""
  echo "This application is still in beta, and still actively changing.  Please test appropriately."
}


echo ""
echo ""
echo "We are installing the MongoHQ CLI to /usr/local/bin/mongohq. The open sourced code for the CLI is available at https://github.com/MongoHQ/mongohq-cli."
echo ""
echo ""
installMongoHQCli

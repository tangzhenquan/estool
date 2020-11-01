#!/bin/bash

INSTALL_DIR=`pwd`"/build/bin"


function build()
{
	echo  "Building $1 ..."

	if [ -f "$1/Makefile" ]; then
		build_cwd=`pwd`
    cd $1
    make install DESTDIR=$2
    if [ "$?" != "0" ]; then
		  echo "$1 failed"
		  exit 1
    else
      echo "$1 done"
    fi
    cd $build_cwd
	else
		echo "not found Makefile "
		exit 1
	fi
}




build_items="cmd/esImport \
cmd/tool \
cmd/esQuery"

if [ ! -z "$1" ]; then
  if 	[ "$1" != "all" ]; then
	  build_items=$1
	fi
fi


if [ ! -z "$2" ]; then
	INSTALL_DIR=$1
fi

if [ ! -d "$INSTALL_DIR" ]; then
  mkdir -p $INSTALL_DIR
fi


for item in $build_items
do
	build $item $INSTALL_DIR
done

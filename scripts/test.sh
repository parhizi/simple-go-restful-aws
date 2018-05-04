#!/usr/bin/env bash

echo "Testing .go files"

cd src/handlers/

for folder in */;
  do
  (cd $folder

    for innerFile in *;
    do
      if [ $innerFile == *".out" ] ; then
        rm $innerFile
      fi

      if [ $innerFile == *".html" ] ; then
        rm $innerFile
      fi
    
    done
      
    go test -coverprofile=cover.out
    go tool cover -html=cover.out -o cover.html

    )
  done

echo "Done."
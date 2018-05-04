#!/usr/bin/env bash
dep init
dep ensure

echo "Compiling functions to bin/handlers/ ..."

rm -rf bin/

cd src/handlers/

for folder in */;
  do
  (cd $folder
    for f in *.go;
    do  
      if [ $f == *"_test.go" ] ; then
  echo "— "$f "Skipped"
        continue;
      fi

      filename="${f%.go}"
    
      if GOOS=linux go build -o "../../../bin/handlers/$filename" ${f}; then
        echo "✓ Compiled $filename"
      else
        echo "✕ Failed to compile $filename!"
        exit 1
      fi
    
    done)
  done

echo "Done."
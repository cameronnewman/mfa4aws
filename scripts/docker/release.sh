
#!/usr/bin/env bash

BIN_FOLDER="bin"
RELEASE_FOLDER="release"
FILES=$BIN_FOLDER/[^.]*

if [[ ! -d $BIN_FOLDER ]]; then
  echo "** bin folder not present, nothing to release"
  exit 1
fi

if [[ -d $RELEASE_FOLDER ]]; then
  rm -rf $RELEASE_FOLDER
fi

mkdir -p $RELEASE_FOLDER

for f in $FILES
do
  echo "Processing $f file..."
  FILE=`basename $f`
  tar -czvf $RELEASE_FOLDER/$FILE.tar.gz --transform 's#^bin/##' $f
  if [ $? -ne 0 ]; then
      echo 'An error has occurred! Aborting the script execution...'
      exit 1
  fi
done

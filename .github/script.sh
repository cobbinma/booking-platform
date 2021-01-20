#!/bin/sh
echo starting...
set -e
go get -u github.com/charypar/monobuild

# get dependencies that need to be built
branch=$(git rev-parse --abbrev-ref HEAD)
if [ "$branch" = "master" ]
then
  dependencies=$(monobuild diff --main-branch)
else
  dependencies=$(monobuild diff --base-branch remotes/origin/master)
fi

# get docker tag from library
get_docker_tag() {
  DOCKERHUB_OWNER=${DOCKERHUB_OWNER}
  IMAGE_PREFIX="booking"
  IMAGE_NAME=$(echo "$1" | sed 's:.*/::')
  IMAGES_TAG="$branch"
  echo "$DOCKERHUB_OWNER"/"$IMAGE_PREFIX"-"$IMAGE_NAME":"$IMAGES_TAG"
}

# inject makefile if not in lib
for dep in $dependencies
do
  dep=$(echo "$dep" | sed 's/\://g')
  FILE="$dep"/Makefile
if [ ! -f "$FILE" ]; then
  echo "copying default makefile to $FILE"
  cp lib/Makefile "$dep"/Makefile
fi
done

# test
for dep in $dependencies
do
  echo testing "$dep"
  dep=$(echo "$dep" | sed 's/\://g')
  make -C "$dep" test
done

# build docker images
for dep in $dependencies
do
  echo building "$dep"
  dep=$(echo "$dep" | sed 's/\://g')
  make -C "$dep" TAG="$(get_docker_tag "$dep")" build
done

# push docker images
if [ "$branch" = "master" ]
then
for dep in $dependencies
do
    echo deploying "$dep"
    dep=$(echo "$dep" | sed 's/\://g')
    make -C "$dep" TAG="$(get_docker_tag "$dep")" deploy
done
fi
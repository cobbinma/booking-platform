#!/bin/sh
echo starting...
set -e
go get -u github.com/charypar/monobuild
branch=$(git rev-parse --abbrev-ref HEAD)
if [ "$branch" = "master" ]
then
changed_libs=$(monobuild diff --main-branch)
else
changed_libs=$(monobuild diff)
fi
get_tag() {
  DOCKERHUB_OWNER="cobbinma"
  IMAGE_PREFIX="booking"
  IMAGE_NAME=$(echo "$1" | sed 's:.*/::')
  IMAGES_TAG="$branch"
  echo "$DOCKERHUB_OWNER"/"$IMAGE_PREFIX"-"$IMAGE_NAME":"$IMAGES_TAG"
}

# test
for lib in $changed_libs
do
  echo testing "$lib"
  # shellcheck disable=SC2039
  lib="${lib//:}"
  make -C "$lib" test
done

# build docker images
for lib in $changed_libs
do
  echo building "$lib"
  # shellcheck disable=SC2039
  lib="${lib//:}"
  docker build "$lib" -t "$(get_tag "$lib")"
done

# push docker images
if [ "$branch" = "master" ]
then
for lib in $changed_libs
do
    echo building "$lib"
    # shellcheck disable=SC2039
    docker push "$(get_tag "${lib//:}")"
done
fi
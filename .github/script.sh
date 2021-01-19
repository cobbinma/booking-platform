#!/bin/sh
echo starting...
set -e
go get -u github.com/charypar/monobuild
branch=$(git rev-parse --abbrev-ref HEAD)
if [ "$branch" = "master" ]
then
dependencies=$(monobuild diff --main-branch)
else
dependencies=$(monobuild diff --base-branch remotes/origin/master)
fi
get_tag() {
  DOCKERHUB_OWNER=${DOCKERHUB_OWNER}
  IMAGE_PREFIX="booking"
  IMAGE_NAME=$(echo "$1" | sed 's:.*/::')
  IMAGES_TAG="$branch"
  echo "$DOCKERHUB_OWNER"/"$IMAGE_PREFIX"-"$IMAGE_NAME":"$IMAGES_TAG"
}

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
  docker build "$dep" -t "$(get_tag "$dep")"
done

# push docker images
if [ "$branch" = "master" ]
then
for dep in $dependencies
do
    echo deploying "$dep"
    dep=$(echo "$dep" | sed 's/\://g')
    docker push "$(get_tag "$dep")"
done
fi
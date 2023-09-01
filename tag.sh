#!/bin/bash

cd "$(dirname "$(readlink -f "$0")")/src" || exit 1

git fetch --tags

# 获取当前 commit 的哈希值
commit_hash=$(git rev-parse HEAD)

tag=$(git tag "--contains=${commit_hash}" 2>/dev/null)

if [ -n "$tag" ]; then
    echo "commit tag: $tag"
	exit
fi

stats=$(git status -s -u)
if [ -n "$stats" ]; then
	echo "dirty workspace"
	echo
	echo "$stats"
	echo
	exit
fi

(cd ../test && make || exit 1)

# 获取最新的 tag
latest_tag=$(git describe --tags --abbrev=0)

# 解析最新 tag 的版本号部分（假设标签格式为 v<major>.<minor>.<patch>）
tag_prefix="v"
version="${latest_tag#"$tag_prefix"}"

if [[ ! "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "tag syntax error: $version"
    exit 1
fi

# 解析主、次、修订版本号
IFS='.' read -ra version_parts <<< "$version"
major=${version_parts[0]}
minor=${version_parts[1]}
patch=${version_parts[2]}

# 增加修订版本号
patch=$((patch + 1))

# 构建新的 tag
new_tag="${tag_prefix}${major}.${minor}.${patch}"

# 打上新的 tag
git push || exit 1
git tag -s "$new_tag" -m "$new_tag"
git push origin "$new_tag"

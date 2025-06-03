#!/bin/bash

set -e

find ./proto/entpb -name '*.proto' -exec \
  sed -i '' 's|option go_package = "example/ent/proto/entpb";|option go_package = "example/gen/entpb";|' {} +

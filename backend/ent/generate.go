package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate ./schema
//go:generate go run -mod=mod entgo.io/contrib/entproto/cmd/entproto -path ./schema
//go:generate bash ./fix_protoc_go_package.sh
//go:generate rm -rf ../../proto/entpb
//go:generate mv ./proto/entpb ../../proto

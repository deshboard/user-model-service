PROTO_PATH = vendor/github.com/deshboard/apis/deshboard

.PHONY: proto

proto: ## Generate code from protocol buffer
	@mkdir -p apis
	protoc -I ${PROTO_PATH} ${PROTO_PATH}/iam/user/v1alpha1/directory.proto ${PROTO_PATH}/iam/user/v1alpha1/repository.proto --go_out=plugins=grpc:apis

envcheck::
	$(call executable_check,protoc,protoc)

PROTO_DIR := proto
GOOGLE_APIS := third_party/googleapis
OUT := pb/user

gen:
	mkdir -p $(OUT)
	protoc \
	  --proto_path=$(PROTO_DIR) \
	  --proto_path=$(GOOGLE_APIS) \
	  --go_out=$(OUT) --go_opt=paths=source_relative \
	  --go-grpc_out=$(OUT) --go-grpc_opt=paths=source_relative \
	  --grpc-gateway_out=$(OUT) --grpc-gateway_opt=paths=source_relative \
	  $(PROTO_DIR)/user.proto


all:
	protoc --java_out=./ types.proto
	protoc --go_out=./ types.proto
	protoc --python_out=./ types.proto

clean:
	rm -rf protoc/*

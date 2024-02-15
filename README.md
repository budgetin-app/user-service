
# Budgetin - User Service

## Overview
User Service is one of the microservices in Budgetin that responsible for user authentication and user authorization by managing the user roles and permissions.

This repository contains an application written in Go that utilizes Protocol Buffers and gRPC for communication. This README provides step-by-step instructions on how to install, compile, and run the application.

## Installation and Setup
Before running the application, ensure you have Go installed on your machine. You can download and install it from the [Go Official Website](https://go.dev/dl/).

Once Go is installed, follow these steps:

1. **Clone the repository**: 
	```bash
	git clone https://github.com/budgetin-app/user-service.git
    ```

2. **Change directory** into the project folder:
    ```bash
    cd user-service
    ```

3. **Install Protocol Buffers Compiler (protoc)**:\
You need to have Protocol Buffers Compiler installed. You can download it from the [Protocol Buffers GitHub releases page](https://github.com/protocolbuffers/protobuf/releases).

4. **Install protoc-gen-go**:\
To generate Go files from .proto files, you need to install protoc-gen-go. You can do this using the following command:
    ```bash
    go install google.golang.org/protobuf/cmd/protoc-gen-go
    ```

## Compilation and Running

1. **Compile the proto file**\
The first step is to compile the `.proto` file using protoc. Navigate to the directory where the `.proto` file is located and execute the following command:
	```bash
	protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ./app/proto/*.proto
	```

2. **Build the generated wire injection file**\
After compiling the `.proto` file, you need to generate the wire injection file. Run the following command:
	```bash
	wire ./config
	```

3. **Build and run the application through `main.go` file**\
Finally, you can build and run the application using the `main.go` file. Execute the following commands:
	```bash
	go build -o main.exe ; ./main
	```

## Contributing

Contributions are welcome! If you find any bugs or have suggestions for improvements, please feel free to open an issue or submit a pull request.

## License
This project is licensed under the [MIT License](LICENSE).

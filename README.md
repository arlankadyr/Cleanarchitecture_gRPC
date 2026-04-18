# AP2 Assignment 2

Medical Scheduling Platform with gRPC.

## Project Structure

```
ap2-assignment2/
├── doctor-service/
│   ├── cmd/doctor-service/main.go          ← entry point
│   ├── internal/
│   │   ├── model/doctor.go
│   │   ├── repository/memory.go
│   │   ├── usecase/doctor.go
│   │   └── transport/grpc/handler.go
│   └── proto/
│       ├── doctor.proto
│       ├── doctor.pb.go
│       └── doctor_grpc.pb.go
├── appointment-service/
│   ├── cmd/appointment-service/main.go
│   ├── internal/
│   │   ├── model/appointment.go
│   │   ├── repository/memory.go
│   │   ├── usecase/appointment.go
│   │   ├── client/doctor_client.go
│   │   └── transport/grpc/handler.go
│   └── proto/
│       ├── appointment.proto
│       ├── appointment.pb.go
│       └── appointment_grpc.pb.go
├── go.mod
├── go.sum
└── README.md




## Prerequisites

- Go 1.21 or later
- (Optional, to regenerate stubs) `protoc` + `protoc-gen-go` + `protoc-gen-go-grpc`

---

## Running Both Services Locally

> **Start Order:** Doctor Service first, then Appointment Service.

**Terminal 1 – Doctor Service (port 50051)**

```bash
go run ./doctor-service/cmd/doctor-service
```

**Terminal 2 – Appointment Service (port 50052)**

```bash
go run ./appointment-service/cmd/appointment-service
```

The Appointment Service reads `DOCTOR_SERVICE_ADDR` from the environment
(default: `localhost:50051`). Override it if needed:

```bash
DOCTOR_SERVICE_ADDR=localhost:50051 go run ./appointment-service/cmd/appointment-service
```

---

## How to Regenerate Proto Stubs

### 1. Install protoc

Download from https://github.com/protocolbuffers/protobuf/releases and add to `PATH`.

### 2. Install Go plugins

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 3. Regenerate

```bash
# Doctor Service
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       doctor-service/proto/doctor.proto

# Appointment Service
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       appointment-service/proto/appointment.proto
```

---

## Services

Doctor Service runs on port 50051.
Appointment Service runs on port 50052.

Both have in-memory storage.

---

## How it works

The Appointment Service calls the Doctor Service via gRPC to check if the doctor exists before creating an appointment.

## Testing with grpcurl

You can run `grpcurl_commands.sh` to test the endpoints if you have `grpcurl` installed.
Or use the `Postman_gRPC_Collection.json` to test via Postman.

#include <chrono>
#include <iostream>
#include <memory>
#include <random>
#include <string>
#include <thread>

#include <grpc/grpc.h>
#include <grpcpp/channel.h>
#include <grpcpp/client_context.h>
#include <grpcpp/create_channel.h>
#include <grpcpp/security/credentials.h>

#include "proto/ultrasound.grpc.pb.h"

int main (int argc, char** argv)
{
    std::string host = "localhost:8080";
    auto channel = grpc::CreateChannel(host, grpc::InsecureChannelCredentials());

    auto stub = ultrasound::UltraScanner::NewStub(channel);

    grpc::ClientContext context;
    ultrasound::CreateGraphRequest request;
    ultrasound::CreateGraphResponse response;

    grpc::Status status = stub->CreateGraph(&context, request, &response);
    if (status.ok())
    {
        std::cout << "ok\n";
        std::cout << response.message() << "\n";
    }
    else
    {
        std::cout << "not good\n";
        std::cout << status.error_message() << "\n";
    }

    return 0;
}

# Development guidelines

We created this project with go-blueprint, we have a Go API (entry point
`cmd/api/main.go`) with gin router and a HTMX front-end (in `cmd/web`).

## Use clean architecture

Organize the code in a clean architecture with the concepts controllers,
commands, use cases, services, repositories, clients and models for
entities/models/requests/responses for application entities and data transfer
object.

A controller builds a command from a request and calls the use cases with the
command.

The use case uses services to implement the business logic.

The services use repositories and clients to create/modify/query data.

## Folder structure

Keep all the code in the `internal` with the following clean code structure,
with some made up examples for user manangement and external Ollama API access:

```
internal
    adapter
        in
            web
                UserUpdateController
                UserUpdateRequest
                UserUpdateReponse
        out
            persistance
                UserRepositoryImpl
                UserRepositoryModel
            clients
                OllamaApiClientImpl
                OllamaMessageModel
    application
        domain
            entities
                User
            use_cases
                UserUpdateCommandImpl
            services
                UserService
        port
            in
                UserUpdateUseCase (interface)
                UserUpdateCommand
            out
                persistance
                    UserRepository (interface)
                clients
                    OllamaApiClient (interface)
    common
    database
        Database (driver)
    server
        Routes (call controllers)
        Server
```

## Dependency injection.

Use dependency injection with the go fx library.

## Tests

When you implement a new feature always add tests.

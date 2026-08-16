# Crossbowhttp 
Crossbowhttp is a lightweight adapter for exposing [crossbow](https://github.com/JuanX-G/crossbow) servers over http.
It includes basic handling of JSON and supports middleware chaining (see [chain](./chain.go)).

## Usage
Adapters just have to implement the `HttpMarshaller` interface. For examples see [json](./json/json.go).
The marshaler allows decoding incoming requests and encoding outgoing responses.
# Gograph API drafts

## Items
- Property:
  - It is a pair of `key string` and `value {}interface`
  - Allowed value types
    - Numbers: Both integer values, with capacity as Java’s Long type, and floating points, with capacity as Java’s Double.
    - Booleans
    - Strings
    - Arrays: Of the above basic types
      - all values in the array must be of the same type
      - type is inferred from the values in the array
      - when storing values, the server has to know the type of the array and that means you can send an empty array only if an array is already stored for the given property; if no array exists already, the server will reject the request

- Node

- Relationship
  - has a direction, with a default to ... ? (one of both, outgoing, incoming)

- Index

- Batch

- GraphDatabaseService
  - Query(cypherQuery string)

***

var service GraphDatabaseService = new(GraphDatabaseService)
err := service.Connect("http://localhost:7474/db/data", "username", "password")
// Is it possible to create and Save empty node, like new(Node)?
var node Node = new(Node).SetType("person").SetProperty("firstname", "Mateusz").SetProperty("age", "99")
service.Save(node) // or explicitly `Create`?
service.Update(node) // or only `Save`?
var nodeId int64 = node.getId()
startNode, err := service.GetNodeById(nodeId)
var props map[string]string = {"age": "100"}
startNode.SetProperties(props) // sets `age` dirty

endNode, err := service.GetNodeById(2)

var rel Relationship = newRelationship(startNode, endNode).SetType("likes")
service.CreateIfAbsent(rel) // or `service.GetOrCreate` as both could return the node anyway

batch := service.Batch()

batch.Create(...)
batch.GetNodeById(...)

// perform operations as if on the `service` directly, though all APIs may not be available

batch.Execute()

***

- Create POST
- Read GET
- Replace PUT
- Update PATCH
- Delete DELETE

Philosophy
- close to the REST API
- can learn it by reading the official Neo4J docs, it's easy, idiomatic to translate API calls into Go calls
- on top of that there may exist some helper functions, to ease the execution of common operations


- both Node and Relationship can implement a set of common interfaces
  - for coping with attached properties (e.g. names `PropertyHolder`)
  - they both have an ID, and URI (or is that URL?) `Identity`
  - they both can have a type `Typer`


node, err := service.Get(Node{id=303})
node, err := service.Post(Node{}) // Create
node, err := service.Post(Node{"properties"=map[string]string{"name":"mako"}})

val, err := service.Get(node.Property('name'))
props, err := service.Get(node.Properties())

err := service.Put(node.Property('name', 'joe')) // Replace
err := service.Put(node.Properties(map[string]string{"name":"joe", "age":33))
err := service.Delete(node.Property('age'))
err := service.Delete(node.Properties()) // deletes all
err := service.Delete(node.Properties(map[string]string{"age":33)) // deletes selected

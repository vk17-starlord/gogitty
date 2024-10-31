package core

// GitObject defines the common interface for all Git objects
type GitObject interface {
	Serialize(repo string) any   // Serialize the object's data
	Deserialize(data []byte) any // Deserialize the data into the object
	Init()                       // Initialize the object (optional)
}

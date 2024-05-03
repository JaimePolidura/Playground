package objects

type ArrayObject struct {
	Object      Object
	ContentType ObjectType
	Content     []byte
}

package outputs

// Output interface defines log output destinations
type Output interface {
	Write(data []byte) error
	Close() error
	GetName() string
}

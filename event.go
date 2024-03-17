package bus

// 40 byte
type event struct {
	subKey subKey
	msg    []byte
}

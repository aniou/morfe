

package memory

type Memory interface {
        Read(address uint32) byte
        Write(address uint32, value byte)

        Shutdown()
        Size() uint32
        Clear()
	Dump(address uint32) []byte
}


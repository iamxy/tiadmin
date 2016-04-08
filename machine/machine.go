package machine

const (
	shortIDLen = 8
)

type Mach interface {
	ID() string
	ShortID() string
	MatchID(ID string) bool
}

type Machine struct {
	MachID         string
	HostName       string
	HostRegion     string
	HostDatacenter string
	PublicIP       string
	IsAlive        bool
	Stat           MachineStat
}

func NewMachine() Mach {
	return &Machine{}
}

func (m *Machine) ID() string {
	return m.MachID
}

func (m *Machine) ShortID() string {
	if len(m.MachID) <= shortIDLen {
		return m.MachID
	}
	return m.MachID[0:shortIDLen]
}

func (m *Machine) MatchID(ID string) bool {
	return m.MachID == ID || m.ShortID() == ID
}

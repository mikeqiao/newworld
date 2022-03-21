package manager

var ModManager *MManager

type MManager struct {
	BaseMod *BaseMod
}

func (m *MManager) Init() {
	m.BaseMod = new(BaseMod)
}

func (m *MManager) Run() {
}

func (m *MManager) Close() {
}

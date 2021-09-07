package marker

type Marker struct {
	Base    int
	Current int
}

func New(base int) *Marker {
	return &Marker{Base: base, Current: base}
}

func (m *Marker) Next() int {
	m.Current += 1
	return m.Current
}

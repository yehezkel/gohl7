package gohl7


type Message struct{
	raw []byte
	segments []*Segment
}

func (m *Message) Segments() ([]*Segment, error){
	return m.segments, nil
}
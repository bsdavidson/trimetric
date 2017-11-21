package trimet

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"time"

	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *StopTimeEvent) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zxvk uint32
	zxvk, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zxvk > 0 {
		zxvk--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "delay":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Delay = nil
			} else {
				if z.Delay == nil {
					z.Delay = new(int32)
				}
				*z.Delay, err = dc.ReadInt32()
				if err != nil {
					return
				}
			}
		case "time":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Time = nil
			} else {
				if z.Time == nil {
					z.Time = new(time.Time)
				}
				*z.Time, err = dc.ReadTime()
				if err != nil {
					return
				}
			}
		case "uncertainty":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Uncertainty = nil
			} else {
				if z.Uncertainty == nil {
					z.Uncertainty = new(int32)
				}
				*z.Uncertainty, err = dc.ReadInt32()
				if err != nil {
					return
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *StopTimeEvent) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "delay"
	err = en.Append(0x83, 0xa5, 0x64, 0x65, 0x6c, 0x61, 0x79)
	if err != nil {
		return err
	}
	if z.Delay == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteInt32(*z.Delay)
		if err != nil {
			return
		}
	}
	// write "time"
	err = en.Append(0xa4, 0x74, 0x69, 0x6d, 0x65)
	if err != nil {
		return err
	}
	if z.Time == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteTime(*z.Time)
		if err != nil {
			return
		}
	}
	// write "uncertainty"
	err = en.Append(0xab, 0x75, 0x6e, 0x63, 0x65, 0x72, 0x74, 0x61, 0x69, 0x6e, 0x74, 0x79)
	if err != nil {
		return err
	}
	if z.Uncertainty == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteInt32(*z.Uncertainty)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *StopTimeEvent) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "delay"
	o = append(o, 0x83, 0xa5, 0x64, 0x65, 0x6c, 0x61, 0x79)
	if z.Delay == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendInt32(o, *z.Delay)
	}
	// string "time"
	o = append(o, 0xa4, 0x74, 0x69, 0x6d, 0x65)
	if z.Time == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendTime(o, *z.Time)
	}
	// string "uncertainty"
	o = append(o, 0xab, 0x75, 0x6e, 0x63, 0x65, 0x72, 0x74, 0x61, 0x69, 0x6e, 0x74, 0x79)
	if z.Uncertainty == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendInt32(o, *z.Uncertainty)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *StopTimeEvent) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zbzg uint32
	zbzg, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zbzg > 0 {
		zbzg--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "delay":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Delay = nil
			} else {
				if z.Delay == nil {
					z.Delay = new(int32)
				}
				*z.Delay, bts, err = msgp.ReadInt32Bytes(bts)
				if err != nil {
					return
				}
			}
		case "time":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Time = nil
			} else {
				if z.Time == nil {
					z.Time = new(time.Time)
				}
				*z.Time, bts, err = msgp.ReadTimeBytes(bts)
				if err != nil {
					return
				}
			}
		case "uncertainty":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Uncertainty = nil
			} else {
				if z.Uncertainty == nil {
					z.Uncertainty = new(int32)
				}
				*z.Uncertainty, bts, err = msgp.ReadInt32Bytes(bts)
				if err != nil {
					return
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *StopTimeEvent) Msgsize() (s int) {
	s = 1 + 6
	if z.Delay == nil {
		s += msgp.NilSize
	} else {
		s += msgp.Int32Size
	}
	s += 5
	if z.Time == nil {
		s += msgp.NilSize
	} else {
		s += msgp.TimeSize
	}
	s += 12
	if z.Uncertainty == nil {
		s += msgp.NilSize
	} else {
		s += msgp.Int32Size
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *StopTimeUpdate) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zbai uint32
	zbai, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zbai > 0 {
		zbai--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "stop_sequence":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.StopSequence = nil
			} else {
				if z.StopSequence == nil {
					z.StopSequence = new(uint32)
				}
				*z.StopSequence, err = dc.ReadUint32()
				if err != nil {
					return
				}
			}
		case "stop_id":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.StopID = nil
			} else {
				if z.StopID == nil {
					z.StopID = new(string)
				}
				*z.StopID, err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "arrival":
			err = z.Arrival.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "departure":
			err = z.Departure.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "schedule_relationship":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.ScheduleRelationship = nil
			} else {
				if z.ScheduleRelationship == nil {
					z.ScheduleRelationship = new(int32)
				}
				*z.ScheduleRelationship, err = dc.ReadInt32()
				if err != nil {
					return
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *StopTimeUpdate) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 5
	// write "stop_sequence"
	err = en.Append(0x85, 0xad, 0x73, 0x74, 0x6f, 0x70, 0x5f, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65)
	if err != nil {
		return err
	}
	if z.StopSequence == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteUint32(*z.StopSequence)
		if err != nil {
			return
		}
	}
	// write "stop_id"
	err = en.Append(0xa7, 0x73, 0x74, 0x6f, 0x70, 0x5f, 0x69, 0x64)
	if err != nil {
		return err
	}
	if z.StopID == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteString(*z.StopID)
		if err != nil {
			return
		}
	}
	// write "arrival"
	err = en.Append(0xa7, 0x61, 0x72, 0x72, 0x69, 0x76, 0x61, 0x6c)
	if err != nil {
		return err
	}
	err = z.Arrival.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "departure"
	err = en.Append(0xa9, 0x64, 0x65, 0x70, 0x61, 0x72, 0x74, 0x75, 0x72, 0x65)
	if err != nil {
		return err
	}
	err = z.Departure.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "schedule_relationship"
	err = en.Append(0xb5, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x5f, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x68, 0x69, 0x70)
	if err != nil {
		return err
	}
	if z.ScheduleRelationship == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteInt32(*z.ScheduleRelationship)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *StopTimeUpdate) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "stop_sequence"
	o = append(o, 0x85, 0xad, 0x73, 0x74, 0x6f, 0x70, 0x5f, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65)
	if z.StopSequence == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendUint32(o, *z.StopSequence)
	}
	// string "stop_id"
	o = append(o, 0xa7, 0x73, 0x74, 0x6f, 0x70, 0x5f, 0x69, 0x64)
	if z.StopID == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendString(o, *z.StopID)
	}
	// string "arrival"
	o = append(o, 0xa7, 0x61, 0x72, 0x72, 0x69, 0x76, 0x61, 0x6c)
	o, err = z.Arrival.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "departure"
	o = append(o, 0xa9, 0x64, 0x65, 0x70, 0x61, 0x72, 0x74, 0x75, 0x72, 0x65)
	o, err = z.Departure.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "schedule_relationship"
	o = append(o, 0xb5, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x5f, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x68, 0x69, 0x70)
	if z.ScheduleRelationship == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendInt32(o, *z.ScheduleRelationship)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *StopTimeUpdate) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zcmr uint32
	zcmr, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zcmr > 0 {
		zcmr--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "stop_sequence":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.StopSequence = nil
			} else {
				if z.StopSequence == nil {
					z.StopSequence = new(uint32)
				}
				*z.StopSequence, bts, err = msgp.ReadUint32Bytes(bts)
				if err != nil {
					return
				}
			}
		case "stop_id":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.StopID = nil
			} else {
				if z.StopID == nil {
					z.StopID = new(string)
				}
				*z.StopID, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "arrival":
			bts, err = z.Arrival.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "departure":
			bts, err = z.Departure.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "schedule_relationship":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.ScheduleRelationship = nil
			} else {
				if z.ScheduleRelationship == nil {
					z.ScheduleRelationship = new(int32)
				}
				*z.ScheduleRelationship, bts, err = msgp.ReadInt32Bytes(bts)
				if err != nil {
					return
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *StopTimeUpdate) Msgsize() (s int) {
	s = 1 + 14
	if z.StopSequence == nil {
		s += msgp.NilSize
	} else {
		s += msgp.Uint32Size
	}
	s += 8
	if z.StopID == nil {
		s += msgp.NilSize
	} else {
		s += msgp.StringPrefixSize + len(*z.StopID)
	}
	s += 8 + z.Arrival.Msgsize() + 10 + z.Departure.Msgsize() + 22
	if z.ScheduleRelationship == nil {
		s += msgp.NilSize
	} else {
		s += msgp.Int32Size
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *TripDescriptor) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zajw uint32
	zajw, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zajw > 0 {
		zajw--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "trip_id":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.TripID = nil
			} else {
				if z.TripID == nil {
					z.TripID = new(string)
				}
				*z.TripID, err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "route_id":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.RouteID = nil
			} else {
				if z.RouteID == nil {
					z.RouteID = new(string)
				}
				*z.RouteID, err = dc.ReadString()
				if err != nil {
					return
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *TripDescriptor) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "trip_id"
	err = en.Append(0x82, 0xa7, 0x74, 0x72, 0x69, 0x70, 0x5f, 0x69, 0x64)
	if err != nil {
		return err
	}
	if z.TripID == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteString(*z.TripID)
		if err != nil {
			return
		}
	}
	// write "route_id"
	err = en.Append(0xa8, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x5f, 0x69, 0x64)
	if err != nil {
		return err
	}
	if z.RouteID == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteString(*z.RouteID)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *TripDescriptor) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "trip_id"
	o = append(o, 0x82, 0xa7, 0x74, 0x72, 0x69, 0x70, 0x5f, 0x69, 0x64)
	if z.TripID == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendString(o, *z.TripID)
	}
	// string "route_id"
	o = append(o, 0xa8, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x5f, 0x69, 0x64)
	if z.RouteID == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendString(o, *z.RouteID)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *TripDescriptor) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zwht uint32
	zwht, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zwht > 0 {
		zwht--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "trip_id":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.TripID = nil
			} else {
				if z.TripID == nil {
					z.TripID = new(string)
				}
				*z.TripID, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "route_id":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.RouteID = nil
			} else {
				if z.RouteID == nil {
					z.RouteID = new(string)
				}
				*z.RouteID, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *TripDescriptor) Msgsize() (s int) {
	s = 1 + 8
	if z.TripID == nil {
		s += msgp.NilSize
	} else {
		s += msgp.StringPrefixSize + len(*z.TripID)
	}
	s += 9
	if z.RouteID == nil {
		s += msgp.NilSize
	} else {
		s += msgp.StringPrefixSize + len(*z.RouteID)
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *TripUpdate) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zcua uint32
	zcua, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zcua > 0 {
		zcua--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "trip":
			err = z.Trip.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "vehicle":
			err = z.Vehicle.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "stop_time_update":
			var zxhx uint32
			zxhx, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.StopTimeUpdates) >= int(zxhx) {
				z.StopTimeUpdates = (z.StopTimeUpdates)[:zxhx]
			} else {
				z.StopTimeUpdates = make([]StopTimeUpdate, zxhx)
			}
			for zhct := range z.StopTimeUpdates {
				err = z.StopTimeUpdates[zhct].DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "timestamp":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Timestamp = nil
			} else {
				if z.Timestamp == nil {
					z.Timestamp = new(time.Time)
				}
				*z.Timestamp, err = dc.ReadTime()
				if err != nil {
					return
				}
			}
		case "delay":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Delay = nil
			} else {
				if z.Delay == nil {
					z.Delay = new(int32)
				}
				*z.Delay, err = dc.ReadInt32()
				if err != nil {
					return
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *TripUpdate) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 5
	// write "trip"
	err = en.Append(0x85, 0xa4, 0x74, 0x72, 0x69, 0x70)
	if err != nil {
		return err
	}
	err = z.Trip.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "vehicle"
	err = en.Append(0xa7, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65)
	if err != nil {
		return err
	}
	err = z.Vehicle.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "stop_time_update"
	err = en.Append(0xb0, 0x73, 0x74, 0x6f, 0x70, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.StopTimeUpdates)))
	if err != nil {
		return
	}
	for zhct := range z.StopTimeUpdates {
		err = z.StopTimeUpdates[zhct].EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "timestamp"
	err = en.Append(0xa9, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70)
	if err != nil {
		return err
	}
	if z.Timestamp == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteTime(*z.Timestamp)
		if err != nil {
			return
		}
	}
	// write "delay"
	err = en.Append(0xa5, 0x64, 0x65, 0x6c, 0x61, 0x79)
	if err != nil {
		return err
	}
	if z.Delay == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteInt32(*z.Delay)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *TripUpdate) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "trip"
	o = append(o, 0x85, 0xa4, 0x74, 0x72, 0x69, 0x70)
	o, err = z.Trip.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "vehicle"
	o = append(o, 0xa7, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65)
	o, err = z.Vehicle.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "stop_time_update"
	o = append(o, 0xb0, 0x73, 0x74, 0x6f, 0x70, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65)
	o = msgp.AppendArrayHeader(o, uint32(len(z.StopTimeUpdates)))
	for zhct := range z.StopTimeUpdates {
		o, err = z.StopTimeUpdates[zhct].MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "timestamp"
	o = append(o, 0xa9, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70)
	if z.Timestamp == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendTime(o, *z.Timestamp)
	}
	// string "delay"
	o = append(o, 0xa5, 0x64, 0x65, 0x6c, 0x61, 0x79)
	if z.Delay == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendInt32(o, *z.Delay)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *TripUpdate) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zlqf uint32
	zlqf, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zlqf > 0 {
		zlqf--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "trip":
			bts, err = z.Trip.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "vehicle":
			bts, err = z.Vehicle.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "stop_time_update":
			var zdaf uint32
			zdaf, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.StopTimeUpdates) >= int(zdaf) {
				z.StopTimeUpdates = (z.StopTimeUpdates)[:zdaf]
			} else {
				z.StopTimeUpdates = make([]StopTimeUpdate, zdaf)
			}
			for zhct := range z.StopTimeUpdates {
				bts, err = z.StopTimeUpdates[zhct].UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "timestamp":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Timestamp = nil
			} else {
				if z.Timestamp == nil {
					z.Timestamp = new(time.Time)
				}
				*z.Timestamp, bts, err = msgp.ReadTimeBytes(bts)
				if err != nil {
					return
				}
			}
		case "delay":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Delay = nil
			} else {
				if z.Delay == nil {
					z.Delay = new(int32)
				}
				*z.Delay, bts, err = msgp.ReadInt32Bytes(bts)
				if err != nil {
					return
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *TripUpdate) Msgsize() (s int) {
	s = 1 + 5 + z.Trip.Msgsize() + 8 + z.Vehicle.Msgsize() + 17 + msgp.ArrayHeaderSize
	for zhct := range z.StopTimeUpdates {
		s += z.StopTimeUpdates[zhct].Msgsize()
	}
	s += 10
	if z.Timestamp == nil {
		s += msgp.NilSize
	} else {
		s += msgp.TimeSize
	}
	s += 6
	if z.Delay == nil {
		s += msgp.NilSize
	} else {
		s += msgp.Int32Size
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *TripUpdatesMsg) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zjfb uint32
	zjfb, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zjfb > 0 {
		zjfb--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "trip_update":
			var zcxo uint32
			zcxo, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.TripUpdates) >= int(zcxo) {
				z.TripUpdates = (z.TripUpdates)[:zcxo]
			} else {
				z.TripUpdates = make([]TripUpdate, zcxo)
			}
			for zpks := range z.TripUpdates {
				err = z.TripUpdates[zpks].DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *TripUpdatesMsg) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "trip_update"
	err = en.Append(0x81, 0xab, 0x74, 0x72, 0x69, 0x70, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.TripUpdates)))
	if err != nil {
		return
	}
	for zpks := range z.TripUpdates {
		err = z.TripUpdates[zpks].EncodeMsg(en)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *TripUpdatesMsg) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "trip_update"
	o = append(o, 0x81, 0xab, 0x74, 0x72, 0x69, 0x70, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65)
	o = msgp.AppendArrayHeader(o, uint32(len(z.TripUpdates)))
	for zpks := range z.TripUpdates {
		o, err = z.TripUpdates[zpks].MarshalMsg(o)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *TripUpdatesMsg) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zeff uint32
	zeff, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zeff > 0 {
		zeff--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "trip_update":
			var zrsw uint32
			zrsw, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.TripUpdates) >= int(zrsw) {
				z.TripUpdates = (z.TripUpdates)[:zrsw]
			} else {
				z.TripUpdates = make([]TripUpdate, zrsw)
			}
			for zpks := range z.TripUpdates {
				bts, err = z.TripUpdates[zpks].UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *TripUpdatesMsg) Msgsize() (s int) {
	s = 1 + 12 + msgp.ArrayHeaderSize
	for zpks := range z.TripUpdates {
		s += z.TripUpdates[zpks].Msgsize()
	}
	return
}

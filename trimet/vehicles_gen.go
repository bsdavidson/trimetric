package trimet

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Position) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "lat":
			z.Latitude, err = dc.ReadFloat32()
			if err != nil {
				return
			}
		case "lng":
			z.Longitude, err = dc.ReadFloat32()
			if err != nil {
				return
			}
		case "bearing":
			z.Bearing, err = dc.ReadFloat32()
			if err != nil {
				return
			}
		case "odometer":
			z.Odometer, err = dc.ReadFloat64()
			if err != nil {
				return
			}
		case "speed":
			z.Speed, err = dc.ReadFloat32()
			if err != nil {
				return
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
func (z *Position) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 5
	// write "lat"
	err = en.Append(0x85, 0xa3, 0x6c, 0x61, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteFloat32(z.Latitude)
	if err != nil {
		return
	}
	// write "lng"
	err = en.Append(0xa3, 0x6c, 0x6e, 0x67)
	if err != nil {
		return err
	}
	err = en.WriteFloat32(z.Longitude)
	if err != nil {
		return
	}
	// write "bearing"
	err = en.Append(0xa7, 0x62, 0x65, 0x61, 0x72, 0x69, 0x6e, 0x67)
	if err != nil {
		return err
	}
	err = en.WriteFloat32(z.Bearing)
	if err != nil {
		return
	}
	// write "odometer"
	err = en.Append(0xa8, 0x6f, 0x64, 0x6f, 0x6d, 0x65, 0x74, 0x65, 0x72)
	if err != nil {
		return err
	}
	err = en.WriteFloat64(z.Odometer)
	if err != nil {
		return
	}
	// write "speed"
	err = en.Append(0xa5, 0x73, 0x70, 0x65, 0x65, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteFloat32(z.Speed)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Position) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "lat"
	o = append(o, 0x85, 0xa3, 0x6c, 0x61, 0x74)
	o = msgp.AppendFloat32(o, z.Latitude)
	// string "lng"
	o = append(o, 0xa3, 0x6c, 0x6e, 0x67)
	o = msgp.AppendFloat32(o, z.Longitude)
	// string "bearing"
	o = append(o, 0xa7, 0x62, 0x65, 0x61, 0x72, 0x69, 0x6e, 0x67)
	o = msgp.AppendFloat32(o, z.Bearing)
	// string "odometer"
	o = append(o, 0xa8, 0x6f, 0x64, 0x6f, 0x6d, 0x65, 0x74, 0x65, 0x72)
	o = msgp.AppendFloat64(o, z.Odometer)
	// string "speed"
	o = append(o, 0xa5, 0x73, 0x70, 0x65, 0x65, 0x64)
	o = msgp.AppendFloat32(o, z.Speed)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Position) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "lat":
			z.Latitude, bts, err = msgp.ReadFloat32Bytes(bts)
			if err != nil {
				return
			}
		case "lng":
			z.Longitude, bts, err = msgp.ReadFloat32Bytes(bts)
			if err != nil {
				return
			}
		case "bearing":
			z.Bearing, bts, err = msgp.ReadFloat32Bytes(bts)
			if err != nil {
				return
			}
		case "odometer":
			z.Odometer, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		case "speed":
			z.Speed, bts, err = msgp.ReadFloat32Bytes(bts)
			if err != nil {
				return
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
func (z *Position) Msgsize() (s int) {
	s = 1 + 4 + msgp.Float32Size + 4 + msgp.Float32Size + 8 + msgp.Float32Size + 9 + msgp.Float64Size + 6 + msgp.Float32Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *RouteType) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zbai int
		zbai, err = dc.ReadInt()
		(*z) = RouteType(zbai)
	}
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z RouteType) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteInt(int(z))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z RouteType) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendInt(o, int(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *RouteType) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zcmr int
		zcmr, bts, err = msgp.ReadIntBytes(bts)
		(*z) = RouteType(zcmr)
	}
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z RouteType) Msgsize() (s int) {
	s = msgp.IntSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *VehicleDescriptor) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "id":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.ID = nil
			} else {
				if z.ID == nil {
					z.ID = new(string)
				}
				*z.ID, err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "label":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Label = nil
			} else {
				if z.Label == nil {
					z.Label = new(string)
				}
				*z.Label, err = dc.ReadString()
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
func (z *VehicleDescriptor) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "id"
	err = en.Append(0x82, 0xa2, 0x69, 0x64)
	if err != nil {
		return err
	}
	if z.ID == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteString(*z.ID)
		if err != nil {
			return
		}
	}
	// write "label"
	err = en.Append(0xa5, 0x6c, 0x61, 0x62, 0x65, 0x6c)
	if err != nil {
		return err
	}
	if z.Label == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteString(*z.Label)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *VehicleDescriptor) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "id"
	o = append(o, 0x82, 0xa2, 0x69, 0x64)
	if z.ID == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendString(o, *z.ID)
	}
	// string "label"
	o = append(o, 0xa5, 0x6c, 0x61, 0x62, 0x65, 0x6c)
	if z.Label == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendString(o, *z.Label)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *VehicleDescriptor) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "id":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.ID = nil
			} else {
				if z.ID == nil {
					z.ID = new(string)
				}
				*z.ID, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "label":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Label = nil
			} else {
				if z.Label == nil {
					z.Label = new(string)
				}
				*z.Label, bts, err = msgp.ReadStringBytes(bts)
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
func (z *VehicleDescriptor) Msgsize() (s int) {
	s = 1 + 3
	if z.ID == nil {
		s += msgp.NilSize
	} else {
		s += msgp.StringPrefixSize + len(*z.ID)
	}
	s += 6
	if z.Label == nil {
		s += msgp.NilSize
	} else {
		s += msgp.StringPrefixSize + len(*z.Label)
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *VehiclePosition) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zhct uint32
	zhct, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zhct > 0 {
		zhct--
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
		case "position":
			err = z.Position.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "current_stop_sequence":
			z.CurrentStopSequence, err = dc.ReadUint32()
			if err != nil {
				return
			}
		case "stop_id":
			z.StopID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "current_status":
			z.CurrentStatus, err = dc.ReadInt32()
			if err != nil {
				return
			}
		case "timestamp":
			z.Timestamp, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "congestion_level":
			z.CongestionLevel, err = dc.ReadInt32()
			if err != nil {
				return
			}
		case "occupancy_status":
			z.OccupancyStatus, err = dc.ReadInt32()
			if err != nil {
				return
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
func (z *VehiclePosition) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 9
	// write "trip"
	err = en.Append(0x89, 0xa4, 0x74, 0x72, 0x69, 0x70)
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
	// write "position"
	err = en.Append(0xa8, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return err
	}
	err = z.Position.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "current_stop_sequence"
	err = en.Append(0xb5, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x6f, 0x70, 0x5f, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteUint32(z.CurrentStopSequence)
	if err != nil {
		return
	}
	// write "stop_id"
	err = en.Append(0xa7, 0x73, 0x74, 0x6f, 0x70, 0x5f, 0x69, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteString(z.StopID)
	if err != nil {
		return
	}
	// write "current_status"
	err = en.Append(0xae, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteInt32(z.CurrentStatus)
	if err != nil {
		return
	}
	// write "timestamp"
	err = en.Append(0xa9, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70)
	if err != nil {
		return err
	}
	err = en.WriteUint64(z.Timestamp)
	if err != nil {
		return
	}
	// write "congestion_level"
	err = en.Append(0xb0, 0x63, 0x6f, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6c, 0x65, 0x76, 0x65, 0x6c)
	if err != nil {
		return err
	}
	err = en.WriteInt32(z.CongestionLevel)
	if err != nil {
		return
	}
	// write "occupancy_status"
	err = en.Append(0xb0, 0x6f, 0x63, 0x63, 0x75, 0x70, 0x61, 0x6e, 0x63, 0x79, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteInt32(z.OccupancyStatus)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *VehiclePosition) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 9
	// string "trip"
	o = append(o, 0x89, 0xa4, 0x74, 0x72, 0x69, 0x70)
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
	// string "position"
	o = append(o, 0xa8, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e)
	o, err = z.Position.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "current_stop_sequence"
	o = append(o, 0xb5, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x6f, 0x70, 0x5f, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65)
	o = msgp.AppendUint32(o, z.CurrentStopSequence)
	// string "stop_id"
	o = append(o, 0xa7, 0x73, 0x74, 0x6f, 0x70, 0x5f, 0x69, 0x64)
	o = msgp.AppendString(o, z.StopID)
	// string "current_status"
	o = append(o, 0xae, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73)
	o = msgp.AppendInt32(o, z.CurrentStatus)
	// string "timestamp"
	o = append(o, 0xa9, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70)
	o = msgp.AppendUint64(o, z.Timestamp)
	// string "congestion_level"
	o = append(o, 0xb0, 0x63, 0x6f, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6c, 0x65, 0x76, 0x65, 0x6c)
	o = msgp.AppendInt32(o, z.CongestionLevel)
	// string "occupancy_status"
	o = append(o, 0xb0, 0x6f, 0x63, 0x63, 0x75, 0x70, 0x61, 0x6e, 0x63, 0x79, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73)
	o = msgp.AppendInt32(o, z.OccupancyStatus)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *VehiclePosition) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zcua uint32
	zcua, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zcua > 0 {
		zcua--
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
		case "position":
			bts, err = z.Position.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "current_stop_sequence":
			z.CurrentStopSequence, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				return
			}
		case "stop_id":
			z.StopID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "current_status":
			z.CurrentStatus, bts, err = msgp.ReadInt32Bytes(bts)
			if err != nil {
				return
			}
		case "timestamp":
			z.Timestamp, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "congestion_level":
			z.CongestionLevel, bts, err = msgp.ReadInt32Bytes(bts)
			if err != nil {
				return
			}
		case "occupancy_status":
			z.OccupancyStatus, bts, err = msgp.ReadInt32Bytes(bts)
			if err != nil {
				return
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
func (z *VehiclePosition) Msgsize() (s int) {
	s = 1 + 5 + z.Trip.Msgsize() + 8 + z.Vehicle.Msgsize() + 9 + z.Position.Msgsize() + 22 + msgp.Uint32Size + 8 + msgp.StringPrefixSize + len(z.StopID) + 15 + msgp.Int32Size + 10 + msgp.Uint64Size + 17 + msgp.Int32Size + 17 + msgp.Int32Size
	return
}

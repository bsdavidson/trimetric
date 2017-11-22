package trimet

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *CalendarDate) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "service_id":
			z.ServiceID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "date":
			z.Date, err = dc.ReadTime()
			if err != nil {
				return
			}
		case "exception_type":
			z.ExceptionType, err = dc.ReadInt()
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
func (z CalendarDate) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "service_id"
	err = en.Append(0x83, 0xaa, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteString(z.ServiceID)
	if err != nil {
		return
	}
	// write "date"
	err = en.Append(0xa4, 0x64, 0x61, 0x74, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteTime(z.Date)
	if err != nil {
		return
	}
	// write "exception_type"
	err = en.Append(0xae, 0x65, 0x78, 0x63, 0x65, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.ExceptionType)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z CalendarDate) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "service_id"
	o = append(o, 0x83, 0xaa, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64)
	o = msgp.AppendString(o, z.ServiceID)
	// string "date"
	o = append(o, 0xa4, 0x64, 0x61, 0x74, 0x65)
	o = msgp.AppendTime(o, z.Date)
	// string "exception_type"
	o = append(o, 0xae, 0x65, 0x78, 0x63, 0x65, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65)
	o = msgp.AppendInt(o, z.ExceptionType)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *CalendarDate) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "service_id":
			z.ServiceID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "date":
			z.Date, bts, err = msgp.ReadTimeBytes(bts)
			if err != nil {
				return
			}
		case "exception_type":
			z.ExceptionType, bts, err = msgp.ReadIntBytes(bts)
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
func (z CalendarDate) Msgsize() (s int) {
	s = 1 + 11 + msgp.StringPrefixSize + len(z.ServiceID) + 5 + msgp.TimeSize + 15 + msgp.IntSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Route) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "id":
			z.RouteID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "agency_id":
			z.AgencyID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "short_name":
			z.ShortName, err = dc.ReadString()
			if err != nil {
				return
			}
		case "long_name":
			z.LongName, err = dc.ReadString()
			if err != nil {
				return
			}
		case "type":
			z.Type, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "url":
			z.URL, err = dc.ReadString()
			if err != nil {
				return
			}
		case "color":
			z.Color, err = dc.ReadString()
			if err != nil {
				return
			}
		case "text_color":
			z.TextColor, err = dc.ReadString()
			if err != nil {
				return
			}
		case "sort_order":
			z.SortOrder, err = dc.ReadInt()
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
func (z *Route) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 9
	// write "id"
	err = en.Append(0x89, 0xa2, 0x69, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteString(z.RouteID)
	if err != nil {
		return
	}
	// write "agency_id"
	err = en.Append(0xa9, 0x61, 0x67, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x69, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteString(z.AgencyID)
	if err != nil {
		return
	}
	// write "short_name"
	err = en.Append(0xaa, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.ShortName)
	if err != nil {
		return
	}
	// write "long_name"
	err = en.Append(0xa9, 0x6c, 0x6f, 0x6e, 0x67, 0x5f, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.LongName)
	if err != nil {
		return
	}
	// write "type"
	err = en.Append(0xa4, 0x74, 0x79, 0x70, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.Type)
	if err != nil {
		return
	}
	// write "url"
	err = en.Append(0xa3, 0x75, 0x72, 0x6c)
	if err != nil {
		return err
	}
	err = en.WriteString(z.URL)
	if err != nil {
		return
	}
	// write "color"
	err = en.Append(0xa5, 0x63, 0x6f, 0x6c, 0x6f, 0x72)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Color)
	if err != nil {
		return
	}
	// write "text_color"
	err = en.Append(0xaa, 0x74, 0x65, 0x78, 0x74, 0x5f, 0x63, 0x6f, 0x6c, 0x6f, 0x72)
	if err != nil {
		return err
	}
	err = en.WriteString(z.TextColor)
	if err != nil {
		return
	}
	// write "sort_order"
	err = en.Append(0xaa, 0x73, 0x6f, 0x72, 0x74, 0x5f, 0x6f, 0x72, 0x64, 0x65, 0x72)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.SortOrder)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Route) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 9
	// string "id"
	o = append(o, 0x89, 0xa2, 0x69, 0x64)
	o = msgp.AppendString(o, z.RouteID)
	// string "agency_id"
	o = append(o, 0xa9, 0x61, 0x67, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x69, 0x64)
	o = msgp.AppendString(o, z.AgencyID)
	// string "short_name"
	o = append(o, 0xaa, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.ShortName)
	// string "long_name"
	o = append(o, 0xa9, 0x6c, 0x6f, 0x6e, 0x67, 0x5f, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.LongName)
	// string "type"
	o = append(o, 0xa4, 0x74, 0x79, 0x70, 0x65)
	o = msgp.AppendInt(o, z.Type)
	// string "url"
	o = append(o, 0xa3, 0x75, 0x72, 0x6c)
	o = msgp.AppendString(o, z.URL)
	// string "color"
	o = append(o, 0xa5, 0x63, 0x6f, 0x6c, 0x6f, 0x72)
	o = msgp.AppendString(o, z.Color)
	// string "text_color"
	o = append(o, 0xaa, 0x74, 0x65, 0x78, 0x74, 0x5f, 0x63, 0x6f, 0x6c, 0x6f, 0x72)
	o = msgp.AppendString(o, z.TextColor)
	// string "sort_order"
	o = append(o, 0xaa, 0x73, 0x6f, 0x72, 0x74, 0x5f, 0x6f, 0x72, 0x64, 0x65, 0x72)
	o = msgp.AppendInt(o, z.SortOrder)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Route) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "id":
			z.RouteID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "agency_id":
			z.AgencyID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "short_name":
			z.ShortName, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "long_name":
			z.LongName, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "type":
			z.Type, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "url":
			z.URL, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "color":
			z.Color, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "text_color":
			z.TextColor, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "sort_order":
			z.SortOrder, bts, err = msgp.ReadIntBytes(bts)
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
func (z *Route) Msgsize() (s int) {
	s = 1 + 3 + msgp.StringPrefixSize + len(z.RouteID) + 10 + msgp.StringPrefixSize + len(z.AgencyID) + 11 + msgp.StringPrefixSize + len(z.ShortName) + 10 + msgp.StringPrefixSize + len(z.LongName) + 5 + msgp.IntSize + 4 + msgp.StringPrefixSize + len(z.URL) + 6 + msgp.StringPrefixSize + len(z.Color) + 11 + msgp.StringPrefixSize + len(z.TextColor) + 11 + msgp.IntSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Stop) DecodeMsg(dc *msgp.Reader) (err error) {
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
			z.ID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "code":
			z.Code, err = dc.ReadString()
			if err != nil {
				return
			}
		case "name":
			z.Name, err = dc.ReadString()
			if err != nil {
				return
			}
		case "desc":
			z.Desc, err = dc.ReadString()
			if err != nil {
				return
			}
		case "lat":
			z.Lat, err = dc.ReadFloat64()
			if err != nil {
				return
			}
		case "lng":
			z.Lon, err = dc.ReadFloat64()
			if err != nil {
				return
			}
		case "zone_id":
			z.ZoneID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "url":
			z.URL, err = dc.ReadString()
			if err != nil {
				return
			}
		case "location_type":
			z.LocationType, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "parent_station":
			z.ParentStation, err = dc.ReadString()
			if err != nil {
				return
			}
		case "direction":
			z.Direction, err = dc.ReadString()
			if err != nil {
				return
			}
		case "position":
			z.Position, err = dc.ReadString()
			if err != nil {
				return
			}
		case "wheelchair_boarding":
			z.WheelchairBoarding, err = dc.ReadInt()
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
func (z *Stop) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 13
	// write "id"
	err = en.Append(0x8d, 0xa2, 0x69, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteString(z.ID)
	if err != nil {
		return
	}
	// write "code"
	err = en.Append(0xa4, 0x63, 0x6f, 0x64, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Code)
	if err != nil {
		return
	}
	// write "name"
	err = en.Append(0xa4, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Name)
	if err != nil {
		return
	}
	// write "desc"
	err = en.Append(0xa4, 0x64, 0x65, 0x73, 0x63)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Desc)
	if err != nil {
		return
	}
	// write "lat"
	err = en.Append(0xa3, 0x6c, 0x61, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteFloat64(z.Lat)
	if err != nil {
		return
	}
	// write "lng"
	err = en.Append(0xa3, 0x6c, 0x6e, 0x67)
	if err != nil {
		return err
	}
	err = en.WriteFloat64(z.Lon)
	if err != nil {
		return
	}
	// write "zone_id"
	err = en.Append(0xa7, 0x7a, 0x6f, 0x6e, 0x65, 0x5f, 0x69, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteString(z.ZoneID)
	if err != nil {
		return
	}
	// write "url"
	err = en.Append(0xa3, 0x75, 0x72, 0x6c)
	if err != nil {
		return err
	}
	err = en.WriteString(z.URL)
	if err != nil {
		return
	}
	// write "location_type"
	err = en.Append(0xad, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.LocationType)
	if err != nil {
		return
	}
	// write "parent_station"
	err = en.Append(0xae, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return err
	}
	err = en.WriteString(z.ParentStation)
	if err != nil {
		return
	}
	// write "direction"
	err = en.Append(0xa9, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Direction)
	if err != nil {
		return
	}
	// write "position"
	err = en.Append(0xa8, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Position)
	if err != nil {
		return
	}
	// write "wheelchair_boarding"
	err = en.Append(0xb3, 0x77, 0x68, 0x65, 0x65, 0x6c, 0x63, 0x68, 0x61, 0x69, 0x72, 0x5f, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.WheelchairBoarding)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Stop) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 13
	// string "id"
	o = append(o, 0x8d, 0xa2, 0x69, 0x64)
	o = msgp.AppendString(o, z.ID)
	// string "code"
	o = append(o, 0xa4, 0x63, 0x6f, 0x64, 0x65)
	o = msgp.AppendString(o, z.Code)
	// string "name"
	o = append(o, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Name)
	// string "desc"
	o = append(o, 0xa4, 0x64, 0x65, 0x73, 0x63)
	o = msgp.AppendString(o, z.Desc)
	// string "lat"
	o = append(o, 0xa3, 0x6c, 0x61, 0x74)
	o = msgp.AppendFloat64(o, z.Lat)
	// string "lng"
	o = append(o, 0xa3, 0x6c, 0x6e, 0x67)
	o = msgp.AppendFloat64(o, z.Lon)
	// string "zone_id"
	o = append(o, 0xa7, 0x7a, 0x6f, 0x6e, 0x65, 0x5f, 0x69, 0x64)
	o = msgp.AppendString(o, z.ZoneID)
	// string "url"
	o = append(o, 0xa3, 0x75, 0x72, 0x6c)
	o = msgp.AppendString(o, z.URL)
	// string "location_type"
	o = append(o, 0xad, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65)
	o = msgp.AppendInt(o, z.LocationType)
	// string "parent_station"
	o = append(o, 0xae, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	o = msgp.AppendString(o, z.ParentStation)
	// string "direction"
	o = append(o, 0xa9, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e)
	o = msgp.AppendString(o, z.Direction)
	// string "position"
	o = append(o, 0xa8, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e)
	o = msgp.AppendString(o, z.Position)
	// string "wheelchair_boarding"
	o = append(o, 0xb3, 0x77, 0x68, 0x65, 0x65, 0x6c, 0x63, 0x68, 0x61, 0x69, 0x72, 0x5f, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67)
	o = msgp.AppendInt(o, z.WheelchairBoarding)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Stop) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
			z.ID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "code":
			z.Code, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "name":
			z.Name, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "desc":
			z.Desc, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "lat":
			z.Lat, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		case "lng":
			z.Lon, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		case "zone_id":
			z.ZoneID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "url":
			z.URL, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "location_type":
			z.LocationType, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "parent_station":
			z.ParentStation, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "direction":
			z.Direction, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "position":
			z.Position, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "wheelchair_boarding":
			z.WheelchairBoarding, bts, err = msgp.ReadIntBytes(bts)
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
func (z *Stop) Msgsize() (s int) {
	s = 1 + 3 + msgp.StringPrefixSize + len(z.ID) + 5 + msgp.StringPrefixSize + len(z.Code) + 5 + msgp.StringPrefixSize + len(z.Name) + 5 + msgp.StringPrefixSize + len(z.Desc) + 4 + msgp.Float64Size + 4 + msgp.Float64Size + 8 + msgp.StringPrefixSize + len(z.ZoneID) + 4 + msgp.StringPrefixSize + len(z.URL) + 14 + msgp.IntSize + 15 + msgp.StringPrefixSize + len(z.ParentStation) + 10 + msgp.StringPrefixSize + len(z.Direction) + 9 + msgp.StringPrefixSize + len(z.Position) + 20 + msgp.IntSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *StopTime) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "trip_id":
			z.TripID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "arrival_time":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.ArrivalTime = nil
			} else {
				if z.ArrivalTime == nil {
					z.ArrivalTime = new(Time)
				}
				err = dc.ReadExtension(z.ArrivalTime)
				if err != nil {
					return
				}
			}
		case "departure_time":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.DepartureTime = nil
			} else {
				if z.DepartureTime == nil {
					z.DepartureTime = new(Time)
				}
				err = dc.ReadExtension(z.DepartureTime)
				if err != nil {
					return
				}
			}
		case "stop_id":
			z.StopID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "stop_sequence":
			z.StopSequence, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "stop_headsign":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.StopHeadsign = nil
			} else {
				if z.StopHeadsign == nil {
					z.StopHeadsign = new(string)
				}
				*z.StopHeadsign, err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "pickup_type":
			z.PickupType, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "drop_off_type":
			z.DropOffType, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "shape_dist_traveled":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.ShapeDistTraveled = nil
			} else {
				if z.ShapeDistTraveled == nil {
					z.ShapeDistTraveled = new(float64)
				}
				*z.ShapeDistTraveled, err = dc.ReadFloat64()
				if err != nil {
					return
				}
			}
		case "timepoint":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Timepoint = nil
			} else {
				if z.Timepoint == nil {
					z.Timepoint = new(int)
				}
				*z.Timepoint, err = dc.ReadInt()
				if err != nil {
					return
				}
			}
		case "continuous_drop_off":
			z.ContinuousDropOff, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "continuous_pickup":
			z.ContinuousPickup, err = dc.ReadInt()
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
func (z *StopTime) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 12
	// write "trip_id"
	err = en.Append(0x8c, 0xa7, 0x74, 0x72, 0x69, 0x70, 0x5f, 0x69, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteString(z.TripID)
	if err != nil {
		return
	}
	// write "arrival_time"
	err = en.Append(0xac, 0x61, 0x72, 0x72, 0x69, 0x76, 0x61, 0x6c, 0x5f, 0x74, 0x69, 0x6d, 0x65)
	if err != nil {
		return err
	}
	if z.ArrivalTime == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteExtension(z.ArrivalTime)
		if err != nil {
			return
		}
	}
	// write "departure_time"
	err = en.Append(0xae, 0x64, 0x65, 0x70, 0x61, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65)
	if err != nil {
		return err
	}
	if z.DepartureTime == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteExtension(z.DepartureTime)
		if err != nil {
			return
		}
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
	// write "stop_sequence"
	err = en.Append(0xad, 0x73, 0x74, 0x6f, 0x70, 0x5f, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.StopSequence)
	if err != nil {
		return
	}
	// write "stop_headsign"
	err = en.Append(0xad, 0x73, 0x74, 0x6f, 0x70, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x73, 0x69, 0x67, 0x6e)
	if err != nil {
		return err
	}
	if z.StopHeadsign == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteString(*z.StopHeadsign)
		if err != nil {
			return
		}
	}
	// write "pickup_type"
	err = en.Append(0xab, 0x70, 0x69, 0x63, 0x6b, 0x75, 0x70, 0x5f, 0x74, 0x79, 0x70, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.PickupType)
	if err != nil {
		return
	}
	// write "drop_off_type"
	err = en.Append(0xad, 0x64, 0x72, 0x6f, 0x70, 0x5f, 0x6f, 0x66, 0x66, 0x5f, 0x74, 0x79, 0x70, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.DropOffType)
	if err != nil {
		return
	}
	// write "shape_dist_traveled"
	err = en.Append(0xb3, 0x73, 0x68, 0x61, 0x70, 0x65, 0x5f, 0x64, 0x69, 0x73, 0x74, 0x5f, 0x74, 0x72, 0x61, 0x76, 0x65, 0x6c, 0x65, 0x64)
	if err != nil {
		return err
	}
	if z.ShapeDistTraveled == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteFloat64(*z.ShapeDistTraveled)
		if err != nil {
			return
		}
	}
	// write "timepoint"
	err = en.Append(0xa9, 0x74, 0x69, 0x6d, 0x65, 0x70, 0x6f, 0x69, 0x6e, 0x74)
	if err != nil {
		return err
	}
	if z.Timepoint == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteInt(*z.Timepoint)
		if err != nil {
			return
		}
	}
	// write "continuous_drop_off"
	err = en.Append(0xb3, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x6f, 0x75, 0x73, 0x5f, 0x64, 0x72, 0x6f, 0x70, 0x5f, 0x6f, 0x66, 0x66)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.ContinuousDropOff)
	if err != nil {
		return
	}
	// write "continuous_pickup"
	err = en.Append(0xb1, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x6f, 0x75, 0x73, 0x5f, 0x70, 0x69, 0x63, 0x6b, 0x75, 0x70)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.ContinuousPickup)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *StopTime) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 12
	// string "trip_id"
	o = append(o, 0x8c, 0xa7, 0x74, 0x72, 0x69, 0x70, 0x5f, 0x69, 0x64)
	o = msgp.AppendString(o, z.TripID)
	// string "arrival_time"
	o = append(o, 0xac, 0x61, 0x72, 0x72, 0x69, 0x76, 0x61, 0x6c, 0x5f, 0x74, 0x69, 0x6d, 0x65)
	if z.ArrivalTime == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = msgp.AppendExtension(o, z.ArrivalTime)
		if err != nil {
			return
		}
	}
	// string "departure_time"
	o = append(o, 0xae, 0x64, 0x65, 0x70, 0x61, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65)
	if z.DepartureTime == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = msgp.AppendExtension(o, z.DepartureTime)
		if err != nil {
			return
		}
	}
	// string "stop_id"
	o = append(o, 0xa7, 0x73, 0x74, 0x6f, 0x70, 0x5f, 0x69, 0x64)
	o = msgp.AppendString(o, z.StopID)
	// string "stop_sequence"
	o = append(o, 0xad, 0x73, 0x74, 0x6f, 0x70, 0x5f, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65)
	o = msgp.AppendInt(o, z.StopSequence)
	// string "stop_headsign"
	o = append(o, 0xad, 0x73, 0x74, 0x6f, 0x70, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x73, 0x69, 0x67, 0x6e)
	if z.StopHeadsign == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendString(o, *z.StopHeadsign)
	}
	// string "pickup_type"
	o = append(o, 0xab, 0x70, 0x69, 0x63, 0x6b, 0x75, 0x70, 0x5f, 0x74, 0x79, 0x70, 0x65)
	o = msgp.AppendInt(o, z.PickupType)
	// string "drop_off_type"
	o = append(o, 0xad, 0x64, 0x72, 0x6f, 0x70, 0x5f, 0x6f, 0x66, 0x66, 0x5f, 0x74, 0x79, 0x70, 0x65)
	o = msgp.AppendInt(o, z.DropOffType)
	// string "shape_dist_traveled"
	o = append(o, 0xb3, 0x73, 0x68, 0x61, 0x70, 0x65, 0x5f, 0x64, 0x69, 0x73, 0x74, 0x5f, 0x74, 0x72, 0x61, 0x76, 0x65, 0x6c, 0x65, 0x64)
	if z.ShapeDistTraveled == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendFloat64(o, *z.ShapeDistTraveled)
	}
	// string "timepoint"
	o = append(o, 0xa9, 0x74, 0x69, 0x6d, 0x65, 0x70, 0x6f, 0x69, 0x6e, 0x74)
	if z.Timepoint == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendInt(o, *z.Timepoint)
	}
	// string "continuous_drop_off"
	o = append(o, 0xb3, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x6f, 0x75, 0x73, 0x5f, 0x64, 0x72, 0x6f, 0x70, 0x5f, 0x6f, 0x66, 0x66)
	o = msgp.AppendInt(o, z.ContinuousDropOff)
	// string "continuous_pickup"
	o = append(o, 0xb1, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x6f, 0x75, 0x73, 0x5f, 0x70, 0x69, 0x63, 0x6b, 0x75, 0x70)
	o = msgp.AppendInt(o, z.ContinuousPickup)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *StopTime) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "trip_id":
			z.TripID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "arrival_time":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.ArrivalTime = nil
			} else {
				if z.ArrivalTime == nil {
					z.ArrivalTime = new(Time)
				}
				bts, err = msgp.ReadExtensionBytes(bts, z.ArrivalTime)
				if err != nil {
					return
				}
			}
		case "departure_time":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.DepartureTime = nil
			} else {
				if z.DepartureTime == nil {
					z.DepartureTime = new(Time)
				}
				bts, err = msgp.ReadExtensionBytes(bts, z.DepartureTime)
				if err != nil {
					return
				}
			}
		case "stop_id":
			z.StopID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "stop_sequence":
			z.StopSequence, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "stop_headsign":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.StopHeadsign = nil
			} else {
				if z.StopHeadsign == nil {
					z.StopHeadsign = new(string)
				}
				*z.StopHeadsign, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "pickup_type":
			z.PickupType, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "drop_off_type":
			z.DropOffType, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "shape_dist_traveled":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.ShapeDistTraveled = nil
			} else {
				if z.ShapeDistTraveled == nil {
					z.ShapeDistTraveled = new(float64)
				}
				*z.ShapeDistTraveled, bts, err = msgp.ReadFloat64Bytes(bts)
				if err != nil {
					return
				}
			}
		case "timepoint":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Timepoint = nil
			} else {
				if z.Timepoint == nil {
					z.Timepoint = new(int)
				}
				*z.Timepoint, bts, err = msgp.ReadIntBytes(bts)
				if err != nil {
					return
				}
			}
		case "continuous_drop_off":
			z.ContinuousDropOff, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "continuous_pickup":
			z.ContinuousPickup, bts, err = msgp.ReadIntBytes(bts)
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
func (z *StopTime) Msgsize() (s int) {
	s = 1 + 8 + msgp.StringPrefixSize + len(z.TripID) + 13
	if z.ArrivalTime == nil {
		s += msgp.NilSize
	} else {
		s += msgp.ExtensionPrefixSize + z.ArrivalTime.Len()
	}
	s += 15
	if z.DepartureTime == nil {
		s += msgp.NilSize
	} else {
		s += msgp.ExtensionPrefixSize + z.DepartureTime.Len()
	}
	s += 8 + msgp.StringPrefixSize + len(z.StopID) + 14 + msgp.IntSize + 14
	if z.StopHeadsign == nil {
		s += msgp.NilSize
	} else {
		s += msgp.StringPrefixSize + len(*z.StopHeadsign)
	}
	s += 12 + msgp.IntSize + 14 + msgp.IntSize + 20
	if z.ShapeDistTraveled == nil {
		s += msgp.NilSize
	} else {
		s += msgp.Float64Size
	}
	s += 10
	if z.Timepoint == nil {
		s += msgp.NilSize
	} else {
		s += msgp.IntSize
	}
	s += 20 + msgp.IntSize + 18 + msgp.IntSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Trip) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zxhx uint32
	zxhx, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zxhx > 0 {
		zxhx--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "id":
			z.ID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "route_id":
			z.RouteID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "service_id":
			z.ServiceID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "direction_id":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.DirectionID = nil
			} else {
				if z.DirectionID == nil {
					z.DirectionID = new(int)
				}
				*z.DirectionID, err = dc.ReadInt()
				if err != nil {
					return
				}
			}
		case "block_id":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.BlockID = nil
			} else {
				if z.BlockID == nil {
					z.BlockID = new(string)
				}
				*z.BlockID, err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "shape_id":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.ShapeID = nil
			} else {
				if z.ShapeID == nil {
					z.ShapeID = new(string)
				}
				*z.ShapeID, err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "headsign":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Headsign = nil
			} else {
				if z.Headsign == nil {
					z.Headsign = new(string)
				}
				*z.Headsign, err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "short_name":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.ShortName = nil
			} else {
				if z.ShortName == nil {
					z.ShortName = new(string)
				}
				*z.ShortName, err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "bikes_allowed":
			z.BikesAllowed, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "wheelchair_accessible":
			z.WheelchairAccessible, err = dc.ReadInt()
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
func (z *Trip) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 10
	// write "id"
	err = en.Append(0x8a, 0xa2, 0x69, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteString(z.ID)
	if err != nil {
		return
	}
	// write "route_id"
	err = en.Append(0xa8, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x5f, 0x69, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteString(z.RouteID)
	if err != nil {
		return
	}
	// write "service_id"
	err = en.Append(0xaa, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteString(z.ServiceID)
	if err != nil {
		return
	}
	// write "direction_id"
	err = en.Append(0xac, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64)
	if err != nil {
		return err
	}
	if z.DirectionID == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteInt(*z.DirectionID)
		if err != nil {
			return
		}
	}
	// write "block_id"
	err = en.Append(0xa8, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x69, 0x64)
	if err != nil {
		return err
	}
	if z.BlockID == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteString(*z.BlockID)
		if err != nil {
			return
		}
	}
	// write "shape_id"
	err = en.Append(0xa8, 0x73, 0x68, 0x61, 0x70, 0x65, 0x5f, 0x69, 0x64)
	if err != nil {
		return err
	}
	if z.ShapeID == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteString(*z.ShapeID)
		if err != nil {
			return
		}
	}
	// write "headsign"
	err = en.Append(0xa8, 0x68, 0x65, 0x61, 0x64, 0x73, 0x69, 0x67, 0x6e)
	if err != nil {
		return err
	}
	if z.Headsign == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteString(*z.Headsign)
		if err != nil {
			return
		}
	}
	// write "short_name"
	err = en.Append(0xaa, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	if z.ShortName == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteString(*z.ShortName)
		if err != nil {
			return
		}
	}
	// write "bikes_allowed"
	err = en.Append(0xad, 0x62, 0x69, 0x6b, 0x65, 0x73, 0x5f, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.BikesAllowed)
	if err != nil {
		return
	}
	// write "wheelchair_accessible"
	err = en.Append(0xb5, 0x77, 0x68, 0x65, 0x65, 0x6c, 0x63, 0x68, 0x61, 0x69, 0x72, 0x5f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x69, 0x62, 0x6c, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.WheelchairAccessible)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Trip) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 10
	// string "id"
	o = append(o, 0x8a, 0xa2, 0x69, 0x64)
	o = msgp.AppendString(o, z.ID)
	// string "route_id"
	o = append(o, 0xa8, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x5f, 0x69, 0x64)
	o = msgp.AppendString(o, z.RouteID)
	// string "service_id"
	o = append(o, 0xaa, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64)
	o = msgp.AppendString(o, z.ServiceID)
	// string "direction_id"
	o = append(o, 0xac, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64)
	if z.DirectionID == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendInt(o, *z.DirectionID)
	}
	// string "block_id"
	o = append(o, 0xa8, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x69, 0x64)
	if z.BlockID == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendString(o, *z.BlockID)
	}
	// string "shape_id"
	o = append(o, 0xa8, 0x73, 0x68, 0x61, 0x70, 0x65, 0x5f, 0x69, 0x64)
	if z.ShapeID == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendString(o, *z.ShapeID)
	}
	// string "headsign"
	o = append(o, 0xa8, 0x68, 0x65, 0x61, 0x64, 0x73, 0x69, 0x67, 0x6e)
	if z.Headsign == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendString(o, *z.Headsign)
	}
	// string "short_name"
	o = append(o, 0xaa, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65)
	if z.ShortName == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendString(o, *z.ShortName)
	}
	// string "bikes_allowed"
	o = append(o, 0xad, 0x62, 0x69, 0x6b, 0x65, 0x73, 0x5f, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64)
	o = msgp.AppendInt(o, z.BikesAllowed)
	// string "wheelchair_accessible"
	o = append(o, 0xb5, 0x77, 0x68, 0x65, 0x65, 0x6c, 0x63, 0x68, 0x61, 0x69, 0x72, 0x5f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x69, 0x62, 0x6c, 0x65)
	o = msgp.AppendInt(o, z.WheelchairAccessible)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Trip) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "id":
			z.ID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "route_id":
			z.RouteID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "service_id":
			z.ServiceID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "direction_id":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.DirectionID = nil
			} else {
				if z.DirectionID == nil {
					z.DirectionID = new(int)
				}
				*z.DirectionID, bts, err = msgp.ReadIntBytes(bts)
				if err != nil {
					return
				}
			}
		case "block_id":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.BlockID = nil
			} else {
				if z.BlockID == nil {
					z.BlockID = new(string)
				}
				*z.BlockID, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "shape_id":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.ShapeID = nil
			} else {
				if z.ShapeID == nil {
					z.ShapeID = new(string)
				}
				*z.ShapeID, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "headsign":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Headsign = nil
			} else {
				if z.Headsign == nil {
					z.Headsign = new(string)
				}
				*z.Headsign, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "short_name":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.ShortName = nil
			} else {
				if z.ShortName == nil {
					z.ShortName = new(string)
				}
				*z.ShortName, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "bikes_allowed":
			z.BikesAllowed, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "wheelchair_accessible":
			z.WheelchairAccessible, bts, err = msgp.ReadIntBytes(bts)
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
func (z *Trip) Msgsize() (s int) {
	s = 1 + 3 + msgp.StringPrefixSize + len(z.ID) + 9 + msgp.StringPrefixSize + len(z.RouteID) + 11 + msgp.StringPrefixSize + len(z.ServiceID) + 13
	if z.DirectionID == nil {
		s += msgp.NilSize
	} else {
		s += msgp.IntSize
	}
	s += 9
	if z.BlockID == nil {
		s += msgp.NilSize
	} else {
		s += msgp.StringPrefixSize + len(*z.BlockID)
	}
	s += 9
	if z.ShapeID == nil {
		s += msgp.NilSize
	} else {
		s += msgp.StringPrefixSize + len(*z.ShapeID)
	}
	s += 9
	if z.Headsign == nil {
		s += msgp.NilSize
	} else {
		s += msgp.StringPrefixSize + len(*z.Headsign)
	}
	s += 11
	if z.ShortName == nil {
		s += msgp.NilSize
	} else {
		s += msgp.StringPrefixSize + len(*z.ShortName)
	}
	s += 14 + msgp.IntSize + 22 + msgp.IntSize
	return
}

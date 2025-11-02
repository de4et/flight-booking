package sro

import (
	"time"

	"flight-booking/internal/model/sro"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// SROToProto converts a domain SRO to protobuf SRO
func SROToProto(s *sro.SRO) *SRO {
	if s == nil {
		return nil
	}

	return &SRO{
		Segments:     segmentsToProto(s.Segments),
		Passengers:   passengersToProto(s.Passengers),
		Class:        string(s.Class),
		Type:         string(s.Type),
		ChannelToken: channelTokenToProto(s.ChannelToken),
		Filters:      filtersToProto(s.Filters),
		Metadata:     metadataToProto(s.Metadata),
	}
}

// ProtoToSRO converts a protobuf SRO to domain SRO
func ProtoToSRO(p *SRO) *sro.SRO {
	if p == nil {
		return nil
	}

	return &sro.SRO{
		Segments:     protoToSegments(p.GetSegments()),
		Passengers:   protoToPassengers(p.GetPassengers()),
		Class:        sro.TravelClass(p.GetClass()),
		Type:         sro.RouteType(p.GetType()),
		ChannelToken: protoToChannelToken(p.GetChannelToken()),
		Filters:      protoToFilters(p.GetFilters()),
		Metadata:     protoToMetadata(p.GetMetadata()),
	}
}

// Helper conversion functions for SRO
func segmentsToProto(segments []sro.Segment) []*Segment {
	protoSegments := make([]*Segment, 0, len(segments))
	for _, segment := range segments {
		protoSegments = append(protoSegments, &Segment{
			From: segment.From,
			To:   segment.To,
			Date: timestamppb.New(segment.Date),
		})
	}
	return protoSegments
}

func protoToSegments(protoSegments []*Segment) []sro.Segment {
	segments := make([]sro.Segment, 0, len(protoSegments))
	for _, protoSegment := range protoSegments {
		if protoSegment != nil {
			var date time.Time
			if protoSegment.GetDate() != nil {
				date = protoSegment.GetDate().AsTime()
			}
			segments = append(segments, sro.Segment{
				From: protoSegment.GetFrom(),
				To:   protoSegment.GetTo(),
				Date: date,
			})
		}
	}
	return segments
}

func passengersToProto(p sro.Passengers) *Passengers {
	return &Passengers{
		Adt: int32(p.ADT),
		Chd: int32(p.CHD),
		Inf: int32(p.INF),
		Src: int32(p.SRC),
		Yth: int32(p.YTH),
		Ins: p.INS,
	}
}

func protoToPassengers(p *Passengers) sro.Passengers {
	if p == nil {
		return sro.Passengers{}
	}
	return sro.Passengers{
		ADT: int(p.GetAdt()),
		CHD: int(p.GetChd()),
		INF: int(p.GetInf()),
		SRC: int(p.GetSrc()),
		YTH: int(p.GetYth()),
		INS: p.GetIns(),
	}
}

func channelTokenToProto(ct sro.ChannelToken) *ChannelToken {
	return &ChannelToken{
		PartnerCode: ct.PartnerCode,
		SourceCode:  ct.SourceCode,
	}
}

func protoToChannelToken(ct *ChannelToken) sro.ChannelToken {
	if ct == nil {
		return sro.ChannelToken{}
	}
	return sro.ChannelToken{
		PartnerCode: ct.GetPartnerCode(),
		SourceCode:  ct.GetSourceCode(),
	}
}

func filtersToProto(f sro.Filters) *Filters {
	return &Filters{
		IsDirectOnly:    f.IsDirectOnly,
		MaxStops:        int32(f.MaxStops),
		WithBaggageOnly: f.WithBaggageOnly,
		Carriers:        f.Carriers,
		CarriersType:    string(f.CarriersType),
		GdsList:         f.GDSList,
		GdsListType:     string(f.GDSListType),
	}
}

func protoToFilters(f *Filters) sro.Filters {
	if f == nil {
		return sro.Filters{}
	}
	return sro.Filters{
		IsDirectOnly:    f.GetIsDirectOnly(),
		MaxStops:        int(f.GetMaxStops()),
		WithBaggageOnly: f.GetWithBaggageOnly(),
		Carriers:        f.GetCarriers(),
		CarriersType:    sro.ListType(f.GetCarriersType()),
		GDSList:         f.GetGdsList(),
		GDSListType:     sro.ListType(f.GetGdsListType()),
	}
}

func metadataToProto(m sro.Metadata) *Metadata {
	return &Metadata{
		IsTest:   m.IsTest,
		Currency: m.Currency,
		Language: m.Language,
		Timeout:  int32(m.Timeout),
	}
}

func protoToMetadata(m *Metadata) sro.Metadata {
	if m == nil {
		return sro.Metadata{}
	}
	return sro.Metadata{
		IsTest:   m.GetIsTest(),
		Currency: m.GetCurrency(),
		Language: m.GetLanguage(),
		Timeout:  int(m.GetTimeout()),
	}
}

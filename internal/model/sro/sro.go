package sro

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidSROToken = fmt.Errorf("invalid SRO token format")

type TravelClass string

const (
	TravelClassE TravelClass = "E"
	TravelClassB TravelClass = "B"
	TravelClassF TravelClass = "F"
	TravelClassW TravelClass = "W"
)

type RouteType string

const (
	RouteTypeOW RouteType = "OW"
	RouteTypeRT RouteType = "RT"
	RouteTypeCX RouteType = "CX"
)

type Segment struct {
	From string    `json:"from"`
	To   string    `json:"to"`
	Date time.Time `json:"date"`
}

type Passengers struct {
	ADT int  `json:"adt"`
	CHD int  `json:"chd"`
	INF int  `json:"inf"`
	SRC int  `json:"src"`
	YTH int  `json:"yth"`
	INS bool `json:"ins"`
}

type ChannelToken struct {
	PartnerCode string `json:"partnerCode"`
	SourceCode  string `json:"sourceCode"`
}

type ListType string

const (
	ListTypeInclude ListType = "include"
	ListTypeExclude ListType = "exclude"
)

type Filters struct {
	IsDirectOnly    bool     `json:"isDirectOnly"`
	MaxStops        int      `json:"maxStops"`
	WithBaggageOnly bool     `json:"withBaggageOnly"`
	Carriers        []string `json:"carriers"`
	CarriersType    ListType `json:"carriersType"`
	GDSList         []string `json:"gdsList"`
	GDSListType     ListType `json:"gdsListType"`
}

type Metadata struct {
	IsTest   bool   `json:"isTest"`
	Currency string `json:"currency"`
	Language string `json:"language"`
	Timeout  int    `json:"timeout"`
}

type SRO struct {
	Segments     []Segment    `json:"segments"`
	Passengers   Passengers   `json:"passengers"`
	Class        TravelClass  `json:"class"`
	Type         RouteType    `json:"type"`
	ChannelToken ChannelToken `json:"channelToken"`
	Filters      Filters      `json:"filters"`
	Metadata     Metadata     `json:"metadata"`
}

func (sro *SRO) IsOW() bool {
	return sro.Type == RouteTypeOW
}

func (sro *SRO) IsRT() bool {
	return sro.Type == RouteTypeRT
}

func (sro *SRO) IsCX() bool {
	return sro.Type == RouteTypeCX
}

func (sro *SRO) GetToken() string {
	var sb strings.Builder

	// **Формат SRO Token**:
	// ```
	// AKV4           - Партнёр (4 символа)
	// 0000           - Источник (4 символа)
	// OW             - Тип (OW/RT/CX)
	// E              - Класс (E/B/F/W)
	// 1              - Количество взрослых
	// 0              - Количество детей
	// 0              - Количество младенцев
	// 0              - Пожилые
	// 0              - Молодёжь
	// 0              - Тестовый режим (0/1)
	// 0              - Только прямые (0/1)
	// 0              - Только с багажом (0/1)
	// 9              - Макс. пересадок (0-9)
	// 0              - Страхование (0/1)
	// MOWLED20241015 - Сегменты (From+To+Date)
	// _I_S7.FS       - Фильтры авиакомпаний (опционально)
	// _GE_1.2        - Фильтры GDS (опционально)
	// _RUB           - Валюта (опционально)
	// _RU            - Язык (опционально)
	// ```

	sb.WriteString(sro.ChannelToken.PartnerCode)
	sb.WriteString(sro.ChannelToken.SourceCode)
	sb.WriteString(string(sro.Type))
	sb.WriteString(string(sro.Class))
	sb.WriteString(fmt.Sprint(sro.Passengers.ADT))
	sb.WriteString(fmt.Sprint(sro.Passengers.CHD))
	sb.WriteString(fmt.Sprint(sro.Passengers.INF))
	sb.WriteString(fmt.Sprint(sro.Passengers.SRC))
	sb.WriteString(fmt.Sprint(sro.Passengers.YTH))
	sb.WriteString(fmt.Sprint(toInt(sro.Metadata.IsTest)))
	sb.WriteString(fmt.Sprint(toInt(sro.Filters.IsDirectOnly)))
	sb.WriteString(fmt.Sprint(toInt(sro.Filters.WithBaggageOnly)))
	sb.WriteString(fmt.Sprint(sro.Filters.MaxStops))
	sb.WriteString(fmt.Sprint(toInt(sro.Passengers.INS)))

	for i := range sro.Segments {
		seg := sro.Segments[i]
		sb.WriteString(fmt.Sprintf("%s%s%s", seg.From, seg.To, seg.Date.Format("20060102")))
	}

	if len(sro.Filters.Carriers) > 0 {
		if sro.Filters.CarriersType == ListTypeInclude {
			sb.WriteString("_I_")
		} else {
			sb.WriteString("_E_")
		}
		sb.WriteString(strings.Join(sro.Filters.Carriers, "."))
	}

	if len(sro.Filters.GDSList) > 0 {
		if sro.Filters.GDSListType == ListTypeInclude {
			sb.WriteString("_GI_")
		} else {
			sb.WriteString("_GE_")
		}
		sb.WriteString(strings.Join(sro.Filters.GDSList, "."))
	}

	if sro.Metadata.Currency != "" {
		sb.WriteString("_" + sro.Metadata.Currency)
	}

	if sro.Metadata.Language != "" {
		sb.WriteString("_" + sro.Metadata.Language)
	}

	return sb.String()
}

func FromToken(token string) (*SRO, error) {
	// min length
	if len(token) < 34 {
		return nil, ErrInvalidSROToken
	}

	if !validatePartnerCode(token[0:4]) {
		// change to specific error?
		return nil, ErrInvalidSROToken
	}
	partnerCode := token[0:4]

	if !validateSourceCode(token[4:8]) {
		return nil, ErrInvalidSROToken
	}
	sourceCode := token[4:8]

	if !validateRouteType(token[8:10]) {
		return nil, ErrInvalidSROToken
	}
	routeType := RouteType(token[8:10])

	if !validateClass(token[10:11]) {
		return nil, ErrInvalidSROToken
	}
	class := TravelClass(token[10:11])

	if !validatePassengerAmount(token[11:12]) {
		return nil, ErrInvalidSROToken
	}
	adt := atoi(token[11:12])

	if !validatePassengerAmount(token[12:13]) {
		return nil, ErrInvalidSROToken
	}
	chd := atoi(token[12:13])

	if !validatePassengerAmount(token[13:14]) {
		return nil, ErrInvalidSROToken
	}
	inf := atoi(token[13:14])

	if !validatePassengerAmount(token[14:15]) {
		return nil, ErrInvalidSROToken
	}
	src := atoi(token[14:15])

	if !validatePassengerAmount(token[15:16]) {
		return nil, ErrInvalidSROToken
	}
	yth := atoi(token[15:16])

	isTest := token[16:17] == "1"
	isDirect := token[17:18] == "1"
	withBaggage := token[18:19] == "1"

	if !validateMaxStops(token[19:20]) {
		return nil, ErrInvalidSROToken
	}
	maxStops := atoi(token[19:20])

	insurance := token[20:21] == "1"

	remaining := token[21:]

	segRegexp := regexp.MustCompile(`([A-Z]{3})([A-Z]{3})(\d{8})`)
	matches := segRegexp.FindAllStringSubmatchIndex(remaining, -1)
	if len(matches) == 0 {
		return nil, ErrInvalidSROToken
	}

	var segments []Segment
	var lastIndex int
	for _, m := range matches {
		from := remaining[m[2]:m[3]]
		to := remaining[m[4]:m[5]]
		dateStr := remaining[m[6]:m[7]]

		date, err := time.Parse("20060102", dateStr)
		if err != nil {
			return nil, ErrInvalidSROToken
		}

		segments = append(segments, Segment{
			From: from,
			To:   to,
			Date: date,
		})
		lastIndex = m[1]
	}

	tail := remaining[lastIndex:]

	var filters Filters
	var metadata Metadata

	for _, part := range strings.Split(tail, "_") {
		if part == "" {
			continue
		}
		switch {
		case strings.HasPrefix(part, "I_"):
			filters.CarriersType = ListTypeInclude
			filters.Carriers = strings.Split(strings.TrimPrefix(part, "I_"), ".")
		case strings.HasPrefix(part, "E_"):
			filters.CarriersType = ListTypeExclude
			filters.Carriers = strings.Split(strings.TrimPrefix(part, "E_"), ".")
		case strings.HasPrefix(part, "G"):
			if strings.HasPrefix(part, "GE_") {
				filters.GDSListType = ListTypeExclude
				filters.GDSList = strings.Split(strings.TrimPrefix(part, "GE_"), ".")
			} else if strings.HasPrefix(part, "GI_") {
				filters.GDSListType = ListTypeInclude
				filters.GDSList = strings.Split(strings.TrimPrefix(part, "GI_"), ".")
			}
		case len(part) == 3:
			metadata.Currency = part
		case len(part) == 2:
			metadata.Language = part
		}
	}

	sro := &SRO{
		Segments: segments,
		Passengers: Passengers{
			ADT: adt,
			CHD: chd,
			INF: inf,
			SRC: src,
			YTH: yth,
			INS: insurance,
		},
		Class: class,
		Type:  routeType,
		ChannelToken: ChannelToken{
			PartnerCode: partnerCode,
			SourceCode:  sourceCode,
		},
		Filters: Filters{
			IsDirectOnly:    isDirect,
			WithBaggageOnly: withBaggage,
			MaxStops:        maxStops,
			Carriers:        filters.Carriers,
			CarriersType:    filters.CarriersType,
			GDSList:         filters.GDSList,
			GDSListType:     filters.GDSListType,
		},
		Metadata: Metadata{
			IsTest:   isTest,
			Currency: metadata.Currency,
			Language: metadata.Language,
		},
	}
	return sro, nil
}

func validatePartnerCode(s string) bool {
	return len(s) == 4 && strings.ToUpper(s) == s
}

func validateSourceCode(s string) bool {
	return validatePartnerCode(s)
}

func validateRouteType(s string) bool {
	return len(s) == 2 && in(s, "OW", "RT", "CX")
}

func validateClass(s string) bool {
	return len(s) == 1 && in(s, "E", "B", "F", "W")
}

func validatePassengerAmount(s string) bool {
	return isDigit(s)
}

func validateMaxStops(s string) bool {
	return isDigit(s)
}

func isDigit(s string) bool {
	return len(s) == 1 && s[0] >= '0' && s[0] <= '9'
}

func in[T comparable](v T, s ...T) bool {
	return slices.Contains(s, v)
}

func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func toInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// TODO - `isInnerFlight()` - определение внутреннего/международного рейса
// - `toArray()` - экспорт параметров для передачи провайдерам

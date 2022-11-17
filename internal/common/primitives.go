package common

import (
	"dbgo/pkg/randoms"
	"math/rand"
	"time"
)

type PrimitiveDataType string

const (
	PrimitiveBoolean   PrimitiveDataType = "boolean"
	PrimitiveCharacter PrimitiveDataType = "character"
	PrimitiveDate      PrimitiveDataType = "date"
	PrimitiveDecimal   PrimitiveDataType = "decimal"
	PrimitiveNumeric   PrimitiveDataType = "numeric"
	PrimitiveUUID      PrimitiveDataType = "uuid"
)

func (p *PrimitiveDataType) Quoted() bool {
	switch *p {
	case PrimitiveCharacter:
		return true
	case PrimitiveDate:
		return true
	case PrimitiveUUID:
		return true
	}

	return false
}

func (p *PrimitiveDataType) DefaultValue(maxSize int) any {
	switch *p {
	case PrimitiveBoolean:
		return randoms.RandomBoolean()
	case PrimitiveCharacter:
		return randoms.RandomStringAlpha(maxSize)
	case PrimitiveDate:
		return randoms.RandomDateUTC(time.Now().Year(), 1970)
	case PrimitiveDecimal:
		return rand.Float32()
	case PrimitiveNumeric:
		if maxSize == 32 {
			return rand.Int31()
		} else if maxSize == 64 {
			return rand.Int63()
		} else {
			return rand.Intn(32767) // assume 16-bit signed
		}
	case PrimitiveUUID:
		return randoms.RandomUUID()
	}

	return nil
}

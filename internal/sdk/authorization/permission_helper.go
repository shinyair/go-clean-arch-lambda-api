package authorization

import (
	"github.com/pkg/errors"
)

const (
	SuperAdminBit     uint64 = 0b1111111111111111111111111111111111111111111111111111111111111111 // max of uint64
	NonePermissionBit uint64 = 0b0

	bitLength int = 8
)

var ErrExceedMaxPermissionBitLenght error = errors.New("exceed permission bit length")

// HasAuthority
//
//	@param userBit
//	@param allowdPermissionBit
//	@return bool
//	@return error
func HasAuthority(userBit uint64, allowdPermissionBit []uint64) (bool, error) {
	if userBit == SuperAdminBit {
		return true, nil
	}
	for _, bit := range allowdPermissionBit {
		ok := (userBit & bit) >= bit
		if ok {
			return true, nil
		}
	}
	return false, nil
}

// GenerateRootBit
//
//	@return uint64
func GenerateRootBit() uint64 {
	return SuperAdminBit
}

// GenerateGrantedBit
//
//	@param grantedIndices
//	@return uint64
//	@return error
func GenerateGrantedBit(grantedIndices []int) (uint64, error) {
	var result uint64 = 0b0
	for _, i := range grantedIndices {
		if i >= bitLength {
			return NonePermissionBit, ErrExceedMaxPermissionBitLenght
		}
		var sub uint64 = 0b1
		sub <<= (bitLength - i)
		result |= sub
	}
	return result, nil
}

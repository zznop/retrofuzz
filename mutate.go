package main;

import(
	"unsafe"
    "math/bits"
)

// MutatorCallback is a function pointer type for mutator routines
type MutatorCallback func(context *RFContext, increment bool)

// MutatorInfo is used as a context for the mutator routines
type MutatorInfo struct {
	index int
	bitPosition int
	currTestCase int
	maxTestCases int
	callback MutatorCallback
}

// ByteFlip flips all 8 bits in a byte
func ByteFlip(context *RFContext, increment bool) {
	index := context.mutatorInfo.index

	b := (*byte)(unsafe.Pointer(&context.inputData[index]))
	*b ^= 0xff

	if increment == false {
		return
	}

	context.mutatorInfo.index++
}

// NibbleFlip flips a 4 bits in a byte starting at bit 0 or bit 4
func NibbleFlip(context *RFContext, increment bool) {
	index := context.mutatorInfo.index
	bitPosition := context.mutatorInfo.bitPosition

	b := (*byte)(unsafe.Pointer(&context.inputData[index]))
	if bitPosition == 0 {
		*b = (*b ^ 0xf0) | (*b & 0x0f)
	} else {
		*b = (*b & 0xf0) | (*b ^ 0x0f)
	}

	if increment == false {
		return
	}

	if (context.mutatorInfo.bitPosition == 0) {
		context.mutatorInfo.bitPosition = 4
	} else {
		context.mutatorInfo.bitPosition = 0
		context.mutatorInfo.index++
	}
}

// BitFlip flips a single bit in a byte
func BitFlip(context *RFContext, increment bool) {
	index := context.mutatorInfo.index
	bitPosition := context.mutatorInfo.bitPosition

	b := (*byte)(unsafe.Pointer(&context.inputData[index]))
	*b ^= 1 << bitPosition

	if increment == false {
		return
	}

	context.mutatorInfo.bitPosition++
	if context.mutatorInfo.bitPosition == 8 {
		context.mutatorInfo.bitPosition = 0
		context.mutatorInfo.index++
	}

	return
}

// InstructionGeneration generates random 68k instructions by simply overwriting the beginning of the code region with
// values ranging from 0x0-0xffff
func InstructionGeneration(context *RFContext, increment bool) {
    var rom *GenesisRom
    var programStart uint32

    rom = (*GenesisRom)(unsafe.Pointer(&context.inputData[0]))
    programStart = bits.ReverseBytes32(rom.offsetProgramStart)
    *(*uint32)(unsafe.Pointer(&context.inputData[programStart])) =
        (uint32)((0xff80 << 16) | context.mutatorInfo.currTestCase)
    context.mutatorInfo.bitPosition = 0
    context.mutatorInfo.index = (int)(programStart)
}

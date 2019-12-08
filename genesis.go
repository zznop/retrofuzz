package main

import(
    "unsafe"
)

// GenesisIVT represents the structure of a Genesis ROM interrupt vector table
type GenesisIVT struct {
    vectPtrBusError uint32
    vectPtrAddressError uint32
    vectPtrIllegalInstruction uint32
    vectPtrDivisionByZero uint32
    vectPtrChkException  uint32
    vectPtrTrapVException uint32
    vectPtrPrivilegeViolation uint32
    vectPtrTraceException uint32
    vectPtrLineAEmulator uint32
    vectPtrLineFEmulator uint32
    vectUnused00 uint32
    vectUnused01 uint32
    vectUnused02 uint32
    vectUnused03 uint32
    vectUnused04 uint32
    vectUnused05 uint32
    vectUnused06 uint32
    vectUnused07 uint32
    vectUnused08 uint32
    vectUnused09 uint32
    vectUnused10 uint32
    vectUnused11 uint32
    vectPtrSpuriousException uint32
    vectPtrIrqL1 uint32
    vectPtrIrqL2 uint32
    vectPtrIrqL3 uint32
    vectPtrIrqL4 uint32
    vectPtrIrqL5 uint32
    vectPtrIrqL6 uint32
    vectPtrIrqL7 uint32
    vectPtrTrap00 uint32
    vectPtrTrap01 uint32
    vectPtrTrap02 uint32
    vectPtrTrap03 uint32
    vectPtrTrap04 uint32
    vectPtrTrap05 uint32
    vectPtrTrap06 uint32
    vectPtrTrap07 uint32
    vectPtrTrap08 uint32
    vectPtrTrap09 uint32
    vectPtrTrap10 uint32
    vectPtrTrap11 uint32
    vectPtrTrap12 uint32
    vectPtrTrap13 uint32
    vectPtrTrap14 uint32
    vectPtrTrap15 uint32
    vectUnused12 uint32
    vectUnused13 uint32
    vectUnused14 uint32
    vectUnused15 uint32
    vectUnused16 uint32
    vectUnused17 uint32
    vectUnused18 uint32
    vectUnused19 uint32
    vectUnused20 uint32
    vectUnused21 uint32
    vectUnused22 uint32
    vectUnused23 uint32
    vectUnused24 uint32
    vectUnused25 uint32
    vectUnused26 uint32
    vectUnused27 uint32
}

// GenesisInfo represents the structure of a Genesis ROM information block
type GenesisInfo struct {
    consoleName [16]byte
    copyright [16]byte
    domesticName [48]byte
    internationalName [48]byte
    serialRevision [14]byte
    checksum uint16
    ioSupport [16]byte
    romStart uint32
    romEnd uint32
    ramStart uint32
    ramEnd uint32
    sramInfo [12]byte
    notes [52]byte
    region [16]byte
}

// GenesisHeader represents the structure of Genesis ROM header fields beginning at byte 0
type GenesisRom struct {
    offsetInitialStack uint32
    offsetProgramStart uint32
    ivt GenesisIVT
    info GenesisInfo
}

// GenesisFixupChecksum calculates the ROM checksum and writes them back to the ROM.
func GenesisFixupChecksum(context *RFContext) {
    const checksumOffset int = 0x18E
    const startOffset int = 0x200
    var checksum uint16
    var currPtr *uint16

    checksum = 0
    for i := startOffset; i < len(context.inputData); i += 2 {
        currPtr = (*uint16)(unsafe.Pointer(&context.inputData[i]))
        checksum += (*currPtr >> 8) | (*currPtr << 8)
    }

    *(*uint16)(unsafe.Pointer(&context.inputData[checksumOffset])) = (checksum << 8) | (checksum >> 8)
}

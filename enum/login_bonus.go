package enum

const (
	LoginBonusReceiveStatusReceived   = 0x00000001
	LoginBonusReceiveStatusReceiving  = 0x00000002
	LoginBonusReceiveStatusUnreceived = 0x00000003
)

const (
	LoginBonusContentGradeNormal = 0x00000001
	LoginBonusContentGradeRare   = 0x00000002
)

const (
	LoginBonusTypeNormal   = 0x00000001
	LoginBonusTypeBeginner = 0x00000002
	LoginBonusTypeEvent2d  = 0x00000003
	LoginBonusTypeEvent3d  = 0x00000004
	LoginBonusTypeBirthday = 0x00000005
	LoginBonusTypeComeback = 0x00000006
)

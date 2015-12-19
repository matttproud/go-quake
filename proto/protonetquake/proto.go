// Package protonetquake defines the netquake protocol, the protocol used in the original game.
package protonetquake

const (
	SVCBad          byte = 0
	SVCNop               = 1
	SVCDisconnect        = 2
	SVCUpdateStat        = 3
	SVCVersion           = 4
	SVCSetView           = 5
	SVCSound             = 6
	SVCTime              = 7
	SVCPrint             = 8
	SVCStuffText         = 9
	SVCSetAngle          = 10
	SVCServerInfo        = 11
	SVCLightStyle        = 12
	SVCUpdateName        = 13
	SVCUpdateFrags       = 14
	SVCClientData        = 15
	SVCStopSound         = 16
	SVCUpdateColors      = 17
	SVCParticle          = 18
	SVCDamage            = 19
	SVCSpawnStatic       = 20
	/* SpawnBinary = 21 */
	SVCSpawnBaseline    = 22
	SVCTempEntity       = 23
	SVCSetPause         = 24
	SVCSignOnNum        = 25
	SVCCenterPrint      = 26
	SVCKilledMonster    = 27
	SVCFoundSecret      = 28
	SVCSpawnStaticSound = 29
	SVCIntermission     = 30
	SVCFinale           = 31
	SVCCDTrack          = 32
	SVCSellScreen       = 33
	SVCCutScene         = 34

	CLCBad           byte = 0
	CLCNop                = 1
	CLCDisconnect         = 2
	CLCMove               = 3
	CLCStringCommand      = 4

	TESpike        byte = 0
	TESuperSpike        = 1
	TEGunshot           = 2
	TEExplosion         = 3
	TETarExplosion      = 4
	TELightning1        = 5
	TELightning2        = 6
	TEWizSpike          = 7
	TEKnigthSpike       = 8
	TELightning3        = 9
	TELavaSplash        = 10
	TETeleport          = 11
	TEExplosion2        = 12
	TEBeam              = 13
	TEImplosion         = 14
	TERailTrail         = 15
)

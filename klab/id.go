package klab

// named klab because these are conventions used in the game itself

// this file mainly extract info from ids
// all functions verified in masterdata.db using SQL

func DefaultSuitMasterIDFromMemberMasterID(memberMasterID int) int {
	// SELECT * FROM m_member_init WHERE suit_m_id != member_m_id * 10000 + 100001001; -> 0
	return 100001001 + memberMasterID*10000
}

func MemberMasterIDFromCardMasterID(cardMasterID int) int {
	// SELECT * FROM m_card WHERE (id / 10000) % 1000 != member_m_id; -> 0
	return (cardMasterID / 10000) % 1000
}

func CardRarityFromCardMasterID(cardMasterID int) int {
	// SELECT * FROM m_card WHERE (id / 100) % 100 != card_rarity_type; -> 0
	return (cardMasterID / 100) % 100
}

func MemberMasterIDFromSuitMasterID(suitMasterID int) int {
	// verified in masterdata.db, all SQL return empty
	if suitMasterID <= 100109 { // special aqours outfit
		// SELECT * FROM m_suit WHERE id <= 100109 AND id % 1000 != member_m_id; -> 0
		return suitMasterID % 1000
	} else if suitMasterID < 100011001 {
		// SELECT * FROM m_suit WHERE id > 100109 AND id < 100011001 AND (id / 100) % 1000 != member_m_id; -> 0
		return (suitMasterID / 100) % 1000
	} else {
		// SELECT * FROM m_suit WHERE id >= 100011001 AND (id / 10000) % 1000 != member_m_id; -> 0
		return (suitMasterID / 10000) % 1000
	}
}

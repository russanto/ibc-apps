package keeper_test

import (
	"github.com/cosmos/ibc-apps/modules/rate-limiting/v10/keeper"
	"github.com/cosmos/ibc-apps/modules/rate-limiting/v10/types"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (s *KeeperTestSuite) TestGetChannelValue() {
	supply := sdkmath.NewInt(100)

	// Mint coins to increase the supply, which will increase the channel value
	err := s.App.BankKeeper.MintCoins(s.Ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(denom, supply)))
	s.Require().NoError(err)

	expected := supply
	actual := s.App.RatelimitKeeper.GetChannelValue(s.Ctx, denom)
	s.Require().Equal(expected, actual)
}

// Adds a rate limit object to the store in preparation for the check rate limit tests
func (s *KeeperTestSuite) SetupCheckRateLimitAndUpdateFlowTest() {
	channelValue := sdkmath.NewInt(100)
	maxPercentSend := sdkmath.NewInt(10)
	maxPercentRecv := sdkmath.NewInt(10)

	s.App.RatelimitKeeper.SetRateLimit(s.Ctx, types.RateLimit{
		Path: &types.Path{
			Denom:             denom,
			ChannelOrClientId: channelId,
		},
		Quota: &types.Quota{
			MaxPercentSend: maxPercentSend,
			MaxPercentRecv: maxPercentRecv,
			DurationHours:  1,
		},
		Flow: &types.Flow{
			Inflow:       sdkmath.ZeroInt(),
			Outflow:      sdkmath.ZeroInt(),
			ChannelValue: channelValue,
		},
	})

	s.App.RatelimitKeeper.RemoveDenomFromBlacklist(s.Ctx, denom)
	s.App.RatelimitKeeper.RemoveWhitelistedAddressPair(s.Ctx, sender, receiver)
}

// Helper function to check the rate limit across a series of transfers
func (s *KeeperTestSuite) processCheckRateLimitAndUpdateFlowTestCase(tc checkRateLimitTestCase) {
	s.SetupCheckRateLimitAndUpdateFlowTest()

	expectedInflow := sdkmath.NewInt(0)
	expectedOutflow := sdkmath.NewInt(0)
	for i, action := range tc.actions {
		if action.addToBlacklist {
			s.App.RatelimitKeeper.AddDenomToBlacklist(s.Ctx, denom)
			continue
		} else if action.removeFromBlacklist {
			s.App.RatelimitKeeper.RemoveDenomFromBlacklist(s.Ctx, denom)
			continue
		}

		if action.addToWhitelist {
			s.App.RatelimitKeeper.SetWhitelistedAddressPair(s.Ctx, types.WhitelistedAddressPair{
				Sender:   sender,
				Receiver: receiver,
			})
			continue
		} else if action.removeFromWhitelist {
			s.App.RatelimitKeeper.RemoveWhitelistedAddressPair(s.Ctx, sender, receiver)
			continue
		}

		amount := sdkmath.NewInt(action.amount)
		packetInfo := keeper.RateLimitedPacketInfo{
			ChannelID: channelId,
			Denom:     denom,
			Amount:    amount,
			Sender:    sender,
			Receiver:  receiver,
		}
		updatedFlow, err := s.App.RatelimitKeeper.CheckRateLimitAndUpdateFlow(s.Ctx, action.direction, packetInfo)

		// Each action optionally errors or skips a flow update
		if action.expectedError != "" {
			s.Require().ErrorContains(err, action.expectedError, tc.name+" - action: #%d - error", i)
		} else {
			s.Require().NoError(err, tc.name+" - action: #%d - no error", i)

			expectedUpdateFlow := !action.skipFlowUpdate
			s.Require().Equal(expectedUpdateFlow, updatedFlow, tc.name+" - action: #%d - updated flow", i)

			if expectedUpdateFlow {
				if action.direction == types.PACKET_RECV {
					expectedInflow = expectedInflow.Add(amount)
				} else {
					expectedOutflow = expectedOutflow.Add(amount)
				}
			}
		}

		// Confirm flow is updated properly (or left as is if the theshold was exceeded)
		rateLimit, found := s.App.RatelimitKeeper.GetRateLimit(s.Ctx, denom, channelId)
		s.Require().True(found)
		s.Require().Equal(expectedInflow.Int64(), rateLimit.Flow.Inflow.Int64(), tc.name+" - action: #%d - inflow", i)
		s.Require().Equal(expectedOutflow.Int64(), rateLimit.Flow.Outflow.Int64(), tc.name+" - action: #%d - outflow", i)
	}
}

func (s *KeeperTestSuite) TestCheckRateLimitAndUpdateFlow_UnidirectionalFlow() {
	testCases := []checkRateLimitTestCase{
		{
			name: "send_under_threshold",
			actions: []action{
				{direction: types.PACKET_SEND, amount: 5},
				{direction: types.PACKET_SEND, amount: 5},
			},
		},
		{
			name: "send_over_threshold",
			actions: []action{
				{direction: types.PACKET_SEND, amount: 5},
				{
					direction: types.PACKET_SEND, amount: 6,
					expectedError: "Outflow exceeds quota",
				},
			},
		},
		{
			name: "recv_under_threshold",
			actions: []action{
				{direction: types.PACKET_RECV, amount: 5},
				{direction: types.PACKET_RECV, amount: 5},
			},
		},
		{
			name: "recv_over_threshold",
			actions: []action{
				{direction: types.PACKET_RECV, amount: 5},
				{
					direction: types.PACKET_RECV, amount: 6,
					expectedError: "Inflow exceeds quota",
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.processCheckRateLimitAndUpdateFlowTestCase(tc)
		})
	}
}

func (s *KeeperTestSuite) TestCheckRateLimitAndUpdatedFlow_BidirectionalFlow() {
	testCases := []checkRateLimitTestCase{
		{
			name: "send_then_recv_under_threshold",
			actions: []action{
				{direction: types.PACKET_SEND, amount: 6},
				{direction: types.PACKET_RECV, amount: 6},
				{direction: types.PACKET_SEND, amount: 6},
				{direction: types.PACKET_RECV, amount: 6},
			},
		},
		{
			name: "recv_then_send_under_threshold",
			actions: []action{
				{direction: types.PACKET_RECV, amount: 6},
				{direction: types.PACKET_SEND, amount: 6},
				{direction: types.PACKET_RECV, amount: 6},
				{direction: types.PACKET_SEND, amount: 6},
			},
		},
		{
			name: "send_then_recv_over_inflow",
			actions: []action{
				{direction: types.PACKET_SEND, amount: 2}, //   -2, Net: -2
				{direction: types.PACKET_RECV, amount: 6}, //   +6, Net: +4
				{direction: types.PACKET_SEND, amount: 2}, //   -2, Net: +2
				{direction: types.PACKET_RECV, amount: 6}, //   +6, Net: +8
				{direction: types.PACKET_SEND, amount: 2}, //   -2, Net: +6
				{direction: types.PACKET_RECV, amount: 6, //    +6, Net: +12 (exceeds threshold)
					expectedError: "Inflow exceeds quota"},
			},
		},
		{
			name: "send_then_recv_over_outflow",
			actions: []action{
				{direction: types.PACKET_SEND, amount: 6}, //   -6, Net: -6
				{direction: types.PACKET_RECV, amount: 2}, //   +2, Net: -4
				{direction: types.PACKET_SEND, amount: 6}, //   -6, Net: -10
				{direction: types.PACKET_RECV, amount: 2}, //   +2, Net: -8
				{direction: types.PACKET_SEND, amount: 6, //    -6, Net: -14 (exceeds threshold)
					expectedError: "Outflow exceeds quota"},
			},
		},
		{
			name: "recv_then_send_over_inflow",
			actions: []action{
				{direction: types.PACKET_RECV, amount: 6}, //   +6, Net: +6
				{direction: types.PACKET_SEND, amount: 2}, //   -2, Net: +4
				{direction: types.PACKET_RECV, amount: 6}, //   +6, Net: +10
				{direction: types.PACKET_SEND, amount: 2}, //   -2, Net: +8
				{direction: types.PACKET_RECV, amount: 6, //    +6, Net: +14 (exceeds threshold)
					expectedError: "Inflow exceeds quota"},
			},
		},
		{
			name: "recv_then_send_over_outflow",
			actions: []action{
				{direction: types.PACKET_RECV, amount: 2}, //  +2, Net: +2
				{direction: types.PACKET_SEND, amount: 6}, //  -6, Net: -4
				{direction: types.PACKET_RECV, amount: 2}, //  +2, Net: -2
				{direction: types.PACKET_SEND, amount: 6}, //  -6, Net: -8
				{direction: types.PACKET_RECV, amount: 2}, //  +2, Net: -6
				{direction: types.PACKET_SEND, amount: 10, //  +6, Net: -12 (exceeds threshold)
					expectedError: "Outflow exceeds quota"},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.processCheckRateLimitAndUpdateFlowTestCase(tc)
		})
	}
}

func (s *KeeperTestSuite) TestCheckRateLimitAndUpdatedFlow_DenomBlacklist() {
	testCases := []checkRateLimitTestCase{
		{
			name: "add_then_remove_from_blacklist", // should succeed
			actions: []action{
				{direction: types.PACKET_RECV, amount: 6},
				{direction: types.PACKET_SEND, amount: 6},
				{addToBlacklist: true},
				{removeFromBlacklist: true},
				{direction: types.PACKET_RECV, amount: 6},
				{direction: types.PACKET_SEND, amount: 6},
			},
		},
		{
			name: "send_recv_blacklist_send",
			actions: []action{
				{direction: types.PACKET_SEND, amount: 6},
				{direction: types.PACKET_RECV, amount: 6},
				{addToBlacklist: true},
				{
					direction: types.PACKET_SEND, amount: 6,
					expectedError: types.ErrDenomIsBlacklisted.Error(),
				},
			},
		},
		{
			name: "send_recv_blacklist_recv",
			actions: []action{
				{direction: types.PACKET_SEND, amount: 6},
				{direction: types.PACKET_RECV, amount: 6},
				{addToBlacklist: true},
				{
					direction: types.PACKET_RECV, amount: 6,
					expectedError: types.ErrDenomIsBlacklisted.Error(),
				},
			},
		},
		{
			name: "recv_send_blacklist_send",
			actions: []action{
				{direction: types.PACKET_RECV, amount: 6},
				{direction: types.PACKET_SEND, amount: 6},
				{addToBlacklist: true},
				{
					direction: types.PACKET_SEND, amount: 6,
					expectedError: types.ErrDenomIsBlacklisted.Error(),
				},
			},
		},
		{
			name: "recv_send_blacklist_recv",
			actions: []action{
				{direction: types.PACKET_RECV, amount: 6},
				{direction: types.PACKET_SEND, amount: 6},
				{addToBlacklist: true},
				{
					direction: types.PACKET_RECV, amount: 6,
					expectedError: types.ErrDenomIsBlacklisted.Error(),
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.processCheckRateLimitAndUpdateFlowTestCase(tc)
		})
	}
}

func (s *KeeperTestSuite) TestCheckRateLimitAndUpdatedFlow_AddressWhitelist() {
	testCases := []checkRateLimitTestCase{
		{
			name: "send_whitelist_send", // should succeed
			actions: []action{
				{direction: types.PACKET_SEND, amount: 6},
				{addToWhitelist: true},
				{direction: types.PACKET_SEND, amount: 6, skipFlowUpdate: true},
			},
		},
		{
			name: "recv_whitelist_recv", // should succeed
			actions: []action{
				{direction: types.PACKET_RECV, amount: 6},
				{addToWhitelist: true},
				{direction: types.PACKET_RECV, amount: 6, skipFlowUpdate: true},
			},
		},
		{
			name: "send_send_whitelist_send", // should succeed
			actions: []action{
				{direction: types.PACKET_SEND, amount: 6},
				{direction: types.PACKET_SEND, amount: 6, expectedError: "Outflow exceeds quota"}, // fails
				{addToWhitelist: true},
				{direction: types.PACKET_SEND, amount: 6, skipFlowUpdate: true}, // succeeds
			},
		},
		{
			name: "recv_recv_whitelist_recv", // should succeed
			actions: []action{
				{direction: types.PACKET_RECV, amount: 6},
				{direction: types.PACKET_RECV, amount: 6, expectedError: "Inflow exceeds quota"}, // fails
				{addToWhitelist: true},
				{direction: types.PACKET_RECV, amount: 6, skipFlowUpdate: true}, // succeeds
			},
		},
		{
			name: "send_recv_send_whitelist_send", // should succeed
			actions: []action{
				{direction: types.PACKET_SEND, amount: 6},
				{direction: types.PACKET_RECV, amount: 6},
				{direction: types.PACKET_SEND, amount: 6},
				{addToWhitelist: true},
				{direction: types.PACKET_SEND, amount: 6, skipFlowUpdate: true},
			},
		},
		{
			name: "recv_send_recv_whitelist_recv", // should succeed
			actions: []action{
				{direction: types.PACKET_RECV, amount: 6},
				{direction: types.PACKET_SEND, amount: 6},
				{direction: types.PACKET_RECV, amount: 6},
				{addToWhitelist: true},
				{direction: types.PACKET_RECV, amount: 6, skipFlowUpdate: true},
			},
		},
		{
			name: "add_then_remove_whitelist_recv",
			actions: []action{
				{direction: types.PACKET_RECV, amount: 6},
				{addToWhitelist: true},
				{removeFromWhitelist: true},
				{direction: types.PACKET_RECV, amount: 6, expectedError: "Inflow exceeds quota"},
			},
		},
		{
			name: "add_then_remove_whitelist_send",
			actions: []action{
				{direction: types.PACKET_SEND, amount: 6},
				{addToWhitelist: true},
				{removeFromWhitelist: true},
				{direction: types.PACKET_SEND, amount: 6, expectedError: "Outflow exceeds quota"},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.processCheckRateLimitAndUpdateFlowTestCase(tc)
		})
	}
}

func (s *KeeperTestSuite) TestUndoSendPacket() {
	// Helper function to check the rate limit outflow amount
	checkOutflow := func(channelId, denom string, expectedAmount sdkmath.Int) {
		rateLimit, found := s.App.RatelimitKeeper.GetRateLimit(s.Ctx, denom, channelId)
		s.Require().True(found, "rate limit should have been found")
		s.Require().Equal(expectedAmount.Int64(), rateLimit.Flow.Outflow.Int64(),
			"outflow - channel: %s, denom: %s", channelId, denom)
	}

	// Create two rate limits
	initialOutflow := sdkmath.NewInt(100)
	packetSendAmount := sdkmath.NewInt(10)
	rateLimit1 := types.RateLimit{
		Path: &types.Path{Denom: denom, ChannelOrClientId: channelId},
		Flow: &types.Flow{Outflow: initialOutflow},
	}
	rateLimit2 := types.RateLimit{
		Path: &types.Path{Denom: "different-denom", ChannelOrClientId: "different-channel"},
		Flow: &types.Flow{Outflow: initialOutflow},
	}
	s.App.RatelimitKeeper.SetRateLimit(s.Ctx, rateLimit1)
	s.App.RatelimitKeeper.SetRateLimit(s.Ctx, rateLimit2)

	// Store a pending packet sequence number of 2 for the first rate limit
	s.App.RatelimitKeeper.SetPendingSendPacket(s.Ctx, channelId, 2)

	// Undo a send of 10 from the first rate limit, with sequence 1
	// If should NOT modify the outflow since sequence 1 was not sent in the current quota
	err := s.App.RatelimitKeeper.UndoSendPacket(s.Ctx, channelId, 1, denom, packetSendAmount)
	s.Require().NoError(err, "no error expected when undoing send packet sequence 1")

	checkOutflow(channelId, denom, initialOutflow)

	// Now undo a send from the same rate limit with sequence 2
	// If should decrement the outflow since 2 is in the current quota
	err = s.App.RatelimitKeeper.UndoSendPacket(s.Ctx, channelId, 2, denom, packetSendAmount)
	s.Require().NoError(err, "no error expected when undoing send packet sequence 2")

	checkOutflow(channelId, denom, initialOutflow.Sub(packetSendAmount))

	// Confirm the outflow of the second rate limit has not been touched
	checkOutflow("different-channel", "different-denom", initialOutflow)

	// Confirm sequence number was removed
	found := s.App.RatelimitKeeper.CheckPacketSentDuringCurrentQuota(s.Ctx, channelId, 2)
	s.Require().False(found, "packet sequence number should have been removed")
}

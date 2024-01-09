package config

import "github.com/pipe01/flydigi-linux/pkg/utils"

var ledUnits = [][][][]byte{
	{
		{
			{50, 50, 0},
			{0, 100, 0},
			{0, 50, 50},
			{0, 0, 100},
			{50, 0, 50},
			{100, 0, 0},
			{0, 116, 255},
			{0, 116, 255},
			{0, 116, 255},
			{0, 116, 255},
		},
		{
			{24, 76, 0},
			{0, 76, 24},
			{0, 0, 100},
			{24, 0, 76},
			{76, 0, 24},
			{76, 24, 0},
			{0, 116, 255},
			{0, 116, 255},
			{0, 116, 255},
			{0, 116, 255},
		},
		{
			{0, 100, 0},
			{0, 50, 50},
			{0, 50, 50},
			{50, 0, 50},
			{100, 0, 0},
			{50, 50, 0},
			{0, 116, 255},
			{0, 116, 255},
			{0, 116, 255},
			{0, 116, 255},
		},
		{
			{100, 0, 0},
			{50, 50, 0},
			{0, 100, 0},
			{0, 50, 50},
			{0, 0, 100},
			{50, 0, 50},
			{0, 116, 255},
			{0, 116, 255},
			{0, 116, 255},
			{0, 116, 255},
		},
		{
			{76, 24, 0},
			{24, 76, 0},
			{0, 76, 24},
			{0, 76, 24},
			{24, 0, 76},
			{76, 0, 24},
			{0, 116, 255},
			{0, 116, 255},
			{0, 116, 255},
			{0, 116, 255},
		},
		{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
	},
	{
		{
			{100, 0, 0},
			{0, 0, 0},
			{0, 100, 0},
			{0, 0, 0},
			{0, 50, 50},
			{0, 0, 0},
			{0, 0, 100},
			{0, 0, 0},
			{50, 0, 50},
			{0, 0, 0},
		},
		{
			{80, 0, 0},
			{0, 0, 0},
			{0, 80, 0},
			{0, 0, 0},
			{0, 40, 40},
			{0, 0, 0},
			{0, 0, 80},
			{0, 0, 0},
			{40, 0, 40},
			{0, 0, 0},
		},
		{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{80, 0, 0},
			{0, 0, 0},
			{0, 80, 0},
			{0, 0, 0},
			{0, 40, 40},
			{0, 0, 0},
			{0, 0, 80},
			{0, 0, 0},
			{40, 0, 40},
			{0, 0, 0},
		},
		{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
	},
	{
		{
			{100, 0, 0},
			{0, 0, 100},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{100, 0, 0},
			{0, 0, 100},
			{0, 0, 100},
			{0, 0, 100},
			{0, 0, 100},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{100, 0, 0},
			{0, 0, 100},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{100, 0, 0},
			{0, 0, 100},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{100, 0, 0},
			{0, 0, 100},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
	},
	{
		{
			{0, 0, 100},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{0, 0, 80},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{0, 0, 80},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
	},
	{
		{
			{0, 0, 100},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{0, 0, 100},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{0, 0, 100},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{0, 0, 100},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{0, 0, 100},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
	},
	{
		{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
	},
}

func getLedGroupList(modeId, ledNum int) []*LedGroup {
	groups := make([]*LedGroup, 0)

	for i := 0; i < 16; i++ {
		if i < ledNum {
			for j := 0; j < 10; j++ {
				groups = append(groups, &LedGroup{
					Units: []*LedUnit{
						{
							R: ledUnits[modeId][i][j][0],
							G: ledUnits[modeId][i][j][1],
							B: ledUnits[modeId][i][j][2],
						},
					},
				})
			}
		} else {
			groups = append(groups, &LedGroup{
				Units: utils.RepeatFunc(func() *LedUnit {
					return &LedUnit{0, 0, 0}
				}, 10),
			})
		}
	}

	return groups
}

var GameHandleName = map[int32]string{
	19: "apex2",
	20: "f1",
	21: "f1",
	22: "f1p",
	23: "f1",
	24: "k1",
	25: "fp1",
	26: "K1SF",
	29: "k1",
	27: "fp1",
	30: "fp1Fate",
	28: "f3",
	80: "f3p",
	81: "f3pip",
	82: "fp2",
	83: "fp2ip",
	84: "k2",
}

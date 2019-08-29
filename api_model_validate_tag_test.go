package libgen

import (
	"github.com/gen-iot/std"
	"testing"
)

func Test_RegisterDevicesRequest_Validate(t *testing.T) {
	req := &RegisterDevicesRequest{
		Devices: []*Device{
			{
				ModelInfo: &ModelInfo{
					Package: "com.pujie88.com",
					Name:    "test",
				},
				Id:       "1",
				Name:     "2",
				Room:     "3",
				Status:   nil,
				MetaData: nil,
			},
		},
	}
	err := std.ValidateStruct(req)
	std.Assert(err == nil, "1 validate should be success")
	req = &RegisterDevicesRequest{
		Devices: nil,
	}
	err = std.ValidateStruct(req)
	std.Assert(err != nil, "2 validate should be failed")
	req = &RegisterDevicesRequest{
		Devices: []*Device{},
	}
	err = std.ValidateStruct(req)
	std.Assert(err != nil, "3 validate should be failed")
	req = &RegisterDevicesRequest{
		Devices: []*Device{
			{
				ModelInfo: &ModelInfo{
					Package: "com.pujie88.iot",
					Name:    "test",
				},
				Id:       "",
				Name:     "",
				Room:     "",
				Status:   nil,
				MetaData: nil,
			},
		},
	}
	err = std.ValidateStruct(req)
	std.Assert(err != nil, "4 validate should be failed")
}

func Test_ListDevicesByRoomRequest_Validate(t *testing.T) {
	req := &ListDevicesByRoomRequest{
		Rooms:          []string{"101"},
		Filter:         nil,
		CategoryFilter: nil,
	}
	err := std.ValidateStruct(req)
	std.Assert(err == nil, "1 validate should be success")
	req = &ListDevicesByRoomRequest{
		Rooms:          []string{},
		Filter:         nil,
		CategoryFilter: nil,
	}
	err = std.ValidateStruct(req)
	std.Assert(err != nil, "2 validate should be failed")

	req = &ListDevicesByRoomRequest{
		Rooms:  []string{"101"},
		Filter: []*ModelInfo{},
	}
	err = std.ValidateStruct(req)
	std.Assert(err != nil, "3 validate should be failed")

	req = &ListDevicesByRoomRequest{
		Rooms: []string{"101"},
		Filter: []*ModelInfo{
			{
				Package: "",
				Name:    "",
			},
		},
	}
	err = std.ValidateStruct(req)
	std.Assert(err != nil, "4 validate should be failed")
	req = &ListDevicesByRoomRequest{
		Rooms: []string{"101"},
		Filter: []*ModelInfo{
			{
				Package: "com.pujie88.iot",
				Name:    "test",
			},
		},
	}
	err = std.ValidateStruct(req)
	std.Assert(err == nil, "5 validate should be success")
}

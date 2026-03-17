package mapper

import (
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "shipment-service/gen/proto/shipment"
	"shipment-service/internal/domain/shipment"
)

func ShipmentToProto(s *shipment.Shipment) *pb.Shipment {
	return &pb.Shipment{
		Id:              s.ID,
		ReferenceNumber: s.ReferenceNumber,
		Origin:          s.Origin,
		Destination:     s.Destination,
		CurrentStatus:   StatusToProto(s.CurrentStatus),
		Driver:          DriverInfoToProto(s.Driver),
		ShipmentAmount:  s.ShipmentAmount,
		DriverRevenue:   s.DriverRevenue,
		CreatedAt:       timestamppb.New(s.CreatedAt),
		UpdatedAt:       timestamppb.New(s.UpdatedAt),
	}
}

func StatusEventToProto(e shipment.StatusEvent) *pb.StatusEvent {
	return &pb.StatusEvent{
		Id:         e.ID,
		ShipmentId: e.ShipmentID,
		Status:     StatusToProto(e.Status),
		Note:       e.Note,
		OccurredAt: timestamppb.New(e.OccurredAt),
	}
}

func DriverInfoToProto(d shipment.DriverInfo) *pb.DriverInfo {
	return &pb.DriverInfo{
		DriverId: d.DriverID,
		UnitId:   d.UnitID,
	}
}

func StatusToProto(s shipment.Status) pb.ShipmentStatus {
	switch s {
	case shipment.StatusPending:
		return pb.ShipmentStatus_SHIPMENT_STATUS_PENDING
	case shipment.StatusPickedUp:
		return pb.ShipmentStatus_SHIPMENT_STATUS_PICKED_UP
	case shipment.StatusInTransit:
		return pb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT
	case shipment.StatusDelivered:
		return pb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED
	case shipment.StatusCancelled:
		return pb.ShipmentStatus_SHIPMENT_STATUS_CANCELLED
	default:
		return pb.ShipmentStatus_SHIPMENT_STATUS_UNSPECIFIED
	}
}

func ProtoToStatus(s pb.ShipmentStatus) (shipment.Status, error) {
	switch s {
	case pb.ShipmentStatus_SHIPMENT_STATUS_PENDING:
		return shipment.StatusPending, nil
	case pb.ShipmentStatus_SHIPMENT_STATUS_PICKED_UP:
		return shipment.StatusPickedUp, nil
	case pb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT:
		return shipment.StatusInTransit, nil
	case pb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED:
		return shipment.StatusDelivered, nil
	case pb.ShipmentStatus_SHIPMENT_STATUS_CANCELLED:
		return shipment.StatusCancelled, nil
	default:
		return "", fmt.Errorf("unknown status: %v", s)
	}
}

func ProtoToDriverInfo(d *pb.DriverInfo) shipment.DriverInfo {
	if d == nil {
		return shipment.DriverInfo{}
	}
	return shipment.DriverInfo{
		DriverID: d.DriverId,
		UnitID:   d.UnitId,
	}
}

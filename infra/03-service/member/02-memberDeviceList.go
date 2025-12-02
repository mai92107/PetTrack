package member

import "context"

func (s *MemberServiceImpl) MemberDeviceList(ctx context.Context, memberId int64) ([]string, error) {
	deviceIds := []string{}
	deviceIds, err := s.memberDeviceRepo.GetMemberDeviceList(ctx, memberId)
	if err != nil {
		return deviceIds, err
	}
	return deviceIds, nil
}

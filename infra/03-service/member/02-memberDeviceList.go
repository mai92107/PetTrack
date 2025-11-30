package member

func (s *MemberServiceImpl) MemberDeviceList(memberId int64) ([]string, error) {
	deviceIds := []string{}
	deviceIds, err := s.memberDeviceRepo.GetMemberDeviceList(memberId)
	if err != nil {
		return deviceIds, err
	}
	return deviceIds, nil
}

package member

import "PetTrack/infra/03-service/member"


var memberService member.MemberServiceImpl

func InitMemberHandler(service member.MemberServiceImpl) {
	memberService = service
}
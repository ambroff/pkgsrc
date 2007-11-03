# $NetBSD: buildlink3.mk,v 1.4 2007/11/03 15:59:18 tnn Exp $

BUILDLINK_DEPTH:=	${BUILDLINK_DEPTH}+
PIDGIN_BUILDLINK3_MK:=	${PIDGIN_BUILDLINK3_MK}+

.if ${BUILDLINK_DEPTH} == "+"
BUILDLINK_DEPENDS+=	pidgin
.endif

BUILDLINK_PACKAGES:=	${BUILDLINK_PACKAGES:Npidgin}
BUILDLINK_PACKAGES+=	pidgin
BUILDLINK_ORDER:=	${BUILDLINK_ORDER} ${BUILDLINK_DEPTH}pidgin

.if ${PIDGIN_BUILDLINK3_MK} == "+"
BUILDLINK_API_DEPENDS.pidgin+=	pidgin>=2.2.2
BUILDLINK_PKGSRCDIR.pidgin?=	../../chat/pidgin
.endif	# PIDGIN_BUILDLINK3_MK

.include "../../chat/libpurple/buildlink3.mk"

BUILDLINK_DEPTH:=	${BUILDLINK_DEPTH:S/+$//}

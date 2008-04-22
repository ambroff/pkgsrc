# $NetBSD: buildlink3.mk,v 1.11 2008/04/22 17:16:54 tnn Exp $

BUILDLINK_DEPTH:=		${BUILDLINK_DEPTH}+
LIBEVENT_BUILDLINK3_MK:=	${LIBEVENT_BUILDLINK3_MK}+

.if !empty(BUILDLINK_DEPTH:M+)
BUILDLINK_DEPENDS+=	libevent
.endif

BUILDLINK_PACKAGES:=	${BUILDLINK_PACKAGES:Nlibevent}
BUILDLINK_PACKAGES+=	libevent
BUILDLINK_ORDER:=	${BUILDLINK_ORDER} ${BUILDLINK_DEPTH}libevent

.if !empty(LIBEVENT_BUILDLINK3_MK:M+)
BUILDLINK_API_DEPENDS.libevent+=libevent>=0.6
BUILDLINK_ABI_DEPENDS.libevent+=libevent-1.4.3* # exact match -- see Makefile
BUILDLINK_PKGSRCDIR.libevent?=	../../wip/libevent
.endif	# LIBEVENT_BUILDLINK3_MK

BUILDLINK_DEPTH:=		${BUILDLINK_DEPTH:S/+$//}

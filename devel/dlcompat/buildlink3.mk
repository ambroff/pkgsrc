# $NetBSD: buildlink3.mk,v 1.7 2004/11/26 09:39:30 jlam Exp $

BUILDLINK_DEPTH:=		${BUILDLINK_DEPTH}+
DLCOMPAT_BUILDLINK3_MK:=	${DLCOMPAT_BUILDLINK3_MK}+

.if !empty(BUILDLINK_DEPTH:M+)
BUILDLINK_DEPENDS+=	dlcompat
.endif

BUILDLINK_PACKAGES:=	${BUILDLINK_PACKAGES:Ndlcompat}
BUILDLINK_PACKAGES+=	dlcompat

.if !empty(DLCOMPAT_BUILDLINK3_MK:M+)
BUILDLINK_DEPENDS.dlcompat+=	dlcompat>=20030629
BUILDLINK_PKGSRCDIR.dlcompat?=	../../devel/dlcompat
.endif  # DLCOMPAT_BUILDLINK3_MK

BUILDLINK_DEPTH:=	${BUILDLINK_DEPTH:S/+$//}

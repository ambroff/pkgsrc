# $NetBSD: buildlink3.mk,v 1.2 2007/12/21 00:19:43 jdolecek Exp $
#
BUILDLINK_DEPTH:=	${BUILDLINK_DEPTH}+
QCA_BUILDLINK3_MK:=	${QCA_BUILDLINK3_MK}+

.if !empty(BUILDLINK_DEPTH:M+)
BUILDLINK_DEPENDS+=	qca2
.endif

BUILDLINK_PACKAGES:=	${BUILDLINK_PACKAGES:Nqca2}
BUILDLINK_PACKAGES+=	qca2
BUILDLINK_ORDER:=	${BUILDLINK_ORDER} ${BUILDLINK_DEPTH}qca2

.if !empty(QCA_BUILDLINK3_MK:M+)
BUILDLINK_API_DEPENDS.qca2+=	qca2>=2.0.0
BUILDLINK_ABI_DEPENDS.qca2?=	qca2>=2.0.0
BUILDLINK_PKGSRCDIR.qca2?=	../../security/qca2
.endif	# QCA2_BUILDLINK3_MK

.include "../../x11/qt4-libs/buildlink3.mk"
.include "../../x11/qt4-tools/buildlink3.mk"

BUILDLINK_DEPTH:=	${BUILDLINK_DEPTH:S/+$//}

# $NetBSD: buildlink3.mk,v 1.12 2013/05/09 07:39:07 adam Exp $

BUILDLINK_TREE+=	mono-addins

.if !defined(MONO_ADDINS_BUILDLINK3_MK)
MONO_ADDINS_BUILDLINK3_MK:=

BUILDLINK_API_DEPENDS.mono-addins+=	mono-addins>=0.3
BUILDLINK_ABI_DEPENDS.mono-addins+=	mono-addins>=1.0nb2
BUILDLINK_PKGSRCDIR.mono-addins?=	../../devel/mono-addins
BUILDLINK_CONTENTS_FILTER.mono-addins+=	${EGREP} '^lib/'

.include "../../lang/mono/buildlink3.mk"
.include "../../x11/gtk-sharp/buildlink3.mk"
.endif # MONO_ADDINS_BUILDLINK3_MK

BUILDLINK_TREE+=	-mono-addins

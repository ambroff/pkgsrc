# $NetBSD: buildlink3.mk,v 1.3 2013/05/09 07:39:17 adam Exp $

BUILDLINK_TREE+=	libdmapsharing3

.if !defined(LIBDMAPSHARING3_BUILDLINK3_MK)
LIBDMAPSHARING3_BUILDLINK3_MK:=

BUILDLINK_API_DEPENDS.libdmapsharing3+=	libdmapsharing3>=2.9.15
BUILDLINK_ABI_DEPENDS.libdmapsharing3?=	libdmapsharing3>=2.9.15nb2
BUILDLINK_PKGSRCDIR.libdmapsharing3?=	../../net/libdmapsharing3

.include "../../net/libsoup24/buildlink3.mk"
.endif	# LIBDMAPSHARING3_BUILDLINK3_MK

BUILDLINK_TREE+=	-libdmapsharing3

# $NetBSD: buildlink3.mk,v 1.6 2004/02/16 21:01:17 jlam Exp $

BUILDLINK_DEPTH:=		${BUILDLINK_DEPTH}+
PKGCONFIG_BUILDLINK3_MK:=	${PKGCONFIG_BUILDLINK3_MK}+

.if !empty(BUILDLINK_DEPTH:M+)
BUILDLINK_DEPENDS+=	pkgconfig
.endif

.if !empty(PKGCONFIG_BUILDLINK3_MK:M+)
BUILDLINK_PACKAGES+=		pkgconfig
BUILDLINK_DEPENDS.pkgconfig+=	pkgconfig>=0.15.0
BUILDLINK_PKGSRCDIR.pkgconfig?=	../../devel/pkgconfig
BUILDLINK_DEPMETHOD.pkgconfig?=	build

PKG_CONFIG_LIBDIR?=	${BUILDLINK_DIR}/lib/pkgconfig
CONFIGURE_ENV+=		PKG_CONFIG=${BUILDLINK_PREFIX.pkgconfig}/bin/pkg-config
CONFIGURE_ENV+=		PKG_CONFIG_LIBDIR=${PKG_CONFIG_LIBDIR:Q}
MAKE_ENV+=		PKG_CONFIG=${BUILDLINK_PREFIX.pkgconfig}/bin/pkg-config
MAKE_ENV+=		PKG_CONFIG_LIBDIR=${PKG_CONFIG_LIBDIR:Q}
.endif	# PKGCONFIG_BUILDLINK3_MK

BUILDLINK_DEPTH:=	${BUILDLINK_DEPTH:S/+$//}

# $NetBSD: buildlink3.mk,v 1.35 2011/01/13 13:36:12 wiz Exp $

BUILDLINK_TREE+=	gimp

.if !defined(GIMP_BUILDLINK3_MK)
GIMP_BUILDLINK3_MK:=

BUILDLINK_API_DEPENDS.gimp+=	gimp>=2.4.0
BUILDLINK_ABI_DEPENDS.gimp+=	gimp>=2.6.11nb3
BUILDLINK_PKGSRCDIR.gimp?=	../../graphics/gimp

pkgbase := gimp
.include "../../mk/pkg-build-options.mk"

.if !empty(PKG_BUILD_OPTIONS.gimp:Maalib)
.include "../../graphics/aalib/buildlink3.mk"
.endif

.if !empty(PKG_BUILD_OPTIONS.gimp:Mmng)
.include "../../graphics/mng/buildlink3.mk"
.endif

.if !empty(PKG_BUILD_OPTIONS.gimp:Msvg)
.include "../../graphics/librsvg/buildlink3.mk"
.endif

.include "../../devel/gettext-lib/buildlink3.mk"
.include "../../fonts/fontconfig/buildlink3.mk"
.include "../../graphics/lcms/buildlink3.mk"
.include "../../graphics/libart/buildlink3.mk"
.include "../../graphics/png/buildlink3.mk"
.include "../../graphics/libwmf/buildlink3.mk"
.include "../../x11/gtk2/buildlink3.mk"
.endif # GIMP_BUILDLINK3_MK

BUILDLINK_TREE+=	-gimp

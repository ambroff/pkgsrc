# $NetBSD: buildlink.mk,v 1.3 2001/06/23 19:26:57 jlam Exp $
#
# This Makefile fragment is included by packages that use mng.
#
# To use this Makefile fragment, simply:
#
# (1) Optionally define BUILDLINK_DEPENDS.mng to the dependency pattern
#     for the version of mng desired.
# (2) Include this Makefile fragment in the package Makefile,
# (3) Add ${BUILDLINK_DIR}/include to the front of the C preprocessor's header
#     search path, and
# (4) Add ${BUILDLINK_DIR}/lib to the front of the linker's library search
#     path.

.if !defined(MNG_BUILDLINK_MK)
MNG_BUILDLINK_MK=	# defined

BUILDLINK_DEPENDS.mng?=	mng>=1.0.0
DEPENDS+=		${BUILDLINK_DEPENDS.mng}:../../graphics/mng

BUILDLINK_PREFIX.mng=	${LOCALBASE}
BUILDLINK_FILES.mng=	include/libmng.h
BUILDLINK_FILES.mng+=	include/libmng_conf.h
BUILDLINK_FILES.mng+=	include/libmng_types.h
BUILDLINK_FILES.mng+=	lib/libmng.*

.include "../../devel/zlib/buildlink.mk"
.include "../../graphics/jpeg/buildlink.mk"
.include "../../graphics/lcms/buildlink.mk"

BUILDLINK_TARGETS.mng=	mng-buildlink
BUILDLINK_TARGETS+=	${BUILDLINK_TARGETS.mng}

pre-configure: ${BUILDLINK_TARGETS.mng}
mng-buildlink: _BUILDLINK_USE

.include "../../mk/bsd.buildlink.mk"

.endif	# MNG_BUILDLINK_MK

# $NetBSD: buildlink2.mk,v 1.3 2003/06/12 19:08:34 wiz Exp $

.if !defined(LAPACK_BUILDLINK2_MK)
LAPACK_BUILDLINK2_MK=	# defined

BUILDLINK_PACKAGES+=		lapack
BUILDLINK_DEPENDS.lapack?=	lapack>=20010201
BUILDLINK_PKGSRCDIR.lapack?=	../../math/lapack
BUILDLINK_DEPMETHOD.lapack?=	build

EVAL_PREFIX+=			BUILDLINK_PREFIX.lapack=lapack
BUILDLINK_PREFIX.lapack_DEFAULT=	${LOCALBASE}
BUILDLINK_FILES.lapack=		lib/liblapack.*

BUILDLINK_TARGETS+=		lapack-buildlink

lapack-buildlink: _BUILDLINK_USE

.endif	# LAPACK_BUILDLINK2_MK

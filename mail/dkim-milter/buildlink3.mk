# $NetBSD: buildlink3.mk,v 1.9 2023/10/24 22:09:42 wiz Exp $

BUILDLINK_TREE+=	dkim-milter

.if !defined(DKIM_MILTER_BUILDLINK3_MK)
DKIM_MILTER_BUILDLINK3_MK:=

BUILDLINK_API_DEPENDS.dkim-milter+=	dkim-milter>=2.8.2
BUILDLINK_ABI_DEPENDS.dkim-milter+=	dkim-milter>=2.8.3nb7
BUILDLINK_PKGSRCDIR.dkim-milter?=	../../mail/dkim-milter
BUILDLINK_DEPMETHOD.dkim-milter?=	build

.include "../../security/openssl/buildlink3.mk"
.include "../../mail/libmilter/buildlink3.mk"
.endif # DKIM_MILTER_BUILDLINK3_MK

BUILDLINK_TREE+=	-dkim-milter

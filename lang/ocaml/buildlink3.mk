# $NetBSD: buildlink3.mk,v 1.6 2004/04/22 09:18:35 tron Exp $

BUILDLINK_DEPTH:=	${BUILDLINK_DEPTH}+
OCAML_BUILDLINK3_MK:=	${OCAML_BUILDLINK3_MK}+

.if !empty(BUILDLINK_DEPTH:M+)
BUILDLINK_DEPENDS+=	ocaml
.endif

BUILDLINK_PACKAGES:=	${BUILDLINK_PACKAGES:Nocaml}
BUILDLINK_PACKAGES+=	ocaml

.if !empty(OCAML_BUILDLINK3_MK:M+)
BUILDLINK_DEPENDS.ocaml+=	ocaml>=3.06nb6
BUILDLINK_PKGSRCDIR.ocaml?=	../../lang/ocaml
BUILDLINK_DEPMETHOD.ocaml?=	build
.endif	# OCAML_BUILDLINK3_MK

BUILDLINK_DEPTH:=	${BUILDLINK_DEPTH:S/+$//}

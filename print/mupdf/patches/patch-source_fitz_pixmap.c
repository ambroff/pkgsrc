$NetBSD: patch-source_fitz_pixmap.c,v 1.1 2017/02/11 09:39:05 leot Exp $

Backport a fix from upstream for CVE-2017-5896:

bug 697515: Fix out of bounds read in fz_subsample_pixmap

Pointer arithmetic for final special case was going wrong.

--- source/fitz/pixmap.c.orig
+++ source/fitz/pixmap.c
@@ -1104,6 +1104,7 @@ fz_subsample_pixmap_ARM(unsigned char *ptr, int w, int h, int f, int factor,
 	"@STACK:r1,<9>,factor,n,fwd,back,back2,fwd2,divX,back4,fwd4,fwd3,divY,back5,divXY\n"
 	"ldr	r4, [r13,#4*22]		@ r4 = divXY			\n"
 	"ldr	r5, [r13,#4*11]		@ for (nn = n; nn > 0; n--) {	\n"
+	"ldr	r8, [r13,#4*17]		@ r8 = back4			\n"
 	"18:				@				\n"
 	"mov	r14,#0			@ r14= v = 0			\n"
 	"sub	r5, r5, r1, LSL #8	@ for (xx = x; xx > 0; x--) {	\n"
@@ -1120,7 +1121,7 @@ fz_subsample_pixmap_ARM(unsigned char *ptr, int w, int h, int f, int factor,
 	"mul	r14,r4, r14		@ r14= v *= divX		\n"
 	"mov	r14,r14,LSR #16		@ r14= v >>= 16			\n"
 	"strb	r14,[r9], #1		@ *d++ = r14			\n"
-	"sub	r0, r0, r8		@ s -= back2			\n"
+	"sub	r0, r0, r8		@ s -= back4			\n"
 	"subs	r5, r5, #1		@ n--				\n"
 	"bgt	18b			@ }				\n"
 	"21:				@				\n"
@@ -1249,6 +1250,7 @@ fz_subsample_pixmap(fz_context *ctx, fz_pixmap *tile, int factor)
 		x += f;
 		if (x > 0)
 		{
+			int back4 = x * n - 1;
 			div = x * y;
 			for (nn = n; nn > 0; nn--)
 			{
@@ -1263,7 +1265,7 @@ fz_subsample_pixmap(fz_context *ctx, fz_pixmap *tile, int factor)
 					s -= back5;
 				}
 				*d++ = v / div;
-				s -= back2;
+				s -= back4;
 			}
 		}
 	}

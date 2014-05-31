$NetBSD: patch-gcc_flags.h,v 1.1 2014/05/31 13:06:25 ryoon Exp $

--- gcc/flags.h.orig	2013-01-10 20:38:27.000000000 +0000
+++ gcc/flags.h
@@ -25,6 +25,11 @@ along with GCC; see the file COPYING3.  
 
 #if !defined(IN_LIBGCC2) && !defined(IN_TARGET_LIBS) && !defined(IN_RTS)
 
+/* Nonzero means warn about any function whose stack usage is larger than N
+   bytes.  The value N is `stack_larger_than_size'.  */
+extern int warn_stack_larger_than;
+extern HOST_WIDE_INT stack_larger_than_size;
+
 /* Names of debug_info_type, for error messages.  */
 extern const char *const debug_type_names[];
 

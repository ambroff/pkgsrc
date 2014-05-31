$NetBSD: patch-gcc_opts.c,v 1.1 2014/05/31 13:06:25 ryoon Exp $

--- gcc/opts.c.orig	2014-05-04 07:05:29.000000000 +0000
+++ gcc/opts.c
@@ -49,6 +49,9 @@ const char *const debug_type_names[] =
   ((strncmp (prefix, string, sizeof prefix - 1) == 0) \
    ? ((string += sizeof prefix - 1), 1) : 0)
 
+int warn_stack_larger_than;
+HOST_WIDE_INT stack_larger_than_size;
+
 void
 set_struct_debug_option (struct gcc_options *opts, location_t loc,
 			 const char *spec)
@@ -468,8 +471,6 @@ static const struct default_options defa
     { OPT_LEVELS_2_PLUS, OPT_fschedule_insns2, NULL, 1 },
 #endif
     { OPT_LEVELS_2_PLUS, OPT_fregmove, NULL, 1 },
-    { OPT_LEVELS_2_PLUS, OPT_fstrict_aliasing, NULL, 1 },
-    { OPT_LEVELS_2_PLUS, OPT_fstrict_overflow, NULL, 1 },
     { OPT_LEVELS_2_PLUS, OPT_freorder_blocks, NULL, 1 },
     { OPT_LEVELS_2_PLUS, OPT_freorder_functions, NULL, 1 },
     { OPT_LEVELS_2_PLUS, OPT_ftree_vrp, NULL, 1 },
@@ -488,6 +489,7 @@ static const struct default_options defa
     { OPT_LEVELS_2_PLUS, OPT_fhoist_adjacent_loads, NULL, 1 },
 
     /* -O3 optimizations.  */
+    { OPT_LEVELS_3_PLUS, OPT_fstrict_aliasing, NULL, 1 },
     { OPT_LEVELS_3_PLUS, OPT_ftree_loop_distribute_patterns, NULL, 1 },
     { OPT_LEVELS_3_PLUS, OPT_fpredictive_commoning, NULL, 1 },
     /* Inlining of functions reducing size is a good idea with -Os
@@ -701,6 +703,8 @@ finish_options (struct gcc_options *opts
 
   if (!opts->x_flag_opts_finished)
     {
+      if (opts->x_flag_pic || opts->x_profile_flag)
+        opts->x_flag_pie = 0;
       if (opts->x_flag_pie)
 	opts->x_flag_pic = opts->x_flag_pie;
       if (opts->x_flag_pic && !opts->x_flag_pie)
@@ -1437,6 +1441,11 @@ common_handle_option (struct gcc_options
       opts->x_warn_frame_larger_than = value != -1;
       break;
 
+    case OPT_Wstack_larger_than_:
+      stack_larger_than_size = value;
+      warn_stack_larger_than = stack_larger_than_size != -1;
+      break;
+
     case OPT_Wstack_usage_:
       opts->x_warn_stack_usage = value;
       opts->x_flag_stack_usage_info = value != -1;

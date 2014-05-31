$NetBSD: patch-gcc_c-family_c-format.c,v 1.1 2014/05/31 13:06:25 ryoon Exp $

--- gcc/c-family/c-format.c.orig	2013-01-10 20:38:27.000000000 +0000
+++ gcc/c-family/c-format.c
@@ -38,6 +38,7 @@ along with GCC; see the file COPYING3.  
    format_type_error.  Target-specific format types do not have
    matching enum values.  */
 enum format_type { printf_format_type, asm_fprintf_format_type,
+		   kprintf_format_type, syslog_format_type,
 		   gcc_diag_format_type, gcc_tdiag_format_type,
 		   gcc_cdiag_format_type,
 		   gcc_cxxdiag_format_type, gcc_gfc_format_type,
@@ -51,6 +52,7 @@ typedef struct function_format_info
   unsigned HOST_WIDE_INT first_arg_num;	/* number of first arg (zero for varargs) */
 } function_format_info;
 
+
 static bool decode_format_attr (tree, function_format_info *, int);
 static int decode_format_type (const char *);
 
@@ -417,6 +419,15 @@ static const format_length_info gcc_diag
   { NO_FMT, NO_FMT, 0 }
 };
 
+static const format_length_info kprintf_length_specs[] =
+{
+  { "h", FMT_LEN_h, STD_C89, NO_FMT, 0 },
+  { "l", FMT_LEN_l, STD_C89, "ll", FMT_LEN_ll, STD_C9L, 0 },
+  { "q", FMT_LEN_ll, STD_EXT, NO_FMT, 0 },
+  { "L", FMT_LEN_L, STD_C89, NO_FMT, 0 },
+  { NO_FMT, NO_FMT, 0 }
+};
+
 /* The custom diagnostics all accept the same length specifiers.  */
 #define gcc_tdiag_length_specs gcc_diag_length_specs
 #define gcc_cdiag_length_specs gcc_diag_length_specs
@@ -597,7 +608,6 @@ static const format_flag_pair strfmon_fl
   { 0, 0, 0, 0 }
 };
 
-
 static const format_char_info print_char_table[] =
 {
   /* C89 conversion specifiers.  */
@@ -641,6 +651,44 @@ static const format_char_info asm_fprint
   { NULL,  0, STD_C89, NOLENGTHS, NULL, NULL, NULL }
 };
 
+static const format_char_info kprint_char_table[] =
+{
+  /* C89 conversion specifiers.  */
+  { "di",  0, STD_C89, { T89_I,   BADLEN,  T89_S,   T89_L,   T9L_LL,  BADLEN,  T99_SST, BADLEN,  BADLEN, BADLEN,  BADLEN,  BADLEN  }, "-wp0 +'I", "i", NULL },
+  { "oxX", 0, STD_C89, { T89_UI,  BADLEN,  T89_US,  T89_UL,  T9L_ULL, BADLEN, T99_ST,  BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN }, "-wp0#",    "i", NULL },
+  { "u",   0, STD_C89, { T89_UI,  BADLEN,  T89_US,  T89_UL,  T9L_ULL, BADLEN, T99_ST,  BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN }, "-wp0'I",   "i", NULL },
+  { "c",   0, STD_C89, { T89_I,   BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN, BADLEN,  BADLEN,  BADLEN  }, "-w",       "", NULL },
+  { "s",   1, STD_C89, { T89_C,   BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN, BADLEN, BADLEN,  BADLEN,   BADLEN,  BADLEN  }, "-wp",      "cR", NULL },
+  { "p",   1, STD_C89, { T89_V,   BADLEN,  BADLEN,  T89_UL,  T9L_LL,  BADLEN,  BADLEN, BADLEN, BADLEN,  BADLEN,   BADLEN,  BADLEN  }, "-wp0",      "c", NULL },
+/* Kernel bitmap formatting */
+  { "b",   0, STD_C89, { T89_I,   BADLEN,  T89_S,  T89_L,  T9L_LL,  BADLEN,  T99_SST,  BADLEN,  BADLEN, BADLEN,   BADLEN,  BADLEN }, "",    "", kprint_char_table + 8   },
+  { NULL,  0, STD_C89, NOLENGTHS, NULL, NULL, NULL },
+/* Kernel bitmap formatting, second part - similar to "s" except for types[] */
+  { "b",   1, STD_C89, { T89_C,   BADLEN,  T89_C,  T89_C,  T89_C,  BADLEN,  T89_C,  BADLEN,  BADLEN, BADLEN, BADLEN, BADLEN }, NULL,		"cR", NULL   }
+};
+
+static const format_char_info syslog_char_table[] =
+{
+  /* C89 conversion specifiers.  */
+  { "di",  0, STD_C89, { T89_I,   T99_SC,  T89_S,   T89_L,   T9L_LL,  TEX_LL,  T99_SST, T99_PD,  T99_IM  }, "-wp0 +'I", "i", NULL  },
+  { "oxX", 0, STD_C89, { T89_UI,  T99_UC,  T89_US,  T89_UL,  T9L_ULL, TEX_ULL, T99_ST,  T99_UPD, T99_UIM }, "-wp0#",    "i", NULL  },
+  { "u",   0, STD_C89, { T89_UI,  T99_UC,  T89_US,  T89_UL,  T9L_ULL, TEX_ULL, T99_ST,  T99_UPD, T99_UIM }, "-wp0'I",   "i", NULL  },
+  { "fgG", 0, STD_C89, { T89_D,   BADLEN,  BADLEN,  T99_D,   BADLEN,  T89_LD,  BADLEN,  BADLEN,  BADLEN  }, "-wp0 +#'", "", NULL   },
+  { "eE",  0, STD_C89, { T89_D,   BADLEN,  BADLEN,  T99_D,   BADLEN,  T89_LD,  BADLEN,  BADLEN,  BADLEN  }, "-wp0 +#",  "", NULL   },
+  { "c",   0, STD_C89, { T89_I,   BADLEN,  BADLEN,  T94_WI,  BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN  }, "-w",       "", NULL   },
+  { "s",   1, STD_C89, { T89_C,   BADLEN,  BADLEN,  T94_W,   BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN  }, "-wp",      "cR", NULL },
+  { "p",   1, STD_C89, { T89_V,   BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN  }, "-w",       "c", NULL  },
+  { "n",   1, STD_C89, { T89_I,   T99_SC,  T89_S,   T89_L,   T9L_LL,  BADLEN,  T99_SST, T99_PD,  T99_IM  }, "",         "W", NULL  },
+  /* C99 conversion specifiers.  */
+  { "F",   0, STD_C99, { T99_D,   BADLEN,  BADLEN,  T99_D,   BADLEN,  T99_LD,  BADLEN,  BADLEN,  BADLEN  }, "-wp0 +#'", "", NULL   },
+  { "aA",  0, STD_C99, { T99_D,   BADLEN,  BADLEN,  T99_D,   BADLEN,  T99_LD,  BADLEN,  BADLEN,  BADLEN  }, "-wp0 +#",  "", NULL   },
+  /* X/Open conversion specifiers.  */
+  { "C",   0, STD_EXT, { TEX_WI,  BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN  }, "-w",       "", NULL   },
+  { "S",   1, STD_EXT, { TEX_W,   BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN  }, "-wp",      "R", NULL  },
+  { "m",   0, STD_C89, { T89_V,   BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN,  BADLEN  }, "-wp",      "", NULL   },
+  { NULL,  0, STD_C89, NOLENGTHS, NULL, NULL, NULL }
+};
+
 static const format_char_info gcc_diag_char_table[] =
 {
   /* C89 conversion specifiers.  */
@@ -817,6 +865,18 @@ static const format_kind_info format_typ
     'w', 0, 'p', 0, 'L', 0,
     NULL, NULL
   },
+  { "kprintf", kprintf_length_specs, kprint_char_table, " +#0-'I", NULL,
+    printf_flag_specs, printf_flag_pairs,
+    FMT_FLAG_ARG_CONVERT|FMT_FLAG_DOLLAR_MULTIPLE|FMT_FLAG_USE_DOLLAR|FMT_FLAG_EMPTY_PREC_OK,
+    'w', 0, 'p', 0, 'L', 0,
+    &integer_type_node, &integer_type_node
+  },
+  { "syslog", printf_length_specs, syslog_char_table, " +#0-'I", NULL,
+    printf_flag_specs, printf_flag_pairs,
+    FMT_FLAG_ARG_CONVERT|FMT_FLAG_DOLLAR_MULTIPLE|FMT_FLAG_USE_DOLLAR|FMT_FLAG_EMPTY_PREC_OK,
+    'w', 0, 'p', 0, 'L', 0,
+    &integer_type_node, &integer_type_node
+  },
   { "gcc_diag",   gcc_diag_length_specs,  gcc_diag_char_table, "q+#", NULL,
     gcc_diag_flag_specs, gcc_diag_flag_pairs,
     FMT_FLAG_ARG_CONVERT,

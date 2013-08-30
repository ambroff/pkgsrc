$NetBSD: patch-mailnews_base_search_src_nsMsgSearchTerm.cpp,v 1.1 2013/08/30 13:00:37 joerg Exp $

--- mailnews/base/search/src/nsMsgSearchTerm.cpp.orig	2013-08-29 23:46:49.000000000 +0000
+++ mailnews/base/search/src/nsMsgSearchTerm.cpp
@@ -197,7 +197,7 @@ nsresult NS_MsgGetStringForAttribute(int
     }
   }
   if (!found)
-    *string = '\0'; // don't leave the string uninitialized
+    *string = NULL; // don't leave the string uninitialized
 
   // we no longer return invalid attribute. If we cannot find the string in the table,
   // then it is an arbitrary header. Return success regardless if found or not

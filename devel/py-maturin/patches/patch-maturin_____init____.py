$NetBSD: patch-maturin_____init____.py,v 1.2 2023/05/07 08:08:22 adam Exp $

Look for correct command with ${PYVERSSUFFIX} appended.

--- maturin/__init__.py.orig	2023-05-07 07:48:01.000000000 +0000
+++ maturin/__init__.py
@@ -55,8 +55,9 @@ def _build_wheel(
 ) -> str:
     # PEP 517 specifies that only `sys.executable` points to the correct
     # python interpreter
+    py_vers = platform.python_version_tuple()
     command = [
-        "maturin",
+        "maturin-" + py_vers[0] + "." + py_vers[1],
         "pep517",
         "build-wheel",
         "-i",

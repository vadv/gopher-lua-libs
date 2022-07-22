// optparse.lua for gopher-lua
package optparse

const lua_optparse = "LS1bWwoKTFVBIE1PRFVMRQoKICBweXRob25pYy5vcHRwYXJzZSAtIEx1YS1iYXNlZCBwYXJ0aWFsIHJlaW1wbGVtZW50YXRpb24gb2YgUHl0aG9uJ3MKICAgICAgb3B0cGFyc2UgWzItM10gY29tbWFuZC1saW5lIHBhcnNpbmcgbW9kdWxlLgoKU1lOT1BTSVMKCiAgbG9jYWwgT3B0aW9uUGFyc2VyID0gcmVxdWlyZSAicHl0aG9uaWMub3B0cGFyc2UiIC4gT3B0aW9uUGFyc2VyCiAgbG9jYWwgb3B0ID0gT3B0aW9uUGFyc2Vye3VzYWdlPSIlcHJvZyBbb3B0aW9uc10gW2d6aXAtZmlsZS4uLl0iLAogICAgICAgICAgICAgICAgICAgICAgICAgICB2ZXJzaW9uPSJmb28gMS4yMyIsIGFkZF9oZWxwX29wdGlvbj1mYWxzZX0KICBvcHQuYWRkX29wdGlvbnsiLWgiLCAiLS1oZWxwIiwgYWN0aW9uPSJzdG9yZV90cnVlIiwgZGVzdD0iaGVscCIsCiAgICAgICAgICAgICAgICAgaGVscD0iZ2l2ZSB0aGlzIGhlbHAifQogIG9wdC5hZGRfb3B0aW9uewogICAgIi1mIiwgIi0tZm9yY2UiLCBkZXN0PSJmb3JjZSIsIGFjdGlvbj0ic3RvcmVfdHJ1ZSIsCiAgICBoZWxwPSJmb3JjZSBvdmVyd3JpdGUgb2Ygb3V0cHV0IGZpbGUifQoKICBsb2NhbCBvcHRpb25zLCBhcmdzID0gb3B0LnBhcnNlX2FyZ3MoKQoKICBpZiBvcHRpb25zLmhlbHAgdGhlbiBvcHQucHJpbnRfaGVscCgpOyBvcy5leGl0KDEpIGVuZAogIGlmIG9wdGlvbnMuZm9yY2UgdGhlbiBwcmludCAnZicgZW5kCiAgZm9yIF8sIG5hbWUgaW4gaXBhaXJzKGFyZ3MpIGRvIHByaW50KG5hbWUpIGVuZAoKREVTQ1JJUFRJT04KCiAgVGhpcyBsaWJyYXJ5IHByb3ZpZGVzIGEgY29tbWFuZC1saW5lIHBhcnNpbmdbMV0gc2ltaWxhciB0byBQeXRob24gb3B0cGFyc2UgWzItM10uCgogIE5vdGU6IFB5dGhvbiBhbHNvIHN1cHBvcnRzIGdldG9wdCBbNF0uCgpTVEFUVVMKCiAgVGhpcyBtb2R1bGUgaXMgZmFpcmx5IGJhc2ljIGJ1dCBjb3VsZCBiZSBleHBhbmRlZC4KCkFQSQoKICBTZWUgc291cmNlIGNvZGUgYW5kIGFsc28gY29tcGFyZSB0byBQeXRob24ncyBkb2NzIFsyLDNdIGZvciBkZXRhaWxzIGJlY2F1c2UKICB0aGUgZm9sbG93aW5nIGRvY3VtZW50YXRpb24gaXMgaW5jb21wbGV0ZS4KCiAgb3B0ID0gT3B0aW9uUGFyc2VyIHt1c2FnZT11c2FnZSwgdmVyc2lvbj12ZXJzaW9uLCBhZGRfaGVscF9vcHRpb249YWRkX2hlbHBfb3B0aW9ufQoKICAgIENyZWF0ZSBjb21tYW5kIGxpbmUgcGFyc2VyLgoKICBvcHQuYWRkX29wdGlvbnN7c2hvcnRmbGFnLCBsb25nZmxhZywgYWN0aW9uPWFjdGlvbiwgbWV0YXZhcj1tZXRhdmFyLCBkZXN0PWRlc3QsIGhlbHA9aGVscH0KCiAgICBBZGQgY29tbWFuZCBsaW5lIG9wdGlvbiBzcGVjaWZpY2F0aW9uLiAgVGhpcyBtYXkgYmUgY2FsbGVkIG11bHRpcGxlIHRpbWVzLgoKICBvcHQucGFyc2VfYXJncygpIC0tPiBvcHRpb25zLCBhcmdzCgogICAgUGVyZm9ybSBhcmd1bWVudCBwYXJzaW5nLgoKREVQRU5ERU5DSUVTCgogIE5vbmUgKG90aGVyIHRoYW4gTHVhIDUuMSBvciA1LjIpCgpSRUZFUkVOQ0VTCgogIFsxXSBodHRwOi8vbHVhLXVzZXJzLm9yZy93aWtpL0NvbW1hbmRMaW5lUGFyc2luZwogIFsyXSBodHRwOi8vZG9jcy5weXRob24ub3JnL2xpYi9vcHRwYXJzZS1kZWZpbmluZy1vcHRpb25zLmh0bWwKICBbM10gaHR0cDovL2Jsb2cuZG91Z2hlbGxtYW5uLmNvbS8yMDA3LzA4L3B5bW90dy1vcHRwYXJzZS5odG1sCiAgWzRdIGh0dHA6Ly9kb2NzLnB5dGhvbi5vcmcvbGliL21vZHVsZS1nZXRvcHQuaHRtbAoKTElDRU5TRQoKICAoYykgMjAwOC0yMDExIERhdmlkIE1hbnVyYS4gIExpY2Vuc2VkIHVuZGVyIHRoZSBzYW1lIHRlcm1zIGFzIEx1YSAoTUlUKS4KCiAgUGVybWlzc2lvbiBpcyBoZXJlYnkgZ3JhbnRlZCwgZnJlZSBvZiBjaGFyZ2UsIHRvIGFueSBwZXJzb24gb2J0YWluaW5nIGEgY29weQogIG9mIHRoaXMgc29mdHdhcmUgYW5kIGFzc29jaWF0ZWQgZG9jdW1lbnRhdGlvbiBmaWxlcyAodGhlICJTb2Z0d2FyZSIpLCB0byBkZWFsCiAgaW4gdGhlIFNvZnR3YXJlIHdpdGhvdXQgcmVzdHJpY3Rpb24sIGluY2x1ZGluZyB3aXRob3V0IGxpbWl0YXRpb24gdGhlIHJpZ2h0cwogIHRvIHVzZSwgY29weSwgbW9kaWZ5LCBtZXJnZSwgcHVibGlzaCwgZGlzdHJpYnV0ZSwgc3VibGljZW5zZSwgYW5kL29yIHNlbGwKICBjb3BpZXMgb2YgdGhlIFNvZnR3YXJlLCBhbmQgdG8gcGVybWl0IHBlcnNvbnMgdG8gd2hvbSB0aGUgU29mdHdhcmUgaXMKICBmdXJuaXNoZWQgdG8gZG8gc28sIHN1YmplY3QgdG8gdGhlIGZvbGxvd2luZyBjb25kaXRpb25zOgoKICBUaGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSBhbmQgdGhpcyBwZXJtaXNzaW9uIG5vdGljZSBzaGFsbCBiZSBpbmNsdWRlZCBpbgogIGFsbCBjb3BpZXMgb3Igc3Vic3RhbnRpYWwgcG9ydGlvbnMgb2YgdGhlIFNvZnR3YXJlLgoKICBUSEUgU09GVFdBUkUgSVMgUFJPVklERUQgIkFTIElTIiwgV0lUSE9VVCBXQVJSQU5UWSBPRiBBTlkgS0lORCwgRVhQUkVTUyBPUgogIElNUExJRUQsIElOQ0xVRElORyBCVVQgTk9UIExJTUlURUQgVE8gVEhFIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZLAogIEZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFORCBOT05JTkZSSU5HRU1FTlQuICBJTiBOTyBFVkVOVCBTSEFMTCBUSEUKICBBVVRIT1JTIE9SIENPUFlSSUdIVCBIT0xERVJTIEJFIExJQUJMRSBGT1IgQU5ZIENMQUlNLCBEQU1BR0VTIE9SIE9USEVSCiAgTElBQklMSVRZLCBXSEVUSEVSIElOIEFOIEFDVElPTiBPRiBDT05UUkFDVCwgVE9SVCBPUiBPVEhFUldJU0UsIEFSSVNJTkcgRlJPTSwKICBPVVQgT0YgT1IgSU4gQ09OTkVDVElPTiBXSVRIIFRIRSBTT0ZUV0FSRSBPUiBUSEUgVVNFIE9SIE9USEVSIERFQUxJTkdTIElOCiAgVEhFIFNPRlRXQVJFLgogIChlbmQgbGljZW5zZSkKCiAtLV1dCgpsb2NhbCBNID0ge19UWVBFPSdtb2R1bGUnLCBfTkFNRT0ncHl0aG9uaWMub3B0cGFyc2UnLCBfVkVSU0lPTj0nMC4zLjIwMTExMTI4J30KCmxvY2FsIGlwYWlycyA9IGlwYWlycwpsb2NhbCB1bnBhY2sgPSB1bnBhY2sKbG9jYWwgaW8gPSBpbwpsb2NhbCB0YWJsZSA9IHRhYmxlCmxvY2FsIG9zID0gb3MKbG9jYWwgYXJnID0gYXJnCgoKbG9jYWwgZnVuY3Rpb24gT3B0aW9uUGFyc2VyKHQpCiAgICBsb2NhbCB1c2FnZSA9IHQudXNhZ2UKICAgIC0tbG9jYWwgdmVyc2lvbiA9IHQudmVyc2lvbgoKICAgIGxvY2FsIG8gPSB7fQogICAgbG9jYWwgb3B0aW9uX2Rlc2NyaXB0aW9ucyA9IHt9CiAgICBsb2NhbCBvcHRpb25fb2YgPSB7fQoKICAgIGZ1bmN0aW9uIG8uZmFpbChzKSAtLSBleHRlbnNpb24KICAgICAgICBpby5zdGRlcnI6d3JpdGUocyAuLiAnXG4nKQogICAgICAgIG9zLmV4aXQoMSkKICAgIGVuZAoKICAgIGZ1bmN0aW9uIG8uYWRkX29wdGlvbihvcHRkZXNjKQogICAgICAgIG9wdGlvbl9kZXNjcmlwdGlvbnNbI29wdGlvbl9kZXNjcmlwdGlvbnMrMV0gPSBvcHRkZXNjCiAgICAgICAgZm9yIF8sdiBpbiBpcGFpcnMob3B0ZGVzYykgZG8KICAgICAgICAgICAgb3B0aW9uX29mW3ZdID0gb3B0ZGVzYwogICAgICAgIGVuZAogICAgZW5kCiAgICBmdW5jdGlvbiBvLnBhcnNlX2FyZ3MoKQogICAgICAgIC0tIGV4cGFuZCBvcHRpb25zIChlLmcuICItLWlucHV0PWZpbGUiIC0+ICItLWlucHV0IiwgImZpbGUiKQogICAgICAgIGxvY2FsIGFyZyA9IHt1bnBhY2soYXJnKX0KICAgICAgICBmb3IgaT0jYXJnLDEsLTEgZG8gbG9jYWwgdiA9IGFyZ1tpXQogICAgICAgICAgICBsb2NhbCBmbGFnLCB2YWwgPSB2Om1hdGNoKCdeKCUtJS0ldyspPSguKiknKQogICAgICAgICAgICBpZiBmbGFnIHRoZW4KICAgICAgICAgICAgICAgIGFyZ1tpXSA9IGZsYWcKICAgICAgICAgICAgICAgIHRhYmxlLmluc2VydChhcmcsIGkrMSwgdmFsKQogICAgICAgICAgICBlbmQKICAgICAgICBlbmQKCiAgICAgICAgbG9jYWwgb3B0aW9ucyA9IHt9CiAgICAgICAgbG9jYWwgYXJncyA9IHt9CiAgICAgICAgbG9jYWwgaSA9IDEKICAgICAgICB3aGlsZSBpIDw9ICNhcmcgZG8gbG9jYWwgdiA9IGFyZ1tpXQogICAgICAgICAgICBsb2NhbCBvcHRkZXNjID0gb3B0aW9uX29mW3ZdCiAgICAgICAgICAgIGlmIG9wdGRlc2MgdGhlbgogICAgICAgICAgICAgICAgbG9jYWwgYWN0aW9uID0gb3B0ZGVzYy5hY3Rpb24KICAgICAgICAgICAgICAgIGxvY2FsIHZhbAogICAgICAgICAgICAgICAgaWYgYWN0aW9uID09ICdzdG9yZScgb3IgYWN0aW9uID09IG5pbCB0aGVuCiAgICAgICAgICAgICAgICAgICAgaSA9IGkgKyAxCiAgICAgICAgICAgICAgICAgICAgdmFsID0gYXJnW2ldCiAgICAgICAgICAgICAgICAgICAgaWYgbm90IHZhbCB0aGVuIG8uZmFpbCgnb3B0aW9uIHJlcXVpcmVzIGFuIGFyZ3VtZW50ICcgLi4gdikgZW5kCiAgICAgICAgICAgICAgICBlbHNlaWYgYWN0aW9uID09ICdzdG9yZV90cnVlJyB0aGVuCiAgICAgICAgICAgICAgICAgICAgdmFsID0gdHJ1ZQogICAgICAgICAgICAgICAgZWxzZWlmIGFjdGlvbiA9PSAnc3RvcmVfZmFsc2UnIHRoZW4KICAgICAgICAgICAgICAgICAgICB2YWwgPSBmYWxzZQogICAgICAgICAgICAgICAgZW5kCiAgICAgICAgICAgICAgICBvcHRpb25zW29wdGRlc2MuZGVzdF0gPSB2YWwKICAgICAgICAgICAgZWxzZQogICAgICAgICAgICAgICAgaWYgdjptYXRjaCgnXiUtJykgdGhlbiBvLmZhaWwoJ2ludmFsaWQgb3B0aW9uICcgLi4gdikgZW5kCiAgICAgICAgICAgICAgICBhcmdzWyNhcmdzKzFdID0gdgogICAgICAgICAgICBlbmQKICAgICAgICAgICAgaSA9IGkgKyAxCiAgICAgICAgZW5kCiAgICAgICAgaWYgb3B0aW9ucy5oZWxwIHRoZW4KICAgICAgICAgICAgby5wcmludF9oZWxwKCkKICAgICAgICAgICAgb3MuZXhpdCgpCiAgICAgICAgZW5kCiAgICAgICAgaWYgb3B0aW9ucy52ZXJzaW9uIHRoZW4KICAgICAgICAgICAgaW8uc3Rkb3V0OndyaXRlKHQudmVyc2lvbiAuLiAiXG4iKQogICAgICAgICAgICBvcy5leGl0KCkKICAgICAgICBlbmQKICAgICAgICByZXR1cm4gb3B0aW9ucywgYXJncwogICAgZW5kCgogICAgbG9jYWwgZnVuY3Rpb24gZmxhZ3Nfc3RyKG9wdGRlc2MpCiAgICAgICAgbG9jYWwgc2ZsYWdzID0ge30KICAgICAgICBsb2NhbCBhY3Rpb24gPSBvcHRkZXNjLmFjdGlvbgogICAgICAgIGZvciBfLGZsYWcgaW4gaXBhaXJzKG9wdGRlc2MpIGRvCiAgICAgICAgICAgIGxvY2FsIHNmbGFnZW5kCiAgICAgICAgICAgIGlmIGFjdGlvbiA9PSBuaWwgb3IgYWN0aW9uID09ICdzdG9yZScgdGhlbgogICAgICAgICAgICAgICAgbG9jYWwgbWV0YXZhciA9IG9wdGRlc2MubWV0YXZhciBvciBvcHRkZXNjLmRlc3Q6dXBwZXIoKQogICAgICAgICAgICAgICAgc2ZsYWdlbmQgPSAjZmxhZyA9PSAyIGFuZCAnICcgLi4gbWV0YXZhcgogICAgICAgICAgICAgICAgICAgICAgICBvciAgJz0nIC4uIG1ldGF2YXIKICAgICAgICAgICAgZWxzZQogICAgICAgICAgICAgICAgc2ZsYWdlbmQgPSAnJwogICAgICAgICAgICBlbmQKICAgICAgICAgICAgc2ZsYWdzWyNzZmxhZ3MrMV0gPSBmbGFnIC4uIHNmbGFnZW5kCiAgICAgICAgZW5kCiAgICAgICAgcmV0dXJuIHRhYmxlLmNvbmNhdChzZmxhZ3MsICcsICcpCiAgICBlbmQKCiAgICBmdW5jdGlvbiBvLnByaW50X2hlbHAoKQogICAgICAgIGlvLnN0ZG91dDp3cml0ZSgiVXNhZ2U6ICIgLi4gdXNhZ2U6Z3N1YignJSVwcm9nJywgYXJnWzBdKSAuLiAiXG4iKQogICAgICAgIGlvLnN0ZG91dDp3cml0ZSgiXG4iKQogICAgICAgIGlvLnN0ZG91dDp3cml0ZSgiT3B0aW9uczpcbiIpCiAgICAgICAgbG9jYWwgbWF4d2lkdGggPSAwCiAgICAgICAgZm9yIF8sb3B0ZGVzYyBpbiBpcGFpcnMob3B0aW9uX2Rlc2NyaXB0aW9ucykgZG8KICAgICAgICAgICAgbWF4d2lkdGggPSBtYXRoLm1heChtYXh3aWR0aCwgI2ZsYWdzX3N0cihvcHRkZXNjKSkKICAgICAgICBlbmQKICAgICAgICBmb3IgXyxvcHRkZXNjIGluIGlwYWlycyhvcHRpb25fZGVzY3JpcHRpb25zKSBkbwogICAgICAgICAgICBpby5zdGRvdXQ6d3JpdGUoIiAgIiAuLiAoJyUtJy4ubWF4d2lkdGguLidzICAnKTpmb3JtYXQoZmxhZ3Nfc3RyKG9wdGRlc2MpKQogICAgICAgICAgICAgICAgICAgIC4uIG9wdGRlc2MuaGVscCAuLiAiXG4iKQogICAgICAgIGVuZAogICAgZW5kCiAgICBpZiB0LmFkZF9oZWxwX29wdGlvbiA9PSBuaWwgb3IgdC5hZGRfaGVscF9vcHRpb24gPT0gdHJ1ZSB0aGVuCiAgICAgICAgby5hZGRfb3B0aW9ueyItLWhlbHAiLCBhY3Rpb249InN0b3JlX3RydWUiLCBkZXN0PSJoZWxwIiwKICAgICAgICAgICAgICAgICAgICAgaGVscD0ic2hvdyB0aGlzIGhlbHAgbWVzc2FnZSBhbmQgZXhpdCJ9CiAgICBlbmQKICAgIGlmIHQudmVyc2lvbiB0aGVuCiAgICAgICAgby5hZGRfb3B0aW9ueyItLXZlcnNpb24iLCBhY3Rpb249InN0b3JlX3RydWUiLCBkZXN0PSJ2ZXJzaW9uIiwKICAgICAgICAgICAgICAgICAgICAgaGVscD0ib3V0cHV0IHZlcnNpb24gaW5mby4ifQogICAgZW5kCiAgICByZXR1cm4gbwplbmQKCgpNLk9wdGlvblBhcnNlciA9IE9wdGlvblBhcnNlcgoKCnJldHVybiBNCg=="

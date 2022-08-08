[admin@localhost ~]$ ldd convstr
	linux-vdso.so.1 =>  (0x00007ffd5a7dd000)
	libiconv.so.2 => /var/lib/libiconv.so.2 (0x00007fb433978000)
	libpthread.so.0 => /lib64/libpthread.so.0 (0x00007fb43375c000)
	libc.so.6 => /lib64/libc.so.6 (0x00007fb43338e000)
	/lib64/ld-linux-x86-64.so.2 (0x00007fb433c76000)


libiconv.so.2 需要加载到动态库环境变量或放入到指定目录

export export LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:/var/lib




#ifndef MSB_CORE_H
#define MSB_CORE_H

/*
    core.h - core syscalls like sync, reboot, etc...

    This code is part of libmsb (the Malino Syscall Bridge).

    Copyleft 2024 malino. This code is licensed under GPL2.
*/

void msb_sync();
long msb_reboot(unsigned int cmd);
long msb_write(unsigned int fd, const char *buf);
long msb_read(unsigned int fd, char *buf, unsigned long count);

#endif
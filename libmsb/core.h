#ifndef MSB_CORE_H
#define MSB_CORE_H

/*
    core.h - core syscalls like sync, reboot, etc...

    This code is part of libmsb (the Malino Syscall Bridge).

    Copyleft 2024 malino. This code is licensed under MPL-2.0.
*/

#include <stdint.h>
#include <unistd.h>
#include <stdbool.h>

void msb_sync();
int msb_reboot(uint32_t cmd);
long msb_write(uint32_t fd, const char *_Nonnull buf);
long msb_read(uint32_t fd, char *_Nonnull buf, uint64_t count);
int msb_mount(const char *_Nonnull source, const char *_Nonnull target, const char *_Nonnull fstype, uint64_t flags, const void *_Nullable data);
int msb_umount2(const char *_Nonnull target, int flags);
pid_t msb_getpid();
int msb_forkexec(const char *_Nonnull path, char *const _Nullable argv[_Nullable], char *const _Nullable envp[_Nullable], bool wait);
long msb_dsc(long rax, long rdi, long rsi, long rdx, long r10, long r8, long r9);

#endif
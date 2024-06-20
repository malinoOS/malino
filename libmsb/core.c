/*
    core.c - core syscalls like sync, reboot, etc...

    This code is part of libmsb (the Malino Syscall Bridge).

    Copyleft 2024 malino. This code is licensed under GPL2.
*/

#include <unistd.h>
#include <string.h>
#include <linux/reboot.h>
#include <sys/syscall.h>
#include "core.h"
#include <stdint.h>
#include <sys/wait.h>
#include <sys/mount.h>
#include <stdbool.h>
#include <errno.h>

void msb_sync() {
    sync();
}

int msb_reboot(uint32_t cmd) {
    if (syscall(SYS_reboot, LINUX_REBOOT_MAGIC1, LINUX_REBOOT_MAGIC2, cmd) == -1)
        return errno;
    return 0;
}

long msb_write(uint32_t fd, const char *_Nonnull buf) {
    long w = write(fd, buf, strlen(buf));
    if (w == -1)
        return -errno;
    return w;
}

long msb_read(uint32_t fd, char *_Nonnull buf, uint64_t count) {
    long r = read(fd, buf, count);
    if (r == -1)
        return -errno;
    return r;
}

int msb_mount(const char *_Nonnull source, const char *_Nonnull target, const char *_Nonnull fstype, uint64_t flags, const void *_Nullable data) {
    if (mount(source, target, fstype, flags, data) == -1)
        return errno;
    return 0;
}

int msb_umount2(const char *_Nonnull target, int flags) {
    if (umount2(target, flags) == -1)
        return errno;
    return 0;
}

pid_t msb_getpid() {
    return getpid();
}

int msb_forkexec(const char *_Nonnull path, char *const _Nullable argv[], char *const _Nullable envp[], bool wait) {
    pid_t pid = syscall(SYS_fork);
    if (pid < 0) {
        return -errno;
    } else if (pid == 0) {
        // Child process
        execve(path, argv, envp);
        return -errno;
    } else {
        // Parent Process
        if (wait) {
            int status;
            if (waitpid(pid, &status, 0) == -1)
                return -errno;

            if (WIFEXITED(status)) 
                return WEXITSTATUS(status);

            else if (WIFSIGNALED(status))
                return WTERMSIG(status);

            else
                return -42; // ENOMSG
        }
        return 0;
    }
}

long msb_dsc(long rax, long rdi, long rsi, long rdx, long r10, long r8, long r9) {
    return syscall(rax,rdi,rsi,rdx,r10,r8,r9);
}
/*
    core.c - core syscalls like sync, reboot, etc...

    This code is part of libmsb (the Malino Syscall Bridge).

    Copyleft 2024 malino. This code is licensed under MPL-2.0.
*/

#include <unistd.h>
#include <string.h>
#include <sys/syscall.h>
#include "core.h"
#include <stdint.h>
#include <sys/wait.h>
#include <sys/mount.h>
#include <stdbool.h>
#include <errno.h>
#include <sys/reboot.h>

void msb_sync() {
    // sync.
    sync();
}

int msb_reboot(uint32_t cmd) {
    // call reboot(2) with command, return errno if fail
    if (reboot(cmd) == -1)
        return errno;
    return 0;
}

long msb_write(uint32_t fd, const char *_Nonnull buf) {
    // call write(2) with args, return -errno if fail
    long w = write(fd, buf, strlen(buf));
    if (w == -1)
        return -errno;
    return w;
}

long msb_read(uint32_t fd, char *_Nonnull buf, uint64_t count) {
    // call read(2) with args, return -errno if fail
    long r = read(fd, buf, count);
    if (r == -1)
        return -errno;
    return r;
}

int msb_mount(const char *_Nonnull source, const char *_Nonnull target, const char *_Nonnull fstype, uint64_t flags, const void *_Nullable data) {
    // call mount(2) with args, return errno if fail
    if (mount(source, target, fstype, flags, data) == -1)
        return errno;
    return 0;
}

int msb_umount2(const char *_Nonnull target, int flags) {
    // call umount2(2) with args, return errno if fail
    if (umount2(target, flags) == -1)
        return errno;
    return 0;
}

pid_t msb_getpid() {
    // just return getpid(). apparently it will never set errno according to `man getpid.2`
    return getpid();
}

int msb_forkexec(const char *_Nonnull path, char *const _Nullable argv[], char *const _Nullable envp[], bool wait) {
    // fork & execute
    pid_t pid = fork();
    if (pid < 0) {
        return -errno;
    } else if (pid == 0) {
        // Child process
        if (execve(path, argv, envp) == -1)
            return -errno;
        return 0;
    } else {
        // Parent Process
        if (wait) {
            int status;
            // wait for it to exit, return -errno if something weird happened
            if (waitpid(pid, &status, 0) == -1)
                return -errno;

            // If the program exited
            if (WIFEXITED(status)) 
                return WEXITSTATUS(status);

            // If it was signaled to death
            else if (WIFSIGNALED(status))
                return WTERMSIG(status);

            // "i don't know what happened here"
            else
                return -42; // -ENOMSG
        }
        // since we can't get the return code, just return 0 since the process spawned successfully
        return 0;
    }
}

long msb_dsc(long rax, long rdi, long rsi, long rdx, long r10, long r8, long r9) {
    // syscall.
    return syscall(rax,rdi,rsi,rdx,r10,r8,r9);
}
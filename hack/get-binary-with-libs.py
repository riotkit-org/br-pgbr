#!/usr/bin/env python3
import os
import subprocess
import re
import sys

ALREADY_COPIED = []


def copy_dependencies(bin_name: str, target_dir: str):
    bin_path = subprocess.check_output(['whereis', bin_name]).decode('utf-8').split(' ')[1].strip()

    try:
        os.mkdir(target_dir)
    except FileExistsError:
        pass

    print(f' >> Copying dependencies for {bin_name}')
    copy_dependencies_for_path(bin_path, target_dir, depth=0)


def copy_dependencies_for_path(bin_path: str, target_dir: str, depth: int = 0):
    if bin_path in ALREADY_COPIED:
        return

    ldd = subprocess.check_output(['ldd', bin_path]).decode('utf-8').split("\n")

    print(ldd)

    for line in ldd:
        is_static_ld = "=>" not in line and "ld-" in line

        if not line.strip() or ("=>" not in line and not is_static_ld):
            continue

        if is_static_ld:
            parsed = re.findall('\s*([\/A-Za-z\-_.0-9]+)\ ', line)
            orig_name = os.path.basename(parsed[0])
        else:
            try:
                parsed = re.findall('=>\s*([\/A-Za-z\-_.0-9]+)\ ', line)
                orig_name = re.findall('\s*([A-Za-z\-_.0-9]+)\s*=>', line)
            except Exception as e:
                print(">> Line caused error: ", line)
                print("   Exception: ", e)
                raise

            orig_name = orig_name[0]

        if parsed[0] == "ldd":
            # ['\tldd (0x7f5a18c9c000)', '\tlibc.musl-x86_64.so.1 => ldd (0x7f5a18c9c000)', '']
            continue

        real_path = subprocess.check_output(['readlink', '-f', parsed[0]]).decode('utf-8').strip()
        print((" " * depth * 2) + f">> Copying {real_path} -> {orig_name}")
        subprocess.check_call(['cp', real_path, target_dir + "/" + orig_name])

        if not is_static_ld:
            copy_dependencies_for_path(real_path, target_dir, depth + 1)

    ALREADY_COPIED.append(bin_path)


if __name__ == '__main__':
    copy_dependencies(sys.argv[1], sys.argv[2])

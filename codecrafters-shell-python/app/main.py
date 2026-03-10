import sys
import shutil
import subprocess

COMMANDS: list[str] = ["exit", "echo", "type"]


def run_a_program(commands: list[str]):
    first = commands[0]
    args = commands[1:]
    first_valid = shutil.which(first)
    if not first_valid:
        return

    print(f"Program was passed {len(args)} (including program name).")
    for arg in args:
        print("Arg", arg)

    subprocess.call([first, *commands[1:]])


def main():
    sys.stdout.write("$ ")
    commands = input().split()
    command = commands[0]

    if command == "exit":
        exit()
    elif command == "echo":
        print(" ".join(commands[1:]))
    elif command == "type":
        sub_command = " ".join(commands[1:])
        if sub_command in COMMANDS:
            print(f"{sub_command} is a shell builtin")
        else:
            path = shutil.which(sub_command)
            if path:
                print(f"{sub_command} is {path}")
            else:
                print(f"{sub_command} not found")
    else:
        run_a_program(commands)
        print(f"{command}: command not found")
    main()


if __name__ == "__main__":
    main()

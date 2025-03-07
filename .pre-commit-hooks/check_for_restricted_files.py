import sys
import os

def main():
    files = sys.argv[1:]
    for file in files:
        if file.endswith('.json') or file.endswith('.jsonl'):
            print(f"Error: Found a .json or .jsonl file: {file}")
            sys.exit(1)

    sys.exit(0)

if __name__ == "__main__":
    main()

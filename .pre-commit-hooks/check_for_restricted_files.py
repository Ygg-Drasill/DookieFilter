import sys
import os

def main():
    files = sys.argv[1:]
    for file in files:
        if file.endswith('.json') or file.endswith('.jsonl') or file.endswith('.csv'):
            print(f"Error: Found a .json, .jsonl or .csv file: {file}")
            sys.exit(1)

    sys.exit(0)

if __name__ == "__main__":
    main()

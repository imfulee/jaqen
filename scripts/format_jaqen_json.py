import json
import sys

if len(sys.argv) != 2:
    print("USAGE: python3 format_jaqen_json.py /path/to/jaqen_config.json")
    sys.exit()

file_name = sys.argv[1]
data = {}

with open(file_name, "r") as file:
    data = json.load(file)

mapping_override_key = "mapping_override"
nations_map = data[mapping_override_key]
ordered_nations_map = dict(sorted(nations_map.items()))
data[mapping_override_key] = ordered_nations_map

with open(file_name, "w") as file:
    json.dump(data, file, ensure_ascii=False, indent=2)

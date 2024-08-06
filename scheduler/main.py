import json
from typing import TypedDict
import sys
from scheduler import balance_dicts


class LessonData(TypedDict):
    lesson: str
    title: str
    duration: int


def get_data() -> list[LessonData]:
    with open("data/lessons.json", "r") as file:
        data: list[dict] = json.load(file)
    result: list[LessonData] = []
    for element in data:
        duration: str = element["duration"]
        colon = duration.find(":")
        hours = int(duration[:colon])
        minutes = int(duration[colon + 1 :])
        result.append(
            {
                "lesson": element["lesson"],
                "title": element["title"],
                "duration": hours * 60 + minutes,
            }
        )
    return result


def print_condensed_data(data: list[list[LessonData]]) -> None:
    for d in data:
        full_duration = sum([e["duration"] for e in d])
        print(f"{full_duration}m - {', '.join([e['lesson'] for e in d])}")


def main(already_complete: int = 7, total_days: int = 13):
    already_complete = int(already_complete)
    total_days = int(total_days)

    data = get_data()[already_complete:]
    condensed = balance_dicts(data, "duration", total_days)
    print_condensed_data(condensed)


if __name__ == "__main__":
    main(*sys.argv[1:])

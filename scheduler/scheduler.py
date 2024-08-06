def balance_ints(original: list[int], length: int) -> list[list[int]]:
    """Returns 2D list with values condensed and balanced to given size.

    Ex. [3, 3, 1, 2, 3, 4], 5 -> [[3], [3], [1, 2], [3], [5]]"""

    average: int = int(sum(original) / length)
    new: list[list[int]] = []

    i = 0
    while i < len(original):
        temp: list[int] = []
        while sum(temp) < average and len(original) > i:
            temp.append(original[i])
            i += 1
        new.append(temp)

    return new


def balance_dicts(original: list[dict], key: str, length: int) -> list[list[dict]]:
    """Expects list of dicts with given key matching an int value.
    Returns 2D list with values condensed and balanced to given size.

    Ex. [{"char": "a", "time": 3}, {"char": "b", "time": 1}, {"char": "c", "time": 2}], 2 ->
    [[{"char": "a", "time": 3}], [{"char": "b", "time": 1}, {"char": "c", "time": 2}]]"""

    average: int = int(sum([d[key] for d in original]) / length)
    new: list[list[dict]] = []

    i = 0
    while i < len(original):
        temp: list[dict] = []
        while sum([d[key] for d in temp]) < average and len(original) > i:
            temp.append(original[i])
            i += 1
        new.append(temp)

    return new

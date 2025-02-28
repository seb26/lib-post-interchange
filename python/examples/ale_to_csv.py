from libale import ALELib

from pathlib import Path
from typing import Iterable, Generator, Union
import argparse
import csv
import io
import logging

logger = logging.getLogger(__name__)

def is_ale_file_ext(filename: str | Path):
    return str(filename).lower()[-3:] == 'ale'

def get_ale_filepaths(
    input: Union[str, Iterable],
    allow_all: bool = False,
    recurse: bool = False,
) -> Generator[Path, None, None]:
    """
    Return Path()s of files and folders that are recognised media files, determined by file extension

    :param input: Path()s in a string or Iterable
    :param allow_all: Skip file extension check and pass all items. Default is False
    :param recurse: deeply recurse all folders found and return all items beneath this path in the filesystem.
                    Default is False and will just return files that are immediate children

    [Originally from get_media_duration]
    """
    def _iterate(filepath: Path):
        if filepath.is_file():
            if is_ale_file_ext(filepath) or allow_all:
                yield filepath
        elif filepath.is_dir():
            for child in filepath.iterdir():
                childpath = Path(child)
                if is_ale_file_ext(childpath) or allow_all:
                    yield childpath
                if recurse is True:
                    if childpath.is_dir():
                        yield from _iterate(childpath)

    if isinstance(input, str):
        filepath = Path(input)
        yield from _iterate(filepath)
    elif isinstance(input, Iterable):
        for input_item in input:
            filepath = Path(input_item)
            yield from _iterate(filepath)

def parse(items: list, user_map: list = None):
    def _build(clip: dict, mappings: list[dict]) -> Generator:
        for mapping in mappings:
            ale_col = mapping.get('ale_col')
            csv_col = mapping.get('csv_col')
            value = clip.get(ale_col)
            if csv_col and value:
                # Good column name and good value
                yield csv_col, value
            elif csv_col and not value:
                # Good column name but no result in the source
                yield csv_col, None
            else:
                logger.debug(f'Skipped during mapping - ALE Col: {ale_col} - CSV Col: {csv_col} - Value: {value}')

    mappings = []
    user_mapped_columns = False
    if user_map:
        user_mapped_columns = True
        logger.debug(f'Input mapping: {user_map}')
        for index, mapping_raw in enumerate(user_map):
            try:
                ale_col, csv_col = mapping_raw.split(':')
                mappings.append(
                    dict(
                        ale_col = ale_col,
                        csv_col = csv_col,
                    )
                )
            except (
                ValueError,
                AttributeError,
            ) as e:
                raise Exception(f'Error interpreting this map item #{index + 1}: {mapping_raw}. Ensure it is using the syntax: ALEColumnName:CSVColumnName')
    if not user_map:
        # Map will be created during ALE parsing
        pass
    logger.debug(f'Mappings: {mappings}')
    logger.debug(f'Input items: {items}')
    ale_files = list( get_ale_filepaths(items, recurse=True) )
    ale_files.sort()

    entries = []
    for ale_file in ale_files:
        logger.debug(f'Parsing {ale_file.name}')
        with open(ale_file, 'r') as f:
            ale = ALELib.parse(f.read())
            if user_mapped_columns is False:
                # If no user map specified, pass all ALE columns through to the CSV unaltered
                mappings_for_this_ale = [ { 'ale_col': col, 'csv_col': col } for col in ale.columns ]
                mappings.extend(mappings_for_this_ale)
            for index, clip in enumerate(ale.clips):
                # Look up the values per column
                try:
                    subtable = dict( _build(clip, mappings) )
                except Exception as e:
                    logger.debug(f'Exception {type(e)} on ALE file {str(ale_file)}: Clip #{index + 1}')
                    logger.debug(e, exc_info=1)
                    continue
                entries.append(subtable)
    # Process columns
    columns = [ mapping['csv_col'] for mapping in mappings ]
    logger.debug(f'Columns: {columns}')
    return entries, columns

def write_csv(file_target, entries: list, columns: list):
    output_csv_file = csv.DictWriter(
        file_target,
        delimiter = ',',
        fieldnames = columns,
    )
    output_csv_file.writeheader()
    for item in entries:
        output_csv_file.writerow(item)
    return output_csv_file

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('items', help='folders or ALE files', nargs='+', action='append')
    parser.add_argument('--map', help='map of ALE column names to CSV column names. Format is ale_col:csv_col, space separated. Escape spaces in the column names with a backslash or quotes. Example: Name:name_csv "Start:Start TC"', nargs='+', action='append')
    parser.add_argument('--debug', help='include debug output', action='store_true', default=False)
    parser.add_argument('-o', help='output to CSV file, otherwise prints to stdout', type=str, required=False)
    args = parser.parse_args()

    if args.debug:
        logging.basicConfig(level=logging.DEBUG)

    # Parse
    # args.items and args.map are returned as double lists - work around this with first index only
    entries, columns = parse(args.items[0], args.map[0])

    # Output
    if args.o:
        logger.debug(f'Writing CSV output to file: {args.o}')
        file_target = open(args.o, 'w', encoding='utf-8')
        write_csv(file_target, entries, columns)
        logger.debug(f'Done.')
    else:
        file_target = io.StringIO()
        write_csv(file_target, entries, columns)
        print( file_target.getvalue().strip() )

if __name__ == '__main__':
    main()
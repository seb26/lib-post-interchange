# Author: Sebastian Reategui

import re

class ALELibDefaults:
    VIDEO_FORMAT = 'CUSTOM'
    AUDIO_FORMAT = '48khz'
    FPS = '24'

class ALELibParseError(Exception):
    pass

class ALE:
    def __init__(self, data: str, table_data: list[dict], video_format: str, audio_format: str, fps: str):
        self.audio_format = audio_format
        self.data = data
        self.fps = fps
        self.clips = table_data
        self.video_format = video_format

    @property
    def columns(self) -> list:
        return list(self.clips[0].keys())

class ALELib:
    @classmethod
    def parse(self, data: str) -> ALE:
        """
        @param data: str - `.ALE` file read into string
        """

        pattern_section_separator = r"(Heading|Column|Data)[\r\n]+"
        pattern_delimiter = '\t'

        # Split ALE data into sections by - Heading, Column, Data
        data_section_raw = re.split(
            pattern_section_separator,
            data,
            flags = re.MULTILINE,
        )
        data_section = [ line.strip() for line in data_section_raw if line ]

        # Heading
        metadata = {}
        if data_section[0] == 'Heading':
            heading_lines = data_section[1].splitlines()
            if heading_lines[0] != 'FIELD_DELIM\tTABS':
                raise ALELibParseError('ERROR - ALELib: parse_ale: unrecognised FIELD_DELIM value. Only TABS is recognised.')
            for index, line in enumerate(heading_lines):
                if not line or line == '':
                    continue
                try:
                    key, value = line.split(pattern_delimiter)
                except ValueError:
                    raise ALELibParseError(f'Unable to parse line {index + 1}')
                metadata[key] = value
        else:
            raise ALELibParseError(f'Unable to parse this file, it does not start with "Heading"')

        # Columns
        if data_section[2] == 'Column':
            columns = [col.strip() for col in data_section[3].strip().split(pattern_delimiter) if col.strip()]
        else:
            raise ALELibParseError('Unable to parse columns.')

        # Data
        table_data = []
        if data_section[4] == 'Data':
            data_rows = [row.strip() for row in data_section[5].strip().splitlines() if row.strip()]
            for row_string in data_rows:
                rows = row_string.split(pattern_delimiter)
                row = dict(zip(columns, rows))
                table_data.append(row)
        else:
            raise ALELibParseError('Unable to parse data fields.')

        return ALE(
            data = data,
            table_data = table_data,
            video_format = metadata['VIDEO_FORMAT'],
            audio_format = metadata['AUDIO_FORMAT'],
            fps = int(metadata['FPS']),
        )
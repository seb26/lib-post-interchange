# ale_to_csv.py

## Usage
```
usage: ale_to_csv.py [-h] [--map MAP [MAP ...]] [--debug] [-o O] items [items ...]

positional arguments:
  items                folders or ALE files

options:
  -h, --help           show this help message and exit
  --map MAP [MAP ...]  map of ALE column names to CSV column names.
                       Format is ale_col:csv_col, space separated.
                       Escape spaces in the column names with a backslash or quotes.
                       Example: Name:name_csv "Start:Start TC"
  --debug              include debug output
  -o                   output to CSV file, otherwise prints to stdout
```

### Notes

- The content of multiple ALE files is output to a single, continuous CSV. To route multiple ALEs each to a single CSV, just run `ale_to_csv.py` multiple times each time specifying only one ALE file. 

## Examples
### Basic import

One ALE file in, and CSV output sent to stdout.

```bash
python scripts/ale_to_csv.py samples/A901R1AA_AVID.ale
```
```csv
Name,Source File,Clip,Duration,Tracks,Start,End,FPS,Original_video,Audio_format,Audio_sr,Audio_bit,Frame_width,Frame_height,Uuid,Sup_version,Exposure_index,Gamma,White_balance,Cc_shift,Look_name,Look_burned_in,Sensor_fps,Shutter_angle,Manufacturer,Camera_model,Camera_sn,Camera_id,Camera_index,Project_fps,Storage_sn,Production,Cinematographer,Operator,Director,Location,Company,User_info1,User_info2,Date_camera,Time_camera,Reel_name,Scene,Take,ASC_SAT,ASC_SOP,Look_user_lut,Lut_file_name,Nd_filterdensity,Focus_distance_unit,Lens_sn,Lens_type,Image_orientation,Image_sharpness,Image_detail,Image_denoising
A001C001_240426_R1AA,A001C001_240426_R1AA.mxf,C001,00:00:29:06,V,03:44:36:21,03:45:06:02,25,ARRIRAW (2202p),,,,3424,2202,B772D724-03CC-11,6.01.02,800,LOG-C,5600,+0,ARRI 709.AML,No,25.000,180.0,ARRI,ALEXA Mini,0021666,R1AA,A,25.000,182501300631,,,,,,,,,20240426,12h58m45s,A001R1AA,,,1.000,(1.000 1.000 1.000)(0.000 0.000 0.000)(1.000 1.000 1.000),No,,0,Imperial,0,Panavision PRIMO_ZOO,0,0,0,0
[...]
```

<br />

### Output to a file

```bash
python scripts/ale_to_csv.py samples/A901R1AA_AVID.ale -o test_output.csv
```

<br />

### User column mapping

User-specified mapping of ALE columns to CSV columns is shown. Useful to include only specific columns, and to determine new columns names that are tailored for the receiving program. 

The format is:

```--map ALECol:CSVCol```

There are two ways escape spaces in column names - either quote it or escape the backspace. 

```bash
python scripts/ale_to_csv.py samples/ --map "Name:Reel Name" Look_name:LUT\ Used
```
```csv
Reel Name,LUT Used
A001C001_240426_R1AA,ARRI 709.AML
A001C002_240426_R1AA,ARRI 709.AML
A901C001_240426_R1AA,ARRI 709.AML
[...]
```



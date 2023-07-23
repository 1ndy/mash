# Mash -- copy one file into another
Mash is a commandline utility for inserting code or yaml into a yaml file. This 
is useful when developing cloudformation templates for AWS, especially when 
using lambdas or quicksight. Operation is simple"
```
     __  ______   _____ __  __
    /  |/  /   | / ___// / / /
   / /|_/ / /| | \__ \/ /_/ /
  / /  / / ___ |___/ / __  /  -- combine code and yaml
 /_/  /_/_/  |_/____/_/ /_/
 Version 0.1

Usage: mash [code|yaml] <file> [into|over] <yaml_file> at <path.seperated.by.dots>

    [code|yaml]        whether <file> should be inserted as code (with a | for multiline) or as yaml
    <file>             the name of the code or yaml to insert into another file
    [into|over]        into will produce a new file, over will overwrite
    <yaml_file>        the yaml file to insert into
    at                 The word at. The design is very human
    <path>             The sequence of keys in <yaml_file> representing the location to insert

```

## Example - inserting code

Suppose we have the following simplified cloudformation template
```
Parameters:
  Name: asdf

Resources:
  Compute:
    Type: AWS::EC2::instance
    Properties:
      Name: 'computer'
      ImageId: 'ami-72d929jd'
  
  Storage:
    Type: AWS::S3::Bucket
    Properties:
      Name: 'storage-bucket'

  Lambda:
    Type: AWS::Lambda::Function
    Properties:
      Code:
        ZipFile:
```
This template sets up a lambda function. Imagine we have a very complex python 
script that needs to run in the lambda, but was developed and tested locally.
```
import pandas as pd
df = pd.read_csv('data.csv')
print(df.head())
```
Copying the python code into the cloudformation template is time consuming and 
error prone. Instead of doing that, use Mash:
```
mash code main.py into cfn.yaml at Resources.Lambda.Properties.Code.ZipFile
```
This will copy the code into the template at the path specified. It will add 
the right amount of indentiation and a character to to enable a multiline 
string in the template:
```
Parameters:
  Name: asdf

Resources:
  Compute:
    Type: AWS::EC2::instance
    Properties:
      Name: 'computer'
      ImageId: 'ami-72d929jd'
  
  Storage:
    Type: AWS::S3::Bucket
    Properties:
      Name: 'storage-bucket'

  Lambda:
    Type: AWS::Lambda::Function
    Properties:
      Code:
        ZipFile: |
          import pandas as pd
          df = pd.read_csv('data.csv')
          print(df.head())

```

## JSON
MAsh doesn't support JSON. I know it's technically a subset of yaml, but parsing it is more complex than parsing yaml. Mash doesn't use the yaml go package, it finds keys with regex and implements a simple Tree data structure to store and query the structure of the yaml file. Adding support for JSON is possible, but for now if you need that functionality, use jq/yq to do the apropriate conversions before mashing.
```
> yq -p json -o yaml fil.json > file.yml
> mash ...
> yq -p yaml -o json output.yml
```

## Tabs
Indenting with tabs is not supported at this time. If your file has tabs, try piping through `sed` first.
# Data structure

The internal data structure is part of a mechanism that allows the developer to easily model the structure of the application they are developing.

## Basic definition

We use three main datasets

*data_type* - Contains definition of supported types

| Name                    |
|-------------------------|
| id                      |
| type                    |
| format                  |
| defailt_validation_rule |
| default_value           | 

*struct* - Contains the definition of the database structure

| Name                  |
|-----------------------|
| id                    |
| parent_id             |
| data_type_id          |
| label                 |

*dataset* - Contains the data itself, is composed of the following elements.

| Name         |
|--------------|
| id           |
| parent_id    |
| struct_id    |
| value        |


## Recursive nature

## Preproc

## Postproc

## Data valitation

## Save data

## Data history

## Abstraction and isolation

## Examples




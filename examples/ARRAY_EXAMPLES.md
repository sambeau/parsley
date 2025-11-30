# Array and String Manipulation Examples

This file demonstrates the array and string manipulation features in Pars.

## Array Indexing

let arr = [10,20,30,40,50];
let first = arr[0];
let third = arr[2];
let last = arr[-1];
let secondLast = arr[-2];

first
third
last
secondLast

## Array Slicing

let numbers = [1,2,3,4,5,6,7,8,9,10];
let slice1 = numbers[2:5];
let slice2 = numbers[0:3];
let slice3 = numbers[5:10];

slice1
slice2
slice3

## Array Concatenation

let arr1 = [1,2,3];
let arr2 = [4,5,6];
let combined = arr1 ++ arr2;

combined

let extended = [1] ++ [2] ++ [3] ++ [4,5,6];
extended

## String Operations

let greeting = "Hello, Pars!";
let char1 = greeting[0];
let char2 = greeting[7];
let lastChar = greeting[-1];

char1
char2
lastChar

let substring1 = greeting[0:5];
let substring2 = greeting[7:11];

substring1
substring2

let fullGreeting = "Hello, " + "World!";
fullGreeting

let name = "Pars";
let message = "Welcome to " + name + "!";
message

## Length Functions

let arrayLength = len(arr);
let stringLength = len(greeting);
let emptyLength = len("");

arrayLength
stringLength
emptyLength

## Complex Example

let data = [100,200,300,400,500];
let doubled = for(x in data) { x * 2 };
let firstThree = doubled[0:3];
let withExtra = firstThree ++ [999];
let result = "Result: " + "length=" ++ len(withExtra);

doubled
firstThree
withExtra

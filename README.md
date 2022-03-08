# ambassadors-app
App for smart-shuffling 9th grade students into pod groups for orientation 

### - [Overview](#overview) 
### - [Validation](#validation)
### - [Nuances](#nuances)
### - [Algorithm](#algorithm-technical-overview)

# Creating pods
### Goal: put students in pods so that each group (ex. Spanish, Graham, Band) is represented in each pod as closely to the percent that it is represented in the population

## Overview
For example, suppose the population had the following diversity: `Gender: female: 45%, male: 45%, nob-binary: 10%`, `School: graham: 30%, crittenden: 30%, blach: 40%`.

Our algorithm would try to make it so that in each pod, each of these groups has the same diversity (`female: 45%, ..., blach: 40%`).

It can take as many types of fields into account. Ex. currently, it takes middle school, gender, language, and group memberships into account, but more could easily be added.
It would make sure that all of these groups are represented as close to the population as possible in every pod.

## Validation
To see how good our algorithm performs, we generate a sample dataset of 600 students into given percentages of each category (gender, school, group memberships, ...).
After our algorithm makes pods, we check how close each group in each pod was represented as it is in the population and calculate the errors.
- An error is the absolute difference between the percent of a group in a pod vs it's percent in the population
    - Ex. if the population was 50% male and a pod was 45% male, the error would be 5%
- We find the average of all errors across every group (male, female, graham, spanish, eld, ...) in every pod
    - ### The average error in our generated dataset is 2%
    - This means that on average, a group such as male or graham would have a diversity in every pod that is only 2% away from the percent in the population!
- ### The highest error is less than 20%
- ### There are errors for every single group in all of the ~60 pods, and less than 10 were worse than 10%

## Nuances
In actual pods, gender is the most visible aspect of diversity, so we thought about increasing the weight of gender in making pods. However given the tiny error, we decided this was not needed and we could leave the algorithm at its pure and simple form. 

ELD 1 and 2 students are not proficient in English, so they are put in separate pods.
The code treats the percentages in their population separately, and still makes their pods as diverse as possible using the same metrics as the other pods.


## Algorithm (technical overview)
Each student can be though of as an array of strings, which are their fields.
A percents object is defined as a map from field to percent (float) of the population that that group is.


To create pods, when you are adding the next student to the pod, you simply choose the student with the maximum weight.
The weight is a function of the percents of the current pod and the percents of the full student population.

### Calculating the weight of a student
To calculate the weight of a student, for each field:
Add the difference between the population percent of the field and the pod percent of the field to the total weight

In pseudo code, this means:

```python
weight = 0
for field in student:
  weight += population[field] - pod[field]
```

This means that if groups are underrepresented in the pod, students with that group have a positive weight contribution from that field, and if groups are overrrepresented in the pod, students with that group have a negative weight contribution from that field.

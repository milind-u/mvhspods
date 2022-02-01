# ambassadors-app
App for smart-shuffling 9th grade students into pod groups for orientation 

<a href="https://docs.google.com/document/d/15hoqYRgNQ6xVI2ec4xbVub2b8ZB5rycKkReRvVbN30Y/edit?usp=sharing">Link to doc</a>

# Smart-shuffling algorithm
### Goal: put students in pods so that each group (ex. Spanish, Graham, Band) is represented in each pod as closely to the percent that it is represented in the population

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

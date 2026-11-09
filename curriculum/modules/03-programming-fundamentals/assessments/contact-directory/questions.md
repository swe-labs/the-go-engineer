# Questions

Answer each question with concrete reasoning and evidence from your Contact Directory implementation.

## Question 1

You implemented a directory using `map[string]Contact`. Why did you choose the phone number as the map key? What would change if you used the email as the key instead?

## Question 2

Your `ListContacts` function returns contacts sorted by name. How did you implement case-insensitive sorting? What happens if two contacts have names that differ only in case (e.g., "alice" and "Alice")?

## Question 3

The `NewContact` function trims whitespace from all fields and returns an error if any field is empty. Why is trimming important? Can you think of a scenario where a user might accidentally include leading/trailing whitespace?

## Question 4

`SearchContacts` searches both name and email fields. Why might searching only one field be insufficient? What performance trade-off does searching both fields introduce for a directory with 10,000 contacts?

## Question 5

Explain the difference between `AddContact` rejecting a duplicate phone vs. silently overwriting it. Which behaviour is safer and why? What if the user's intent was to update an existing contact?

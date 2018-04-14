# Snipper

Snipper is a proof of concept for applying Lift Snippet semantics to YAML
templating.

Snipper takes in a template and then any number of transformers to create a
final YAML file. YAML transformers are other YAML files whose keys conform to
selector syntax and values indicate what input to the provided selector should
be used.

## Getting started

Take the following example:

```yaml
spec:
  name: awesome-dude
  annotations:
    name: awesome-dude
    team: awesome-dude
  metadata:
    name: awesome-dude
    team: awesome-dude
  alerters:
  - test
```

Someone obviously got lazy with copy and paste when writing this YAML. As usual,
copy and paste works fine when you must, but anyone with a half decent
templating engine will avoid it when possible. So, let's use Snipper to solve a
very silly problem: let's change all of these to their correct values using
Snipper transformers.

The first thing we want is to change the name at the top level. To do so we
could define a name transformer like so:

```yaml
---
transformer:
- 'spec.name': very-important-production-application
```

Then run it with snipper.

```
$ snipper template.yaml name_transformer.yaml
```

And see that, indeed, `spec.name` has changed, but the other names remain
unaffected. However, we can use selector syntax to get the result we want. Let's
add another line to the transformer:

```yaml
---
transformer:
- 'spec.name': very-important-production-application
- 'spec.*.name': very-important-production-application
```

This tells snipper to apply to the top level name, _and_ to any name below a
direct child of `spec`. But what if we could combine this into a single rule?
Well, it turns out we can.

```yaml
transformer:
- 'spec.**.name': very-important-production-application
```

This will cause snipper to change the value of all children of spec that have
the key `name` to our desired value.

Now, what if we wanted to define a separate transformer for teams and alerting?
Well we might define a `teams_transformer.yaml` like so:

```yaml
transformer:
- 'spec.**.team': Data Science
  'alerters':
  - pagerduty
```

Then we could apply _both_ of our transformers from the snipper command line:

```
$ snipper template.yaml name_transformer.yaml team_transformer.yaml
```

That's all well and good, but what if I need to add extra alerters that are
unrelated to the team as well? That's easy. Selector notation also supports
operations that change the transformation from a simple replace to support
append or prepend operations.

```yaml
transformer:
- 'alerters+':
  - systems_dashboard
```

When added to the existing transformers this results in the following alerters:

```yaml
alerters:
- pagerduty
- systems_dashboard
```

You could have alernately made that a prepend operation by writing `+alerters`
as your selector.

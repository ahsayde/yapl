# Statements

## Match

The `match` statement provides a way to filter resources which will be validated by the policy.

It can be a [condition](#condition) or a [logical condition](#logical-condition).


#### Examples

- `match` statement with single condition:

```yaml
match:
 field: kind
 operator: equal
 value: Deployment
```

- `match` statement with logical conditions:

```yaml
match:
  and:
  - field: kind
    operator: equal
    value: Deployment
  - field: metadata.namespace
    operator: equal
    value: default
```

## Exclude

The `exclude` statement is used to exclude resources from being validated by the policy.

It can be a [condition](#condition) or a [logical conditions](#logical-conditions).


#### Examples

- `exclude` statement with single condition:

```yaml
exclude:
 field: metadata.namespace
 operator: equal
 value: kube-system
```

- `exclude` statement with logical conditions:

```yaml
exclude:
  or:
  - field: kind
    operator: equal
    value: Secret
  - field: metadata.namespace
    operator: equal
    value: kube-system
```

## Rules

A `yapl` policy can contain one or more [rules](#rules). Each rule consist of a [condition](#condition) and a `result`


```yaml
rules:
- condition:
    field: metadata.name
    equal: hasPrefix
    value: app
  result:
    msg: resource name must starts with 'app'
```

### Conditional Rules

to add a condition when a rule is evaluated you can use when field to define a condition when a rule can be evaluated

```yaml
rules:
- when:
    < condition or logical condition >
  condition:
    < condition >
  result:
    < result >
```

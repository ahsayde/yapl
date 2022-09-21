
# Summary
- [Structure](#structure)
- [Concepts](#concepts)
  - [Condition](#condition)
  - [Logical Condition](#logical-condition)
  - [Expressions](#expressions)
- [Statements](#statements)
    - [Match](#match)
    - [Exclude](#exclude)
    - [Rules](#rules)
      - [Conditional Rules](#conditional-rules)
- [Operators](#operators)
- [Context](#context)
- [Functions](#functions)

# Structure

 ![Structure](structure.png)

# Concepts

## Condition

## Logical Condition

Logical condition combine the result of multiple [conditions](#condition) to produce a single result.

Logical operators `and`, `or` and `not` are used to define the relationship of conditions. Logical condition can has multiple levels.

### Examples

```yaml
  and:
  - < condition >
  - < condition >
```

```yaml
  or:
  - < condition >
  - < condition >
```

```yaml
  not:
    < condition >
```

```yaml
  and:
  - < condition >
  - < condition >
  - or:
    - < condition >
    - < condition >
    - not:
        < condition >
```

## Expressions


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

## Operators

Here is the list of the availabe operators which can be used in [conditions](#condition).

| Operator    | Alias    | Description                                                               | Field Value | Operator Value |
|-------------|----------|---------------------------------------------------------------------------|-------------|----------------|
| `equal`     | `eq`     | Checks field's value equal provided value                                 | `any`       | `any`          |
| `hasPrefix` |          | Checks whether field's value begins with prefix                           | `string`    | `string`       |
| `hasSuffix` |          | Checks whether field's value ends with suffix                             | `string`    | `string`       |
| `regex`     |          | Checks whether field's value matches the provided regex                   | `string`    | `string`       |
| `minValue`  | `min`    | Checks whether field's value is greater than or equals provided value     | `number`    | `number`       |
| `maxValue`  | `max`    | Checks whether field's value is less than or equals the provided value    | `number`    | `number`       |
| `in`        |          | Checks whether field's value in the provided value                        | `any`       | `array`        |
| `contains`  |          | Checks whether field's value contains the provided value                  | `array`     | `any`          |
| `length`    | `len`    | Checks whether field's value length equals the provided value             | `array`     | `integer`      |
| `minLength` | `minlen` | Checks whether field's value has minimum length equals the provided value | `array`     | `integer`      |
| `maxLength` | `maxlen` | Checks whether field's value has maximum length equals the provided value | `array`     | `integer`      |


## Context

Contexts are a way to access information about current rule

### Input

Input contex allow you to access the resource object using ```${ .Input.< field json path > }``` expression.

#### Example

```yaml
rules:
- condition:
    field: metadata.name
    operator: hasPrefix
    value: app
  result:
    msg: resource name ${ .Input.metadata.name } must start with prefix 'app'
```

### Parameters

You can access parameters passed during the evaluation of input using expression`${ .Params.<variable name> }`.

#### Example

```yaml
exclude:
  field: metadata.namespace
  operator: in
  value: ${ .Params.excluded_namespaces }
```

### Environment Variables 

You can access environment variable value by using expression `${ .Env.<variable name> }`

#### Example


```yaml
rules:
- condition:
    field: request.body
    operator: maxLength
    value: ${ .Env.MAX_BODY_SIZE }
  result: request body must not exceed ${ .Env.MAX_BODY_SIZE } 
```

### Current Condition

`Cond` context allow you to access all the information of the current condition.

> This context is only availabe on field `rules.result`.


| Key                              | Type              | Description                                           |
|----------------------------------|-------------------|-------------------------------------------------------|
| `Cond.Field.Value`               | `any`             | The value of the field                                |
| `Cond.Field.Index`               | `integer`         | The index of the field if field's parent is an array  |
| `Cond.Field.Parent`              | `[field object]`  | The parent of the field                               |
| `Cond.Operator`                  | `string`          | Condition's operator                                  |
| `Cond.Value`                     | `string`          | Condition's value                                     |


#### Examples

```yaml
rules:
- condition:
    field: metadata.name
    operator: hasPrefix
    value: app
  result:
    msg: resource name must starts with ${ .Cond.Value } but found ${ .Cond.Field.Value }
    key: ${ .Field.Path }
```

```yaml
rules:
- condition:
    field: spec.containers[*].image
    operator: hasPrefix
    value: my-registry
  result:
    msg: container ${ .Cond.Field.Parent.Value.name } uses invalid image registry 
```


## Functions

Here is the list of the builtin functions which can be used in [expressions](#expressions)

### String

- `split`
- `lower`
- `upper`
- `title`
- `join`
- `trim`
- `trimLeft`
- `trimRigh`
- `trimPrefix`
- `trimSuffix`
- `replace`
- `replaceAll`

### Math

- `round`
- `ceil`
- `abs`
- `floor`
- `max`
- `min`

### Date & Time

- `date`
- `now`
- `year`
- `month`
- `weekday`
- `day`
- `hour`
- `minute`
- `second`

### Type Conversion

- `bool`

<!-- ## Objects

### Condition

| Field         | Type     |  Notes   |         Description                                               |
|---------------|----------|----------|-------------------------------------------------------------------|
|  `field`      | `string` | optional | the json path of the field to check                               |
|  `expr`       | `string` | optional | expression to evaulate                                            |
|  `operator`   | `string` | required | condition's operator. Available operators are [here](#operators)  |
|  `value`      |  `any`   | required | value to compare field's value with                               |


### logical Condition

| Field       | Type                |  Notes    |         Description                                                          |
|-------------|---------------------------------|----------|-------------------------------------------------------------------|
|  `and`      | [`array[logical Condition]`]() | optional | the json path of the field to check                               |
|  `or`       | [`array[logical Condition]`]() | optional | expression to evaulate                                            |
|  `not`      | [`logical Condition`]()            | optional | condition's operator. Available operators are [here](#operators)  |


### Rule

| Field         | Type                  |  Notes   |         Description                                               |
|---------------|-----------------------|----------|-------------------------------------------------------------------|
|  `condition`  | [`Condition`]()       | required | the json path of the field to check                               |
|  `when`       | [`logical Condition`]()| optional | expression to evaulate                                            |
|  `result`     | `any`                 | required | condition's operator. Available operators are [here](#operators)  |

 -->

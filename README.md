
# YAPL
*YAML as Policy Language*

[![codecov](https://codecov.io/gh/ahsayde/yapl/branch/master/graph/badge.svg?token=N5CJSZBHNF)](https://codecov.io/gh/ahsayde/yapl)
![example workflow](https://github.com/ahsayde/yapl/actions/workflows/main.yml/badge.svg)


- [Syntax](#syntax)
- [Usage](#usage)
- [Contexts](#contexts)
- [Operators](#operators)


## Playground

You can play with `YAPL` with this [playground](https://ahsayde.github.io/yapl-playground/) tool.


## Syntax
 

### Selecting Resources

The `match` and `exclude` are optional declarations used to filter resources which will be validated by the rules.

They have the same syntax, and can be used together to include and exclude resources.

#### Examples

For example this code will only validate resources when `request.method` is equal `GET`

```yaml
match:
  field: request.method
  operator: equal
  value: GET
```

this code will validate all resources except when the `request.method` is equal `GET`

```yaml
exclude:
  field: request.method
  operator: equal
  value: GET
```

this code will validate all resources when `request.method` equals `GET` and `request.endpoint` not equals `/api`

```yaml
match:
  field: request.method
  operator: equal
  value: GET

exclude:
  field: request.endpoint
  operator: equal
  value: /api
```

the previous example can be done using only `match` statement as follow:

```yaml
match:
  and:
  - field: request.method
    operator: equal
    value: GET
  - not:
      field: request.endpoint
      operator: equal
      value: /api
```

### Validating Resources

Policies can have multiple rules. Each rule has its own condition and result object.

rule's condition doesn't support nested conditions like `match` and `exclude`, It's only one level

```yaml
rules:
- condition: 
    field: metadata.name
    equal: hasPrefix
    value: app
  result: container name must starts with app
```

the `result` field could be an object, see this example

```yaml
rules:
- condition: 
    field: metadata.name
    equal: hasPrefix
    value: app
  result:
    msg: container name must starts with app
```

#### Conditional Rule

You can add conditional rule by adding field `when` which implements the `condition` object

In this example, instead for writing two policies for each container type. Ypu can add two conditional rules

```yaml
rules:
- when:
    field: image
    equal: equal
    value: my-app-image
  condition: 
    field: metadata.name
    equal: hasPrefix
    value: app
  result: app container's name must starts with app

- when:
    field: image
    equal: equal
    value: my-db-image
  condition: 
    field: metadata.name
    equal: hasPrefix
    value: db
  result: db container's name must starts with db
```

### Context

Contexts are a way to access information about current rule

### Available Contexts

| Name                           | Type     | Description                           |         Availability            |
|--------------|-----------------|--------------------------------------------------|---------------------------------|
| [`Input`](#input-context)       | `object` | The input to be checked              |          all fields             |
| [`Params`](#params-context)     | `object` | Parameters passed during evaulation  |          all fields             |
| [`Env`](#env-context)           | `object` | exported environment variable        |          all fields             |
| [`Cond`](#condition-context)    | `object` | Current condition information,       |  only on `rules.result` field  |


#### `Input` Context

Input contex allow you to access the resource object

```yaml
rules:
- condition:
    field: user.role
    operator: equal
    value: admin
  result: user ${ .Input.user.name } doesn't has access
```

#### `Params` Context

You can access parameters passed during the evaluation of input using `${ .Params.<variable name> }` expression. For example

```yaml
rules:
- condition:
    field: request.body
    operator: maxLength
    value: ${ .Params.max_body_size }
  result: request body must not exceed ${ .Params.max_body_size } 
```

#### `Env` context

You can access any environment variable value by using expression `${ .Env.<variable name> }`

```yaml
rules:
- condition:
    field: request.body
    operator: maxLength
    value: ${ .Env.MAX_BODY_SIZE }
  result: request body must not exceed ${ .Env.MAX_BODY_SIZE } 
```

#### `Cond` context

`Cond` context allow you to access all the information of the current field.

This context is only availabe on field `rules.result`.


| Key                              | Type              | Description                                           |
|----------------------------------|-------------------|-------------------------------------------------------|
| `Cond.Field.Value`               | `any`             | The value of the field                                |
| `Cond.Field.Index`               | `integer`         | The index of the field if field's parent is an array  |
| `Cond.Field.Parent`              | `[field object]`  | The parent of the field                               |
| `Cond.Operator`                  | `string`          | Condition's operator                                  |
| `Cond.Value`                     | `string`          | Condition's value                                     |

Example

```yaml
rules:
  - condition:
      field: metadata.labels.app
      operator: hasPrefix
      value: app
    result:
      msg: resource app label must starts with ${ .Cond.Value } but found ${ .Cond.Field.Value }
      key: ${ .Field.Path }
```



## Usage

```go
input := map[string]interface{}{
  "users": []interface{}{
    map[string]interface{}{
      "name": "bob",
      "role": "member",
    },
  },
}

params := map[string]interface{}{
  "role": "admin",
}

raw, err := ioutil.ReadFile("policy.yaml")
	if err != nil {
		panic(err)
	}

policy, err := yapl.Parse(raw)
if err != nil {
  panic(err)
}

result, err := policy.Eval(input, params)
if err != nil {
  panic(err)
}

```


## Operators

| Operator    | Alias    | Description                                                               | Field Value | Operator Value |
|-------------|----------|---------------------------------------------------------------------------|-------------|----------------|
| `exists`    |          | checks whether field is exist or not                                      | `any`       |                |
| `equal`     | `eq`     | asserts field's value equal provided value                                | `any`       | `any`          |
| `hasPrefix` |          | checks whether field's value begins with prefix                           | `string`    | `string`       |
| `hasSuffix` |          | checks whether field's value ends with suffix                             | `string`    | `string`       |
| `regex`     |          | checks whether field's value matches the provided regex                   | `string`    | `string`       |
| `minValue`  | `min`    | checks whether field's value is greater than or equals provided value     | `number`    | `number`       |
| `maxValue`  | `max`    | checks whether field's value is less than or equals the provided value    | `number`    | `number`       |
| `in`        |          | checks whether field's value in the provided value                        | `any`       | `array`        |
| `contains`  |          | checks whether field's value contains the provided value                  | `array`     | `any`          |
| `length`    | `len`    | checks whether field's value length equals the provided value             | `array`     | `integer`      |
| `minLength` | `minlen` | checks whether field's value has minimum length equals the provided value | `array`     | `integer`      |
| `maxLength` | `maxlen` | checks whether field's value has maximum length equals the provided value | `array`     | `integer`      |

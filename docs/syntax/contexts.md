
# Context

Contexts are a way to access information about policy runtime

### Available Contexts
- [`.Input`](#input)
- [`.Params`](#parameters)
- [`.Env`](#environment-variables)
- [`.Cond`](#current-condition)

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


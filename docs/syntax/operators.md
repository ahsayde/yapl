# Operators

## Built-in Operators

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


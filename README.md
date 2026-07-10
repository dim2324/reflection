## 🌟 Ключевые особенности

- 🧠 **Рефлексия** — обход полей структуры через `reflect.Type` и `reflect.Value`
- 🏷️ **Теги структуры** — извлечение правил через `tag.Get("validate")`
- 📐 **Множественные правила** — разбор через `strings.Split` с разделителями `;` и `=`
- 🔤 **Поддержка Unicode** — корректная работа с кириллицей, иероглифами и эмодзи через конвертацию в руны
- 📊 **Разные типы данных** — строки, целые числа (`int`, `int64`), числа с плавающей точкой (`float64`)
- ❌ **Fail-fast подход** — возврат ошибки при первом же несоответствии
- 🧩 **Регулярные выражения** — проверка паттернов через стандартную библиотеку `regexp`

Ожидаемый вывод программы:

Validation error: field Name length must be at least 3 characters
Validation error: field Age must not exceed 65
Validation error: field Email does not match pattern ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$
Validation passed successfully
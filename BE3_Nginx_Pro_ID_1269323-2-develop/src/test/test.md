# Проверка балансировки

После применения конфига Nginx с заголовком `X-Upstream-Addr` открой
`https://be3.shop`, затем инструменты разработчика (F12) → Console.
Выполни скрипт:

```js
const counts = {};

for (let i = 0; i < 40; i++) {
  const response = await fetch('/api/v1/products');
  const backend = response.headers.get('X-Upstream-Addr');
  counts[backend] = (counts[backend] || 0) + 1;
}

console.table(counts);
```

Ожидаемый результат: три адреса бэкендов с распределением примерно
`20:10:10`, соответствующим весам `2:1:1`.

Запросы выполняются последовательно. Без токена API возвращает `401`;
эти ответы не кешируются и позволяют проверить выбор бэкенда для каждого запроса.
Если вместо адреса появляется `null`, проверь, применён ли конфиг с заголовком
`X-Upstream-Addr`.

# Проверка кеширования

После применения конфига с заголовком `X-Cache-Status` получи действующий
`access_token` через Swagger. На странице `https://be3.shop` открой F12 → Console
и выполни скрипт. В появившееся поле вставь токен без префикса `Bearer`.

```js
{
  const token = prompt('Вставь access_token без Bearer');
  const url = `/api/v1/products?demo=${crypto.randomUUID()}`;
  const results = [];

  for (let i = 1; i <= 2; i++) {
    const response = await fetch(url, {
      headers: { Authorization: `Bearer ${token}` },
      cache: 'no-store'
    });
    await response.text();

    results.push({
      запрос: i,
      HTTP: response.status,
      кеш: response.headers.get('X-Cache-Status'),
      бэкенд: response.headers.get('X-Upstream-Addr')
    });
  }

  console.table(results);
}
```

Ожидаемый результат:

| Запрос | HTTP | Кеш | Бэкенд |
|---|---|---|---|
| 1 | 200 | MISS | Адрес сервера |
| 2 | 200 | HIT | `null` — обращения к бэкенду не было |

Уникальный параметр `demo` позволяет повторять демонстрацию с новым ключом кеша.
`cache: 'no-store'` отключает кеш браузера; ответы сохраняет кеш Nginx.
Запросы нужно выполнить в пределах минуты — срока свежести ответа в кеше.
При `401` проверь токен: кешируются только успешные ответы `200`.

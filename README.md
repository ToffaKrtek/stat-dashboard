# stat-dashboard
Сборщик метрик сервера с постройкой графиков на golang

### DAEMON
```
stat-dashboard -server=server_name # Название отображается в графике как заголовок (${server}-${current_date})
```
```
stat-dashboard -period=10 # Периодичность запуска сборщика метрик (в минутах)
```
```
stat-dashboard -repeat=week # Периодичность цикла "сбор"-"создание графика"-"обнуление данных" [day, week, hour, infinity, debug]
  # debug - 1 минута и period каждые 2 секунды
  # infinity - без создания графика и обнуления данных period каждые 2 секунды
```
```
stat-dashboard -cli #отображение графика в терминале в реальном времени (TODO::вынести в cli клиент)
```

# Simple fast fuzzy

Fuzzy library written in go.

This library is a simple fuzzy search with unicode normalization and an arbitrary score system.

## Test the library

```bash
git clone https://github.com/eregnier/sffuzzy
cd ssfuzzy
make test-trace
```

## Usage

`go get github.com/eregnier/sffuzzy`

Usage samples are in [test.go](test.go)

A minimal usage code below:

```go
  //One shot search
  names := []string{"super man", "super noel", "super du"}
  results := fuzzy.SearchOnce("perdu", &names, fuzzy.Options{Sort: true, Limit: 5, Normalize: true})
```

```go
  //Use search cache for performance
  names := []string{"super man", "super noel", "super du"}
  options := fuzzy.Options{Sort: true, Limit: 5, Normalize: true}
  cacheTargets := fuzzy.Prepare(&names, options)
  results := fuzzy.Search("perdu", cacheTargets, options)
```

## Options

```go
  options := fuzzy.Options{Sort: true, Normalize: true, Limit: 10}
```

This options structure have the following options

**Prop**|**Type**|**Description**
:-----:|:-----:|:-----:
`Sort`|bool|Orders result depending on results score
`Normalize`|bool|Handles searches in texts with special characters. Make search more flexible / less strict
`Limit`|int|Define how many results are kept un search return value

## Performances

The given sample file is a flat csv loaded as string list.

Multi thread performances are worse with my code, so I reverted

On a AMD 3600x and on a single core I fuzzy search a text in 40 ms in `SearchOnce` mode

When I build cache with a `Prepare(data *[]string)`, and then I run a `Search`, the prepare takes about 36ms and the Search about 4ms on the data sample of ~155K lines for 320Ko.

## Execution

The following code [test.go](test.go)

Have the following output

```bash
2020/04/29 02:07:57 TestMinimalSearch &{[{super du 8 5 1} {super man 3 3 4} {super noel 3 3 5}] 8}
2020/04/29 02:07:57 TestMinimalSearchCache &{[{super du 8 5 1} {super man 3 3 4} {super noel 3 3 5}] 8}
2020/04/29 02:07:57  + Cache search, first search is slower.
2020/04/29 02:07:57  🕑 Duration: 9.435114ms
2020/04/29 02:07:57  + Cached searches
2020/04/29 02:07:57  🕑 Duration: 3.883137ms
2020/04/29 02:07:57 [{San Francisco;United States 13 10 16} {South San Francisco;United States 12 10 22} {St. Francis;United States 12 10 19}]
2020/04/29 02:07:57  🕑 Duration: 2.100023ms
2020/04/29 02:07:57 [{Mumbai;India 11 6 0} {Mumbwa;Zambia 8 6 6} {Mount Gambier;Australia 8 6 16}]
2020/04/29 02:07:57  🕑 Duration: 3.992907ms
2020/04/29 02:07:57 [{Hong Kong;Hong Kong 16 8 0} {Xiangkhoang;Laos 13 8 3} {Mokhotlong;Lesotho 12 8 7}]
2020/04/29 02:07:57  🕑 Duration: 2.312343ms
2020/04/29 02:07:57 [{Agadez;Niger 11 6 0} {Várzea Grande;Brazil 8 6 7} {Altagracia de Orituco;Venezuela 8 6 21}]
2020/04/29 02:07:57  🕑 Duration: 2.034594ms
2020/04/29 02:07:57 [{Palmas;Brazil 10 5 0} {La Palma;Panama 10 5 0} {Las Palmas de Gran Canaria;Spain 10 5 0}]
2020/04/29 02:07:57  🕑 Duration: 4.029126ms
2020/04/29 02:07:57 [{Sucre;Bolivia 20 12 0} {Quime;Bolivia 13 7 0} {Villazón;Bolivia 13 7 0}]
2020/04/29 02:07:57  🕑 Duration: 3.910717ms
2020/04/29 02:07:57 [{Ibb;Yemen 16 8 0} {Ibb;Yemen 16 8 0} {Dhamār;Yemen 11 5 0}]
2020/04/29 02:07:57  🕑 Duration: 3.826816ms
2020/04/29 02:07:57 [{West View;United States 16 8 0} {Westview;United States 16 8 0} {Viera West;United States 14 8 3}]
2020/04/29 02:07:57  + Search all at once
2020/04/29 02:07:57  🕑 Duration: 44.714400ms
2020/04/29 02:07:57 Print plain unmarshaled json results
2020/04/29 02:07:57 [
  {
    "target": "Ōsaka;Japan",
    "score": 13,
    "matchCount": 10,
    "typos": 1
  },
  {
    "target": "Northwest Harborcreek;United States",
    "score": 5,
    "matchCount": 5,
    "typos": 29
  },
  {
    "target": "Oshakati;Namibia",
    "score": 5,
    "matchCount": 5,
    "typos": 11
  },
  {
    "target": "Colombo;Sri Lanka",
    "score": 5,
    "matchCount": 5,
    "typos": 11
  },
  {
    "target": "Coxsackie;United States",
    "score": 5,
    "matchCount": 5,
    "typos": 17
  }
]
PASS
ok  	_/home/utopman/sources/sffuzzy	0.083s
```

## Licence

[MIT](LICENCE.md)
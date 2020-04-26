# Simple fast fuzzy

Fuzzy library written in go.

This library is a simple fuzzy search with unicode normalization and an arbitrary score system.

## Test the library

```bash
git clone https://github.com/eregnier/simple-fast-fuzzy
cd simple-fast-fuzzy
go get
make
```

## Usage

Usage samples are in [test.go](test.go)

## Performances

The given sample file is a flat csv loaded as string list.

Multi thread performances are worse with my code, so I reverted

On a AMD 3600x and on a single core I fuzzy search a text in 40 ms in `SearchOnce` mode

When I build cache with a `Prepare(data *[]string)`, and then I run a `Search`, the prepare takes about 36ms and the Search about 4ms on the data sample of ~155K lines for 320Ko.

## Execution

The following code [test.go](test.go)

Have the following output

```bash
duration ms> 44
[
  {
    "target": "Ōsaka;Japan",
    "score": 10,
    "matchCount": 10,
    "typos": 1,
    "complete": true
  },
  {
    "target": "Yuzhno-Sakhalinsk;Russia",
    "score": 5,
    "matchCount": 5,
    "typos": 5,
    "complete": false
  },
  {
    "target": "Oshakati;Namibia",
    "score": 5,
    "matchCount": 5,
    "typos": 5,
    "complete": false
  },
  {
    "target": "Makedonska Kamenica;Macedonia",
    "score": 5,
    "matchCount": 5,
    "typos": 5,
    "complete": false
  },
  {
    "target": "Zhosaly;Kazakhstan",
    "score": 5,
    "matchCount": 5,
    "typos": 5,
    "complete": false
  },
  {
    "target": "Osakarovka;Kazakhstan",
    "score": 5,
    "matchCount": 5,
    "typos": 5,
    "complete": false
  },
  {
    "target": "Mombasa;Kenya",
    "score": 4,
    "matchCount": 4,
    "typos": 5,
    "complete": false
  },
  {
    "target": "Korsakov;Russia",
    "score": 4,
    "matchCount": 4,
    "typos": 5,
    "complete": false
  },
  {
    "target": "Moses Lake;United States",
    "score": 4,
    "matchCount": 4,
    "typos": 5,
    "complete": false
  },
  {
    "target": "P’yŏngsan;Korea, North",
    "score": 4,
    "matchCount": 4,
    "typos": 5,
    "complete": false
  }
]

 + Perform cache search, first search is slower.
duration ms> 11
 + Perform cached searches
duration ms> 2
[{San Francisco;Argentina 9 9 5 false} {San Francisco de Macorís;Dominican Republic 9 9 5 false} {San Francisco;United States 9 9 5 false} {San Francisco;El Salvador 9 9 5 false} {San Fernando;Philippines 8 8 5 false}]
duration ms> 1
[{Mumbai;India 7 6 0 true} {Mayumba;Gabon 5 5 5 false} {Mumbwa;Zambia 5 5 5 false} {Capenda Camulemba;Angola 5 5 5 false} {Namutumba;Uganda 5 5 5 false}]

```

## Licence

[MIT](LICENCE.md)
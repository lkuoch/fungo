# FUNGO

Fungo is a high level language interpreter written in go

Some examples

```
let map = fn(arr, f) {
	let iter = fn(arr, accumulated) {
		if (len(arr) == 0) {
			accumulated
		} else {
			iter(rest(arr), push(accumulated, f(first(arr))))
		}
	};

	iter(arr, []);
};

let a = [1, 2, 3, 4];
let double = fn(x) { x * 2};

map(a, double)
// ^ Prints: [2, 4, 6, 8]

let reduce = fn(arr, initial, f) {
	let iter = fn(arr, result) {
		if (len(arr) == 0) {
			result
		}
	} else {
		iter(rest(arr), f(result, first(arr)));
	}
};

let sum = fn(arr) {
	reduce(arr, 0, fn(initial, el) { initial + el } )
}

sum([1, 2, 3, 3, 5, 6])
// ^ Prints: 15
```

// Package bitset implements bitsets. A bitset is a set of non-negative
// integers represented using a bit for each integer.
//
// There are three implementations:
//
// Use Dense for the common case where you know the largest value that could be in the set,
// and you want to use a sequence of bits. Addition, removal and membership tests on a
// Dense bitset are very fast, and memory is proportional to the largest possible value
// (one bit per possible value).
//
// Use Sparse for bitsets whose values can come from a wide range. Sparse
// bitsets take more time per operation, but can use less memory than a Dense
// bitset if the set contains relatively few elements drawn for a large range. For example,
// A Dense bitset with the elements 1000, 2000, ...., 1_000_000 would occupy 125K, while
// a Sparse one would take
// XXXXXXXXXXXXXXXX TODO
//
// Set64 is a faster version of Dense when the largest possible value is 63.
package bitset

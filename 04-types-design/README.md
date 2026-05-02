# Section 04: Types and Design

This section covers the Go type system, focusing on data structure definition, behavioral abstraction, and logic composition. Learners will transition from procedural function-based programming to domain-oriented design using structs, methods, and interfaces.

## Technical Objectives

- **Data Modeling**: Defining complex state using named structs and understanding memory alignment.
- **Behavioral Attachment**: Implementing methods with value and pointer receivers to encapsulate logic.
- **Abstraction**: Using interfaces to define behavioral contracts and achieve polymorphism.
- **Composition**: Utilizing struct and interface embedding as the primary mechanism for logic reuse.
- **Text Processing**: Understanding internal string representations (runes/bytes) and text manipulation.

## Zero-Magic Machine Boundary

In Section 04, we establish the following technical constraints:

- **Structs** are contiguous memory blocks; "fields" are just offsets from a base pointer.
- **Methods** are syntactic sugar for functions where the "receiver" is passed as the first parameter.
- **Interfaces** are two-word headers (itab + data pointer); method calls involve dynamic dispatch via the itab.
- **Composition** is field delegation at the compiler level, not runtime method lookup.

## Curriculum Map

### Track 1: TI (Types & Interfaces)
Foundational lessons on type definition and polymorphism.
- [TI.1] `1-struct`: Contiguous data grouping and memory layout.
- [TI.2] `2-methods`: Encapsulating behavior via receivers.
- [TI.7] `7-receiver-sets`: Method set rules and interface compatibility.
- [TI.3] `3-interfaces`: Abstract behavioral contracts.
- [TI.4] `4-interface-embedding`: Composing contracts from smaller interfaces.
- [TI.6] `6-type-switch`: Runtime type inspection and assertions.
- [TI.11] `11-dynamic-typing-with-any`: Working with the empty interface.
- [TI.5] `5-stringer`: Implementing the standard string conversion interface.
- [TI.8] `8-custom-errors`: Defining domain-specific error types.
- [TI.9] `9-generics`: Type-safe generic programming.
- [TI.14] `14-complex-generic-constraints`: Advanced generic bounds.
- [TI.15] `15-generic-data-structures`: Implementing generic containers.
- [TI.12] `12-functional-options`: The functional options pattern for configuration.
- [TI.13] `13-method-values`: Using methods as first-class functions.
- [TI.10] `10-payroll-processor`: Capstone project for Track 1.

### Track 2: CO (Composition)

Logic sharing through structural delegation and method promotion.

- [CO.1] [16-composition](./16-composition/README.md): Modeling "has-a" relationships via named fields.
- [CO.2] [17-embedding](./17-embedding/README.md): Anonymous fields and method set promotion.
- [CO.3] [18-bank-account-project](./18-bank-account-project/README.md): Case study in shadowing and promotion.

### Track 3: ST (Strings & Text)

Internal representation and high-level manipulation of text data.

- [ST.1] [19-strings](./19-strings/README.md): Immutability, byte-slice representation, and 'strings.Builder'.
- [ST.2] [20-formatting](./20-formatting/README.md): Template-driven output and error wrapping with 'fmt'.
- [ST.3] [21-unicode](./21-unicode/README.md): UTF-8 encoding, runes, and safe character iteration.
- [ST.4] [22-regex](./22-regex/README.md): Pattern matching and extraction via the RE2 engine.
- [ST.5] [23-text-template](./23-text-template/README.md): Decoupling presentation from data logic.
- [ST.6] [24-config-parser-project](./24-config-parser-project/README.md): Integrating text tools into a production configuration pipeline.

## Prerequisites

- `FE.1-10` Functions & Errors track.
- Basic understanding of memory addressability (Section 00).

## Next Step

After completing this section, proceed to [Section 05: Packages and I/O](../05-packages-io/README.md) to apply these types within library boundaries and external system interactions.

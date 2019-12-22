--- title
go string rune byte.md
--- description
what's the difference anyways
--- main


how do you convert the idea of a **symbol**
to something the computer can represent?

## def

as humans see it

- string: _immutable_ sequence of 0+ character
- string literal: _immutable_ **utf8** sequence of 0+ characters
- character: valid sequence of 1+ runes **(code points)**
- rune: valid sequence of 1+ bytes
- byte: 8 bits

### examples

the process of turning ideas into bits and bytes a computer can read

- _ab_ is:
  - 1 string **length 2**
  - 2 characters
  - 2 runes
  - 2 bytes
- _日本_ is:
  - 1 string **length 6**
  - 2 character
  - 2 runes
  - 6 bytes
- _👩‍🏫_ is:
  - 1 string **length 11**
  - 1 character
  - 3 runes
  - 11 bytes

# Go ethereum key finder

Number of possibles public address:
- 40 hex character = 16 ^ 40 = 1.5e+48

Number of unique active address: 5e+7
> https://etherscan.io/chart/address

_Chance to randomly generate a key matching an existing wallet 1.5e+48 / 5e+7 = 3e+40_

The current performance of the program is 1billion key per minute: 1e+9

>  nb_possible_public_address / nb_unique_active_address / nb_generated_per_minute = time_to_break_one_key_in_minutes

It gives us 1.5e+48 / 5e+7 / 1e+9 = 3e+31

Per month: 3e+31 / 4.3e+4 = 7e+26
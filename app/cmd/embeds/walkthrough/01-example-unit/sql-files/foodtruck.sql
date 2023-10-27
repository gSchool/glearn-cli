-- create a database called 'food_trucks' then run:
-- psql -d food_trucks -f foodtruck.sql

-- QUESTIONS

-- Select all users please

-- Select all ids and avatars of users

-- Select users who were referred

-- Select the users who referred other users

-- Select truck names with their user's username

-- Select user emails with the count of the number of trucks they own
-- https://www.postgresql.org/docs/current/tutorial-agg.html
-- https://www.postgresql.org/docs/current/functions-aggregate.html

-- Select users who do not own trucks

-- Select truck names with each of ther menu item names and prices ordered by price from cheapest to most expensive
-- display the prices as dollar values with decimal places https://www.postgresql.org/docs/current/typeconv-func.html
-- https://www.postgresql.org/docs/current/queries-order.html

-- Select users first and last names with the average cost of their menu items, ordered most expensive to least

-- Select the truck with the most expensive burger
-- https://www.postgresql.org/docs/current/queries-limit.html

-- Select the truck with the cheapest burger

-- Select the count of trucks in each food truck category available

-- Select content of reviews which contain content with the word 'good' without case sensitivity
-- https://www.postgresql.org/docs/current/functions-matching.html

-- select the food truck name with their average review score, ordered by average desc, and rounded to two decimals
-- https://www.postgresql.org/docs/current/functions-math.html

-- One of the truck owners might be deliberately trying to tank other truck owner's review ratings!
-- select users with their id and first name and their average review score rounded to two decimals and the count of their reviews
-- is anyone trying to game the system unfairly?

-- select the food truck name with their average review score, only exclude scores given by people trying to cheat ratings
-- if a user has given 5 or more reviews with an average review score of 1, do not count their review ratings

DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS truck_menu_items;
DROP TABLE IF EXISTS menu_items;
DROP TABLE IF EXISTS trucks;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id INTEGER NOT NULL UNIQUE,
    first TEXT NOT NULL,
    last TEXT NOT NULL,
    avatar TEXT NOT NULL,
    email TEXT NOT NULL,
    username TEXT NOT NULL,
    referrer_id INTEGER REFERENCES users("id") ON DELETE CASCADE
);

CREATE TABLE trucks (
    id INTEGER NOT NULL UNIQUE,
    name TEXT NOT NULL,
    website TEXT NOT NULL,
    category TEXT NOT NULL check(category = 'American' or category = 'Asian' or category = 'French' or category = 'Mediterranean' or category = 'Indian' or category = 'Italian' or category = 'Latin'),
    vegetarian_friendly BOOLEAN NOT NULL,
    owner_id INTEGER NOT NULL REFERENCES users("id") ON DELETE CASCADE
);

CREATE TABLE menu_items (
    id INTEGER NOT NULL UNIQUE,
    name TEXT NOT NULL,
    calories INTEGER NOT NULL
);

CREATE TABLE reviews (
    id INTEGER NOT NULL UNIQUE,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    reviewer_id INTEGER REFERENCES users("id") ON DELETE CASCADE,
    rating INTEGER NOT NULL check(rating = 1 or rating = 2 or rating = 3 or rating = 4 or rating = 5),
    truck_id INTEGER NOT NULL REFERENCES trucks("id") ON DELETE CASCADE
);

CREATE TABLE truck_menu_items (
    truck_id INTEGER NOT NULL REFERENCES trucks("id") ON DELETE CASCADE,
    menu_item_id INTEGER NOT NULL REFERENCES menu_items("id") ON DELETE CASCADE,
    price INTEGER NOT NULL
);

INSERT INTO users VALUES
  (1, 'John', 'Smith', '1780000.jpeg', 'John@gmail.com', 'Jsmith', null),
  (2, 'Dave', 'Jones', '178030988.jpeg', 'Dave@gmail.com', 'Djones', null),
  (3, 'Patrick', 'Lacquer', '17800.jpeg', 'Patrick@gmail.com', 'Placquer', null),
  (4, 'Abbie', 'Schmabbie', '10000.jpeg', 'Abbie@gmail.com', 'Aschabbie', null),
  (5, 'David', 'Agarwal', '170000.jpeg', 'David@gmail.com', 'Dagarwal', null),
  (6, 'Susie', 'Chen', '60000.jpeg', 'Susie@gmail.com', 'Schen', null),
  (7, 'Matt', 'Gvido', '14780000.jpeg', 'Matt@gmail.com', 'Mvido', null),
  (8, 'Zuirch', 'Hernández', '1780320948.jpeg', 'Zurich@gmail.com', 'Zhernandez',5),
  (9, 'Will', 'Smith', '17800000089.jpeg', 'Will@gmail.com', 'Wsmith',4),
  (10, 'Yuri', 'Mikhailov', '17800000293908.jpeg', 'Yuri@gmail.com', 'Ymikhailov',3)
  ;

INSERT INTO trucks VALUES
  (1, 'Trucky', 'https://trucky.com', 'American', true, 1),
  (2, 'Necks', 'https://necks.com', 'Mediterranean', true, 2),
  (3, '100% Human Food', 'https://humanfood.com', 'French', false, 2),
  (4, 'Questionably Low Cost', 'https://epa.gov', 'Mediterranean', true, 3),
  (5, 'B-b-bistro', 'https://b-b-b-bitro.com', 'Italian', true, 3),
  (6, 'NaN Naan', 'https://nan-naan.net', 'Indian', true, 4),
  (7, 'White Cheese Queso', 'https://queso-white-cheese.com', 'Latin', false, 4),
  (8, 'Delicious Burgers', 'https://delicious-burgers.com', 'American', false, 4),
  (9, 'Quantifiable Flavor', 'https://delicious-burgers.com', 'American', false, 4),
  (10, 'Dumpling in French Onion Soup', 'https://dumplingsoup.com', 'French', true, 1),
  (11, 'Hello World', 'https://example.com', 'Italian', false, 1)
  ;

INSERT INTO menu_items VALUES
  (1, 'Burger', 300),
  (2, 'Fries', 150),
  (3, 'Milk Shake', 400),
  (4, 'Neck Meat', 130),
  (5, 'Vegetarian Neck Meat', 120),
  (6, 'Naan', 80),
  (7, 'Garlic Naan', 80),
  (8, 'Human Chow', 500),
  (9, 'Bachelor Chow', 500),
  (10, 'Bachelorette Chow', 500),
  (11, 'Cole Slaw', 100),
  (12, 'Raw Salmon', 305),
  (13, 'Poke Bowl', 505),
  (14, 'Simple Salad', 105),
  (15, 'Polymorphic Sandwich', 1005),
  (16, 'Hard Bread', 50),
  (17, 'Chips & Dip', 150),
  (18, 'Guac', 200),
  (19, 'Large Burrito', 1100),
  (20, '100 Units of Flavor', 100),
  (21, '200 Units of Flavor', 200),
  (22, '400 Units of Flavor', 400),
  (23, 'Dumplings', 400),
  (24, 'French Onion Soup', 300),
  (25, 'Baguette', 110),
  (26, 'NOT A NUMBER', 100),
  (27, 'Fish Flakes', 500)
  ;

INSERT INTO truck_menu_items VALUES
  (1, 1, 500), -- trucky, burger
  (1, 2, 300), -- trucky, fries
  (1, 9, 800), -- trucky, Bachelor Chow
  (2, 4, 550), -- Necks, neck meat
  (2, 5, 450), -- Necks, veggie neck meat
  (3, 8, 600), -- 100% human food, hc
  (3, 9, 600), -- 100% human food, bc
  (3, 10, 600), -- 100% human food, bettec
  (3, 1, 600), -- 100% human food, burger
  (4, 11, 100), -- low cost, slaw
  (4, 12, 100), -- low cost, salmon
  (4, 1, 100), -- low cost, burger
  (4, 16, 100), -- low cost, hard bread
  (5, 1, 900), -- bistro, burger
  (5, 13, 800), -- bistro, poke bowl
  (5, 14, 400), -- bistro, salad
  (5, 15, 1200), -- bistro, sandwich
  (6, 6, 400), -- naan, naan
  (6, 7, 400), -- naan, garlic nan
  (6, 26, 700), -- naan, NOT A NUMBER
  (7, 17, 400), -- queso, chips
  (7, 18, 400), -- queso, guac
  (7, 18, 900), -- queso, burrito
  (8, 1, 1200), -- dburger, burger
  (8, 2, 800), -- dburger, fries
  (9, 20, 800), -- qf, 100
  (9, 21, 800), -- qf, 200
  (9, 22, 800), -- qf, 400
  (9, 1, 800), -- qf, burger
  (10, 23, 400), -- french onion, dumpling
  (10, 24, 700), -- french onion, soup
  (10, 25, 200) -- french onion, baguette
  ;

INSERT INTO reviews VALUES
  (1, 'Delicious and debateably nutritious', 'I once visited this truck with elation at the idea of having such a sumptuous feast, only to find a puppy claw lingering in my coleslaw. Still ate it, it was good.', 5, 4, 5),
  (2, 'Dissatisfied', 'The owner of this food truck led me to believe that his macaroni noodles were gluten free. They were not, I can taste gluten like a snake can taste heat.', 2, 1, 4),
  (3, 'Sometimes I wonder...', 'You know, I just really dont get it, why is there a dumpling in this french onion soup? I went to this food truck with hopes to solve the mystery of why on earth anyone would ever do this. Im still perplexed. Its like putting a basketball player on the baseball field, entertaining to watch, but really not that good.', 5, 2, 10),
  (4, 'Xs and Os', 'Love love love! I wanted to eat a second helping just because the first melted on my tounge in the most delectable manner. I mean, it was melted cheese, so it was already kind of melted, but WOW.', 8, 5, 7),
  (5, 'Yuck', 'If flavor is actually quantifiable, this place gets a -10 on the richter scale. Ive had flour more flavorful than this. If I had to come down from my flavor tower one more time to eat this trash, I would probably skip the stairs and just jump', 2, 1, 9),
  (6, 'meh.jpg', 'I really wasnt sure what to expect out of this human food thing. Well, I cant say that I understand what I really ate. It had taste, and texture, and both were fine, but... what the hell was that? Would I eat it again? Yeah. Would I pay $15 for it again? Nah.', 6, 3, 3),
  (7, 'Terrible service', 'How hard is it to serve people burgers? I mean seriously. I am so indignant right now. I sat at the window while this bozo with his butt crack hanging out cooked burgers for like 2 whole minutes before acknowledging my presence. I am a paying customer, bozo.', 2, 1, 8),
  (8, 'BMTSSBMTSS', 'THE AMBIANCE OF THIS PLACE IS KILLER! SO UPBEAT! SO HIP! AND THE FOOD WAS GOOD TO BOOT! TOTALLY DIG THE EDM MUSIC PLAYING TOO! DEADMAU5 IS MY JAM!', 3, 4, 5),
  (9, '.', 'Bad', 2, 1, 1),
  (10, 'Unfortunate', 'Bad hombres over there.', 2, 1, 6),
  (11, 'NaaS', 'THAT NAAN TOTALL IS A NUMBER, I GIVE IT A 4 ON THE MY-TASTEBUD-SCALE, TASTEBUS WHICH I OWN AND ARE HUMAN', 7, 4, 6),
  (12, 'grumpy-cat.gif', 'I cant even', 2, 1, 10),
  (13, 'bad service, bad food', 'the title says it all', 2, 1, 7),
  (14, 'Pretty good I thik', 'enjoyed!', 6, 4, 7),
  (15, 'Delectable', 'rarely are the tastebuds titillated to such perfection. I have recommended to all my friends and I recommend you do the same', 10, 5, 7),
  (16, 'Desirable', 'Flat bread, fair prices, excellent puns; you cannot go wrong here. Every menu item is incredible. Once I ordered every item on the menu twice, so gOoD', 10, 5, 6),
  (18, 'Disgusting', 'I did not enjoy this at all, which is a shame because the service was quite good. The operator was informed and chatty, but my meal somehow came out cold!', 10, 2, 4),
  (19, 'Fine', 'this is fine', 9, 3, 4),
  (20, 'listen up punk', 'certain sql queries seem to use asyndetons while other rely on polysyndetons and I am not sure it is clear which is which plz help', 8, 2, 4),
  (21, 'good', 'Good', 7, 4, 5),
  (22, 'good', 'I found it GOOD, I am a human', 7, 4, 3),
  (23, 'good', '[ 170.555779] Call Trace:\n [ 170.558221] [<ffffffff816045b6>] dump_stack+0x19/0x1b\n [ 170.563346] [<ffffffff8106e29b>] warn_slowpath_common+0x6b/0xb0\n [ 170.569336] [<ffffffff8106e3ea>] warn_slowpath_null+0x1a/0x20\n [ 170.575153] [<ffffffff814a352d>] cpufreq_update_policy+0x1dd/0x1f0\n [ 170.581403] [<ffffffff814a3540>] ? cpufreq_update_policy+0x1f0/0x1f0\n [ 170.587827] [<ffffffff8136d40b>] cpufreq_set_cur_state.part.3+0x8c/0x95\n [ 170.594510] [<ffffffff8136d4b5>] processor_set_cur_state+0xa1/0xdb\n [ 170.600761] [<ffffffff8148a1e5>] thermal_cdev_update+0x95/0xb0\n [ 170.606664] [<ffffffff8148c849>] step_wise_throttle+0x59/0x90\n [ 170.612480] [<ffffffff8148aaeb>] handle_thermal_trip+0x5b/0x160', 7, 4, 11),
  (24, 'serious foodie', 'The exterior decor is very warm and welcoming, and the staff greeted me with aplomb', 5, 5, 6)
  ;

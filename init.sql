CREATE DATABASE cust_data
WITH 
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TEMPLATE = template0;


\connect cust_data;

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS "Users" (
  "user_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "username" varchar(255) NOT NULL,
  "email" varchar(255) NOT NULL,
  "password_hash" varchar(255) NOT NULL,
  "plaid_access_token" varchar(255) NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "Allocations" (
  "allocation_type" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "allocation_description" varchar(255) NOT NULL,
  "allocation_factor" decimal NOT NULL,
  "user_id" UUID NOT NULL REFERENCES "Users"("user_id")
);

CREATE TABLE IF NOT EXISTS "Category" (
  "category_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "plaid_category_primary_descriptor" varchar(255) NOT NULL,
  "plaid_category_detailed_descriptor" varchar(255) NOT NULL,
  "category_name" varchar(255) NOT NULL,
  "category_description" varchar(255) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "Expenses" (
  "expense_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "expense_description" varchar(255) NOT NULL,
  "expense_amount" decimal NOT NULL,
  "expense_category" UUID NOT NULL REFERENCES "Category"("category_id"),
  "user_id" UUID NOT NULL REFERENCES "Users"("user_id"),
  "allocation_type" UUID NOT NULL REFERENCES "Allocations"("allocation_type"),
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS "Income" (
  "income_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "income_description" varchar(255) NOT NULL,
  "income_amount" decimal NOT NULL,
  "income_frequency" varchar(255) NOT NULL,
  "user_id" UUID NOT NULL REFERENCES "Users"("user_id"),
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);




INSERT INTO "Users" ("user_id", "username", "email", "password_hash", "plaid_access_token") VALUES ('ed1bec4c-0a1b-4783-b47f-16ba0650b821', 'admin', 'admin@smartsplit.com', '$2a$10$nLavVuPde6DTLfHwkoxKkOOYfUt/QZrIg2Uq0W5HcyetavCl7ND12', 'access-sandbox-5423b0c9-2019-4f5e-bddd-2b41e52e5651'); --acess token user
--INSERT INTO "Users" ("user_id", "username", "email", "password_hash", "plaid_access_token") VALUES ('ed1bec4c-0a1b-4783-b47f-16ba0650b821', 'admin', 'admin@smartsplit.com', '$2a$10$nLavVuPde6DTLfHwkoxKkOOYfUt/QZrIg2Uq0W5HcyetavCl7ND12', ''); -- No access token user

INSERT INTO "Allocations" ("allocation_type", "allocation_description", "allocation_factor", "user_id") VALUES ('9f3c76e9-9d43-4480-a56d-a176b783f24d', 'Needs', 0.5, 'ed1bec4c-0a1b-4783-b47f-16ba0650b821');
INSERT INTO "Allocations" ("allocation_type", "allocation_description", "allocation_factor", "user_id") VALUES ('ac184cdf-b7ff-4eb9-b757-628770d566fb', 'Debts and Repayment', 0.1, 'ed1bec4c-0a1b-4783-b47f-16ba0650b821');
INSERT INTO "Allocations" ("allocation_type", "allocation_description", "allocation_factor", "user_id") VALUES ('f981f988-5be8-4a9b-bb39-392dd646ddbd', 'Wants', 0.3, 'ed1bec4c-0a1b-4783-b47f-16ba0650b821');
INSERT INTO "Allocations" ("allocation_type", "allocation_description", "allocation_factor", "user_id") VALUES ('184906a8-94f8-459e-b654-88e42d246579', 'Savings', 0.1, 'ed1bec4c-0a1b-4783-b47f-16ba0650b821');

INSERT INTO "Category" ("category_id", "plaid_category_primary_descriptor", "plaid_category_detailed_descriptor", "category_name","category_description") VALUES ('1ae53e57-8b82-45f2-a8cd-94d43932ab54', 'RENT_AND_UTILITIES', 'RENT_AND_UTILITIES_RENT', 'Rent','Payment, Rent');
INSERT INTO "Category" ("category_id", "plaid_category_primary_descriptor", "plaid_category_detailed_descriptor", "category_name","category_description") VALUES ('c2f89cdc-5ff7-46fc-92b3-f14bbdec7404', 'RENT_AND_UTILITIES', 'RENT_AND_UTILITIES_GAS_AND_ELECTRICITY', 'Utilities','Electric, Utilities');
INSERT INTO "Category" ("category_id", "plaid_category_primary_descriptor", "plaid_category_detailed_descriptor", "category_name","category_description") VALUES ('b0332b93-42a8-4eac-936f-ab8d859e9ce3', 'FOOD_AND_DRINK', 'FOOD_AND_DRINK_GROCERIES', 'Groceries','Groceries, Food');
INSERT INTO "Category" ("category_id", "plaid_category_primary_descriptor", "plaid_category_detailed_descriptor", "category_name","category_description") VALUES ('a9f5ae3e-d44e-49ac-82c5-1bc13c354f58', 'FOOD_AND_DRINK', 'FOOD_AND_DRINK_RESTAURANT', 'Restaurants','Restaurants, Dining');
INSERT INTO "Category" ("category_id", "plaid_category_primary_descriptor", "plaid_category_detailed_descriptor", "category_name","category_description") VALUES ('e4b43a8e-b330-413d-80fa-0dd727bb7bb5', 'TRANSPORTATION', 'TRANSPORTATION_GAS','Gas','Oil and gas');
INSERT INTO "Category" ("category_id", "plaid_category_primary_descriptor", "plaid_category_detailed_descriptor", "category_name","category_description") VALUES ('2d899e52-09dd-4815-b045-d41fe26a94e2', 'RENT_AND_UTILITIES', 'RENT_AND_UTILITIES_TELEPHONE', 'Utilities', 'Utilities, Phone Bill');
INSERT INTO "Category" ("category_id", "plaid_category_primary_descriptor", "plaid_category_detailed_descriptor", "category_name","category_description") VALUES ('c91064d5-ac80-40a2-92db-b0d82ac7f4e5', 'LOAN_PAYMENTS', 'LOAN_PAYMENTS_STUDENT_LOAN_PAYMENT',  'Loan Payment', 'Student Loan Payment');
INSERT INTO "Category" ("category_id", "plaid_category_primary_descriptor", "plaid_category_detailed_descriptor", "category_name","category_description") VALUES ('615b48c8-7b43-42ff-b147-6580219235e6', 'LOAN_PAYMENTS', 'LOAN_PAYMENTS_CREDIT_CARD_PAYMENT',  'Credit Card', 'Credit Card Payment');
INSERT INTO "Category" ("category_id", "plaid_category_primary_descriptor", "plaid_category_detailed_descriptor", "category_name","category_description") VALUES ('ff8b02f5-d11f-4344-89c9-d6880227746a', 'RENT_AND_UTILITIES', 'RENT_AND_UTILITIES_INTERNET_AND_CABLE', 'Internet Service', 'Internet, Services');
INSERT INTO "Category" ("category_id", "plaid_category_primary_descriptor", "plaid_category_detailed_descriptor", "category_name","category_description") VALUES ('fa5a2c7b-3674-4522-8e38-078fa1da0585', 'GENERAL_SERVICES', 'GENERAL_SERVICES_STORAGE',  'Storage Service', 'Storage, Services');
INSERT INTO "Category" ("category_id", "plaid_category_primary_descriptor", "plaid_category_detailed_descriptor", "category_name","category_description") VALUES ('3be2fb20-9cbc-42a5-bcc4-ddb98e93a194', 'PERSONAL_CARE', 'PERSONAL_CARE_GYMS_AND_FITNESS_CENTERS',  'Gym', 'Gyms and Fitness Centers, Recreation');
INSERT INTO "Category" ("category_id", "plaid_category_primary_descriptor", "plaid_category_detailed_descriptor", "category_name","category_description") VALUES ('aac6042d-9026-4dd3-b34e-23eb68598b90', 'ENTERTAINMENT', 'ENTERTAINMENT_OTHER_ENTERTAINMENT', 'Social Clubs', 'Social Clubs, Arts and Entertainment');
INSERT INTO "Category" ("category_id", "plaid_category_primary_descriptor", "plaid_category_detailed_descriptor", "category_name","category_description") VALUES ('bb29d737-319e-4585-864c-a00701b2a231', 'ENTERTAINMENT', 'ENTERTAINMENT_TV_AND_MOVIES', 'Subscription Service', 'Subscription, Services');
INSERT INTO "Category" ("category_id", "plaid_category_primary_descriptor", "plaid_category_detailed_descriptor", "category_name","category_description") VALUES ('30bb0179-90db-4f99-a0bb-45abbbff8efb', 'GENERAL_SERVICES', 'GENERAL_SERVICES_INSURANCE',  'Insurance', 'Insurance, Services');
INSERT INTO "Category" ("category_id", "plaid_category_primary_descriptor", "plaid_category_detailed_descriptor", "category_name","category_description") VALUES ('367e8925-3682-42f0-bbce-beb0eee2fbab', 'MEDICAL', 'MEDICAL_PHARMACIES_AND_SUPPLEMENTS',  'Pharmacy', 'Pharmacy, Health');
INSERT INTO "Category" ("category_id", "plaid_category_primary_descriptor", "plaid_category_detailed_descriptor", "category_name","category_description") VALUES ('fdf97f4c-9739-43b9-acc7-2fb074bfb012', '', '',  'Savings', 'Emergency Funds, Savings');


INSERT INTO "Income" ("income_id", "income_description", "income_amount", "user_id", "income_frequency") VALUES ('8693199c-c2c6-4f75-8138-bda75a186581', 'Salary', 5000, 'ed1bec4c-0a1b-4783-b47f-16ba0650b821', 'Semi-Monthly');



INSERT INTO "Expenses" ("expense_description", "expense_amount", "expense_category", "user_id", "allocation_type") VALUES ('Rent', 1500, '1ae53e57-8b82-45f2-a8cd-94d43932ab54', 'ed1bec4c-0a1b-4783-b47f-16ba0650b821', '9f3c76e9-9d43-4480-a56d-a176b783f24d');
INSERT INTO "Expenses" ("expense_description", "expense_amount", "expense_category", "user_id", "allocation_type") VALUES ('FPL', 100, 'c2f89cdc-5ff7-46fc-92b3-f14bbdec7404', 'ed1bec4c-0a1b-4783-b47f-16ba0650b821', '9f3c76e9-9d43-4480-a56d-a176b783f24d');
INSERT INTO "Expenses" ("expense_description", "expense_amount", "expense_category", "user_id", "allocation_type") VALUES ('Groceries', 500, 'b0332b93-42a8-4eac-936f-ab8d859e9ce3', 'ed1bec4c-0a1b-4783-b47f-16ba0650b821', '9f3c76e9-9d43-4480-a56d-a176b783f24d');
INSERT INTO "Expenses" ("expense_description", "expense_amount", "expense_category", "user_id", "allocation_type") VALUES ('Dining out', 400, 'a9f5ae3e-d44e-49ac-82c5-1bc13c354f58', 'ed1bec4c-0a1b-4783-b47f-16ba0650b821', 'f981f988-5be8-4a9b-bb39-392dd646ddbd');
INSERT INTO "Expenses" ("expense_description", "expense_amount", "expense_category", "user_id", "allocation_type") VALUES ('Vehicle gas', 120, 'e4b43a8e-b330-413d-80fa-0dd727bb7bb5', 'ed1bec4c-0a1b-4783-b47f-16ba0650b821', '9f3c76e9-9d43-4480-a56d-a176b783f24d');
INSERT INTO "Expenses" ("expense_description", "expense_amount", "expense_category", "user_id", "allocation_type") VALUES ('Internet', 80, 'ff8b02f5-d11f-4344-89c9-d6880227746a', 'ed1bec4c-0a1b-4783-b47f-16ba0650b821', '9f3c76e9-9d43-4480-a56d-a176b783f24d');
INSERT INTO "Expenses" ("expense_description", "expense_amount", "expense_category", "user_id", "allocation_type") VALUES ('Phone Bill', 150, '2d899e52-09dd-4815-b045-d41fe26a94e2', 'ed1bec4c-0a1b-4783-b47f-16ba0650b821', '9f3c76e9-9d43-4480-a56d-a176b783f24d');
INSERT INTO "Expenses" ("expense_description", "expense_amount", "expense_category", "user_id", "allocation_type") VALUES ('Storage', 100, 'fa5a2c7b-3674-4522-8e38-078fa1da0585', 'ed1bec4c-0a1b-4783-b47f-16ba0650b821', '9f3c76e9-9d43-4480-a56d-a176b783f24d');
INSERT INTO "Expenses" ("expense_description", "expense_amount", "expense_category", "user_id", "allocation_type") VALUES ('Gym', 50, '3be2fb20-9cbc-42a5-bcc4-ddb98e93a194', 'ed1bec4c-0a1b-4783-b47f-16ba0650b821', 'f981f988-5be8-4a9b-bb39-392dd646ddbd');
INSERT INTO "Expenses" ("expense_description", "expense_amount", "expense_category", "user_id", "allocation_type") VALUES ('TV subscriptions', 30, 'bb29d737-319e-4585-864c-a00701b2a231', 'ed1bec4c-0a1b-4783-b47f-16ba0650b821', 'f981f988-5be8-4a9b-bb39-392dd646ddbd');
INSERT INTO "Expenses" ("expense_description", "expense_amount", "expense_category", "user_id", "allocation_type") VALUES ('Car insurance', 120, '30bb0179-90db-4f99-a0bb-45abbbff8efb', 'ed1bec4c-0a1b-4783-b47f-16ba0650b821', '9f3c76e9-9d43-4480-a56d-a176b783f24d');
INSERT INTO "Expenses" ("expense_description", "expense_amount", "expense_category", "user_id", "allocation_type") VALUES ('Credit card payment', 50, '615b48c8-7b43-42ff-b147-6580219235e6', 'ed1bec4c-0a1b-4783-b47f-16ba0650b821', 'ac184cdf-b7ff-4eb9-b757-628770d566fb');
INSERT INTO "Expenses" ("expense_description", "expense_amount", "expense_category", "user_id", "allocation_type") VALUES ('Student loan payment', 350, 'c91064d5-ac80-40a2-92db-b0d82ac7f4e5', 'ed1bec4c-0a1b-4783-b47f-16ba0650b821', 'ac184cdf-b7ff-4eb9-b757-628770d566fb');















# How to check permission by bitmap

## Scenario background
### For example, we have 3 services,
sales service(current service):
  - order feature
  - shop feature
  - headquarters management
  - branch shop management
  - delivery feature

notification service(external service)
  - email notification feature

tracker service(external service):
  - sales trackers
  - order tracker
  - shop tracker
  - delivery tracker
  - notification tracker
  - email notification tracker

### Related features across services
here are 2 related features:
- order feature in sales service(current service)
- order tracker feature in tracker service(external service)
and we define the permission of users to order feature in sales service like this:
  - can read orders:
    (users have permssion of sales service && users have permission of order feature in sales service)
    || (users have permission of tracker service && users have permission of order tracker feature in tracker service)
  - can write orders:
    users have permssion of sales service && have permission of order feature in sales service

## Assign permission index
let go back to the 3 services, and assign an index for each item in the services.

sales service(current service):                     index: 0
  - order feature                                   index: 1
  - shop feature                                    index: 2
  - headquarters management                         index: 3
  - branch shop management                          index: 4
  - delivery feature                                index: 5

notification service(external service):             index: 0
  - email notification feature                      index: 1

tracker service(external service):                  index: 0
  - sales trackers                                  index: 1
  - order tracker                                   index: 2
  - shop tracker                                    index: 3
  - delivery tracker                                index: 4
  - notification tracker                            index: 5
  - email notification tracker                      index: 6

## Calculate bitmap by granted permission indicies
Considering that notification service is not related to current service, we can ignore it.
then we design a 8-bit bitmap to store if users have the permission of each item in a service.

- scenario1:
  - desc: users have permssion of sales service && users have permission of order feature in sales service
  - bitmap: 0b11000000

- scenario2:
  - desc: users have permission of tracker service && users have permission of order tracker feature under sales tracker in tracker service
  - bitmap: 0b11100000

then the read/write order permission check comes to be:
  - read: (user_sales_permission_bit & required_sales_permission_bit(0b11000000) >= required_sales_permission_bit(0b11000000))
    || (user_tracker_permission_bit & required_tracker_permission_bit(0b11100000) >= required_tracker_permission_bit(0b11100000))
  - write: user_sales_permission_bit & required_sales_permission_bit(0b11000000) >= required_sales_permission_bit(0b11000000)

## Check permission by bitmap
#### permission check example 1:
user has permissions of order feature and delivery feature in sales service, then user_sales_permission_bit is 0b11000100
0b11000100 & 0b11000000 = 0b11000000 = 192

#### permission check example 2:
user has permissions of headquarters management feature and delivery feature in sales service, then user_sales_permission_bit is 0b10110100
0b10110100 & 0b11000000 = 0b10000100
0b10000100(132) < 0b11000000(192)

#### permission check that across services:
in our case, we need to check 2 services, then we can set users permission bit like 0b1100010011000000.
first 8 bits are permissions of sales service, last 8 bits are permissions of tracker service
final example:
  - read permission check: user_permission_bit > 0b1100000000000000 || (user_permission_bit & 0b0000000011111111) > 0b0000000011100000
  - write permission check: user_permission_bit > 0b1100000000000000

as we have uint64 type, where we can store 64bit, so we can concat 8 services together(8 bits per service)

s1      s2      s3      s4      s5      s6      s7      s8
1000000010000000100000001000000010000000100000001000000010000000